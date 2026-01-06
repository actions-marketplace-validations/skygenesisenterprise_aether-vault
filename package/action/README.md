<div align="center">

# 🔐 Aether Vault Action

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![Go](https://img.shields.io/badge/Go-1.25+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![GitHub Action](https://img.shields.io/badge/GitHub_Action-Verified-green?style=for-the-badge&logo=github)](https://github.com/marketplace) [![OIDC](https://img.shields.io/badge/OIDC-Enabled-orange?style=for-the-badge)](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect)

**🚀 Enterprise-Ready GitHub Action for Secure Aether Vault Authentication**

A comprehensive GitHub Action that provides **secure, zero-knowledge authentication** with Aether Vault using **GitHub OIDC tokens**. Features enterprise-grade security, multi-platform support, and seamless integration into CI/CD pipelines.

[🚀 Quick Start](#-quick-start) • [📋 Features](#-features) • [🛠️ Tech Stack](#️-tech-stack) • [📁 Architecture](#-architecture) • [🔧 Usage](#-usage) • [🤝 Contributing](#-contributing)

[![GitHub stars](https://img.shields.io/github/stars/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/stargazers) [![GitHub forks](https://img.shields.io/github/forks/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/network) [![GitHub issues](https://img.shields.io/github/issues/github/skygenesisenterprise/aether-vault)](https://github.com/skygenesisenterprise/aether-vault/issues)

</div>

---

## 🌟 What is Aether Vault Action?

**Aether Vault Action** is a **secure, enterprise-grade GitHub Action** that provides **zero-knowledge authentication** with Aether Vault using **GitHub OIDC tokens**. It eliminates the need for static secrets while maintaining the highest security standards for CI/CD pipelines.

### 🎯 Our Security-First Vision

- **🔐 Zero-Knowledge Authentication** - No static secrets stored in repository
- **🚀 GitHub OIDC Integration** - Native token exchange with GitHub Actions
- **⚡ Multi-Platform Support** - Linux, macOS (amd64, arm64) binaries
- **🛡️ Enterprise-Grade Security** - Role-based access, short-lived tokens
- **📋 Policy Enforcement** - Security policy checks with detailed reporting
- **🔗 Seamless Integration** - Drop-in replacement for existing workflows
- **📊 Comprehensive Auditing** - Detailed logs and correlation IDs
- **🎨 Developer-Friendly** - Simple configuration, clear error messages

---

## 🆕 What's New - Latest Features

### 🎯 **Major Enhancements in v1.0+**

#### 🔐 **Enhanced Security Model** (NEW)

- ✅ **GitHub OIDC Authentication** - Native token exchange without static secrets
- ✅ **Short-Lived Tokens** - Ephemeral Vault tokens with configurable TTL
- ✅ **Role-Based Access** - Fine-grained permissions through Vault roles
- ✅ **Policy Enforcement** - Security policy checks with violation reporting
- ✅ **Audit Trail** - Complete logging with correlation IDs

#### 🚀 **Multi-Platform Architecture** (IMPROVED)

- ✅ **Cross-Platform Binaries** - Linux and macOS support
- ✅ **Multi-Architecture** - amd64 and arm64 binaries included
- ✅ **Composite Action** - Efficient GitHub Actions implementation
- ✅ **Go-Based Backend** - High-performance native Go binary

#### 📊 **Enhanced User Experience** (IMPROVED)

- ✅ **Intuitive Configuration** - Simple, well-documented inputs
- ✅ **Clear Error Messages** - Helpful debugging information
- ✅ **Comprehensive Examples** - Real-world usage patterns
- ✅ **Detailed Documentation** - Complete API and usage guides

---

## 📊 Current Status

> **✅ Production Ready**: Enterprise-grade security with GitHub OIDC integration.

### ✅ **Currently Implemented**

#### 🔐 **Core Security Features**

- ✅ **GitHub OIDC Authentication** - Complete token exchange implementation
- ✅ **Zero-Knowledge Model** - No static secrets in repository
- ✅ **Multi-Platform Support** - Linux/macOS binaries (amd64/arm64)
- ✅ **Policy Enforcement** - Security checks with detailed reporting
- ✅ **Audit Logging** - Structured logs with correlation IDs

#### 🛠️ **Technical Implementation**

- ✅ **Go Backend Binary** - High-performance native implementation
- ✅ **Composite Action** - Efficient GitHub Actions integration
- ✅ **Multi-Architecture** - Cross-platform binary compilation
- ✅ **Error Handling** - Comprehensive error reporting
- ✅ **Input Validation** - Secure configuration validation

#### 📋 **Developer Experience**

- ✅ **Simple Configuration** - Intuitive input parameters
- ✅ **Clear Documentation** - Comprehensive usage guides
- ✅ **Real-World Examples** - Practical implementation patterns
- ✅ **Debugging Support** - Verbose logging and troubleshooting

### 🔄 **In Development**

- **Windows Platform Support** - Additional platform coverage
- **Advanced Policy Engine** - Enhanced security rule evaluation
- **Performance Optimizations** - Caching and connection pooling
- **Extended Audit Features** - Enhanced reporting capabilities

### 📋 **Planned Features**

- **Multi-Vault Support** - Connect to multiple Vault instances
- **Secret Injection** - Automatic secret population in workflows
- **Webhook Integration** - Real-time security event notifications
- **Advanced Analytics** - Security metrics and insights

---

## 🚀 Quick Start

### 📋 Prerequisites

- **GitHub Repository** with OIDC enabled
- **Aether Vault Server** with GitHub OIDC configuration
- **Vault Role** configured for your repository
- **GitHub Actions** workflow permissions

### 🔧 Basic Usage

1. **Enable OIDC in your GitHub repository**

2. **Configure Vault role for your repository**

3. **Add the action to your workflow**

```yaml
name: Deploy with Security
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      id-token: write # Required for OIDC
      contents: read

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Authenticate with Aether Vault
        uses: skygenesisenterprise/aether-vault@v1
        with:
          vault-url: ${{ secrets.VAULT_URL }}
          auth-method: github-oidc
          role: my-app-role
          policy-mode: enforce
```

### 🎯 **Advanced Configuration**

```yaml
- name: Advanced Aether Vault Authentication
  uses: skygenesisenterprise/aether-vault@v1
  with:
    vault-url: ${{ secrets.VAULT_URL }}
    auth-method: github-oidc
    role: production-deploy-role
    policy-mode: enforce
    audience: aether-vault
    allow-token-output: false
  env:
    LOG_LEVEL: debug
```

---

## 🛠️ Tech Stack

### 🔐 **Security Layer**

```
GitHub OIDC + Aether Vault
├── 🎯 JWT Token Exchange (GitHub Actions)
├── 🔐 Role-Based Authentication (Vault)
├── 📋 Policy Enforcement (Security Rules)
├── 📊 Audit Logging (Structured JSON)
└── 🚀 Zero-Knowledge Model (No Static Secrets)
```

### ⚙️ **Implementation Layer**

```
Go 1.25+ + GitHub Actions
├── 🐹 Native Go Binary (High Performance)
├── 🔗 Composite Action (Efficient Integration)
├── 🌐 HTTP Client (Resty Library)
├── 📝 Structured Logging (Logrus)
└── 🔧 Input Validation (Go Validation)
```

### 🏗️ **Platform Support**

```
Multi-Platform Architecture
├── 🐧 Linux (amd64, arm64)
├── 🍎 macOS (amd64, arm64)
├── 🪟 Windows (Planned)
└── 🐳 Docker (Container Support)
```

---

## 📁 Architecture

### 🏗️ **Action Structure**

```
package/action/
├── action.yml                 # GitHub Action definition
├── cmd/
│   └── main.go                # Go binary entry point
├── internal/
│   ├── auth/                  # OIDC authentication
│   │   └── oidc.go           # GitHub token exchange
│   ├── config/                # Configuration management
│   │   └── config.go         # Environment validation
│   ├── vault/                 # Vault API client
│   │   └── client.go         # API communication
│   ├── github/                # GitHub context
│   │   └── context.go        # Runtime information
│   └── output/                # Output management
│       └── manager.go        # GitHub outputs
├── bin/                       # Pre-compiled binaries
├── go.mod                     # Go modules
├── go.sum                     # Dependencies checksum
├── Makefile                   # Build automation
├── README.md                  # This documentation
├── USAGE_EXAMPLES.md          # Usage examples
└── LICENSE                    # MIT License
```

### 🔄 **Authentication Flow**

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   GitHub Actions │    │   Aether Vault   │    │   Policy Engine  │
│   (OIDC Provider)│◄──►│   (Token Store)  │◄──►│   (Security)     │
│                  │    │                  │    │                 │
│ • JWT Token     │    │ • Role Mapping   │    │ • Policy Rules   │
│ • Repository    │    │ • Token Exchange │    │ • Violation Check│
│ • Workflow      │    │ • Short-Lived    │    │ • Audit Report  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
            │                       │                       │
            ▼                       ▼                       ▼
     ┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
     │  Go Binary      │    │  Vault API       │    │  Security Output│
     │  (Authentication)│   │  (Communication) │   │  (Results)      │
     │                 │    │                  │    │                 │
     │ • Token Validation│ │ • HTTP Requests  │    │ • Status Report  │
     │ • Error Handling │ │ • Response Parse  │    │ • Violation Details│
     │ • Logging        │ │ • Retry Logic     │    │ • Audit Trail    │
     └─────────────────┘    └──────────────────┘    └─────────────────┘
```

---

## 🔧 Usage

### 📋 **Input Parameters**

| Parameter            | Required | Default        | Description                                       |
| -------------------- | -------- | -------------- | ------------------------------------------------- | -------- |
| `vault-url`          | ✅ Yes   | -              | Aether Vault server URL                           |
| `auth-method`        | ❌ No    | `github-oidc`  | Authentication method                             |
| `role`               | ❌ No    | -              | Vault role for authentication                     |
| `policy-mode`        | ❌ No    | `enforce`      | Policy enforcement mode (`enforce`                | `audit`) |
| `audience`           | ❌ No    | `aether-vault` | OIDC audience for token exchange                  |
| `allow-token-output` | ❌ No    | `false`        | Allow vault token in outputs (security-sensitive) |

### 📊 **Output Parameters**

| Parameter     | Description                        |
| ------------- | ---------------------------------- | ------------ |
| `status`      | Policy check status (`success`     | `violation`) |
| `report-id`   | Audit report ID for correlation    |
| `vault-token` | Vault token (if allowed by policy) |

### 🎯 **Usage Examples**

#### **Basic Authentication**

```yaml
- name: Authenticate with Aether Vault
  uses: skygenesisenterprise/aether-vault@v1
  with:
    vault-url: ${{ secrets.VAULT_URL }}
    role: my-app-role
```

#### **Policy Enforcement**

```yaml
- name: Security Policy Check
  uses: skygenesisenterprise/aether-vault@v1
  with:
    vault-url: ${{ secrets.VAULT_URL }}
    role: security-check-role
    policy-mode: enforce
    audience: my-audience
```

#### **Audit Mode**

```yaml
- name: Security Audit
  uses: skygenesisenterprise/aether-vault@v1
  with:
    vault-url: ${{ secrets.VAULT_URL }}
    role: audit-role
    policy-mode: audit
    allow-token-output: true
```

#### **Debug Mode**

```yaml
- name: Debug Authentication
  uses: skygenesisenterprise/aether-vault@v1
  with:
    vault-url: ${{ secrets.VAULT_URL }}
    role: debug-role
  env:
    LOG_LEVEL: debug
```

---

## 🔧 Development

### 🎯 **Build Commands**

```bash
# Build all platforms
make build

# Build specific platform
make linux-amd64
make linux-arm64
make darwin-amd64
make darwin-arm64

# Build current platform
make build-local

# Run tests
make test

# Run security checks
make security

# Format code
make fmt

# Lint code
make lint

# Create release package
make release
```

### 📋 **Development Workflow**

```bash
# Setup development environment
cd package/action
go mod tidy

# Run tests
go test -v ./...

# Build binary
make build-local

# Test locally
export VAULT_URL="https://vault.dev.local"
export AUTH_METHOD="github-oidc"
./bin/aether-vault-linux-amd64
```

### 🐛 **Testing**

```bash
# Unit tests
go test ./internal/...

# Integration tests (requires Vault)
VAULT_URL="https://vault.test.local" go test ./...

# Security checks
make security

# Code quality
make check
```

---

## 🔐 Security Considerations

### 🛡️ **Token Security**

- **No Static Secrets** - Uses GitHub OIDC tokens exclusively
- **Short-Lived Tokens** - Vault tokens have configurable TTL
- **Role-Based Access** - Fine-grained permissions through Vault roles
- **Secure Token Handling** - Tokens never logged or exposed unnecessarily

### 🔒 **OIDC Security**

- **JWT Validation** - Signature and claims verification
- **Repository Scoping** - Tokens scoped to specific repository
- **Audience Verification** - Ensures tokens are for intended purpose
- **Expiration Handling** - Automatic token refresh and expiration

### 🌐 **Network Security**

- **HTTPS-Only** - All communication encrypted
- **Certificate Validation** - Proper SSL/TLS certificate verification
- **Request Timeouts** - Configurable timeout limits
- **Retry Logic** - Resilient handling of transient failures

---

## 🚨 Troubleshooting

### 🔧 **Common Issues**

#### **OIDC Token Not Available**

```bash
# Ensure repository has OIDC enabled
# Check GitHub repository settings > Actions > General
# Verify "Actions" is not disabled
```

#### **Vault Authentication Failed**

```bash
# Verify Vault URL accessibility
curl -v $VAULT_URL/v1/sys/health

# Check role configuration in Vault
vault read auth/github/role/your-role-name

# Validate audience setting
# Ensure audience matches Vault configuration
```

#### **Policy Violation**

```bash
# Enable debug logging
env:
  LOG_LEVEL: debug

# Check policy details in logs
# Review violation report in action outputs
```

### 🐛 **Debug Mode**

```yaml
- name: Debug Authentication
  uses: skygenesisenterprise/aether-vault@v1
  with:
    vault-url: ${{ secrets.VAULT_URL }}
    role: debug-role
  env:
    LOG_LEVEL: debug
    VAULT_DEBUG: true
```

---

## 🤝 Contributing

We're looking for contributors to help enhance this security-focused GitHub Action! Whether you're experienced with Go, GitHub Actions, security, or Vault integration, there's a place for you.

### 🎯 **How to Get Started**

1. **Fork the repository** and create a feature branch
2. **Read the documentation** and understand the security model
3. **Join discussions** about security enhancements and features
4. **Start small** - Documentation, tests, or security improvements
5. **Follow our security-first guidelines** and code standards

### 🏗️ **Areas Needing Help**

- **Go Development** - Core binary enhancements, security features
- **GitHub Actions Experts** - Workflow optimization, best practices
- **Security Specialists** - OIDC implementation, Vault integration
- **Documentation** - Security guides, usage examples, API docs
- **Testing** - Unit tests, integration tests, security testing
- **Platform Support** - Windows binaries, additional architectures

### 📝 **Contribution Process**

1. **Security First** - All changes must maintain security standards
2. **Create a branch** with a descriptive name
3. **Implement changes** following Go best practices
4. **Test thoroughly** including security testing
5. **Submit a pull request** with security-focused description
6. **Address feedback** from maintainers and security review

---

## 📞 Support & Community

### 💬 **Get Help**

- 📖 **[Documentation](README.md)** - Complete usage guide
- 📋 **[Usage Examples](USAGE_EXAMPLES.md)** - Real-world patterns
- 🐛 **[GitHub Issues](https://github.com/skygenesisenterprise/aether-vault/issues)** - Bug reports and feature requests
- 💡 **[GitHub Discussions](https://github.com/skygenesisenterprise/aether-vault/discussions)** - Security questions and ideas
- 📧 **Email** - security@skygenesisenterprise.com

### 🐛 **Reporting Security Issues**

For security vulnerabilities, please email us directly at **security@skygenesisenterprise.com** rather than opening public issues.

---

## 📊 Project Status

| Component                    | Status         | Technology        | Security          | Notes                        |
| ---------------------------- | -------------- | ----------------- | ----------------- | ---------------------------- |
| **Core Authentication**      | ✅ Working     | GitHub OIDC       | **Enterprise**    | Complete token exchange      |
| **Multi-Platform Support**   | ✅ Working     | Go Binaries       | **Hardened**      | Linux/macOS (amd64/arm64)    |
| **Policy Enforcement**       | ✅ Working     | Vault API         | **Strict**        | Security rule validation     |
| **Audit Logging**            | ✅ Working     | Structured JSON   | **Comprehensive** | Correlation IDs              |
| **Error Handling**           | ✅ Working     | Go Error Patterns | **Secure**        | No token leakage             |
| **Documentation**            | ✅ Working     | Markdown          | **Complete**      | Usage guides and examples    |
| **Testing Suite**            | 🔄 In Progress | Go Testing        | **Security**      | Unit and integration tests   |
| **Windows Support**          | 📋 Planned     | Go Compilation    | **Hardened**      | Additional platform coverage |
| **Performance Optimization** | 📋 Planned     | Go Profiling      | **Enhanced**      | Caching and pooling          |

---

## 🏆 Sponsors & Partners

**Development led by [Sky Genesis Enterprise](https://skygenesisenterprise.com)**

We're looking for security-focused sponsors and partners to help enhance this open-source security tool.

[🤝 Become a Sponsor](https://github.com/sponsors/skygenesisenterprise)

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

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

- **Sky Genesis Enterprise** - Project leadership and security expertise
- **GitHub Actions Team** - Excellent OIDC integration platform
- **Vault Team** - Secure secrets management solution
- **Go Community** - High-performance, security-focused programming language
- **Open Source Security Community** - Tools, libraries, and security best practices
- **Enterprise Security Experts** - Guidance and security review

---

<div align="center">

### 🚀 **Join Us in Building a More Secure CI/CD Future!**

[⭐ Star This Repo](https://github.com/skygenesisenterprise/aether-vault) • [🐛 Report Issues](https://github.com/skygenesisenterprise/aether-vault/issues) • [💡 Security Discussions](https://github.com/skygenesisenterprise/aether-vault/discussions)

---

**🔐 Enterprise-Grade Security with Zero-Knowledge Authentication!**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) security team**

_Building secure, enterprise-ready GitHub Actions with GitHub OIDC integration_

</div>
