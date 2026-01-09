<div align="center">

# ⚙️ Aether Vault Configuration Guide

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![Go](https://img.shields.io/badge/Go-1.25+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![Viper](https://img.shields.io/badge/Viper-Config-green?style=for-the-badge&logo=go)](https://github.com/spf13/viper) [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)

**🔥 Complete Configuration Guide for Enterprise Authentication & Vault Server**

Comprehensive configuration documentation for Aether Vault Server, covering all aspects of setup, environment variables, security settings, and deployment configuration.

[🚀 Quick Setup](#-quick-setup) • [🌍 Environment Variables](#-environment-variables) • [📝 Configuration Files](#-configuration-files) • [🔐 Security Configuration](#-security-configuration) • [🗄️ Database Configuration](#️-database-configuration) • [🚀 Production Setup](#-production-setup) • [🔧 Advanced Settings](#-advanced-settings)

</div>

---

## 🚀 Quick Setup

### 📋 **Prerequisites**

- **Go** 1.25.0 or higher
- **PostgreSQL** 15.0 or higher
- **Git** for version control
- **Make** (included with most systems)

### 🔧 **Fastest Setup**

```bash
# 1. Clone and setup
git clone https://github.com/skygenesisenterprise/aether-vault.git
cd aether-vault/server

# 2. Quick start with defaults
make quick-start

# 3. Server will be running on http://localhost:8080
```

### 📝 **Manual Setup**

```bash
# 1. Install dependencies
go mod download
go mod tidy

# 2. Copy environment template
cp .env.example .env

# 3. Edit configuration
nano .env

# 4. Setup database
make db-migrate

# 5. Start server
make go-server
```

---

## 🌍 Environment Variables

### 🎯 **Required Variables**

| Variable                        | Description           | Default    | Example                            |
| ------------------------------- | --------------------- | ---------- | ---------------------------------- |
| `VAULT_JWT_SECRET`              | JWT signing secret    | _required_ | `your-super-secret-jwt-key-here`   |
| `VAULT_SECURITY_ENCRYPTION_KEY` | Secret encryption key | _required_ | `your-32-character-encryption-key` |

### 🖥️ **Server Configuration**

| Variable                     | Description             | Default       | Example      |
| ---------------------------- | ----------------------- | ------------- | ------------ |
| `VAULT_SERVER_HOST`          | Server bind address     | `0.0.0.0`     | `127.0.0.1`  |
| `VAULT_SERVER_PORT`          | Server port             | `8080`        | `3000`       |
| `VAULT_SERVER_ENVIRONMENT`   | Environment mode        | `development` | `production` |
| `VAULT_SERVER_READ_TIMEOUT`  | Read timeout (seconds)  | `30`          | `60`         |
| `VAULT_SERVER_WRITE_TIMEOUT` | Write timeout (seconds) | `30`          | `60`         |

### 🗄️ **Database Configuration**

| Variable                  | Description       | Default     | Example              |
| ------------------------- | ----------------- | ----------- | -------------------- |
| `VAULT_DATABASE_HOST`     | Database host     | `localhost` | `db.example.com`     |
| `VAULT_DATABASE_PORT`     | Database port     | `5432`      | `5432`               |
| `VAULT_DATABASE_USER`     | Database user     | `vault`     | `vault_user`         |
| `VAULT_DATABASE_PASSWORD` | Database password | _empty_     | `secure_db_password` |
| `VAULT_DATABASE_DBNAME`   | Database name     | `vault`     | `production_vault`   |
| `VAULT_DATABASE_SSLMODE`  | SSL mode          | `disable`   | `require`            |

### 🔐 **Security Configuration**

| Variable                        | Description       | Default  | Example  |
| ------------------------------- | ----------------- | -------- | -------- |
| `VAULT_SECURITY_KDF_ITERATIONS` | PBKDF2 iterations | `100000` | `200000` |
| `VAULT_SECURITY_SALT_LENGTH`    | Salt length       | `32`     | `64`     |

### 🎟️ **JWT Configuration**

| Variable               | Description                | Default    | Example                |
| ---------------------- | -------------------------- | ---------- | ---------------------- |
| `VAULT_JWT_EXPIRATION` | Token expiration (seconds) | `3600`     | `7200`                 |
| `VAULT_JWT_SECRET`     | JWT signing secret         | _required_ | `your-jwt-secret-here` |

### 📊 **Audit Configuration**

| Variable                 | Description          | Default | Example |
| ------------------------ | -------------------- | ------- | ------- |
| `VAULT_AUDIT_ENABLED`    | Enable audit logging | `true`  | `true`  |
| `VAULT_AUDIT_LOG_LEVEL`  | Log level            | `info`  | `debug` |
| `VAULT_AUDIT_LOG_FORMAT` | Log format           | `json`  | `text`  |

---

## 📝 Configuration Files

### 🌍 **.env.example**

```bash
# =============================================================================
# Aether Vault Server Environment Configuration
# =============================================================================
# Copy this file to .env and update with your values
# =============================================================================

# Server Configuration
VAULT_SERVER_HOST=0.0.0.0
VAULT_SERVER_PORT=8080
VAULT_SERVER_ENVIRONMENT=development
VAULT_SERVER_READ_TIMEOUT=30
VAULT_SERVER_WRITE_TIMEOUT=30

# Database Configuration
VAULT_DATABASE_HOST=localhost
VAULT_DATABASE_PORT=5432
VAULT_DATABASE_USER=vault
VAULT_DATABASE_PASSWORD=your_secure_password
VAULT_DATABASE_DBNAME=vault
VAULT_DATABASE_SSLMODE=disable

# Security Configuration
VAULT_SECURITY_ENCRYPTION_KEY=your-32-character-encryption-key
VAULT_SECURITY_KDF_ITERATIONS=100000
VAULT_SECURITY_SALT_LENGTH=32

# JWT Configuration
VAULT_JWT_SECRET=your-super-secret-jwt-key-here-must-be-very-long
VAULT_JWT_EXPIRATION=3600

# Audit Configuration
VAULT_AUDIT_ENABLED=true
VAULT_AUDIT_LOG_LEVEL=info
VAULT_AUDIT_LOG_FORMAT=json

# Development Options
# VAULT_DATABASE_PASSWORD=dev_password
# VAULT_SECURITY_ENCRYPTION_KEY=dev_encryption_key_32_chars
# VAULT_JWT_SECRET=dev_jwt_secret_please_change_in_production
```

### 📄 **config.yaml** (Optional)

```yaml
# =============================================================================
# Aether Vault Server Configuration File
# =============================================================================
# This file is optional - environment variables take precedence
# =============================================================================

server:
  host: "0.0.0.0"
  port: 8080
  environment: "development"
  read_timeout: 30
  write_timeout: 30

database:
  host: "localhost"
  port: 5432
  user: "vault"
  password: ""
  dbname: "vault"
  sslmode: "disable"

security:
  encryption_key: ""
  kdf_iterations: 100000
  salt_length: 32

jwt:
  secret: ""
  expiration: 3600

audit:
  enabled: true
  log_level: "info"
  log_format: "json"
```

---

## 🔐 Security Configuration

### 🔑 **Generate Secure Keys**

#### JWT Secret

```bash
# Generate a secure JWT secret (32+ characters)
openssl rand -base64 32

# Or using Go
go run -c 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'
```

#### Encryption Key

```bash
# Generate a secure encryption key (32 characters)
openssl rand -hex 16

# Or using Go
go run -c 'package main; import ("crypto/rand"; "encoding/hex"; "fmt"); func main() { b := make([]byte, 16); rand.Read(b); fmt.Println(hex.EncodeToString(b)) }'
```

### 🛡️ **Security Best Practices**

#### **Production Security Checklist**

- [ ] Use strong, unique secrets (32+ characters)
- [ ] Enable SSL/TLS for all communications
- [ ] Use SSL mode `require` or `verify-full` for database
- [ ] Set environment to `production`
- [ ] Enable audit logging
- [ ] Use secure key derivation (100,000+ iterations)
- [ ] Regularly rotate secrets and keys
- [ ] Monitor and review audit logs

#### **Environment-Specific Security**

```bash
# Development - Relaxed security
VAULT_SERVER_ENVIRONMENT=development
VAULT_DATABASE_SSLMODE=disable
VAULT_SECURITY_KDF_ITERATIONS=10000

# Staging - Enhanced security
VAULT_SERVER_ENVIRONMENT=staging
VAULT_DATABASE_SSLMODE=require
VAULT_SECURITY_KDF_ITERATIONS=50000

# Production - Maximum security
VAULT_SERVER_ENVIRONMENT=production
VAULT_DATABASE_SSLMODE=verify-full
VAULT_SECURITY_KDF_ITERATIONS=100000
```

---

## 🗄️ Database Configuration

### 🐘 **PostgreSQL Setup**

#### **Install PostgreSQL**

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS
brew install postgresql
brew services start postgresql

# CentOS/RHEL
sudo yum install postgresql-server postgresql-contrib
sudo postgresql-setup initdb
sudo systemctl start postgresql
```

#### **Create Database and User**

```bash
# Connect to PostgreSQL
sudo -u postgres psql

# Create database
CREATE DATABASE vault;

# Create user
CREATE USER vault_user WITH PASSWORD 'secure_password';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE vault TO vault_user;

# Exit
\q
```

#### **Database Connection String**

```bash
# Format: postgres://user:password@host:port/dbname?sslmode=mode
postgres://vault_user:secure_password@localhost:5432/vault?sslmode=require
```

### 🔧 **Database Optimization**

#### **PostgreSQL Configuration** (postgresql.conf)

```ini
# Memory Configuration
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 4MB
maintenance_work_mem = 64MB

# Connection Configuration
max_connections = 100
shared_preload_libraries = 'pg_stat_statements'

# Logging Configuration
log_statement = 'all'
log_min_duration_statement = 1000
log_checkpoints = on
log_connections = on
log_disconnections = on
```

#### **Connection Pooling**

```go
// Application-side connection pool (configured in server)
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

---

## 🚀 Production Setup

### 🏗️ **Production Configuration**

#### **Environment Variables**

```bash
# Production Server Configuration
VAULT_SERVER_HOST=0.0.0.0
VAULT_SERVER_PORT=443
VAULT_SERVER_ENVIRONMENT=production
VAULT_SERVER_READ_TIMEOUT=60
VAULT_SERVER_WRITE_TIMEOUT=60

# Production Database Configuration
VAULT_DATABASE_HOST=db.production.com
VAULT_DATABASE_PORT=5432
VAULT_DATABASE_USER=vault_prod
VAULT_DATABASE_PASSWORD=super_secure_db_password
VAULT_DATABASE_DBNAME=vault_production
VAULT_DATABASE_SSLMODE=verify-full

# Production Security Configuration
VAULT_SECURITY_ENCRYPTION_KEY=32-character-production-encryption-key
VAULT_SECURITY_KDF_ITERATIONS=100000
VAULT_SECURITY_SALT_LENGTH=64

# Production JWT Configuration
VAULT_JWT_SECRET=super-long-secure-jwt-secret-for-production-use-only
VAULT_JWT_EXPIRATION=7200

# Production Audit Configuration
VAULT_AUDIT_ENABLED=true
VAULT_AUDIT_LOG_LEVEL=info
VAULT_AUDIT_LOG_FORMAT=json
```

#### **Nginx Configuration** (SSL Termination)

```nginx
server {
    listen 443 ssl http2;
    server_name vault.example.com;

    # SSL Configuration
    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;

    # Security Headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";

    # Proxy to Aether Vault
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 🔧 **Systemd Service**

```ini
# /etc/systemd/system/aether-vault.service
[Unit]
Description=Aether Vault Server
After=network.target postgresql.service

[Service]
Type=simple
User=vault
Group=vault
WorkingDirectory=/opt/aether-vault
ExecStart=/opt/aether-vault/aether-vault
Restart=always
RestartSec=10
Environment=VAULT_SERVER_ENVIRONMENT=production
Environment=VAULT_DATABASE_HOST=localhost
Environment=VAULT_DATABASE_USER=vault
Environment=VAULT_DATABASE_PASSWORD=/run/secrets/vault_db_password

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/aether-vault

[Install]
WantedBy=multi-user.target
```

---

## 🔧 Advanced Settings

### 📊 **Performance Configuration**

#### **Go Runtime Configuration**

```go
// Advanced server configuration in main.go
runtime.GOMAXPROCS(runtime.NumCPU())

// HTTP Server configuration
server := &http.Server{
    Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
    Handler:      router.GetEngine(),
    ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
    WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
    IdleTimeout:  120 * time.Second,
    MaxHeaderBytes: 1 << 20, // 1MB
}
```

#### **Gin Configuration**

```go
// Production Gin configuration
if cfg.Server.Environment == "production" {
    gin.SetMode(gin.ReleaseMode)
}

// Enable JSON pretty print in development
if cfg.Server.Environment == "development" {
    gin.DisableConsoleColor()
}
```

### 🗄️ **Database Advanced Settings**

#### **GORM Configuration**

```go
// Advanced GORM configuration
gormConfig := &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
    PrepareStmt:              true,
    DisableForeignKeyConstraintWhenMigrating: true,
}

// Database connection with advanced settings
db, err := gorm.Open(postgres.Open(dsn), gormConfig)
if err != nil {
    return nil, fmt.Errorf("failed to connect to database: %w", err)
}

// Advanced connection pool settings
sqlDB, err := db.DB()
if err != nil {
    return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
}

sqlDB.SetMaxIdleConns(25)
sqlDB.SetMaxOpenConns(200)
sqlDB.SetConnMaxLifetime(time.Hour)
sqlDB.SetConnMaxIdleTime(30 * time.Minute)
```

### 🔐 **Security Advanced Configuration**

#### **Rate Limiting Configuration**

```go
// Advanced rate limiting per endpoint
rateLimitMiddleware := middleware.NewRateLimitMiddleware(
    100,  // requests
    60,   // per minute
    10,   // burst size
    time.Minute,
)

// Rate limiting by user
userRateLimitMiddleware := middleware.NewUserRateLimitMiddleware(
    1000, // requests per user
    time.Hour,
)
```

#### **CORS Configuration**

```go
// Advanced CORS configuration
corsConfig := cors.Config{
    AllowOrigins:     []string{"https://app.example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
    ExposeHeaders:    []string{"X-RateLimit-Limit", "X-RateLimit-Remaining"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
```

---

## 🛠️ Configuration Validation

### ✅ **Startup Validation**

The server validates configuration on startup:

```go
// Required in production
if cfg.Server.Environment == "production" {
    if cfg.Database.Host == "" {
        panic("Database host is required in production")
    }
    if cfg.JWT.Secret == "" {
        panic("JWT secret is required in production")
    }
    if cfg.Security.EncryptionKey == "" {
        panic("Encryption key is required in production")
    }
}

// Validate values
if len(cfg.Security.EncryptionKey) != 32 {
    panic("Encryption key must be exactly 32 characters")
}

if len(cfg.JWT.Secret) < 32 {
    panic("JWT secret must be at least 32 characters")
}
```

### 🔍 **Configuration Health Check**

```bash
# Check configuration health
make health

# Expected output
✅ Server: OK (localhost:8080)
✅ Database: OK (connected)
✅ JWT Secret: OK (configured)
✅ Encryption Key: OK (configured)
✅ Audit Logging: OK (enabled)
```

---

## 📝 Configuration Examples

### 🏢 **Enterprise Configuration**

```bash
# Enterprise-grade configuration
VAULT_SERVER_HOST=0.0.0.0
VAULT_SERVER_PORT=443
VAULT_SERVER_ENVIRONMENT=production

# High-availability database cluster
VAULT_DATABASE_HOST=db-cluster.internal
VAULT_DATABASE_PORT=5432
VAULT_DATABASE_USER=vault_enterprise
VAULT_DATABASE_PASSWORD=vault_enterprise_password
VAULT_DATABASE_DBNAME=vault_enterprise
VAULT_DATABASE_SSLMODE=verify-full

# Enterprise security
VAULT_SECURITY_ENCRYPTION_KEY=enterprise-encryption-key-32
VAULT_SECURITY_KDF_ITERATIONS=200000
VAULT_SECURITY_SALT_LENGTH=64

# Enterprise JWT
VAULT_JWT_SECRET=enterprise-jwt-secret-for-production-use-only
VAULT_JWT_EXPIRATION=3600

# Enterprise audit
VAULT_AUDIT_ENABLED=true
VAULT_AUDIT_LOG_LEVEL=info
VAULT_AUDIT_LOG_FORMAT=json
```

### 🧪 **Development Configuration**

```bash
# Development configuration
VAULT_SERVER_HOST=127.0.0.1
VAULT_SERVER_PORT=8080
VAULT_SERVER_ENVIRONMENT=development

# Local database
VAULT_DATABASE_HOST=localhost
VAULT_DATABASE_PORT=5432
VAULT_DATABASE_USER=vault_dev
VAULT_DATABASE_PASSWORD=dev_password
VAULT_DATABASE_DBNAME=vault_dev
VAULT_DATABASE_SSLMODE=disable

# Development security (lower security for ease of use)
VAULT_SECURITY_ENCRYPTION_KEY=dev-encryption-key-32-chars
VAULT_SECURITY_KDF_ITERATIONS=10000
VAULT_SECURITY_SALT_LENGTH=16

# Development JWT
VAULT_JWT_SECRET=dev-jwt-secret-not-for-production
VAULT_JWT_EXPIRATION=86400

# Development audit
VAULT_AUDIT_ENABLED=true
VAULT_AUDIT_LOG_LEVEL=debug
VAULT_AUDIT_LOG_FORMAT=text
```

---

## 🔧 Troubleshooting

### ❌ **Common Configuration Issues**

#### **Database Connection Failed**

```bash
# Check database connection
psql -h localhost -U vault -d vault

# Common solutions
# 1. Check database is running
sudo systemctl status postgresql

# 2. Check connection parameters
#    - Host is correct
#    - Port is correct
#    - User exists
#    - Password is correct
#    - Database exists

# 3. Check SSL mode
#    - disable for local development
#    - require/verify-full for production
```

#### **JWT Secret Issues**

```bash
# Generate new JWT secret
openssl rand -base64 32

# Update environment
export VAULT_JWT_SECRET="your-new-secret-here"

# Restart server
make go-server
```

#### **Encryption Key Issues**

```bash
# Check key length
echo -n "your-key" | wc -c
# Should be exactly 32 characters

# Generate new key
openssl rand -hex 16

# Update environment
export VAULT_SECURITY_ENCRYPTION_KEY="your-new-32-character-key"
```

### 🔍 **Debug Configuration**

```bash
# Show configuration
make config

# Debug configuration loading
go run main.go --debug-config

# Check environment variables
env | grep VAULT_

# Validate configuration file
go run main.go --validate-config
```

---

## 🔗 Related Documentation

- [📖 Server Documentation](./README.md)
- [🏗️ Architecture Guide](./architecture.md)
- [📚 API Documentation](./api.md)
- [🚀 Deployment Guide](./deployment.md)
- [🔒 Security Guide](./security.md)

---

<div align="center">

### ⚙️ **Configure Your Aether Vault Server for Maximum Performance and Security!**

[📖 Full Documentation](../../README.md) • [🐛 Report Issues](https://github.com/skygenesisenterprise/aether-vault/issues) • [💡 Configuration Help](https://github.com/skygenesisenterprise/aether-vault/discussions)

---

**🔧 Production-Ready Configuration for Enterprise Authentication & Vault Management**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

</div>
