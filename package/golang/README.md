<div align="center">

# 🐹 Aether Vault Go SDK

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![Go](https://img.shields.io/badge/Go-1.25+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![Go Report Card](https://goreportcard.com/badge/github.com/skygenesisenterprise/aether-vault/package/golang?style=for-the-badge)](https://goreportcard.com/report/github.com/skygenesisenterprise/aether-vault/package/golang)

**🔥 Native Go Client Library for Aether Vault**

A high-performance, type-safe Go SDK for interacting with Aether Vault - the modern secrets management and identity platform. This package provides native Go integration with comprehensive features for authentication, secrets management, TOTP, identity management, policies, and audit logging.

[🚀 Quick Start](#-quick-start) • [📋 Features](#-features) • [🛠️ Installation](#️-installation) • [📚 Documentation](#-documentation) • [🔧 API Reference](#-api-reference) • [🤝 Contributing](#-contributing)

[![GoDoc](https://pkg.go.dev/badge/github.com/skygenesisenterprise/aether-vault/package/golang.svg?style=for-the-badge)](https://pkg.go.dev/github.com/skygenesisenterprise/aether-vault/package/golang)

</div>

---

## 🌟 What is Aether Vault Go SDK?

**Aether Vault Go SDK** is the official Go client library for Aether Vault, providing native Go integration with type safety, high performance, and comprehensive feature coverage for modern secrets management and identity operations.

### 🎯 Key Features

- **🚀 High Performance** - Native Go implementation with efficient HTTP client
- **🔐 Complete Authentication** - JWT token management and automatic refresh
- **🗝️ Secrets Management** - Full CRUD operations for secrets and keys
- **🔢 TOTP Support** - Time-based One-Time Password generation and validation
- **👤 Identity Management** - User and group management capabilities
- **📋 Policy Engine** - Access control and policy enforcement
- **📊 Audit Logging** - Comprehensive audit trail and compliance features
- **🔄 Context Support** - Proper Go context handling for cancellation and timeouts
- **🛡️ Error Handling** - Structured error types with detailed error codes
- **⚡ Concurrent Safe** - Thread-safe client for concurrent operations
- **🔧 Configurable** - Flexible configuration with timeouts, retries, and TLS options

---

## 📋 Features Overview

### 🔐 **Authentication Client**

- JWT token management with automatic refresh
- Login/logout operations
- Token validation and renewal
- Session management

### 🗝️ **Secrets Management**

- Create, read, update, delete secrets
- Secret versioning and history
- Encrypted storage and retrieval
- Bulk operations support

### 🔢 **TOTP Client**

- Generate TOTP tokens
- Validate TOTP codes
- TOTP configuration management
- Backup code handling

### 👤 **Identity Management**

- User and group operations
- Role assignments
- Profile management
- Directory synchronization

### 📋 **Policy Engine**

- Policy creation and management
- Access rule evaluation
- Role-based access control
- Policy inheritance

### 📊 **Audit Logging**

- Comprehensive audit trails
- Event tracking and filtering
- Compliance reporting
- Log aggregation support

---

## 🚀 Quick Start

### 📋 Prerequisites

- **Go** 1.25.0 or higher
- **Aether Vault Server** running and accessible
- Valid API credentials (if required)

### 🔧 Installation

```bash
# Get the package
go get github.com/skygenesisenterprise/aether-vault/package/golang

# Import in your Go code
import "github.com/skygenesisenterprise/aether-vault/package/golang"
```

### ⚡ Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/skygenesisenterprise/aether-vault/package/golang"
)

func main() {
    // Create vault client
    vault, err := vault.New(vault.NewConfig(
        "https://your-vault-server.com",
        "your-api-token",
    ))
    if err != nil {
        log.Fatal("Failed to create vault client:", err)
    }
    defer vault.Close()

    ctx := context.Background()

    // Check vault health
    if err := vault.Health(ctx); err != nil {
        log.Fatal("Vault is not healthy:", err)
    }

    // Get vault version
    version, err := vault.Version(ctx)
    if err != nil {
        log.Fatal("Failed to get version:", err)
    }

    fmt.Printf("Connected to Aether Vault v%s\n", version)

    // Store a secret
    secret, err := vault.Secrets.Store(ctx, &secrets.StoreRequest{
        Path:  "myapp/database",
        Value: "database-connection-string",
        Metadata: map[string]string{
            "environment": "production",
            "owner":       "my-team",
        },
    })
    if err != nil {
        log.Fatal("Failed to store secret:", err)
    }

    fmt.Printf("Secret stored with ID: %s\n", secret.ID)

    // Retrieve the secret
    retrieved, err := vault.Secrets.Get(ctx, secret.ID)
    if err != nil {
        log.Fatal("Failed to retrieve secret:", err)
    }

    fmt.Printf("Retrieved secret: %s\n", retrieved.Value)
}
```

### 🔐 Authentication Example

```go
func authenticateExample() {
    vault, err := vault.New(vault.NewConfig(
        "https://your-vault-server.com",
        "", // Empty token for initial auth
    ))
    if err != nil {
        log.Fatal(err)
    }
    defer vault.Close()

    ctx := context.Background()

    // Login with credentials
    authResp, err := vault.Auth.Login(ctx, &auth.LoginRequest{
        Username: "admin@example.com",
        Password: "your-password",
    })
    if err != nil {
        log.Fatal("Login failed:", err)
    }

    fmt.Printf("Authenticated! Token expires at: %v\n", authResp.ExpiresAt)

    // Update client with token
    vault.GetConfig().Token = authResp.Token

    // Validate current token
    valid, err := vault.Auth.Validate(ctx)
    if err != nil {
        log.Fatal("Token validation failed:", err)
    }

    if valid {
        fmt.Println("Token is valid!")
    }
}
```

### 🔢 TOTP Example

```go
func totpExample() {
    vault, err := vault.New(vault.NewConfig(
        "https://your-vault-server.com",
        "your-api-token",
    ))
    if err != nil {
        log.Fatal(err)
    }
    defer vault.Close()

    ctx := context.Background()

    // Generate a new TOTP secret
    totpSecret, err := vault.TOTP.Generate(ctx, &totp.GenerateRequest{
        Username: "user@example.com",
        Issuer:   "My App",
    })
    if err != nil {
        log.Fatal("Failed to generate TOTP:", err)
    }

    fmt.Printf("TOTP Secret: %s\n", totpSecret.Secret)
    fmt.Printf("QR Code URL: %s\n", totpSecret.QRCodeURL)

    // Validate a TOTP code
    valid, err := vault.TOTP.Validate(ctx, &totp.ValidateRequest{
        Username: "user@example.com",
        Code:     "123456", // User's TOTP code
    })
    if err != nil {
        log.Fatal("Failed to validate TOTP:", err)
    }

    if valid {
        fmt.Println("TOTP code is valid!")
    }
}
```

---

## 🛠️ Configuration

### 🔧 Configuration Options

```go
config := &vault.Config{
    Endpoint:   "https://your-vault-server.com",
    Token:      "your-api-token",
    Timeout:    30 * time.Second,
    RetryCount: 3,
    UserAgent:  "my-app/1.0.0",
    Debug:      true,
    TLSConfig: &config.TLSConfig{
        InsecureSkipVerify: false,
        CertFile:          "/path/to/cert.pem",
        KeyFile:           "/path/to/key.pem",
        CAFile:            "/path/to/ca.pem",
    },
    Headers: map[string]string{
        "X-Custom-Header": "custom-value",
    },
}

vault, err := vault.New(config)
```

### 🔄 Default Configuration

```go
// Use default configuration with custom endpoint and token
vault, err := vault.New(vault.NewConfig(
    "https://your-vault-server.com",
    "your-api-token",
))

// Or start from defaults and customize
config := vault.DefaultConfig()
config.Timeout = 60 * time.Second
config.RetryCount = 5
config.Debug = true

vault, err := vault.New(config)
```

---

## 📚 API Reference

### 🔧 Main Client

```go
type Vault struct {
    Auth     *auth.AuthClient     // Authentication operations
    Secrets  *secrets.SecretsClient // Secrets management
    TOTP     *totp.TOTPClient     // TOTP operations
    Identity *identity.IdentityClient // Identity management
    Policies *policies.PolicyClient // Policy management
    Audit    *audit.AuditClient    // Audit logging
}

// Create new vault client
func New(cfg *Config) (*Vault, error)

// Create default configuration
func DefaultConfig() *Config

// Create configuration with endpoint and token
func NewConfig(endpoint, token string) *Config

// Health check
func (v *Vault) Health(ctx context.Context) error

// Get vault version
func (v *Vault) Version(ctx context.Context) (string, error)

// Get vault information
func (v *Vault) Info(ctx context.Context) (map[string]interface{}, error)

// Close client and cleanup resources
func (v *Vault) Close() error
```

### 🔐 Authentication Client

```go
// Login with credentials
func (a *AuthClient) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)

// Logout current session
func (a *AuthClient) Logout(ctx context.Context) error

// Validate current token
func (a *AuthClient) Validate(ctx context.Context) (bool, error)

// Refresh authentication token
func (a *AuthClient) Refresh(ctx context.Context) (*RefreshResponse, error)

// Get current user info
func (a *AuthClient) Whoami(ctx context.Context) (*UserInfo, error)
```

### 🗝️ Secrets Client

```go
// Store a new secret
func (s *SecretsClient) Store(ctx context.Context, req *StoreRequest) (*Secret, error)

// Retrieve a secret
func (s *SecretsClient) Get(ctx context.Context, id string) (*Secret, error)

// Update a secret
func (s *SecretsClient) Update(ctx context.Context, id string, req *UpdateRequest) (*Secret, error)

// Delete a secret
func (s *SecretsClient) Delete(ctx context.Context, id string) error

// List secrets
func (s *SecretsClient) List(ctx context.Context, filter *ListFilter) ([]*Secret, error)

// Get secret versions
func (s *SecretsClient) Versions(ctx context.Context, id string) ([]*SecretVersion, error)
```

### 🔢 TOTP Client

```go
// Generate new TOTP secret
func (t *TOTPClient) Generate(ctx context.Context, req *GenerateRequest) (*TOTPSecret, error)

// Validate TOTP code
func (t *TOTPClient) Validate(ctx context.Context, req *ValidateRequest) (bool, error)

// Enable TOTP for user
func (t *TOTPClient) Enable(ctx context.Context, req *EnableRequest) error

// Disable TOTP for user
func (t *TOTPClient) Disable(ctx context.Context, req *DisableRequest) error

// Get backup codes
func (t *TOTPClient) BackupCodes(ctx context.Context, req *BackupCodesRequest) ([]string, error)
```

---

## 🔄 Advanced Usage

### 🌐 Context Management

```go
// Use context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

secret, err := vault.Secrets.Get(ctx, secretID)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Request timed out")
    } else {
        log.Printf("Failed to get secret: %v", err)
    }
}
```

### 🔄 Retry Logic

```go
// Configure retry behavior
config := vault.DefaultConfig()
config.RetryCount = 5
config.Timeout = 15 * time.Second

