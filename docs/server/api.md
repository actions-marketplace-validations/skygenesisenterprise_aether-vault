<div align="center">

# 📚 Aether Vault API Documentation

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![API Version](https://img.shields.io/badge/API-v1.0-green?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault) [![Go](https://img.shields.io/badge/Go-1.25+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![Gin](https://img.shields.io/badge/Gin-REST-lightgrey?style=for-the-badge&logo=go)](https://gin-gonic.com/)

**🔥 Complete REST API Reference for Enterprise Authentication & Vault Management**

Comprehensive API documentation for the Aether Vault Server - an enterprise-grade authentication and secret management system with advanced security features.

[🔗 Base URL](#-base-url) • [🔐 Authentication](#-authentication) • [👤 Users](#-users) • [🗄️ Secrets](#️-secrets) • [🔐 TOTP](#-totp-2fa) • [🆔 Identity](#-identity) • [📊 Audit](#-audit) • [⚙️ System](#️-system)

</div>

---

## 🌐 Base URL

```
Development: http://localhost:8080
Production:  https://your-domain.com
API Version: /api/v1
```

All API endpoints are prefixed with `/api/v1/`.

---

## 🔐 Authentication

### Overview

The Aether Vault API uses JWT (JSON Web Tokens) for authentication. Some endpoints require authentication, while others are public for system access.

### Authentication Headers

```http
Authorization: Bearer <jwt_token>
Content-Type: application/json
X-Request-ID: <correlation_id>
```

### Public Endpoints

- `POST /api/v1/auth/login` - User authentication
- `GET /api/v1/system/health` - Health check
- `GET /api/v1/system/version` - Version information

### Protected Endpoints

All other endpoints require a valid JWT token in the `Authorization` header.

---

## 🔑 Authentication Endpoints

### POST /api/v1/auth/login

Authenticates a user and returns a JWT token.

**Request:**

```json
{
  "email": "user@example.com",
  "password": "secure_password",
  "totp_code": "123456"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "user": {
      "id": "uuid-here",
      "email": "user@example.com",
      "role": "user"
    }
  }
}
```

**Status Codes:**

- `200 OK` - Authentication successful
- `401 Unauthorized` - Invalid credentials
- `403 Forbidden` - Account locked or disabled
- `429 Too Many Requests` - Rate limit exceeded

### POST /api/v1/auth/logout

Logs out a user and invalidates the JWT token.

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "success": true,
  "message": "Successfully logged out"
}
```

### GET /api/v1/auth/session

Retrieves current session information.

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid-here",
      "email": "user@example.com",
      "role": "user",
      "created_at": "2025-01-09T10:00:00Z",
      "last_login": "2025-01-09T10:30:00Z"
    },
    "session": {
      "expires_at": "2025-01-09T11:30:00Z",
      "created_at": "2025-01-09T10:30:00Z"
    }
  }
}
```

---

## 👤 User Management Endpoints

All user endpoints require authentication.

### GET /api/v1/users

Retrieves a list of users (requires admin privileges).

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**

- `page` (int, optional) - Page number (default: 1)
- `limit` (int, optional) - Items per page (default: 20)
- `search` (string, optional) - Search by email or name

**Response:**

```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "uuid-here",
        "email": "user@example.com",
        "role": "user",
        "created_at": "2025-01-09T10:00:00Z",
        "updated_at": "2025-01-09T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 1,
      "pages": 1
    }
  }
}
```

### GET /api/v1/users/:id

Retrieves a specific user by ID.

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid-here",
    "email": "user@example.com",
    "role": "user",
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:00:00Z"
  }
}
```

### POST /api/v1/users

Creates a new user (requires admin privileges).

**Headers:** `Authorization: Bearer <token>`

**Request:**

```json
{
  "email": "newuser@example.com",
  "password": "secure_password",
  "role": "user"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid-here",
    "email": "newuser@example.com",
    "role": "user",
    "created_at": "2025-01-09T10:00:00Z"
  }
}
```

### PUT /api/v1/users/:id

Updates a user (requires admin privileges or own user).

**Headers:** `Authorization: Bearer <token>`

**Request:**

```json
{
  "email": "updated@example.com",
  "role": "admin"
}
```

### DELETE /api/v1/users/:id

Deletes a user (requires admin privileges).

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

---

## 🗄️ Secret Management Endpoints

All secret endpoints require authentication.

### GET /api/v1/secrets

Retrieves a list of secrets accessible to the user.

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**

- `page` (int, optional) - Page number
- `limit` (int, optional) - Items per page
- `search` (string, optional) - Search by name or description

**Response:**

```json
{
  "success": true,
  "data": {
    "secrets": [
      {
        "id": "uuid-here",
        "name": "database-password",
        "description": "Production database password",
        "created_at": "2025-01-09T10:00:00Z",
        "updated_at": "2025-01-09T10:00:00Z"
      }
    ]
  }
}
```

### POST /api/v1/secrets

Creates a new secret.

**Headers:** `Authorization: Bearer <token>`

**Request:**

```json
{
  "name": "api-key",
  "value": "secret-api-key-value",
  "description": "External API key",
  "tags": ["production", "external"]
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid-here",
    "name": "api-key",
    "description": "External API key",
    "tags": ["production", "external"],
    "created_at": "2025-01-09T10:00:00Z"
  }
}
```

### GET /api/v1/secrets/:id

Retrieves a specific secret (decrypts the value).

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid-here",
    "name": "api-key",
    "value": "secret-api-key-value",
    "description": "External API key",
    "tags": ["production", "external"],
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:00:00Z"
  }
}
```

