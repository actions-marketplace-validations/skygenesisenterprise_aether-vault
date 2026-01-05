<div align="center">

# 🚀 Aether Vault Node.js SDK

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![TypeScript](https://img.shields.io/badge/TypeScript-5-blue?style=for-the-badge&logo=typescript)](https://www.typescriptlang.org/) [![Node.js](https://img.shields.io/badge/Node.js-18+-green?style=for-the-badge&logo=node.js)](https://nodejs.org/) [![Next.js](https://img.shields.io/badge/Next.js-16-black?style=for-the-badge&logo=next.js)](https://nextjs.org/)

**🔐 Official SDK for Aether Vault - Centralized secrets, TOTP, and identity management for modern applications.**

This SDK provides a comprehensive, type-safe interface for interacting with Aether Vault API, eliminating the need for raw `/api/v1/*` fetch calls throughout your codebase.

[🚀 Quick Start](#-quick-start) • [📋 Features](#-features) • [📖 Installation](#-installation) • [🎨 Usage Examples](#-usage-examples) • [📚 API Reference](#-api-reference) • [🏗️ Architecture](#-architecture) • [🛠️ Error Handling](#-error-handling) • [⚙️ Configuration](#-configuration) • [🤝 Contributing](#-contributing)

[![GitHub stars](https://img.shields.io/github/stars/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/stargazers) [![GitHub forks](https://img.shields.io/github/forks/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/network) [![GitHub issues](https://img.shields.io/github/issues/github/skygenesisenterprise/aether-vault)](https://github.com/skygenesisenterprise/aether-vault/issues)

</div>

---

## 🌟 What is Aether Vault Node.js SDK?

**Aether Vault Node.js SDK** is the official TypeScript client library for the Aether Vault secure secrets management platform. It provides a **clean, intuitive, and type-safe interface** that completely abstracts away HTTP complexity while giving developers full access to all Aether Vault capabilities.

### 🎯 Our Mission

- **🔐 Centralized Authentication** - JWT, session, and token-based auth with automatic refresh
- **🔒 Secure Secrets Management** - Complete CRUD operations with encryption and rotation
- **⚡ High-Performance TOTP** - Time-based One-Time Password generation and verification
- **👤 Identity Management** - User profiles, roles, and session control
- **🌐 Complete Auditing** - Full audit trails and compliance logging
- **🛡️ Type Safety** - Full TypeScript support with strict mode and comprehensive types
- **🔄 Environment Flexibility** - Support for local, cloud, and appliance deployments
- **🎨 Developer Experience** - Rich autocompletion, error handling, and debugging support

---

## ✨ Features

### 🔐 **Authentication & Security**

- **Multiple Auth Methods** - JWT tokens, bearer tokens, session-based auth
- **Automatic Token Refresh** - Seamless token renewal with configurable refresh logic
- **Password Management** - Secure password change, reset, and recovery flows
- **Two-Factor Authentication** - TOTP setup, verification, and management
- **Session Management** - View and revoke active user sessions

### 🔒 **Secrets Management**

- **Complete CRUD Operations** - Create, read, update, delete, and search secrets
- **Secure Storage** - Client-side encryption with server-side decryption
- **Secret Rotation** - Automatic and manual secret value rotation
- **Flexible Metadata** - Tags, descriptions, and custom metadata support
- **Expiration Handling** - Automatic secret expiration and archival

### ⚡ **TOTP (Two-Factor)**

- **QR Code Generation** - Setup TOTP with QR codes for mobile apps
- **Code Generation** - Generate time-based codes for verification
- **Backup Codes** - Generate and manage backup codes for recovery
- **Multiple Algorithms** - Support for SHA1, SHA256, and SHA512
- **Custom Time Steps** - Configurable time periods (30s, 60s, 90s)

### 👤 **Identity & User Management**

- **User Profiles** - Complete user information and profile management
- **Role-Based Access** - Hierarchical permissions and role management
- **Session Control** - View, manage, and revoke user sessions
- **Account Security** - Password policies, 2FA enforcement, email verification

### 🌐 **Audit & Compliance**

- **Comprehensive Logging** - All operations automatically logged with full context
- **Advanced Filtering** - Filter by user, action, resource, date range, and more
- **Export Capabilities** - Export audit logs to CSV, JSON, or custom formats
- **Real-time Monitoring** - Live audit trail updates and notifications

### 🛠️ **Type Safety & Developer Experience**

- **Full TypeScript Support** - Complete type definitions for all API responses
- **Strict Type Checking** - Compile-time error prevention
- **Rich Autocompletion** - Full IntelliSense support in IDEs
- **Typed Errors** - Comprehensive error types with specific error codes
- **Promise-Based API** - Modern async/await patterns throughout

### 🔄 **Multi-Environment Support**

- **Development Mode** - Optimized for local development with debugging
- **Production Mode** - Optimized for production with minimal logging
- **Environment Variables** - Support for multiple deployment configurations
- **Dynamic Endpoints** - Easy switching between local, cloud, and appliance deployments

---

## 📖 Installation

### Prerequisites

- **Node.js** 18.0.0 or higher
- **TypeScript** 5.0 or higher (recommended)
- **pnpm** 9.0.0 or higher (recommended)

### Installation

```bash
# npm
npm install aether-vault

# yarn
yarn add aether-vault

# pnpm (recommended)
pnpm add aether-vault
```

### Verification

```bash
# Verify installation
npm list aether-vault
```

---

## 🚀 Quick Start

### Basic Setup

```typescript
import { createVaultClient } from "aether-vault";

// Create vault client
const vault = createVaultClient({
  baseURL: "http://localhost:8080", // or "https://vault.example.com/api/v1"
  auth: {
    type: "session", // Uses browser cookies for web apps
    // Other options:
    // type: "jwt",
    // token: "your-jwt-token",
    // type: "bearer",
    // token: "your-bearer-token",
  },
  timeout: 10000, // Request timeout in milliseconds
  retry: true, // Enable automatic retries
  maxRetries: 3, // Maximum retry attempts
});
```

### First Authentication

```typescript
// Login user
const session = await vault.auth.login({
  username: "user@example.com",
  password: "securePassword123",
});

console.log("Logged in as:", session.user.displayName);
console.log("Token expires:", session.expiresAt);

// Check current session
const currentSession = await vault.auth.session();
console.log("Session valid:", currentSession.valid);
```

### Secret Management

```typescript
// Create a secret
const secret = await vault.secrets.create({
  name: "Database Connection",
  value: "postgresql://user:pass@localhost:5432/mydb",
  description: "Production database connection",
  type: "database",
  tags: "production,critical",
  expiresAt: new Date("2025-12-31"),
});

// List all secrets
const { secrets } = await vault.secrets.list({
  pageSize: 50,
  tags: "production",
});

// Get a specific secret
const dbSecret = await vault.secrets.getValue("Database_Connection");
console.log("Database URL:", dbSecret);

// Update a secret
const updated = await vault.secrets.update("Database_Connection", {
  description: "Updated database connection",
});

// Delete a secret
await vault.secrets.delete("old-api-key");
```

### TOTP Management

```typescript
// Generate TOTP for a service
const totp = await vault.totp.generate({
  name: "GitHub 2FA",
  description: "Two-factor authentication for GitHub",
  account: "user@example.com",
});

console.log("Scan QR code with:", totp.qrCode);
console.log("Backup codes:", totp.backupCodes);

// Generate a code
const { code, remainingSeconds } = await vault.totp.generate("github-2fa");
console.log("Your code:", code, "Valid for", remainingSeconds, "seconds");
```

### User Identity

```typescript
// Get current user profile
const user = await vault.identity.me();
console.log("Current user:", user.email, user.displayName);

// Get user policies
const policies = await vault.identity.policies();
console.log("User policies:", policies);
```

### Audit Logging

```typescript
// Get audit logs with filtering
const auditLogs = await vault.audit.list({
  dateFrom: new Date("2025-01-01"),
  action: "secret_access",
  pageSize: 100,
});

// Get failed authentication attempts
const failedLogins = await vault.audit.getFailedAuth({
  dateFrom: new Date(Date.now() - 24 * 60 * 60 * 1000), // Last 24 hours
  pageSize: 50,
});

// Export audit logs to CSV
const csvData = await vault.audit.exportToCSV({
  dateFrom: new Date("2025-01-01"),
  dateTo: new Date("2025-01-31"),
});
```

### System Health

```typescript
// Check system health
const health = await vault.system.health();
console.log("System status:", health.status);

// Get version information
const version = await vault.system.version();
console.log("Aether Vault version:", version.version);

// Comprehensive system status
const status = await vault.system.status();
console.log("System healthy:", status.healthy);
console.log("Version:", status.version.version);
```

---

## 📚 Usage Examples

### Complete Workflow Example

```typescript
import { createVaultClient } from "aether-vault";

async function completeWorkflow() {
  const vault = createVaultClient({
    baseURL: "http://localhost:8080",
    auth: { type: "session" },
  });

  try {
    console.log("🚀 Starting Aether Vault workflow...");

    // 1. Check system health
    const health = await vault.system.health();
    if (health.status !== "healthy") {
      throw new Error("Aether Vault is not healthy");
    }
    console.log("✅ System is healthy");

    // 2. Authenticate user
    const session = await vault.auth.login({
      username: "user@example.com",
      password: "securePassword123",
    });
    console.log("✅ Authenticated as:", session.user.displayName);

    // 3. Create a secret
    const secret = await vault.secrets.create({
      name: "API Key",
      value: "sk_live_1234567890",
      type: "api_key",
      tags: "production,critical",
    });
    console.log("✅ Created secret:", secret.name);

    // 4. Generate TOTP
    const totp = await vault.totp.generate({
      name: "Mobile App 2FA",
    });
    console.log("✅ Generated TOTP:", totp.name);

    // 5. List secrets
    const { secrets } = await vault.secrets.list();
    console.log("✅ Found", secrets.total, "secrets");

    // 6. Get user profile
    const user = await vault.identity.me();
    console.log("✅ Current user:", user.email);

    // 7. Check audit logs
    const auditLogs = await vault.audit.list({ pageSize: 10 });
    console.log("✅ Recent audit entries:", auditLogs.entries.length);

    // 8. Logout
    await vault.auth.logout();
    console.log("✅ Logged out successfully");

    console.log("🎉 Workflow completed successfully!");
  } catch (error) {
    console.error("❌ Workflow failed:", error);
  }
}

// Run the workflow
completeWorkflow();
```

### Next.js Integration Example

```typescript
// app/lib/vault.ts
import { createVaultClient } from "aether-vault";

export const vault = createVaultClient({
  baseURL: process.env.NEXT_PUBLIC_VAULT_URL || "/api/v1",
  auth: {
    type: "session", // For web applications
  },
});

export async function getServerSecrets() {
  try {
    return await vault.secrets.list();
  } catch (error) {
    console.error("Failed to fetch secrets:", error);
    throw error;
  }
}
```

```typescript
// app/components/secrets-manager.tsx
"use client";

import { useState, useEffect } from "react";
import { vault } from "@/lib/vault";

export function SecretsManager() {
  const [secrets, setSecrets] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadSecrets() {
      try {
        const response = await vault.secrets.list({ pageSize: 20 });
        setSecrets(response.secrets);
      } catch (error) {
        console.error("Failed to load secrets:", error);
      } finally {
        setLoading(false);
      }
    }

    loadSecrets();
  }, []);

  const handleCreateSecret = async (name: string, value: string) => {
    try {
      await vault.secrets.create({ name, value });
      const response = await vault.secrets.list({ pageSize: 20 });
      setSecrets(response.secrets);
    } catch (error) {
      console.error("Failed to create secret:", error);
    }
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Secrets Manager</h1>
      <button onClick={() => handleCreateSecret("NEW_SECRET", "value")}>
        Create Secret
      </button>
      <ul>
        {secrets.map((secret) => (
          <li key={secret.id}>
            <strong>{secret.name}</strong> - {secret.description}
          </li>
        ))}
      </ul>
    </div>
  );
}
```

---

## 📚 API Reference

### Core Methods

#### Authentication (`vault.auth.*`)

```typescript
// Login with credentials
vault.auth.login(credentials: AuthCredentials): Promise<AuthSession>

// Logout current session
vault.auth.logout(): Promise<void>

// Get current session
vault.auth.session(): Promise<{ user: UserIdentity; valid: boolean }>

// Register new user
vault.auth.register(userData: RegisterData): Promise<UserIdentity>

// Change password
vault.auth.changePassword(data: PasswordChangeData): Promise<void>

// Request password reset
vault.auth.forgotPassword(email: string): Promise<void>

// Reset password with token
vault.auth.resetPassword(data: PasswordResetData): Promise<void>

// Validate token
vault.auth.validate(): Promise<boolean>

// Get current user
vault.auth.getCurrentUser(): Promise<UserIdentity>
```

#### Secrets Management (`vault.secrets.*`)

```typescript
// List all secrets (with pagination)
vault.secrets.list(params?: SecretFilterParams): Promise<SecretListResponse>

// Get specific secret by ID
vault.secrets.get(id: string): Promise<Secret>

// Create new secret
vault.secrets.create(secret: SecretInput): Promise<Secret>

// Update existing secret
vault.secrets.update(id: string, updates: SecretUpdate): Promise<Secret>

// Delete secret
vault.secrets.delete(id: string): Promise<void>
```

#### TOTP Management (`vault.totp.*`)

```typescript
// List TOTP entries
vault.totp.list(params?: TOTPFilterParams): Promise<TOTPListResponse>

// Generate new TOTP
vault.totp.generate(data: TOTPCreateInput): Promise<TOTPEntry>

// Generate TOTP code
vault.totp.generate(id: string): Promise<TOTPCode>

// Update TOTP entry
vault.totp.update(id: string, updates: Partial<TOTPCreateInput>): Promise<TOTPEntry>

// Delete TOTP entry
vault.totp.delete(id: string): Promise<void>
```

#### Identity Management (`vault.identity.*`)

```typescript
// Get current user profile
vault.identity.me(): Promise<UserIdentity>

// Get user policies
vault.identity.policies(): Promise<Policy[]>
```

#### Audit Logging (`vault.audit.*`)

```typescript
// List audit entries
vault.audit.list(filter?: AuditFilter): Promise<AuditListResponse>

// Get specific audit entry
vault.audit.getEntry(id: string): Promise<AuditEntry>

// Get audit entries for specific user
vault.audit.getUserEntries(userId: string, options?: Omit<AuditFilter, "userId">): Promise<AuditListResponse>

// Get failed authentication attempts
vault.audit.getFailedAuth(options?: Omit<AuditFilter, "action">): Promise<AuditListResponse>

// Get secret access logs
vault.audit.getSecretAccess(options?: Omit<AuditFilter, "resource">): Promise<AuditListResponse>

// Export audit logs to CSV
vault.audit.exportToCSV(filter?: AuditFilter): Promise<string>
```

#### System Operations (`vault.system.*`)

```typescript
// Check system health
vault.system.health(): Promise<HealthResponse>

// Get version information
vault.system.version(): Promise<VersionResponse>

// Check system readiness
vault.system.ready(): Promise<boolean>

// Get system metrics
vault.system.metrics(): Promise<SystemMetrics>

// Get comprehensive system status
vault.system.status(): Promise<SystemStatus>
```

---

## 🛠️ Error Handling

The SDK provides comprehensive, typed error handling for all API operations:

### Error Types

```typescript
import {
  VaultError,
  VaultAuthError,
  VaultPermissionError,
  VaultNotFoundError,
  VaultServerError,
  VaultNetworkError,
} from "aether-vault";
```

### Error Handling Patterns

```typescript
import { VaultError, VaultAuthError } from "aether-vault";

try {
  const secret = await vault.secrets.get("non-existent-secret");
  console.log("Secret value:", secret.value);
} catch (error) {
  if (error instanceof VaultNotFoundError) {
    console.log("Secret not found - create it first");
  } else if (error instanceof VaultAuthError) {
    console.log("Authentication failed - please login again");
    window.location.href = "/login";
  } else if (error instanceof VaultPermissionError) {
    console.log("Permission denied - insufficient access rights");
  } else if (error instanceof VaultError) {
    console.log("Vault error:", error.message);
  } else {
    console.error("Unexpected error:", error);
  }
}
```

---

## ⚙️ Configuration

### Client Configuration Options

```typescript
interface VaultConfig {
  // API base URL - no trailing slash
  baseURL: string;

  // Authentication configuration
  auth: AuthConfig;

  // Request timeout (milliseconds)
  timeout?: number;

  // Retry configuration
  retry?: boolean;
  maxRetries?: number;
  retryDelay?: number;

  // Custom headers
  headers?: Record<string, string>;

  // Debug mode for development
  debug?: boolean;
}
```

### Authentication Types

```typescript
// Session-based (recommended for web apps)
interface SessionAuthConfig {
  type: "session";
}

// JWT-based (recommended for API services)
interface JwtAuthConfig {
  type: "jwt";
  token: string;
  jwt?: {
    autoRefresh?: boolean;
    refreshEndpoint?: string;
    refreshFn?: (token: string) => Promise<string>;
  };
}

// Bearer token-based
interface BearerAuthConfig {
  type: "bearer";
  token: string;
}

// No authentication
interface NoAuthConfig {
  type: "none";
}
```

### Environment Variables

```bash
# Local development
NEXT_PUBLIC_VAULT_URL=http://localhost:8080

# Cloud deployment
NEXT_PUBLIC_VAULT_URL=https://api.aethervault.com

# Appliance deployment
NEXT_PUBLIC_VAULT_URL=https://vault.company.internal
```

---

## 🏗️ Architecture

### Modular Design

```
src/
├── core/              # HTTP client, configuration, and error handling
│   ├── client.ts       # Main HTTP client with fetch/isomorphic support
│   ├── config.ts       # Configuration management and validation
│   └── errors.ts       # Typed error definitions
├── auth/              # Authentication domain client
│   └── auth.client.ts   # Authentication operations
├── secrets/           # Secrets management domain client
│   └── secrets.client.ts # Secrets CRUD operations
├── totp/              # TOTP domain client
│   └── totp.client.ts   # TOTP generation and verification
├── identity/          # Identity domain client
│   └── identity.client.ts # User profile and session management
├── audit/              # Audit domain client
│   └── audit.client.ts   # Audit log operations
├── policies/            # Policy management client
│   └── policy.client.ts   # Access control policies
├── system/             # System operations client
│   └── system.client.ts   # Health checks and metrics
└── types/              # Type definitions
│   ├── vault.ts          # Core API types
│   └── index.ts         # Type exports
└── index.ts            # Main SDK entry point
```

### HTTP Client Architecture

- **Isomorphic Design** - Works in Node.js, browsers, and Electron
- **Automatic Authentication** - Token management and refresh
- **Request/Response Interceptors** - Automatic error handling and retries
- **Comprehensive Error Types** - Specific error classes for different failure modes
- **Type Safety** - Full TypeScript support with strict mode

### Security Features

- **Secure Token Storage** - In-memory or configurable storage
- **Automatic Token Refresh** - Background token renewal
- **Rate Limiting Awareness** - Respects server rate limits
- **Request/Response Validation** - Input validation and response verification
- **HTTPS Support** - Secure communication with API endpoints

---

## 🛠️ Development

### TypeScript Configuration

```json
// tsconfig.json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "noImplicitReturns": true,
    "noUnusedLocals": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "forceConsistentCasingInFileNames": true,
    "moduleResolution": "node",
    "declaration": true,
    "outDir": "./dist",
    "rootDir": "./src"
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "**/*.test.ts"]
}
```

### Development Workflow

```bash
# Install dependencies
pnpm install

# Run in development mode
pnpm dev

# Build for production
pnpm build

# Run tests
pnpm test

# Type checking
pnpm typecheck
```

### Debugging

```typescript
const vault = createVaultClient({
  baseURL: "http://localhost:8080",
  debug: true, // Enable debug logging
});

// All requests and responses will be logged
```

---

## 🤝 Contributing

We welcome contributions to help improve the Aether Vault Node.js SDK! Whether you're experienced with TypeScript, API clients, authentication systems, or security-focused development, there's a place for you.

### How to Get Started

1. **Fork the repository** and create a feature branch
2. **Check the issues** for tasks that need help
3. **Join discussions** about architecture and features
4. **Start small** - Documentation, tests, or minor features
5. **Follow our guidelines** and commit standards

### Areas Needing Help

- **Core SDK Development** - API clients, error handling, HTTP layer
- **TypeScript Types** - Type definitions and interfaces
- **Documentation** - API docs, examples, and guides
- **Testing** - Unit tests and integration tests
- **Security** - Authentication and secure communication
- **Developer Experience** - Debugging, autocompletion, and tooling

### Contribution Guidelines

- **Make-First Workflow** - Use `pnpm` commands for all operations
- **TypeScript Strict Mode** - All code must pass strict type checking
- **Component Structure** - Follow established patterns for client organization
- **API Design** - RESTful endpoints with proper HTTP methods
- **Error Handling** - Comprehensive error handling and logging
- **Security First** - Validate all inputs and implement proper authentication

### Development Commands

```bash
# Install dependencies
pnpm install

# Run tests
pnpm test

# Type checking
pnpm typecheck

# Build library
pnpm build

# Run examples
pnpm dev examples/basic-usage.ts
```

---

## 📄 License

MIT License - see [LICENSE](../../../LICENSE) file for details.

---

## 🚀 Support & Community

### Get Help

- 📖 **[Documentation](../docs/)** - Comprehensive guides and API docs
- 📋 **[Issue Tracker](https://github.com/skygenesisenterprise/aether-vault/issues)** - Bug reports and feature requests
- 💡 **[GitHub Discussions](https://github.com/skygenesisenterprise/aether-vault/discussions)** - General questions and ideas
- 📧 **Email** - support@skygenesisenterprise.com

### Reporting Issues

When reporting bugs, please include:

- **Clear description** of the problem
- **Steps to reproduce**
- **Environment information** - Node.js version, TypeScript version, OS, etc.
- **Error logs** or screenshots
- **Expected vs actual behavior**
- **SDK version** being used

---

## 📊 Project Status

| Component                | Status     | Technology                 | Evolution    |
| ------------------------ | ---------- | -------------------------- | ------------ |
| **Core SDK**             | ✅ Working | TypeScript + Native Fetch  | **Enhanced** |
| **Authentication**       | ✅ Working | JWT + Session + Refresh    | **Complete** |
| **Secrets Management**   | ✅ Working | Full CRUD + Encryption     | **Complete** |
| **TOTP Support**         | ✅ Working | QR Codes + Verification    | **Complete** |
| **Identity Management**  | ✅ Working | Profiles + Session Control | **Complete** |
| **Audit Logging**        | ✅ Working | Full Audit + Export        | **Complete** |
| **System Operations**    | ✅ Working | Health + Metrics           | **Complete** |
| **Type Safety**          | ✅ Working | Strict Mode + Types        | **Complete** |
| **Error Handling**       | ✅ Working | Typed Errors + Catching    | **Complete** |
| **Documentation**        | ✅ Working | Complete Examples          | **Enhanced** |
| **Multi-Environment**    | ✅ Working | Local/Cloud/Appliance      | **Enhanced** |
| **Developer Experience** | ✅ Working | Autocompletion + Debugging | **Enhanced** |

---

<div align="center">

### 🚀 **Ready for Production Use**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

_Building the future of secure secrets management for modern applications_

</div>