vault, err := vault.New(config)
```

### 🔧 Custom HTTP Client

```go
// Configure custom HTTP transport
transport := &http.Transport{
    MaxIdleConns:        10,
    IdleConnTimeout:     30 * time.Second,
    DisableCompression:  true,
}

config := vault.DefaultConfig()
config.HTTPClient = &http.Client{
    Transport: transport,
    Timeout:   config.Timeout,
}

vault, err := vault.New(config)
```

### 📊 Error Handling

```go
// Structured error handling
secret, err := vault.Secrets.Get(ctx, secretID)
if err != nil {
    var vaultErr *errors.Error
    if errors.As(err, &vaultErr) {
        switch vaultErr.Code {
        case errors.ErrCodeNotFound:
            log.Printf("Secret not found: %s", secretID)
        case errors.ErrCodeUnauthorized:
            log.Printf("Unauthorized access to secret: %s", secretID)
        case errors.ErrCodeInternal:
            log.Printf("Internal server error: %v", vaultErr.Message)
        default:
            log.Printf("Vault error: %s", vaultErr.Error())
        }
    } else {
        log.Printf("Unexpected error: %v", err)
    }
    return
}
```

---

## 🧪 Testing

### 🔧 Unit Testing

```go
func TestSecretOperations(t *testing.T) {
    // Create test client with mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Mock response
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id":"test-secret","value":"test-value"}`))
    }))
    defer server.Close()

    config := vault.NewConfig(server.URL, "test-token")
    vault, err := vault.New(config)
    require.NoError(t, err)
    defer vault.Close()

    // Test secret retrieval
    secret, err := vault.Secrets.Get(context.Background(), "test-secret")
    assert.NoError(t, err)
    assert.Equal(t, "test-secret", secret.ID)
    assert.Equal(t, "test-value", secret.Value)
}
```

### 🔄 Integration Testing

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Configure with test vault instance
    endpoint := os.Getenv("VAULT_TEST_ENDPOINT")
    token := os.Getenv("VAULT_TEST_TOKEN")

    if endpoint == "" || token == "" {
        t.Skip("VAULT_TEST_ENDPOINT and VAULT_TEST_TOKEN must be set for integration tests")
    }

    config := vault.NewConfig(endpoint, token)
    vault, err := vault.New(config)
    require.NoError(t, err)
    defer vault.Close()

    // Test real operations
    err = vault.Health(context.Background())
    assert.NoError(t, err)

    version, err := vault.Version(context.Background())
    assert.NoError(t, err)
    assert.NotEmpty(t, version)
}
```

