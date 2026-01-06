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
import { Environment } from "./core/config.types.js";

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
 * Configuration options for createVaultClient.
 */
export interface CreateVaultClientOptions {
  /** Configuration file path (default: "vault.config.ts") */
  configPath?: string;
  /** Environment to use (default: from NODE_ENV or VAULT_ENV) */
  environment?: Environment;
  /** Enable environment variable overrides (default: true) */
  enableEnvOverrides?: boolean;
  /** Strict validation mode (default: false) */
  strict?: boolean;
  /** Override configuration (takes precedence over file/env) */
  config?: Partial<VaultConfig>;
}

/**
 * Creates and returns a configured Aether Vault client.
 * This is the recommended way to initialize the SDK.
 *
 * Automatically loads configuration from vault.config.ts if present,
 * with fallback to environment variables.
 *
 * @param optionsOrConfig - Either configuration options for auto-loading or a VaultConfig object
 * @returns Configured AetherVaultClient instance
 *
 * @example
 * ```typescript
 * import { createVaultClient } from "aether-vault";
 *
 * // Auto-load from vault.config.ts or environment variables
 * const vault = createVaultClient();
 *
 * // Or specify custom options
 * const vault = createVaultClient({
 *   environment: "production",
 *   configPath: "./config/vault.config.ts"
 * });
 *
 * // Or provide explicit config (legacy behavior)
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
export async function createVaultClient(
  optionsOrConfig?: CreateVaultClientOptions | VaultConfig,
): Promise<AetherVaultClient> {
  // Check if this is a VaultConfig (legacy behavior) or options object
  const isLegacyConfig =
    optionsOrConfig &&
    "baseURL" in optionsOrConfig &&
    "auth" in optionsOrConfig;

  if (isLegacyConfig) {
    // Legacy behavior: use provided config directly
    return new AetherVaultClient(optionsOrConfig as VaultConfig);
  }

  // New behavior: auto-load configuration
  const options = (optionsOrConfig as CreateVaultClientOptions) || {};

  try {
    // Import the config loader dynamically to avoid circular dependencies
    const { loadVaultConfig } = await import("./core/config.loader.js");

    // Load configuration from file and/or environment
    const loaderOptions: any = {};
    if (options.configPath) loaderOptions.configPath = options.configPath;
    if (options.environment) loaderOptions.environment = options.environment;
    if (options.enableEnvOverrides !== undefined)
      loaderOptions.enableEnvOverrides = options.enableEnvOverrides;
    if (options.strict !== undefined) loaderOptions.strict = options.strict;

    let config = await loadVaultConfig(loaderOptions);

    // Apply any explicit overrides
    if (options.config) {
      config = { ...config, ...options.config };
    }

    return new AetherVaultClient(config);
  } catch (error) {
    throw new Error(
      `Failed to create Aether Vault client: ${error instanceof Error ? error.message : String(error)}`,
    );
  }
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
