/**
 * Aether Vault Node.js SDK - Usage Examples
 *
 * This file demonstrates how to use the Aether Vault SDK
 * without exposing any API paths to the end user.
 */

import { createVaultClient } from "../src/index.js";

/**
 * Complete authentication flow example
 */
async function authenticationExample() {
  // Initialize vault client
  const vault = createVaultClient({
    baseURL: "http://localhost:8080",
    auth: {
      type: "session",
    },
  });

  try {
    // Login user
    const session = await vault.auth.login({
      username: "user@example.com",
      password: "securePassword123",
    });

    console.log("✅ User logged in successfully");
    console.log("User:", session.user.firstName, session.user.lastName);
    console.log("Token expires:", session.expiresAt);

    // Check current session
    const currentSession = await vault.auth.session();
    console.log("✅ Session valid:", currentSession.valid);

    // Get current user information
    const user = await vault.identity.me();
    console.log("✅ Current user:", user.email, user.displayName);

    // Logout user
    await vault.auth.logout();
    console.log("✅ User logged out successfully");
  } catch (error) {
    console.error("❌ Authentication failed:", error);
  }
}

/**
 * Secrets management example
 */
async function secretsExample() {
  const vault = createVaultClient({
    baseURL: "http://localhost:8080",
    auth: {
      type: "session",
    },
  });

  // Login first
  await vault.auth.login({
    username: "user@example.com",
    password: "securePassword123",
  });

  try {
    // Create a new secret
    const secret = await vault.secrets.create({
      name: "Database Connection",
      description: "Production database connection string",
      value: "postgresql://user:pass@localhost:5432/mydb",
      type: "database",
      tags: "production,database",
      expiresAt: new Date("2025-12-31"),
    });

    console.log("✅ Secret created:", secret.id, secret.name);

    // List all secrets
    const secretsList = await vault.secrets.list();
    console.log("✅ Found", secretsList.total, "secrets");

    // Get specific secret
    const retrievedSecret = await vault.secrets.get(secret.id);
    console.log("✅ Retrieved secret:", retrievedSecret.name);

    // Update secret
    const updatedSecret = await vault.secrets.update(secret.id, {
      description: "Updated database connection string",
      tags: "production,database,updated",
    });

    console.log("✅ Secret updated:", updatedSecret.description);

    // Delete secret
    await vault.secrets.delete(secret.id);
    console.log("✅ Secret deleted successfully");
  } catch (error) {
    console.error("❌ Secrets operation failed:", error);
  }
}

/**
 * TOTP management example
 */
async function totpExample() {
  const vault = createVaultClient({
    baseURL: "http://localhost:8080",
    auth: {
      type: "session",
    },
  });

  // Login first
  await vault.auth.login({
    username: "user@example.com",
    password: "securePassword123",
  });

  try {
    // Create TOTP entry
    const totp = await vault.totp.create({
      name: "GitHub 2FA",
      description: "Two-factor authentication for GitHub",
      algorithm: "SHA1",
      digits: 6,
      period: 30,
    });

    console.log("✅ TOTP created:", totp.id, totp.name);

    // List all TOTP entries
    const totpList = await vault.totp.list();
    console.log("✅ Found", totpList.total, "TOTP entries");

    // Generate TOTP code
    const code = await vault.totp.generate(totp.id);
    console.log("✅ Generated TOTP code:", code.code);
    console.log("⏰ Code expires in:", code.remainingSeconds, "seconds");
  } catch (error) {
    console.error("❌ TOTP operation failed:", error);
  }
}

/**
 * Policy management example
 */
async function policiesExample() {
  const vault = createVaultClient({
    baseURL: "http://localhost:8080",
    auth: {
      type: "session",
    },
  });

  // Login first
  await vault.auth.login({
    username: "admin@example.com",
    password: "adminPassword123",
  });

  try {
    // List current policies
    const policies = await vault.policies.list();
    console.log("✅ Found", policies.total, "policies");

    // Create a new policy
    const newPolicy = await vault.policies.create({
      name: "Allow Secret Read",
      description: "Allows reading of user's own secrets",
      resource: "secret",
      actions: ["read"],
      effect: "allow",
      priority: 100,
    });

    console.log("✅ Policy created:", newPolicy.id, newPolicy.name);

    // Get policies for specific resource
    const secretPolicies = await vault.policies.getByResource("secret");
    console.log("✅ Found", secretPolicies.length, "secret policies");

    // Evaluate policy
    const evaluation = await vault.policies.evaluate("secret", "read", {
      userId: "user-123",
      resourceId: "secret-456",
    });

    console.log("✅ Policy evaluation:", evaluation.allowed);
    if (!evaluation.allowed) {
      console.log("❌ Access denied:", evaluation.reason);
    }
  } catch (error) {
    console.error("❌ Policy operation failed:", error);
  }
}

/**
 * Audit logging example
 */