---

## 📦 Examples Repository

Complete examples are available in the [examples](examples/) directory:

- **[Basic Usage](examples/basic/)** - Getting started guide
- **[Authentication](examples/auth/)** - Authentication flows
- **[Secrets Management](examples/secrets/)** - Secrets operations
- **[TOTP Setup](examples/totp/)** - TOTP configuration
- **[Error Handling](examples/errors/)** - Error handling patterns
- **[Concurrency](examples/concurrency/)** - Concurrent operations

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### 🎯 How to Contribute

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### 🏗️ Development Setup

```bash
# Clone the repository
git clone https://github.com/skygenesisenterprise/aether-vault.git
cd aether-vault/package/golang

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench ./...

# Lint code
golangci-lint run
```

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 🏆 Acknowledgments

- **Aether Vault** - The underlying vault platform
- **Go Community** - Excellent programming language and ecosystem
- **Contributors** - All contributors who help make this project better

---

## 📞 Support & Community

### 💬 **Get Help**

- 📖 **[Documentation](../../docs/)** - Comprehensive guides
- 📦 **[Go Package Reference](https://pkg.go.dev/github.com/skygenesisenterprise/aether-vault/package/golang)** - GoDoc documentation
- 🐛 **[GitHub Issues](https://github.com/skygenesisenterprise/aether-vault/issues)** - Bug reports and feature requests
- 💡 **[GitHub Discussions](https://github.com/skygenesisenterprise/aether-vault/discussions)** - General questions and ideas

### 🐛 **Reporting Issues**

When reporting bugs, please include:

- Go version and OS
- Aether Vault version
- Complete error message and stack trace
- Minimal reproduction case
- Expected vs actual behavior

---

<div align="center">

### 🚀 **Production-Ready Go SDK for Aether Vault!**

[⭐ Star This Repo](https://github.com/skygenesisenterprise/aether-vault) • [🐛 Report Issues](https://github.com/skygenesisenterprise/aether-vault/issues) • [💡 Start a Discussion](https://github.com/skygenesisenterprise/aether-vault/discussions)

---

**🔧 Native Go Integration with Type Safety and High Performance!**

**Made with ❤️ by the [Aether Vault Team](https://aether-vault.com)**

_Built for Go developers who demand performance, reliability, and comprehensive features_

</div>
