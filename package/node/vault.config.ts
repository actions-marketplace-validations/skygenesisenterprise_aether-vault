/**
 * Aether Vault Configuration Template
 *
 * This is the main configuration file for Aether Vault SDK.
 * Copy this file to your project root as vault.config.ts and customize as needed.
 *
 * The SDK will automatically load this configuration when using createVaultClient().
 */

import {
  VaultConfigFile,
  CompleteVaultConfig,
  Environment,
  AuthMethod,
  LogLevel,
  ExtendedAuthConfig,
} from "./src/core/config.types.js";

/**
 * Helper function to create authentication configuration based on environment variables.
 */
function getAuthConfig(): ExtendedAuthConfig {
  const authMethod = (process.env.VAULT_AUTH_METHOD as AuthMethod) || "token";

  switch (authMethod) {
    case "token":
      return {
        method: "token",
        token: process.env.VAULT_TOKEN || "your-vault-token-here",
      };

    case "app-role":
      return {
        method: "app-role",
        roleId: process.env.VAULT_ROLE_ID || "your-role-id",
        secretId: process.env.VAULT_SECRET_ID || "your-secret-id",
      };

    case "oidc":
      return {
        method: "oidc",
        token: process.env.VAULT_OIDC_TOKEN || "your-oidc-token",
        role: process.env.VAULT_OIDC_ROLE || undefined,
      };

    case "jwt":
      return {
        method: "jwt",
        token: process.env.VAULT_TOKEN || "your-jwt-token",
        issuer: process.env.VAULT_JWT_ISSUER,
        audience: process.env.VAULT_JWT_AUDIENCE,
      };

    case "bearer":
      return {
        method: "bearer",
        token: process.env.VAULT_TOKEN || "your-bearer-token",
      };

    case "session":
      return {
        method: "session",
        cookieName: process.env.VAULT_SESSION_COOKIE,
      };

    case "none":
      return {
        method: "none",
      };

    default:
      throw new Error(`Unsupported auth method: ${authMethod}`);
  }
}

/**
 * Default Vault configuration.
 * This serves as the base configuration and will be overridden by environment-specific settings.
 */
const defaultConfig: CompleteVaultConfig = {
  // Vault endpoint - required
  endpoint:
    process.env.VAULT_ENDPOINT ||
    "https://vault.skygenesisenterprise.com/api/v1",

  // Optional namespace for multi-tenant setups
  namespace: process.env.VAULT_NAMESPACE || undefined,

  // Authentication configuration
  auth: getAuthConfig(),

  // Retry configuration
  retry: {
    retries: parseInt(process.env.VAULT_RETRY_RETRIES || "3"),
    delay: parseInt(process.env.VAULT_RETRY_DELAY || "1000"),
    backoff: true,
    maxDelay: 10000,
  },

  // Logging configuration
  logging: {
    level: (process.env.VAULT_LOG_LEVEL as LogLevel) || "info",
    console: true,
    http: process.env.VAULT_DEBUG === "true",
  },

  // Feature flags
  features: {
    autoRenewToken: process.env.VAULT_AUTO_RENEW_TOKEN === "true",
    auditEnabled: process.env.VAULT_AUDIT_ENABLED === "true",
    metricsEnabled: false,
    cachingEnabled: true,
    tracingEnabled: false,
  },

  // Request timeout
  timeout: parseInt(process.env.VAULT_TIMEOUT || "30000"),

  // API version
  apiVersion: process.env.VAULT_API_VERSION || "v1",
};

/**
 * Environment-specific configurations.
 * These will override the default config based on the current environment.
 */
const environments: Record<Environment, Partial<CompleteVaultConfig>> = {
  development: {
    endpoint: "http://localhost:8200/api/v1",
    auth: {
      method: "token",
      token: "dev-token",
    },
    logging: {
      level: "debug",
      console: true,
      http: true,
    },
    retry: {
      retries: 1,
      delay: 500,
    },
    features: {
      auditEnabled: false,
      tracingEnabled: true,
    },
  },

  production: {
    endpoint: "https://vault.production.com/api/v1",
    namespace: "prod",
    logging: {
      level: "warn",
      console: false,
      http: false,
    },
    retry: {
      retries: 5,
      delay: 2000,
      backoff: true,
      maxDelay: 30000,
    },
    features: {
      autoRenewToken: true,
      auditEnabled: true,
      metricsEnabled: true,
      cachingEnabled: true,
      tracingEnabled: false,
    },
    timeout: 60000,
  },

  staging: {
    endpoint: "https://vault.staging.com/api/v1",
    namespace: "staging",
    logging: {
      level: "info",
      console: true,
      http: false,
    },
    retry: {
      retries: 3,
      delay: 1500,
    },
    features: {
      autoRenewToken: true,
      auditEnabled: true,
      metricsEnabled: false,
      cachingEnabled: true,
      tracingEnabled: true,
    },
  },

  test: {
    endpoint: "http://localhost:8200/api/v1",
    auth: {
      method: "token",
      token: "test-token",
    },
    logging: {
      level: "error",
      console: false,
      http: false,
    },
    retry: {
      retries: 0,
      delay: 0,
    },
    features: {
      auditEnabled: false,
      metricsEnabled: false,
      cachingEnabled: false,
      tracingEnabled: false,
    },
    timeout: 5000,
  },
};

/**
 * Complete configuration file export.
 * This is what the SDK will load when importing the configuration.
 */
export const vaultConfig: VaultConfigFile = {
  default: defaultConfig,
  environments,
};

/**
 * Helper function to get configuration for a specific environment.
 * Useful for testing or manual configuration loading.
 */
export function getConfigForEnvironment(env: Environment): CompleteVaultConfig {
  const envConfig = environments[env] || {};
  return {
    ...defaultConfig,
    ...envConfig,
    // Deep merge for nested objects
    auth: { ...defaultConfig.auth, ...envConfig.auth },
    retry: { ...defaultConfig.retry, ...envConfig.retry },
    logging: { ...defaultConfig.logging, ...envConfig.logging },
    features: { ...defaultConfig.features, ...envConfig.features },
  };
}

/**
 * Get current environment from process.env.NODE_ENV or VAULT_ENV.
 */
export function getCurrentEnvironment(): Environment {
  const env = process.env.VAULT_ENV || process.env.NODE_ENV || "development";
  switch (env) {
    case "production":
    case "prod":
      return "production";
    case "staging":
    case "stage":
      return "staging";
    case "test":
      return "test";
    default:
      return "development";
  }
}

/**
 * Get configuration for the current environment.
 */
export function getCurrentConfig(): CompleteVaultConfig {
  return getConfigForEnvironment(getCurrentEnvironment());
}

// Export current config as default for easy importing
export default getCurrentConfig();
