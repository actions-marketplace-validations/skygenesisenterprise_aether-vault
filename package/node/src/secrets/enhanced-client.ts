import {
  SecurityLevel,
  SecretMetadata,
  EnhancedSecret,
  CreateEnhancedSecretRequest,
} from "./security-levels.js";
import { SecretClassifier } from "./classifier.js";
import {
  PolicyEngine,
  AccessContext,
  AccessDecision,
} from "./policy-engine.js";
import { MultiLevelEncryption, EncryptionMetadata } from "./encryption.js";
import { VaultClient } from "../core/client.js";

/**
 * Enhanced secrets client with fine-grained management.
 */
export class EnhancedSecretsClient {
  private policyEngine: PolicyEngine;
  private encryption: MultiLevelEncryption;

  constructor(private readonly httpClient: VaultClient) {
    this.policyEngine = new PolicyEngine();
    this.encryption = new MultiLevelEncryption();
  }

  /**
   * Creates a new enhanced secret with automatic classification.
   */
  public async createEnhanced(
    secret: CreateEnhancedSecretRequest,
  ): Promise<EnhancedSecret> {
    // Auto-classify the secret
    const classifiedMetadata = SecretClassifier.classify(secret);

    // Create encryption metadata
    const { encrypted, metadata: encryptionMetadata } =
      await this.encryption.encrypt(
        secret.value,
        secret.name,
        classifiedMetadata.securityLevel,
      );

    // Create the secret request
    const createRequest = {
      name: secret.name,
      value: encrypted, // Store encrypted value
      description: secret.description,
      tags: secret.tags,
      expiresAt: secret.expiresAt,
      metadata: {
        ...classifiedMetadata,
        encryption: encryptionMetadata,
        version: {
          current: 1,
          total: 1,
          lastRotated: new Date().toISOString(),
          nextRotation: this.calculateNextRotation(classifiedMetadata),
        },
        accessStats: {
          totalAccess: 0,
          uniqueUsers: 0,
          failedAttempts: 0,
        },
        status: "active" as const,
        checksum: this.calculateChecksum(encrypted, encryptionMetadata),
      },
    };

    // Call the API
    const response = await this.httpClient.post<EnhancedSecret>(
      "/secrets/enhanced",
      createRequest,
    );

    return response;
  }

  /**
   * Gets an enhanced secret with access control verification.
   */
  public async getEnhanced(
    id: string,
    context: AccessContext,
    includeValue: boolean = false,
  ): Promise<EnhancedSecret> {
    // Get the secret metadata first (without value)
    const secret = await this.httpClient.get<EnhancedSecret>(
      `/secrets/enhanced/${id}`,
      {
        includeValue: "false",
      },
    );

    // Check access permissions
    const accessDecision = this.policyEngine.evaluateAccess(
      id,
      secret.metadata,
      context,
    );

    if (!accessDecision.granted) {
      throw new Error(`Access denied: ${accessDecision.reason}`);
    }

    // Log the access attempt
    await this.logAccess(id, context, accessDecision.granted);

    if (includeValue) {
      // Get the encrypted value
      const encryptedSecret = await this.httpClient.get<{ value: string }>(
        `/secrets/${id}/value`,
      );

      // Decrypt the value
      const decryptedValue = await this.encryption.decrypt(
        encryptedSecret.value,
        secret.metadata.encryption as any,
      );

      secret.value = decryptedValue;
    }

    // Update access statistics
    await this.updateAccessStats(id, context.userId);

    return secret;
  }

