<div align="center">

# 📦 Aether Vault Package Ecosystem

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![Go](https://img.shields.io/badge/Go-1.21+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![TypeScript](https://img.shields.io/badge/TypeScript-5-blue?style=for-the-badge&logo=typescript)](https://www.typescriptlang.org/) [![Node.js](https://img.shields.io/badge/Node.js-18+-green?style=for-the-badge&logo=node.js)](https://nodejs.org/) [![Docker](https://img.shields.io/badge/Docker-Ready-blue?style=for-the-badge&logo=docker)](https://www.docker.com/)

**🔥 Comprehensive Multi-Language SDK Ecosystem for Aether Vault Integration**

A complete package ecosystem providing native SDKs, CLI tools, and integration packages for seamless Aether Vault deployment across multiple platforms and languages.

[🚀 Quick Start](#-quick-start) • [📦 Available Packages](#-available-packages) • [🛠️ Tech Stack](#️-tech-stack) • [📁 Architecture](#-architecture) • [🤝 Contributing](#-contributing)

[![GitHub stars](https://img.shields.io/github/stars/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/stargazers) [![GitHub forks](https://img.shields.io/github/forks/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/network)

</div>

---

## 🌟 What is the Aether Vault Package Ecosystem?

The **Aether Vault Package Ecosystem** is a comprehensive collection of SDKs, tools, and integration packages designed to provide seamless access to Aether Vault functionality across multiple programming languages and platforms.

### 🎯 Our Vision

- **🚀 Multi-Language Support** - Native SDKs for Go, TypeScript/Node.js, and more
- **🛠️ Developer Tools** - CLI tools and utilities for streamlined development
- **🔗 Platform Integration** - GitHub Apps, Docker containers, and cloud deployment
- **📦 Modular Design** - Use only what you need with minimal dependencies
- **🔐 Security First** - Enterprise-grade authentication and authorization
- **🌐 Universal Compatibility** - Works across different environments and platforms

---

## 📦 Available Packages

### 🐹 **@aether-vault/golang** - Go SDK & CLI

**Purpose**: Native Go SDK and command-line tools for Aether Vault integration.

**Key Features**:

- ✅ High-performance Go client library
- ✅ Complete CLI tool suite for vault management
- ✅ Type-safe Go structs and interfaces
- ✅ Go modules support with semantic versioning
- ✅ Middleware and authentication handling
- ✅ Comprehensive audit logging

**Quick Start**:

```go
import "github.com/skygenesisenterprise/aether-vault/package/golang"

client := golang.NewClient(&golang.Config{
    BaseURL: "https://vault.example.com",
    APIKey:  "your-api-key",
})

// Authenticate and access vault
err := client.Auth.Login("username", "password")
if err != nil {
    log.Fatal(err)
}

// Retrieve secrets
secret, err := client.Secrets.Get("database/credentials")
```

**Installation**:

```bash
go get github.com/skygenesisenterprise/aether-vault/package/golang
```

---

### 📦 **@aether-vault/node** - Node.js/TypeScript SDK

**Purpose**: Universal TypeScript SDK for Node.js and browser environments.

**Key Features**:

- ✅ Universal client (Node.js + Browser support)
- ✅ TypeScript strict mode with full type definitions
- ✅ Authentication and session management
- ✅ Next.js integration hooks and utilities
- ✅ Comprehensive examples and documentation
- ✅ Promise-based API with async/await support

**Quick Start**:

```typescript
import { AetherVaultClient } from "@aether-vault/node";

const client = new AetherVaultClient({
  baseURL: "https://vault.example.com",
  apiKey: "your-api-key",
});

// Authenticate
await client.auth.login({
  username: "user@example.com",
  password: "secure-password",
});

// Access vault services
const secret = await client.secrets.get("api/keys");
const audit = await client.audit.list();
```

**Installation**:

```bash
npm install @aether-vault/node
# or
pnpm add @aether-vault/node
# or
yarn add @aether-vault/node
```

---

### 🐙 **@aether-vault/github** - GitHub Integration

**Purpose**: GitHub Marketplace application for vault automation and CI/CD integration.

**Key Features**:

- ✅ Verified GitHub Marketplace app
- ✅ Automated secret scanning and detection
- ✅ PR commenting and workflow integration
- ✅ Vault client integration for GitHub Actions
- ✅ Webhook security with HMAC-SHA256
- ✅ Docker deployment ready

**Quick Start**:

```typescript
import { GitHubApp } from "@aether-vault/github";

const app = new GitHubApp({
  appId: process.env.GITHUB_APP_ID,
  privateKey: process.env.GITHUB_PRIVATE_KEY,
  vaultEndpoint: process.env.VAULT_ENDPOINT,
});

// Handle webhook events
app.webhook.on("pull_request.opened", async (event) => {
  await app.scanner.scanRepository(event.repository);
  await app.prCommenter.addComment(
    event.pull_request,
    "🔒 Secret scan completed. No vulnerabilities detected.",
  );
});
```

**Installation**:

```bash
npm install @aether-vault/github
```

---

### 🐳 **@aether-vault/docker** - Docker Runtime

**Purpose**: Docker-based runtime environment for containerized vault deployments.

**Key Features**:

- ✅ Lightweight Docker runtime for vault services
- ✅ Container orchestration and management
- ✅ Runtime health monitoring and logging
- ✅ Secure injection of vault credentials
- ✅ Multi-platform container support
- ✅ Kubernetes integration ready

**Quick Start**:

```bash
# Build and run the Docker runtime
docker build -t aether-vault-runtime package/docker/
docker run -d \
  --name vault-runtime \
  -p 8080:8080 \
  -e VAULT_ENDPOINT=https://vault.example.com \
  aether-vault-runtime
```

**Usage**:

```go
// Inside containerized applications
import "github.com/skygenesisenterprise/aether-vault/package/docker/cmd/aether-runtime"

// Runtime automatically handles:
// - Vault authentication
// - Secret injection
// - Health monitoring
// - Audit logging
```

---

### ⚡ **@aether-vault/cli** - Command Line Interface

**Purpose**: Universal CLI tool for vault management and administration.

**Key Features**:

- ✅ Cross-platform CLI (Windows, macOS, Linux)
- ✅ Interactive shell and batch mode
- ✅ Vault initialization and management
- ✅ User and policy administration
- ✅ Audit log inspection
- ✅ Cloud and local runtime support

**Quick Start**:

```bash
# Install the CLI
npm install -g @aether-vault/cli

# Initialize vault
aether-vault init --endpoint https://vault.example.com

# Login and manage vault
aether-vault auth login
aether-vault secrets list
aether-vault policies create admin-policy.json

# Runtime management
aether-vault runtime start --docker
aether-vault runtime status
```

**Installation**:

```bash
npm install -g @aether-vault/cli
# or download binary from releases
```

---

## 🛠️ Tech Stack

### 🎨 **Package Technologies**

```
Multi-Language Ecosystem
├── 🐹 Go SDK (Native Go)
│   ├── High-Performance HTTP Client
│   ├── CLI Tools & Utilities
│   ├── Go Modules Support
│   └── Type-Safe Structs
├── 📦 Node.js SDK (TypeScript)
│   ├── Universal Client (Node.js + Browser)
│   ├── Next.js Integration Hooks
│   ├── Promise-Based API
│   └── Full Type Definitions
├── 🐙 GitHub App (TypeScript + Fastify)
│   ├── Webhook Security
│   ├── Secret Scanning
│   ├── PR Automation
│   └── CI/CD Integration
├── 🐳 Docker Runtime (Go)
│   ├── Container Orchestration
│   ├── Runtime Monitoring
│   ├── Secret Injection
│   └── Health Management
└── ⚡ CLI Tools (Go)
    ├── Cross-Platform Support
    ├── Interactive Shell
    ├── Vault Administration
    └── Runtime Management
```

### 🔧 **Shared Infrastructure**

```
Common Package Foundation
├── 🔐 Authentication & Authorization
│   ├── JWT Token Management
│   ├── Multi-Factor Auth Support
│   ├── Session Management
│   └── Security Middleware
├── 📊 Audit & Logging
│   ├── Structured Logging
│   ├── Audit Trail
│   ├── Event Tracking
│   └── Compliance Reporting
├── 🛡️ Security Features
│   ├── Encryption at Rest
│   ├── Secure Transport
│   ├── Input Validation
│   └── Rate Limiting
├── 🌐 Network & Transport
│   ├── HTTP/HTTPS Clients
│   ├── Connection Pooling
│   ├── Retry Logic
│   └── Circuit Breakers
└── 📦 Package Management
    ├── Semantic Versioning
    ├── Dependency Management
    ├── Build Automation
    └── Release Engineering
```

---

## 📁 Architecture

### 🏗️ **Package Ecosystem Structure**

```
package/
├── golang/                     # 🐹 Go SDK & CLI Tools
│   ├── client/               # HTTP client implementation
│   ├── auth/                 # Authentication handlers
│   ├── secrets/              # Secret management
│   ├── audit/                # Audit logging
│   ├── middleware/           # HTTP middleware
│   └── vault.go              # Main SDK entry point
├── node/                      # 📦 Node.js/TypeScript SDK
│   ├── src/
│   │   ├── core/             # Core client functionality
│   │   ├── auth/             # Authentication client
│   │   ├── secrets/          # Secret management
│   │   ├── audit/            # Audit client
│   │   ├── nextjs/           # Next.js integration
│   │   └── types/            # TypeScript definitions
│   ├── examples/             # Usage examples
│   └── package.json          # Node.js package configuration
├── github/                    # 🐙 GitHub Integration App
│   ├── src/
│   │   ├── webhook/          # Webhook handlers
│   │   ├── scanner/          # Secret scanning
│   │   ├── prCommenter/      # PR automation
│   │   ├── vaultClient/      # Vault integration
│   │   └── auth/             # GitHub authentication
│   ├── Dockerfile            # Container configuration
│   └── docker-compose.yml     # Development environment
├── docker/                    # 🐳 Docker Runtime
│   ├── cmd/
│   │   └── aether-runtime/   # Runtime entry point
│   ├── internal/
│   │   ├── runtime/          # Runtime management
│   │   ├── injector/         # Secret injection
│   │   ├── auth/             # Authentication client
│   │   └── vault/            # Vault client
│   └── Dockerfile            # Runtime container
├── cli/                       # ⚡ Command Line Interface
│   ├── cmd/                  # CLI commands
│   ├── internal/
│   │   ├── client/           # API client
│   │   ├── config/           # Configuration management
│   │   ├── runtime/          # Runtime detection
│   │   └── ui/               # User interface
│   └── main.go               # CLI entry point
└── README.md                 # This file
```

### 🔄 **Integration Patterns**

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Applications  │    │   Package SDKs    │    │   Aether Vault  │
│                 │◄──►│                   │◄──►│   Server        │
│  • Go Apps      │    │  • Go SDK         │    │                 │
│  • Node.js Apps │    │  • Node.js SDK    │    │  • Auth Service │
│  • Web Frontend │    │  • GitHub App     │    │  • Secret Store │
│  • CLI Tools    │    │  • Docker Runtime │    │  • Audit Log    │
│  • CI/CD Pipelines│   │  • CLI Tools      │    │  • Policy Engine│
└─────────────────┘    └──────────────────┘    └─────────────────┘
            │                       │                       │
            ▼                       ▼                       ▼
     Native SDKs          Multi-Language Support     Centralized Vault
     Type Safety          Universal Compatibility    Enterprise Security
     Performance          Platform Integration       Compliance & Audit
```

---

## 🚀 Quick Start

### 📋 **Prerequisites**

- **Aether Vault Server** running and accessible
- **API Credentials** with appropriate permissions
- **Platform Tools**:
  - Go 1.21+ (for Go SDK)
  - Node.js 18+ (for Node.js SDK)
  - Docker (for Docker runtime)
  - Git (for source management)

### 🔧 **Installation Guide**

#### **Go SDK Installation**

```bash
# Install the Go SDK
go get github.com/skygenesisenterprise/aether-vault/package/golang

# Import in your Go application
import "github.com/skygenesisenterprise/aether-vault/package/golang"
```

#### **Node.js SDK Installation**

```bash
# Install via npm
npm install @aether-vault/node

# or via pnpm
pnpm add @aether-vault/node

# or via yarn
yarn add @aether-vault/node
```

#### **GitHub App Installation**

```bash
# Install the GitHub App package
npm install @aether-vault/github

# Set up environment variables
export GITHUB_APP_ID="your-app-id"
export GITHUB_PRIVATE_KEY="path-to-private-key.pem"
export VAULT_ENDPOINT="https://vault.example.com"
```

#### **CLI Installation**

```bash
# Install globally via npm
npm install -g @aether-vault/cli

# Or download binary from GitHub releases
curl -L https://github.com/skygenesisenterprise/aether-vault/releases/latest/download/aether-vault-linux-amd64 -o aether-vault
chmod +x aether-vault
sudo mv aether-vault /usr/local/bin/
```

### 🌐 **Basic Usage Examples**

#### **Go SDK Example**

```go
package main

import (
    "fmt"
    "log"
    "github.com/skygenesisenterprise/aether-vault/package/golang"
)

func main() {
    // Initialize client
    client, err := golang.NewClient(&golang.Config{
        BaseURL: "https://vault.example.com",
        APIKey:  "your-api-key",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Authenticate
    if err := client.Auth.Login("admin", "password"); err != nil {
        log.Fatal(err)
    }

    // Get a secret
    secret, err := client.Secrets.Get("database/credentials")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Secret: %s\n", secret.Value)
}
```

#### **Node.js SDK Example**

```typescript
import { AetherVaultClient } from "@aether-vault/node";

async function main() {
  // Initialize client
  const client = new AetherVaultClient({
    baseURL: "https://vault.example.com",
    apiKey: "your-api-key",
  });

  // Authenticate
  await client.auth.login({
    username: "admin",
    password: "password",
  });

  // Get a secret
  const secret = await client.secrets.get("database/credentials");
  console.log("Secret:", secret.value);
}

main().catch(console.error);
```

#### **CLI Usage Example**

```bash
# Initialize configuration
aether-vault init --endpoint https://vault.example.com

# Authenticate
aether-vault auth login

# List secrets
aether-vault secrets list

# Get a specific secret
aether-vault secrets get database/credentials

# Manage runtime
aether-vault runtime start --docker
aether-vault runtime status
```

---

## 🤝 Contributing

We welcome contributions to the Aether Vault Package Ecosystem! Whether you're interested in improving existing SDKs, adding new language support, or enhancing documentation, there's a place for you.

### 🎯 **How to Get Started**

1. **Choose a package** - Go, Node.js, GitHub, Docker, or CLI
2. **Read the package-specific README** - Understand conventions
3. **Fork the repository** and create a feature branch
4. **Follow our coding standards** - Go fmt, Prettier, ESLint
5. **Add tests** - Ensure comprehensive test coverage
6. **Submit a pull request** with clear description

### 🏗️ **Areas Needing Help**

- **Go SDK Development** - Performance optimization, new features
- **Node.js SDK Enhancement** - Browser compatibility, React integration
- **GitHub App Features** - Advanced workflow automation
- **CLI Tool Expansion** - Additional commands and functionality
- **Docker Runtime** - Kubernetes integration, monitoring
- **Documentation** - Examples, tutorials, API docs
- **Testing** - Unit tests, integration tests, E2E tests

### 📝 **Development Guidelines**

- **Language-Specific Standards** - Follow Go and TypeScript best practices
- **Semantic Versioning** - Use proper version management
- **API Consistency** - Maintain consistent interfaces across SDKs
- **Security First** - Validate all inputs and handle errors securely
- **Performance** - Optimize for speed and resource usage
- **Documentation** - Keep docs updated with code changes

---

## 📞 Support & Community

### 💬 **Get Help**

- 📖 **[Package Documentation](package/)** - Detailed guides for each package
- 🐛 **[GitHub Issues](https://github.com/skygenesisenterprise/aether-vault/issues)** - Bug reports and feature requests
- 💡 **[GitHub Discussions](https://github.com/skygenesisenterprise/aether-vault/discussions)** - General questions and ideas
- 📧 **Email** - packages@skygenesisenterprise.com

### 🐛 **Reporting Issues**

When reporting bugs, please include:

- Package name and version
- Clear description of the problem
- Steps to reproduce
- Environment information
- Error logs or stack traces
- Expected vs actual behavior

---

## 📊 Package Status

| Package                  | Status    | Language   | Platform        | Notes                    |
| ------------------------ | --------- | ---------- | --------------- | ------------------------ |
| **@aether-vault/golang** | ✅ Stable | Go         | Native          | High-performance SDK     |
| **@aether-vault/node**   | ✅ Stable | TypeScript | Node.js/Browser | Universal client         |
| **@aether-vault/github** | ✅ Stable | TypeScript | GitHub          | Verified Marketplace app |
| **@aether-vault/docker** | ✅ Stable | Go         | Docker          | Container runtime        |
| **@aether-vault/cli**    | ✅ Stable | Go         | CLI             | Cross-platform tool      |

---

## 📄 License

All packages in this ecosystem are licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2025 Sky Genesis Enterprise

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

---

## 🙏 Acknowledgments

- **Sky Genesis Enterprise** - Project leadership and package ecosystem development
- **Go Community** - Excellent programming language and tooling
- **TypeScript Team** - Type-safe JavaScript development
- **GitHub** - Platform integration and marketplace support
- **Docker Team** - Container platform and runtime
- **Node.js Community** - Server-side JavaScript ecosystem
- **Open Source Contributors** - Code, feedback, and improvements

---

<div align="center">

### 🚀 **Choose Your Package and Start Building with Aether Vault!**

[📦 Go SDK](package/golang/) • [📦 Node.js SDK](package/node/) • [🐙 GitHub App](package/github/) • [🐳 Docker Runtime](package/docker/) • [⚡ CLI Tools](package/cli/)

---

**🔧 Comprehensive Multi-Language SDK Ecosystem for Enterprise Vault Integration**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

_Building secure, scalable, and developer-friendly vault integration packages_

</div>
