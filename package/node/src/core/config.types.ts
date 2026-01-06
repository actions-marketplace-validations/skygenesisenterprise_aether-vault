/**
 * Extended configuration interfaces for Aether Vault SDK.
 * Provides comprehensive configuration options including environment-based overrides.
 */

/**
 * Environment names for multi-environment configuration.
 */
export type Environment = "development" | "production" | "staging" | "test";

/**
 * Authentication methods supported by Aether Vault.
 */
export type AuthMethod =
  | "token"
  | "app-role"
  | "oidc"
  | "jwt"
  | "bearer"
  | "session"
  | "none";

/**
 * Logging levels for the SDK.
 */
export type LogLevel = "error" | "warn" | "info" | "debug" | "trace";

/**
 * Retry configuration options.
 */
export interface RetryConfig {
  /** Number of retry attempts */
  retries: number;
  /** Delay between retries in milliseconds */
  delay: number;
  /** Enable exponential backoff */
  backoff?: boolean;
  /** Maximum delay for exponential backoff */
  maxDelay?: number;
}

/**
 * Logging configuration.
 */
export interface LoggingConfig {
  /** Logging level */
  level: LogLevel;
  /** Enable console logging */
  console?: boolean;
  /** Custom logger function */
  logger?: (level: LogLevel, message: string, meta?: any) => void;
  /** Enable request/response logging */
  http?: boolean;
}

/**
 * Feature flags for future functionality.
 */
export interface FeaturesConfig {
  /** Enable automatic token renewal */
  autoRenewToken?: boolean;
  /** Enable audit logging */
  auditEnabled?: boolean;
  /** Enable metrics collection */
  metricsEnabled?: boolean;
  /** Enable caching */
  cachingEnabled?: boolean;
  /** Enable request tracing */
  tracingEnabled?: boolean;
}

/**
 * Token-based authentication configuration.
 */
export interface TokenAuthConfig {
  method: "token";
  /** Vault token */
  token: string;
}

/**
 * AppRole authentication configuration.
 */
export interface AppRoleAuthConfig {
  method: "app-role";
  /** Role ID */
  roleId: string;
  /** Secret ID */
  secretId: string;
}

/**
 * OIDC authentication configuration.
 */
export interface OidcAuthConfig {
  method: "oidc";
  /** OIDC token */
  token: string;
  /** OIDC role */
  role?: string;
}

/**
 * Extended authentication configuration supporting all Vault auth methods.
 */
export type ExtendedAuthConfig =
  | TokenAuthConfig
  | AppRoleAuthConfig
  | OidcAuthConfig
  | { method: "jwt"; token: string; issuer?: string; audience?: string }
  | { method: "bearer"; token: string }
  | { method: "session"; cookieName?: string }
  | { method: "none" };

/**
 * Complete Vault configuration with all options.
 */
export interface CompleteVaultConfig {
  /** Vault endpoint URL */
  endpoint: string;
  /** Vault namespace */
  namespace?: string;
  /** Authentication configuration */
  auth: ExtendedAuthConfig;
  /** Retry configuration */
  retry: RetryConfig;
  /** Logging configuration */
  logging: LoggingConfig;
  /** Feature flags */
  features: FeaturesConfig;
  /** Request timeout in milliseconds */
  timeout?: number;
  /** Custom headers */
  headers?: Record<string, string>;
  /** API version */
  apiVersion?: string;
}

/**
 * Environment-specific configuration.
 */
export interface EnvironmentConfig {
  /** Development environment configuration */
  development?: Partial<CompleteVaultConfig>;
  /** Production environment configuration */
  production?: Partial<CompleteVaultConfig>;
  /** Staging environment configuration */
  staging?: Partial<CompleteVaultConfig>;
  /** Test environment configuration */
  test?: Partial<CompleteVaultConfig>;
}

/**
 * Configuration file structure.
 */
export interface VaultConfigFile {
  /** Default configuration */
  default: CompleteVaultConfig;
  /** Environment-specific configurations */
  environments?: EnvironmentConfig;
}

/**
 * Environment variable mapping.
 */
export interface EnvironmentVariables {
  /** Vault endpoint */
  VAULT_ENDPOINT?: string;
  /** Vault namespace */
  VAULT_NAMESPACE?: string;
  /** Auth method */
  VAULT_AUTH_METHOD?: string;
  /** Vault token */
  VAULT_TOKEN?: string;
  /** AppRole ID */
  VAULT_ROLE_ID?: string;
  /** AppRole Secret ID */
  VAULT_SECRET_ID?: string;
  /** OIDC token */
  VAULT_OIDC_TOKEN?: string;
  /** OIDC role */
  VAULT_OIDC_ROLE?: string;
  /** Environment */
  VAULT_ENV?: string;
  /** Log level */
  VAULT_LOG_LEVEL?: string;
  /** Retry attempts */
  VAULT_RETRY_RETRIES?: string;
  /** Retry delay */
  VAULT_RETRY_DELAY?: string;
  /** Timeout */
  VAULT_TIMEOUT?: string;
  /** API version */
  VAULT_API_VERSION?: string;
  /** Enable debug */
  VAULT_DEBUG?: string;
  /** Enable auto token renewal */
  VAULT_AUTO_RENEW_TOKEN?: string;
  /** Enable audit */
  VAULT_AUDIT_ENABLED?: string;
}

/**
 * Configuration loader options.
 */
export interface ConfigLoaderOptions {
  /** Configuration file path */
  configPath?: string;
  /** Environment to use */
  environment?: Environment;
  /** Enable environment variable overrides */
  enableEnvOverrides?: boolean;
  /** Strict mode (throw on missing required fields) */
  strict?: boolean;
}

/**
 * Loaded configuration result.
 */
export interface LoadedConfig {
  /** Final configuration */
  config: CompleteVaultConfig;
  /** Source of configuration values */
  sources: {
    file?: string;
    environment?: string;
    envVars?: string[];
  };
}