  /**
   * Lists enhanced secrets with filtering and access control.
   */
  public async listEnhanced(
    context: AccessContext,
    params?: {
      page?: number;
      pageSize?: number;
      securityLevel?: SecurityLevel[];
      category?: string[];
      ownerId?: string[];
    },
  ): Promise<{
    secrets: EnhancedSecret[];
    total: number;
    page: number;
    pageSize: number;
  }> {
    // Get all secrets that match the filters
    const allSecrets = await this.httpClient.get<{
      secrets: EnhancedSecret[];
      total: number;
    }>("/secrets/enhanced", params);

    // Filter secrets based on access control
    const accessibleSecrets = allSecrets.secrets.filter((secret) => {
      const decision = this.policyEngine.evaluateAccess(
        secret.id,
        secret.metadata,
        context,
      );
      return decision.granted;
    });

    return {
      secrets: accessibleSecrets,
      total: accessibleSecrets.length,
      page: params?.page || 1,
      pageSize: params?.pageSize || 20,
    };
  }

  /**
   * Updates an enhanced secret with access control.
   */
  public async updateEnhanced(
    id: string,
    updates: Partial<{
      value: string;
      description: string;
      tags: string[];
      expiresAt: string;
      status: EnhancedSecret["status"];
      metadata: Partial<SecretMetadata>;
    }>,
    context: AccessContext,
  ): Promise<EnhancedSecret> {
    // Get current secret
    const currentSecret = await this.getEnhanced(id, context, false);

    // Check update permissions
    const updateDecision = this.policyEngine.evaluateAccess(
      id,
      currentSecret.metadata,
      {
        ...context,
        requestMethod: "PUT",
      },
    );

    if (!updateDecision.granted) {
      throw new Error(`Update denied: ${updateDecision.reason}`);
    }

    let encryptedValue: string | undefined;
    let encryptionMetadata: EncryptionMetadata | undefined;

    // Handle value update with re-encryption
    if (updates.value) {
      const { encrypted, metadata } = await this.encryption.encrypt(
        updates.value,
        id,
        currentSecret.metadata.securityLevel,
      );
      encryptedValue = encrypted;
      encryptionMetadata = metadata;
    }

    // Re-classify if metadata is being updated
    let updatedMetadata = updates.metadata;
    if (updates.metadata && Object.keys(updates.metadata).length > 0) {
      const fullUpdateRequest: CreateEnhancedSecretRequest = {
        name: currentSecret.name,
        value: updates.value || currentSecret.value || "",
        metadata: { ...currentSecret.metadata, ...updates.metadata },
      };

      updatedMetadata = SecretClassifier.classify(fullUpdateRequest);
    }

    // Prepare update request
    const updateRequest = {
      ...updates,
      value: encryptedValue,
      metadata: {
        ...currentSecret.metadata,
        ...updatedMetadata,
        ...(encryptionMetadata && { encryption: encryptionMetadata }),
        ...(updates.value && {
          version: {
            ...currentSecret.version,
            current: currentSecret.version.current + 1,
            total: currentSecret.version.total + 1,
            lastRotated: new Date().toISOString(),
          },
        }),
      },
    };

    const response = await this.httpClient.put<EnhancedSecret>(
      `/secrets/enhanced/${id}`,
      updateRequest,
    );

    return response;
  }

  /**
   * Rotates a secret value.
   */
  public async rotate(
    id: string,
    newValue?: string,
    context?: AccessContext,
  ): Promise<EnhancedSecret> {
    const currentSecret = await this.getEnhanced(
      id,
      context || {
        userId: "system",
        userRoles: ["system"],
        userAttributes: {},
        timestamp: new Date().toISOString(),
        sourceIp: "127.0.0.1",
        requestMethod: "POST",
        resource: `/secrets/${id}/rotate`,
      },
    );

    // Check rotation permissions
    const rotateDecision = this.policyEngine.evaluateAccess(
      id,
      currentSecret.metadata,
      {
        ...(context || ({} as AccessContext)),
        requestMethod: "POST",
        resource: `/secrets/${id}/rotate`,
      },
    );

    if (!rotateDecision.granted) {
      throw new Error(`Rotation denied: ${rotateDecision.reason}`);
    }

    // Generate new value if not provided
    const valueToRotate = newValue || this.generateNewValue(currentSecret);

    // Rotate the encryption key
    await this.encryption.rotateKey(id);

    // Encrypt with new key
    const { encrypted, metadata: encryptionMetadata } =
      await this.encryption.encrypt(
        valueToRotate,
        id,
        currentSecret.metadata.securityLevel,
      );

    const updateRequest = {
      value: encrypted,
      metadata: {
        ...currentSecret.metadata,
        encryption: encryptionMetadata,
        version: {
          current: currentSecret.version.current + 1,
          total: currentSecret.version.total + 1,
          lastRotated: new Date().toISOString(),
          nextRotation: this.calculateNextRotation(currentSecret.metadata),
        },
      },
    };

    const response = await this.httpClient.post<EnhancedSecret>(
      `/secrets/${id}/rotate`,
      updateRequest,
    );

    return response;
  }

