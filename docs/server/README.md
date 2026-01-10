<div align="center">

# 🚀 Aether Vault Server

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![Go](https://img.shields.io/badge/Go-1.25+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![Gin](https://img.shields.io/badge/Gin-1.9+-lightgrey?style=for-the-badge&logo=go)](https://gin-gonic.com/) [![GORM](https://img.shields.io/badge/GORM-1.25+-green?style=for-the-badge&logo=go)](https://gorm.io/) [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/) [![JWT](https://img.shields.io/badge/JWT-Verified-green?style=for-the-badge&logo=jwt)](https://jwt.io/)

**🔥 Enterprise-Grade Authentication & Vault Server - Complete Security Architecture with Advanced Identity Management**

A comprehensive authentication and vault management server built for enterprise security. Features **complete authentication system**, **encrypted secret management**, **TOTP 2FA support**, **comprehensive audit logging**, and **enterprise-grade security architecture**.

[🚀 Quick Start](#-quick-start) • [📋 Features](#-features) • [🛠️ Tech Stack](#️-tech-stack) • [🔐 Security](#-security) • [📁 Architecture](#-architecture) • [📚 Documentation](#-documentation) • [🤝 Contributing](#-contributing)

[![GitHub stars](https://img.shields.io/github/stars/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/stargazers) [![GitHub forks](https://img.shields.io/github/forks/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/network) [![GitHub issues](https://img.shields.io/github/issues/github/skygenesisenterprise/aether-vault)](https://github.com/skygenesisenterprise/aether-vault/issues)

</div>

---

## 🌟 What is Aether Vault Server?

**Aether Vault Server** is a production-ready authentication and vault management system designed for enterprise security. It provides a complete ecosystem for identity management, secret storage, and security compliance with advanced features typically found in commercial solutions.

### 🎯 **Core Capabilities**

- **🔐 Complete Authentication** - JWT tokens with TOTP 2FA and session management
- **🗄️ Secret Management** - AES-256 encrypted storage with access controls
- **🌐 Network Management** - Multi-protocol connectivity testing and monitoring
- **🛡️ Enterprise Security** - Rate limiting, input validation, comprehensive audit logging
- **📊 Compliance Ready** - Complete audit trails and compliance reporting
- **⚡ High Performance** - Go-based with PostgreSQL and connection pooling
- **🏗️ Scalable Architecture** - Clean architecture with horizontal scaling support

---

## 🚀 Quick Start

### 📋 **Prerequisites**

- **Go** 1.25.0 or higher
- **PostgreSQL** 15.0 or higher
- **Docker** (optional, for containerized deployment)
- **Make** (included with most systems)

### 🔧 **Installation & Setup**

#### **1. Clone & Setup**

```bash
git clone https://github.com/skygenesisenterprise/aether-vault.git
cd aether-vault/server
```

#### **2. Quick Start (Recommended)**

```bash
# One-command setup and start
make quick-start
```

#### **3. Manual Setup**

```bash
# Install dependencies
go mod download && go mod tidy

# Environment setup
cp .env.example .env
# Edit .env with your configuration

# Database setup
make db-migrate

# Start server
make go-server
```

### 🌐 **Access Points**

Once running, access the server at:

- **API Server**: [http://localhost:8080](http://localhost:8080)
- **Health Check**: [http://localhost:8080/api/v1/system/health](http://localhost:8080/api/v1/system/health)
- **API Documentation**: [http://localhost:8080/api/v1/system/version](http://localhost:8080/api/v1/system/version)

### 🎯 **Essential Commands**

```bash
# 🚀 Development
make quick-start          # Setup and start server
make go-dev               # Development mode with hot reload
make go-server            # Start production server

# 🗄️ Database
make db-migrate           # Run migrations
make db-seed              # Seed development data
make db-studio            # Open database admin tool

# 🔧 Code Quality
make go-test              # Run tests
make go-lint              # Run linter
make go-fmt               # Format code

# 🐳 Docker
make docker-build         # Build Docker image
make docker-run            # Run with Docker Compose

# 🛠️ Utilities
make help                 # Show all commands
make health               # Check service health
make status               # Show project status
```

---

## 📋 Features

### 🔐 **Authentication & Authorization**

- **JWT Token System** - Secure token-based authentication with rotation
- **TOTP 2FA** - Time-based one-time password support
- **Session Management** - Secure session handling with refresh tokens
- **Password Security** - bcrypt with configurable security policies
- **Access Control** - Role-based permissions with fine-grained control

### 🗄️ **Secret Management**

- **AES-256 Encryption** - Military-grade encryption for secret storage
- **Access Policies** - Configurable rules for secret access
- **Version Control** - Track changes and maintain secret history
- **Audit Trail** - Complete logging of all secret operations
- **Key Management** - Secure key derivation and rotation

### 🛡️ **Security Infrastructure**

- **Rate Limiting** - Configurable limits per endpoint and user
- **Input Validation** - Comprehensive request validation and sanitization
- **Security Headers** - CORS, CSP, HSTS, and other security headers
- **Request Tracking** - Correlation IDs for debugging and auditing
- **IP Controls** - Configurable IP access controls

### 🌐 **Network Management**

- **Multi-Protocol Support** - HTTP, HTTPS, SSH, FTP, SFTP, WebDAV, SMB, NFS, RSYNC, Git, Custom
- **Connectivity Testing** - Real-time protocol testing with latency measurement
- **Network Monitoring** - Status tracking and health monitoring for network endpoints
- **Protocol Validation** - Security validation and whitelist/blacklist support
- **Connection Limits** - Configurable concurrent connection limits and timeouts

### 📊 **Enterprise Audit & Compliance**

- **Complete Audit Logging** - All operations logged with full context
- **Security Event Tracking** - Failed logins, access attempts, policy violations
- **Real-time Monitoring** - Live audit stream and configurable alerting
- **Compliance Reporting** - Generate reports for security audits
- **Log Retention** - Configurable retention policies and archiving

---

## 🛠️ Tech Stack

### ⚙️ **Backend Architecture**

```
Go 1.25+ + Gin Framework
├── 🗄️ GORM + PostgreSQL (Database Layer)
├── 🔐 JWT Authentication (Complete Implementation)
├── 🛡️ Security Middleware (Rate Limiting, CORS, Headers)
├── 🌐 HTTP Router (Gin Router with Middleware)
├── 🌐 Network Management (Multi-Protocol Support & Testing)
├── 📊 Structured Logging (Context-aware logging)
├── 🔐 TOTP 2FA (Time-based One-Time Password)
├── 🗄️ Secret Management (AES-256 Encrypted Storage)
├── 📊 Audit Logging (Complete Event Tracking)
└── 🏗️ Policy Engine (Access Control Rules)
```

### 🗄️ **Data Layer**

```
PostgreSQL + GORM
├── 🏗️ Schema Management (Auto-migration)
├── 🔍 Query Builder (Type-Safe Queries)
├── 🔄 Connection Pooling (Performance Optimization)
├── 👤 User Models (Complete Authentication Models)
├── 🗄️ Secret Models (Encrypted Storage Models)
├── 📊 Audit Models (Complete Audit Trail)
├── 📈 Seed Scripts (Development Data)
└── 🔐 Security Models (TOTP, Sessions, Policies)
```

### 🔐 **Security Stack**

```
Enterprise Security Architecture
├── 🛡️ Authentication (JWT + TOTP 2FA)
├── 🔐 Encryption (AES-256 Secret Storage)
├── 🚦 Rate Limiting (Configurable Limits)
├── 🔍 Input Validation (Comprehensive Validation)
├── 🌐 Network Security (Protocol Validation & Monitoring)
├── 🌐 CORS (Cross-Origin Resource Sharing)
├── 📋 Security Headers (CSP, HSTS, X-Frame-Options)
├── 📊 Audit Logging (Complete Event Tracking)
├── 🏗️ Policy Engine (Access Control Rules)
└── 🚨 Security Monitoring (Real-time Alerts)
```

---

## 🔐 Security

### 🎯 **Enterprise-Grade Security Architecture**

Aether Vault implements defense-in-depth security with multiple protection layers:

```
Security Layers
├── 🌐 Network Security
│   ├── HTTPS/TLS Encryption
│   ├── Security Headers
│   └── Rate Limiting
├── 🔐 Authentication Security
│   ├── JWT Token Management
│   ├── TOTP 2FA Support
│   └── Session Management
├── 🗄️ Data Security
│   ├── AES-256 Encryption
│   ├── Secure Password Hashing
│   └── Database Encryption
├── 📊 Audit Security
│   ├── Complete Event Logging
│   ├── Security Monitoring
│   └── Compliance Reporting
└── 🏗️ Application Security
    ├── Input Validation
    ├── Access Control
    └── Policy Enforcement
```

### 🔑 **Key Security Features**

- **Zero Trust Architecture** - Never trust, always verify
- **Encryption Everywhere** - Data encrypted at rest and in transit
- **Complete Audit Trail** - All actions logged and traceable
- **Multi-Factor Authentication** - JWT + TOTP 2FA support
- **Enterprise Compliance** - GDPR, SOC 2, ISO 27001 ready
- **Real-time Monitoring** - Security event tracking and alerting

---

## 📁 Architecture

### 🏗️ **Clean Architecture Design**

```
aether-vault/server/
├── cmd/                     # CLI Entry Points
│   └── server/
│       └── main.go         # Main server entry point
├── src/                     # Source Code
│   ├── config/             # Configuration Management
│   │   └── config.go       # Server, database, security config
│   ├── controllers/        # HTTP Request Handlers
│   │   ├── auth.go         # Authentication endpoints
│   │   ├── user.go         # User management endpoints
│   │   ├── secret.go       # Secret management endpoints
│   │   ├── totp.go         # TOTP 2FA endpoints
│   │   ├── audit.go        # Audit logging endpoints
│   │   ├── system.go       # System health and metrics
│   │   ├── identity.go     # Identity management endpoints
│   │   ├── network.go      # Network management endpoints
│   │   └── policy.go       # Policy management endpoints
│   ├── middleware/         # HTTP Middleware Stack
│   │   ├── auth.go         # JWT authentication middleware
│   │   ├── security.go     # Security headers and validation
│   │   ├── ratelimit.go    # Rate limiting middleware
│   │   ├── audit.go        # Audit logging middleware
│   │   ├── network.go      # Network protocol validation middleware
│   │   ├── user.go         # User context middleware
│   │   └── utils.go        # Utility middleware functions
│   ├── model/              # Data Models & DTOs
│   │   ├── user.go         # User model and structs
│   │   ├── secret.go       # Secret management models
│   │   ├── totp.go         # TOTP configuration models
│   │   ├── audit.go        # Audit log models
│   │   ├── network.go      # Network configuration models
│   │   ├── policy.go       # Policy and rule models
│   │   └── dto.go          # Data Transfer Objects
│   ├── routes/             # Route Definitions
│   │   └── routes.go       # API route configuration
│   ├── services/           # Business Logic Layer
│   │   ├── auth.go         # Authentication service logic
│   │   ├── user.go         # User management service
│   │   ├── secret.go       # Secret management service
│   │   ├── totp.go         # TOTP/2FA service logic
│   │   ├── audit.go        # Audit logging service
│   │   ├── network.go      # Network management service
│   │   ├── policy.go       # Policy enforcement service
│   │   └── system.go       # System monitoring service
│   └── utils/              # Utility Functions
│       ├── crypto.go       # Cryptographic utilities
│       └── logger.go       # Logging utilities
├── utils/                  # Shared Utilities
│   ├── crypto.go           # Cryptographic functions
│   └── logger.go           # Logging configuration
├── main.go                 # Server Entry Point
├── go.mod                  # Go Modules File
├── go.sum                  # Go Modules Checksum
├── .env.example            # Environment Template
├── Dockerfile              # Docker Configuration
├── docker-compose.yml      # Docker Compose Setup
└── README.md               # Server Documentation
```

### 🔄 **Request Flow Architecture**

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │───►│  Middleware │───►│ Controllers │───►│  Services   │
│   Request   │    │   Stack     │    │             │    │   Layer     │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │                   │
       ▼                   ▼                   ▼                   ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Security  │    │   Auth &    │    │   Business  │    │   Data      │
│   Headers   │    │   Validation │    │   Logic     │    │   Access    │
│   CORS      │    │   Rate Limit│    │   Processing │    │   Layer     │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                                                                 │
                                                                 ▼
                                                    ┌─────────────┐
                                                    │   GORM      │
                                                    │  PostgreSQL │
                                                    │  Database   │
                                                    │   Layer     │
                                                    └─────────────┘
                                                                 │
                                                                 ▼
                                                    ┌─────────────┐
                                                    │   Audit     │
                                                    │  Logging    │
                                                    │   System    │
                                                    │   Layer     │
                                                    └─────────────┘
```

---

## 📚 Documentation

### 📖 **Complete Documentation Set**

- **[📚 API Documentation](./api.md)** - Complete REST API reference
- **[🏗️ Architecture Guide](./architecture.md)** - Detailed system architecture
- **[⚙️ Configuration Guide](./configuration.md)** - Setup and configuration
- **[🚀 Deployment Guide](./deployment.md)** - Production deployment
- **[🔒 Security Guide](./security.md)** - Security architecture and best practices

### 🎯 **Quick Links**

| Topic             | Description                          | Link                                   |
| ----------------- | ------------------------------------ | -------------------------------------- |
| **API Reference** | Complete REST API documentation      | [📚 API Docs](./api.md)                |
| **Architecture**  | System design and architecture       | [🏗️ Architecture](./architecture.md)   |
| **Configuration** | Setup and environment configuration  | [⚙️ Configuration](./configuration.md) |
| **Deployment**    | Production deployment guide          | [🚀 Deployment](./deployment.md)       |
| **Security**      | Security features and best practices | [🔒 Security](./security.md)           |

---

## 📊 Project Status

### ✅ **Currently Implemented**

| Component                 | Status     | Technology            | Features                            |
| ------------------------- | ---------- | --------------------- | ----------------------------------- |
| **Go Backend Server**     | ✅ Working | Go + Gin              | High-performance API server         |
| **Authentication System** | ✅ Working | JWT + TOTP 2FA        | Complete enterprise auth            |
| **Secret Management**     | ✅ Working | AES-256 + GORM        | Encrypted storage with policies     |
| **Security Middleware**   | ✅ Working | Custom Gin Middleware | Rate limiting, validation, headers  |
| **Database Layer**        | ✅ Working | GORM + PostgreSQL     | Auto-migrations + enterprise models |
| **Audit System**          | ✅ Working | Custom Logging        | Complete event tracking             |
| **Policy Engine**         | ✅ Working | Custom Rules Engine   | Configurable access control         |
| **TOTP 2FA**              | ✅ Working | Custom TOTP           | Time-based one-time passwords       |
| **Session Management**    | ✅ Working | JWT + Refresh Tokens  | Secure session handling             |
| **Rate Limiting**         | ✅ Working | Custom Middleware     | Configurable limits and bursts      |
| **Docker Deployment**     | ✅ Working | Multi-Stage           | Production-ready containers         |

### 🔄 **In Development**

- **Advanced Policy Engine** - Complex rule evaluation and dynamic policies
- **Multi-Tenant Support** - Organization-based isolation and management
- **API Documentation** - Comprehensive OpenAPI/Swagger documentation
- **Performance Optimization** - Caching, connection pooling, query optimization
- **Testing Suite** - Unit and integration tests with coverage reporting

### 📋 **Planned Features**

- **OAuth2/OpenID Connect** - Standard federation protocols
- **LDAP/Active Directory Integration** - Enterprise directory services
- **WebAuthn/FIDO2** - Hardware-based authentication
- **Advanced Analytics** - Usage patterns and security insights
- **High Availability** - Clustering, failover, and load balancing

---

## 💻 Development

### 🎯 **Development Workflow**

```bash
# New developer setup
make quick-start

# Daily development
make go-dev               # Start development server
make go-fmt               # Format code
make go-lint              # Check code quality
make go-test              # Run tests

# Database changes
make db-migrate           # Apply migrations
make db-seed              # Seed development data
make db-studio            # Browse database

# Before committing
make go-fmt               # Format code
make go-lint              # Check code quality
make go-test              # Run tests
make go-vet               # Static analysis

# Production deployment
make go-build-linux       # Build for production
make docker-build         # Create Docker image
make docker-run           # Deploy
```

### 📋 **Development Guidelines**

- **Go Best Practices** - Follow Go conventions and idiomatic patterns
- **Security First** - Validate all inputs and implement proper authentication
- **Error Handling** - Comprehensive error handling and logging
- **Testing** - Write unit tests for all business logic
- **Documentation** - Maintain comprehensive API documentation
- **Code Quality** - Use gofmt, golangci-lint, and go vet regularly
- **Database Design** - Use proper indexing and constraints
- **API Design** - RESTful endpoints with proper HTTP methods and status codes

---

## 🤝 Contributing

We're looking for contributors to help build this comprehensive enterprise authentication and vault server! Whether you're experienced with Go, security, authentication systems, database design, or enterprise infrastructure, there's a place for you.

### 🎯 **How to Get Started**

1. **Fork the repository** and create a feature branch
2. **Check the issues** for tasks that need help
3. **Join discussions** about architecture and features
4. **Start small** - Documentation, tests, or minor features
5. **Follow our code standards** and commit guidelines

### 🏗️ **Areas Needing Help**

- **Go Backend Development** - API endpoints, business logic, security features
- **Security Specialists** - Authentication, encryption, audit systems, TOTP
- **Database Design** - Schema development, migrations, optimization
- **Enterprise Integration** - LDAP, OAuth2, SAML, federation protocols
- **DevOps Engineers** - Docker, deployment, CI/CD, monitoring
- **Security Experts** - Penetration testing, security audits, compliance
- **Documentation** - API docs, security guides, deployment tutorials
- **Testing** - Unit tests, integration tests, security testing

### 📝 **Contribution Process**

1. **Choose an area** - Core authentication, secret management, or security features
2. **Read the documentation** - Understand the architecture and conventions
3. **Create a branch** with a descriptive name following our standards
4. **Implement your changes** following Go best practices and security guidelines
5. **Test thoroughly** - Include unit tests and security considerations
6. **Submit a pull request** with clear description and testing instructions
7. **Address feedback** from maintainers and security review

---

## 📞 Support & Community

### 💬 **Get Help**

- 📖 **[Documentation](.)** - Comprehensive guides and API docs
- 🐛 **[GitHub Issues](https://github.com/skygenesisenterprise/aether-vault/issues)** - Bug reports and feature requests
- 💡 **[GitHub Discussions](https://github.com/skygenesisenterprise/aether-vault/discussions)** - General questions and ideas
- 📧 **Email** - support@skygenesisenterprise.com

### 🐛 **Reporting Issues**

When reporting bugs, please include:

- Clear description of the problem
- Steps to reproduce
- Environment information (Go version, PostgreSQL version, OS, etc.)
- Error logs or stack traces
- Expected vs actual behavior
- Security considerations (if applicable)

---

## 🏆 Sponsors & Partners

**Development led by [Sky Genesis Enterprise](https://skygenesisenterprise.com)**

We're looking for sponsors and partners to help accelerate development of this open-source enterprise authentication and vault server project.

[🤝 Become a Sponsor](https://github.com/sponsors/skygenesisenterprise)

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](../../LICENSE) file for details.

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

- **Sky Genesis Enterprise** - Project leadership and security architecture
- **Go Community** - High-performance programming language and ecosystem
- **Gin Framework** - Lightweight HTTP web framework
- **GORM Team** - Modern Go database library
- **PostgreSQL Team** - Powerful relational database
- **JWT Community** - Secure token-based authentication standard
- **Docker Team** - Container platform and tools
- **Open Source Community** - Tools, libraries, and security inspiration

---

<div align="center">

### 🚀 **Join Us in Building the Future of Enterprise Authentication & Vault Management!**

[⭐ Star This Repo](https://github.com/skygenesisenterprise/aether-vault) • [🐛 Report Issues](https://github.com/skygenesisenterprise/aether-vault/issues) • [💡 Start a Discussion](https://github.com/skygenesisenterprise/aether-vault/discussions)

---

**🔧 Enterprise-Grade Security with Advanced Authentication & Complete Vault Management!**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

_Building secure, scalable authentication and vault management solutions for the enterprise_

</div>
