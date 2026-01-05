/**
 * Main entry point for the Aether Vault Node.js SDK.
 * Provides the createVaultClient function and exports all SDK components.
 */

import { VaultConfig } from "./core/config.js";
import { VaultClient } from "./core/client.js";
import { AuthClient } from "./auth/auth.client.js";
import { SecretsClient } from "./secrets/secrets.client.js";
import { TotpClient } from "./totp/totp.client.js";
import { IdentityClient } from "./identity/identity.client.js";
import { AuditClient } from "./audit/audit.client.js";
import { PolicyClient } from "./policies/policy.client.js";
import { SystemClient } from "./system/system.client.js";

/**
 * Main Aether Vault client class.
 * Provides access to all SDK functionality through organized domain clients.
 */
export class AetherVaultClient {
  private readonly httpClient: VaultClient;

  /**
   * Authentication client for handling authentication operations.
   */
  public readonly auth: AuthClient;

  /**
   * Secrets client for managing secrets.
   */
  public readonly secrets: SecretsClient;

  /**
   * TOTP client for managing Time-based One-Time Passwords.
   */
  public readonly totp: TotpClient;

  /**
   * Identity client for managing user identity operations.
   */
  public readonly identity: IdentityClient;

  /**
   * Audit client for accessing audit logs.
   */
  public readonly audit: AuditClient;

  /**
   * Policy client for managing access policies.
   */
  public readonly policies: PolicyClient;

  /**
   * System client for health checks and system info.
   */
  public readonly system: SystemClient;

  /**
   * Creates a new AetherVaultClient instance.
   *
   * @param config - Client configuration
   */
  constructor(config: VaultConfig) {
    this.httpClient = new VaultClient(config);

    // Initialize domain clients
    this.auth = new AuthClient(this.httpClient, config.auth);
    this.secrets = new SecretsClient(this.httpClient);
    this.totp = new TotpClient(this.httpClient);
    this.identity = new IdentityClient(this.httpClient);
    this.audit = new AuditClient(this.httpClient);
    this.policies = new PolicyClient(this.httpClient);
    this.system = new SystemClient(this.httpClient);
  }

  /**
   * Updates the authentication token.
   *
   * @param token - New authentication token
   */
  public updateToken(token: string): void {
    this.httpClient.updateToken(token);
  }

  /**
   * Gets the current client configuration.
   *
   * @returns Current client configuration (read-only)
   */
  public getConfig() {
    return this.httpClient.getConfig();
  }
}

/**
 * Creates and returns a configured Aether Vault client.
 * This is the recommended way to initialize the SDK.
 *
 * @param config - Client configuration object
 * @returns Configured AetherVaultClient instance
 *
 * @example
 * ```typescript
 * import { createVaultClient } from "aether-vault";
 *
 * const vault = createVaultClient({
 *   baseURL: "/api/v1",
 *   auth: {
 *     type: "session",
 *   },
 * });
 *
 * await vault.secrets.list();
 * await vault.secrets.get("DATABASE_URL");
 * await vault.totp.generate("github");
 * ```
 */
export function createVaultClient(config: VaultConfig): AetherVaultClient {
  return new AetherVaultClient(config);
}

// Re-export types and classes for convenience
export {
  VaultConfig,
  AuthConfig,
  JwtAuthConfig,
  SessionAuthConfig,
} from "./core/config.js";
export {
  VaultError,
  VaultAuthError,
  VaultPermissionError,
  VaultNotFoundError,
  VaultServerError,
  VaultNetworkError,
} from "./core/errors.js";
export {
  Service,
  ServiceType,
  HealthState,
  ServiceRegistrationRequest,
  ServiceRegistrationResponse,
  ServiceListResponse,
} from "./types/index.js";
export { AuthClient } from "./auth/auth.client.js";
export { SecretsClient } from "./secrets/secrets.client.js";
export { TotpClient } from "./totp/totp.client.js";
export { IdentityClient } from "./identity/identity.client.js";
export { AuditClient } from "./audit/audit.client.js";
export { PolicyClient } from "./policies/policy.client.js";
export { SystemClient } from "./system/system.client.js";