  /**
   * Archives an enhanced secret.
   */
  public async archive(
    id: string,
    context: AccessContext,
  ): Promise<EnhancedSecret> {
    const secret = await this.getEnhanced(id, context, false);

    const archiveDecision = this.policyEngine.evaluateAccess(
      id,
      secret.metadata,
      {
        ...context,
        requestMethod: "POST",
        resource: `/secrets/${id}/archive`,
      },
    );

    if (!archiveDecision.granted) {
      throw new Error(`Archive denied: ${archiveDecision.reason}`);
    }

    const response = await this.httpClient.post<EnhancedSecret>(
      `/secrets/enhanced/${id}/archive`,
    );

    return response;
  }

  /**
   * Sets up an access policy for secrets.
   */
  public setPolicy(policy: any): void {
    this.policyEngine.addPolicy(policy);
  }

  /**
   * Gets access decision for a secret.
   */
  public checkAccess(
    id: string,
    secretMetadata: SecretMetadata,
    context: AccessContext,
  ): AccessDecision {
    return this.policyEngine.evaluateAccess(id, secretMetadata, context);
  }

  /**
   * Calculates next rotation date based on rotation policy.
   */
  private calculateNextRotation(metadata: SecretMetadata): string {
    if (!metadata.rotationPolicy.enabled) {
      return "";
    }

    const intervalDays = metadata.rotationPolicy.intervalDays || 90;
    const nextRotation = new Date();
    nextRotation.setDate(nextRotation.getDate() + intervalDays);

    return nextRotation.toISOString();
  }

  /**
   * Calculates checksum for integrity verification.
   */
  private calculateChecksum(
    value: string,
    encryption: EncryptionMetadata,
  ): string {
    const crypto = require("crypto");
    return crypto
      .createHash("sha256")
      .update(value)
      .update(encryption.checksum)
      .digest("hex");
  }

  /**
   * Logs access attempt for audit purposes.
   */
  private async logAccess(
    secretId: string,
    context: AccessContext,
    granted: boolean,
  ): Promise<void> {
    const logEntry = {
      secretId,
      userId: context.userId,
      sourceIp: context.sourceIp,
      timestamp: context.timestamp,
      granted,
      userAgent: context.userAgent,
      deviceId: context.deviceId,
    };

    await this.httpClient.post("/audit/secret-access", logEntry);
  }

  /**
   * Updates access statistics for a secret.
   */
  private async updateAccessStats(
    secretId: string,
    userId: string,
  ): Promise<void> {
    await this.httpClient.post(`/secrets/${secretId}/stats/access`, { userId });
  }

  /**
   * Generates a new value for automatic rotation.
   */
  private generateNewValue(secret: EnhancedSecret): string {
    const crypto = require("crypto");

    // Generate based on category and current value length
    const length = secret.value ? Math.max(32, secret.value.length) : 32;

    if (secret.metadata.category === "api_key") {
      return crypto
        .randomBytes(length)
        .toString("base64")
        .replace(/[^a-zA-Z0-9]/g, "")
        .substring(0, length);
    }

    if (secret.metadata.category === "encryption_key") {
      return crypto.randomBytes(length / 2).toString("hex");
    }

    // Default: alphanumeric with special characters
    const chars =
      "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*";
    let result = "";
    for (let i = 0; i < length; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
  }
}