async function auditExample() {
  const vault = createVaultClient({
    baseURL: "http://localhost:8080",
    auth: {
      type: "session",
    },
  });

  // Login first
  await vault.auth.login({
    username: "admin@example.com",
    password: "adminPassword123",
  });

  try {
    // Get audit logs
    const auditLogs = await vault.audit.list({
      page: 1,
      pageSize: 50,
      sortBy: "createdAt",
      sortOrder: "desc",
    });

    console.log("✅ Found", auditLogs.total, "audit entries");

    // Get failed authentication attempts
    const failedLogins = await vault.audit.getFailedAuth({
      dateFrom: new Date(Date.now() - 24 * 60 * 60 * 1000), // Last 24h
      pageSize: 100,
    });

    console.log("✅ Found", failedLogins.total, "failed login attempts");

    // Get secret access logs
    const secretAccess = await vault.audit.getSecretAccess({
      dateFrom: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000), // Last 7 days
      sortBy: "createdAt",
      sortOrder: "desc",
    });

    console.log("✅ Found", secretAccess.total, "secret access events");

    // Export audit logs to CSV
    const csvData = await vault.audit.exportToCSV({
      dateFrom: new Date("2025-01-01"),
      dateTo: new Date("2025-01-31"),
      action: "login",
    });

    console.log("✅ Audit logs exported to CSV");
    console.log("📊 CSV data length:", csvData.length, "characters");
  } catch (error) {
    console.error("❌ Audit operation failed:", error);
  }
}

/**
 * System health and monitoring example
 */
async function systemExample() {
  const vault = createVaultClient({
    baseURL: "http://localhost:8080",
  });

  try {
    // Check system health
    const health = await vault.system.health();
    console.log("✅ System status:", health.status);
    console.log("📊 Database status:", health.database);
    console.log("🔢 Version:", health.version);

    // Get version information
    const version = await vault.system.version();
    console.log("📦 Build version:", version.version);
    console.log("⏰ Build time:", version.buildTime);
    console.log("🔗 Git commit:", version.gitCommit);
    console.log("🐹 Go version:", version.goVersion);

    // Check if system is ready
    const isReady = await vault.system.ready();
    console.log("✅ System ready:", isReady);

    // Get comprehensive system status
    const status = await vault.system.status();
    console.log("🏥 System healthy:", status.healthy);
    console.log("🚀 System ready:", status.ready);

    if (status.ready) {
      console.log("✅ Aether Vault is fully operational!");
    }
  } catch (error) {
    console.error("❌ System check failed:", error);
  }
}

/**
 * Complete workflow example
 */
async function completeWorkflowExample() {
  const vault = createVaultClient({
    baseURL: "http://localhost:8080",
    auth: {
      type: "session",
    },
    timeout: 10000,
  });

  try {
    console.log("🚀 Starting Aether Vault workflow...");

    // 1. Check system health
    const health = await vault.system.health();
    if (health.status !== "healthy") {
      throw new Error("Aether Vault is not healthy");
    }
    console.log("✅ System is healthy");

    // 2. Authenticate
    const session = await vault.auth.login({
      username: "user@example.com",
      password: "securePassword123",
    });
    console.log("✅ Authenticated as:", session.user.email);

    // 3. Get user identity
    const user = await vault.identity.me();
    console.log(
      "✅ User:",
      user.displayName || `${user.firstName} ${user.lastName}`,
    );

    // 4. Create a secret
    const secret = await vault.secrets.create({
      name: "API Key",
      description: "External API authentication key",
      value: "sk-live-1234567890",
      type: "api_key",
      tags: "production,external",
    });
    console.log("✅ Created secret:", secret.name);

    // 5. Create TOTP
    const totp = await vault.totp.create({
      name: "Bank App",
      description: "Banking application 2FA",
    });
    console.log("✅ Created TOTP:", totp.name);

    // 6. Check audit logs
    const auditLogs = await vault.audit.list({ pageSize: 10 });
    console.log("✅ Found", auditLogs.total, "recent audit entries");

    // 7. Logout
    await vault.auth.logout();
    console.log("✅ Logged out successfully");

    console.log("🎉 Workflow completed successfully!");
  } catch (error) {
    console.error("❌ Workflow failed:", error);
  }
}

/**
 * Dynamic endpoint configuration example
 */
async function dynamicEndpointExample() {
  // Local development
  const localVault = createVaultClient({
    baseURL: "http://localhost:8080",
    auth: {
      type: "session",
    },
  });

  // Cloud deployment
  const cloudVault = createVaultClient({
    baseURL: "https://api.aethervault.com",
    auth: {
      type: "jwt",
      token: "your-jwt-token",
    },
  });

  // Appliance deployment
  const applianceVault = createVaultClient({
    baseURL: "https://vault.company.internal",
    auth: {
      type: "bearer",
      token: "your-bearer-token",
    },
  });

  console.log("✅ Vault clients configured for different environments");

  // Use whichever client is needed
  const currentVault = localVault; // Switch based on environment

  try {
    const health = await currentVault.system.health();
    console.log("✅ Health check:", health.status);
  } catch (error) {
    console.error("❌ Health check failed:", error);
  }
}

// Export examples for easy testing
export {
  authenticationExample,
  secretsExample,
  totpExample,
  policiesExample,
  auditExample,
  systemExample,
  completeWorkflowExample,
  dynamicEndpointExample,
};

// Run examples if this file is executed directly
if (require.main === module) {
  console.log("🚀 Running Aether Vault SDK Examples\n");

  // Uncomment the example you want to run
  // await authenticationExample();
  // await secretsExample();
  // await totpExample();
  // await policiesExample();
  // await auditExample();
  // await systemExample();
  // await completeWorkflowExample();
  // await dynamicEndpointExample();

  console.log("\n✨ Examples available in this file:");
  console.log("  - authenticationExample()");
  console.log("  - secretsExample()");
  console.log("  - totpExample()");
  console.log("  - policiesExample()");
  console.log("  - auditExample()");
  console.log("  - systemExample()");
  console.log("  - completeWorkflowExample()");
  console.log("  - dynamicEndpointExample()");
}