### PUT /api/v1/secrets/:id

Updates an existing secret.

**Headers:** `Authorization: Bearer <token>`

**Request:**

```json
{
  "name": "updated-api-key",
  "value": "new-secret-value",
  "description": "Updated description"
}
```

### DELETE /api/v1/secrets/:id

Deletes a secret.

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "success": true,
  "message": "Secret deleted successfully"
}
```

---

## 🔐 TOTP 2FA Endpoints

All TOTP endpoints require authentication.

### GET /api/v1/totp

Retrieves TOTP configurations for the user.

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "success": true,
  "data": {
    "totp_configs": [
      {
        "id": "uuid-here",
        "name": "Primary Device",
        "issuer": "Aether Vault",
        "created_at": "2025-01-09T10:00:00Z"
      }
    ]
  }
}
```

### POST /api/v1/totp

Creates a new TOTP configuration.

**Headers:** `Authorization: Bearer <token>`

**Request:**

```json
{
  "name": "Mobile Device",
  "issuer": "Aether Vault"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid-here",
    "name": "Mobile Device",
    "issuer": "Aether Vault",
    "secret": "JBSWY3DPEHPK3PXP",
    "qr_code": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
    "backup_codes": ["123456", "789012", "345678"],
    "created_at": "2025-01-09T10:00:00Z"
  }
}
```

### POST /api/v1/totp/:id/generate

Generates a new TOTP code for verification.

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "success": true,
  "data": {
    "code": "123456",
    "expires_in": 30
  }
}
```

---

## 🆔 Identity Management Endpoints

All identity endpoints require authentication.

### GET /api/v1/identity/me

Retrieves the current user's identity information.

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid-here",
    "email": "user@example.com",
    "role": "user",
    "permissions": ["read:secrets", "write:secrets"],
    "created_at": "2025-01-09T10:00:00Z",
    "last_login": "2025-01-09T10:30:00Z"
  }
}
```

### GET /api/v1/identity/policies

Retrieves access policies for the current user.

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "success": true,
  "data": {
    "policies": [
      {
        "id": "uuid-here",
        "name": "Default User Policy",
        "permissions": ["read:own_secrets", "write:own_secrets"],
        "resources": ["secrets"],
        "created_at": "2025-01-09T10:00:00Z"
      }
    ]
  }
}
```

---

## 📊 Audit Logging Endpoints

All audit endpoints require authentication.

### GET /api/v1/audit/logs

Retrieves audit logs (requires admin privileges).

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**

- `page` (int, optional) - Page number
- `limit` (int, optional) - Items per page
- `user_id` (string, optional) - Filter by user ID
- `action` (string, optional) - Filter by action type
- `start_date` (string, optional) - Start date filter (ISO 8601)
- `end_date` (string, optional) - End date filter (ISO 8601)

**Response:**

```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": "uuid-here",
        "user_id": "user-uuid",
        "action": "secret:created",
        "resource": "secrets/secret-uuid",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0...",
        "created_at": "2025-01-09T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 1,
      "pages": 1
    }
  }
}
```

---

## ⚙️ System Endpoints

System endpoints are public and do not require authentication.

### GET /api/v1/system/health

Returns the health status of the server and its dependencies.

**Response:**

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2025-01-09T10:00:00Z",
    "version": "1.0.0",
    "uptime": "2h30m45s",
    "checks": {
      "database": "healthy",
      "redis": "healthy",
      "storage": "healthy"
    }
  }
}
```

### GET /api/v1/system/version

Returns version information about the server.

**Response:**

```json
{
  "success": true,
  "data": {
    "version": "1.0.0",
    "build": "2025-01-09T10:00:00Z",
    "git_commit": "abc123def456",
    "go_version": "1.25.5",
    "gin_version": "1.9.1"
  }
}
```

---

## 📝 Response Format

All API responses follow a consistent format:

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data
  },
  "message": "Optional success message"
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": "Additional error details"
  },
  "request_id": "uuid-here"
}
```

---

## 🚨 HTTP Status Codes

| Status                      | Description             | Usage                        |
| --------------------------- | ----------------------- | ---------------------------- |
| `200 OK`                    | Success                 | Successful GET, PUT, DELETE  |
| `201 Created`               | Created                 | Successful POST              |
| `400 Bad Request`           | Invalid Request         | Invalid JSON, missing fields |
| `401 Unauthorized`          | Authentication Required | Invalid or missing token     |
| `403 Forbidden`             | Access Denied           | Insufficient permissions     |
| `404 Not Found`             | Resource Not Found      | Invalid resource ID          |
| `409 Conflict`              | Conflict                | Resource already exists      |
| `422 Unprocessable Entity`  | Validation Error        | Invalid input data           |
| `429 Too Many Requests`     | Rate Limited            | Too many requests            |
| `500 Internal Server Error` | Server Error            | Unexpected server error      |

---

## 🔄 Rate Limiting

The API implements rate limiting to prevent abuse:

- **Default Limit**: 100 requests per minute per IP
- **Authenticated Users**: Higher limits based on role
- **Admin Users**: No rate limiting

Rate limit headers are included in responses:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1641763200
```

---

## 📝 Request IDs

All requests receive a unique correlation ID for tracking:

```http
X-Request-ID: uuid-here
```

Include this ID in support requests for faster troubleshooting.

---

## 🔒 Security Considerations

### HTTPS in Production

Always use HTTPS in production environments. HTTP is only supported for development.

### Token Security

- JWT tokens expire after 1 hour by default
- Use refresh tokens for extended sessions
- Invalidate tokens on logout
- Store tokens securely on the client side

### Input Validation

All inputs are validated and sanitized:

- SQL injection protection
- XSS protection
- CSRF protection
- Input length limits

### Audit Logging

All API calls are logged for security:

- Request/response data
- User identification
- IP addresses
- Timestamps

---

## 🧪 Testing the API

### Using curl

```bash
# Health check
curl http://localhost:8080/api/v1/system/health

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Get secrets (with token)
curl -X GET http://localhost:8080/api/v1/secrets \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Using Postman

Import the Postman collection from the repository for easy API testing.

---

## 📚 SDK and Libraries

### Official Libraries

- **Go SDK**: `github.com/skygenesisenterprise/aether-vault-go`
- **JavaScript SDK**: `@skygenesisenterprise/aether-vault-js`
- **Python SDK**: `aether-vault-python`

### Third-Party Libraries

Community-maintained libraries are available on GitHub.

---

## 🔗 Related Documentation

- [📖 Server Documentation](./README.md)
- [🏗️ Architecture Guide](./architecture.md)
- [⚙️ Configuration Guide](./configuration.md)
- [🚀 Deployment Guide](./deployment.md)
- [🔒 Security Guide](./security.md)

---

<div align="center">

### 🚀 **Start Building with the Aether Vault API Today!**

[📖 Full Documentation](../../README.md) • [🐛 Report Issues](https://github.com/skygenesisenterprise/aether-vault/issues) • [💡 Get Help](https://github.com/skygenesisenterprise/aether-vault/discussions)

---

**🔧 Enterprise-Grade Authentication & Vault Management API**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

</div>
