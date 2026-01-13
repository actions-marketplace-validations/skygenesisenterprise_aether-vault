<div align="center">

# 🚀 Aether Vault CLI - Capability Commands Reference

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![Go](https://img.shields.io/badge/Go-1.25+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![Cobra](https://img.shields.io/badge/Cobra-1.8+-lightgrey?style=for-the-badge&logo=go)](https://github.com/spf13/cobra) [![Viper](https://img.shields.io/badge/Viper-1.16+-green?style=for-the-badge&logo=go)](https://github.com/spf13/viper) [![DevOps](https://img.shields.io/badge/DevOps-Ready-orange?style=for-the-badge&logo=devops)](https://www.devops.com/)

**🔐 Complete Capability Command Reference - Master CBAC System**

Comprehensive reference for all `vault capability` commands, including request, validation, management, and revocation of cryptographic capabilities. Capabilities are the core of Aether Vault's Capability-Based Access Control (CBAC) system.

[📋 Command Overview](#-command-overview) • [🔐 Capability Request](#-capability-request) • [✅ Capability Validate](#-capability-validate) • [📋 Capability List](#-capability-list) • [🗑️ Capability Revoke](#-capability-revoke) • [📊 Capability Status](#-capability-status) • [🎯 Use Cases](#-use-cases) • [🔍 Error Handling](#-error-handling) • [🏆 Best Practices](#-best-practices)

[![GitHub stars](https://img.shields.io/github/stars/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/stargazers) [![GitHub forks](https://img.shields.io/github/forks/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/network)

</div>

---

## 📋 Command Overview

### 🎯 **Command Structure**

```bash
vault capability [subcommand] [flags]
```

### 📚 **Available Subcommands**

| Subcommand | Description                  | Purpose                                  |
| ---------- | ---------------------------- | ---------------------------------------- |
| `request`  | Request new capability       | Create time-limited access tokens        |
| `validate` | Validate existing capability | Check if capability is still valid       |
| `list`     | List capabilities            | Display and filter existing capabilities |
| `revoke`   | Revoke capability            | Invalidate capability immediately        |
| `status`   | Show system status           | Display capability engine status         |

### 🔄 **Capability Lifecycle**

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Request       │    │   Validate       │    │   Revoke        │
│   (Create)      │◄──►│   (Check)        │◄──►│   (Invalidate)  │
│  Resource Access │    │  Status Check    │    │  Immediate Stop │
│  Time-Limited    │    │  Constraint Eval │    │  Security Action│
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
   Policy Evaluation        Usage Tracking        Audit Logging
   Constraint Checking       Cache Management      Security Events
   Cryptographic Signing     Expiration Handling   Revocation Trail
```

---

## 🔐 Capability Request

### 🎯 **Syntax**

```bash
vault capability request [flags]
```

### 📋 **Required Flags**

| Flag         | Type        | Description                                |
| ------------ | ----------- | ------------------------------------------ |
| `--resource` | string      | Resource path (e.g., `secret:/db/primary`) |
| `--action`   | stringSlice | Action(s) to grant (e.g., `read`, `write`) |

### ⚙️ **Optional Flags**

| Flag            | Type   | Default       | Description                    |
| --------------- | ------ | ------------- | ------------------------------ |
| `--ttl`         | int64  | 300           | Time-to-live in seconds        |
| `--max-uses`    | int    | 100           | Maximum number of uses         |
| `--identity`    | string | auto-detected | Requesting identity            |
| `--purpose`     | string | -             | Purpose of the request         |
| `--constraints` | string | -             | Constraints in JSON format     |
| `--context`     | string | -             | Request context in JSON format |

### 💡 **Examples**

#### **🎯 Basic Capability Request**

```bash
vault capability request \
  --resource "secret:/db/primary" \
  --action read
```

#### **⏰ Capability with Custom TTL**

```bash
vault capability request \
  --resource "secret:/api/production" \
  --action read,write \
  --ttl 600 \
  --purpose "API access for production deployment"
```

#### **🔒 Capability with Constraints**

```bash
vault capability request \
  --resource "secret:/sensitive/data" \
  --action read \
  --constraints '{"ipAddresses": ["10.0.0.100"], "timeWindow": {"hours": [9,10,11,12,13,14,15,16,17]}}'
```

#### **🐳 Capability with Runtime Context**

```bash
vault capability request \
  --resource "secret:/docker/registry" \
  --action read \
  --context '{"runtime": {"type": "docker", "id": "container123"}, "sourceIP": "10.0.0.100"}'
```

#### **🏢 Enterprise Capability Request**

```bash
vault capability request \
  --resource "secret:/enterprise/database" \
  --action read,write \
  --ttl 1800 \
  --max-uses 50 \
  --identity "production-service" \
  --purpose "Critical business operation" \
  --constraints '{
    "ipAddresses": ["10.0.0.0/8"],
    "timeWindow": {"hours": [9,10,11,12,13,14,15,16,17], "days": ["mon","tue","wed","thu","fri"]},
    "environment": {"type": "production", "region": "us-east-1"}
  }' \
  --context '{
    "runtime": {"type": "kubernetes", "pod": "prod-service-xyz", "namespace": "production"},
    "sourceIP": "10.0.1.100",
    "user": {"id": "service-account", "role": "production-operator"}
  }'
```

### 📤 **Response Format**

#### **📋 Table Format**

```
Capability Request Result:
  Status: granted
  Request ID: req_1234567890_abcdef
  Processing Time: 45ms

Capability Details:
  ID: cap_1234567890_ghijkl
  Type: read
  Resource: secret:/db/primary
  Actions: read
  Identity: app123
  Issuer: aether-vault-agent
  TTL: 300 seconds
  Max Uses: 100
  Issued At: 2025-01-13T10:00:00Z
  Expires At: 2025-01-13T10:05:00Z

Policy Evaluation:
  Decision: allow
  Applied Policies: ["database-access", "app-policy"]
  Reasoning: Request matches database access policy for app identity

Constraints:
  IP Addresses: ["10.0.0.100"]
  Time Window: {"hours": [9,10,11,12,13,14,15,16,17]}
  Environment: {"type": "production"}
```

#### **📄 JSON Format**

```json
{
  "capability": {
    "id": "cap_1234567890_ghijkl",
    "type": "read",
    "resource": "secret:/db/primary",
    "actions": ["read"],
    "identity": "app123",
    "issuer": "aether-vault-agent",
    "issued_at": "2025-01-13T10:00:00Z",
    "expires_at": "2025-01-13T10:05:00Z",
    "ttl": 300,
    "max_uses": 100,
    "used_count": 0,
    "signature": "base64-encoded-signature",
    "constraints": {
      "ipAddresses": ["10.0.0.100"],
      "timeWindow": { "hours": [9, 10, 11, 12, 13, 14, 15, 16, 17] },
      "environment": { "type": "production" }
    }
  },
  "status": "granted",
  "message": "Capability granted successfully",
  "request_id": "req_1234567890_abcdef",
  "processing_time": "45ms",
  "policy_result": {
    "decision": "allow",
    "applied_policies": ["database-access", "app-policy"],
    "reasoning": "Request matches database access policy for app identity"
  }
}
```

---

## ✅ Capability Validate

### 🎯 **Syntax**

```bash
vault capability validate [capability-id] [flags]
```

### 📋 **Arguments**

| Argument        | Type   | Description                      |
| --------------- | ------ | -------------------------------- |
| `capability-id` | string | ID of the capability to validate |

### ⚙️ **Optional Flags**

| Flag        | Type   | Default | Description                       |
| ----------- | ------ | ------- | --------------------------------- |
| `--context` | string | -       | Validation context in JSON format |
| `--verbose` | bool   | false   | Show detailed validation info     |

### 💡 **Examples**

#### **✅ Basic Validation**

```bash
vault capability validate cap_1234567890_ghijkl
```

#### **🔍 Validation with Context**

```bash
vault capability validate cap_1234567890_ghijkl \
  --context '{"sourceIP": "10.0.0.100", "runtime": {"type": "docker", "id": "container123"}}'
```

#### **📊 Verbose Validation**

```bash
vault capability validate cap_1234567890_ghijkl --verbose
```

### 📤 **Response Format**

#### **✅ Successful Validation**

```
Capability Validation Result:
  Valid: true
  Validation Time: 12ms

Context:
  cache_hit: false
  constraints_satisfied: true
  policy_evaluations: 2
  signature_valid: true

Capability Info:
  ID: cap_1234567890_ghijkl
  Resource: secret:/db/primary
  Actions: read
  Uses Remaining: 95
  Expires In: 4m48s
```

#### **❌ Failed Validation**

```
Capability Validation Result:
  Valid: false
  Validation Time: 8ms

Errors:
  EXPIRED: Capability expired at 2025-01-13T10:05:00Z
    Field: expires_at

Warnings:
  USAGE_HIGH: Capability has used 80% of max uses
    Field: used_count
```

#### **📄 JSON Format**

```json
{
  "valid": true,
  "validation_time": "12ms",
  "context": {
    "cache_hit": false,
    "constraints_satisfied": true,
    "policy_evaluations": 2,
    "signature_valid": true
  },
  "capability_info": {
    "id": "cap_1234567890_ghijkl",
    "resource": "secret:/db/primary",
    "actions": ["read"],
    "uses_remaining": 95,
    "expires_in": "4m48s"
  }
}
```

---

## 📋 Capability List

### 🎯 **Syntax**

```bash
vault capability list [flags]
```

### ⚙️ **Optional Flags**

| Flag         | Type   | Default | Description                                 |
| ------------ | ------ | ------- | ------------------------------------------- |
| `--identity` | string | -       | Filter by identity                          |
| `--type`     | string | -       | Filter by capability type                   |
| `--status`   | string | -       | Filter by status (active, expired, revoked) |
| `--limit`    | int    | 50      | Limit number of results                     |
| `--offset`   | int    | 0       | Offset for pagination                       |
| `--format`   | string | table   | Output format (table, json, yaml)           |

### 💡 **Examples**

#### **📋 List All Capabilities**

```bash
vault capability list
```

#### **👤 Filter by Identity**

```bash
vault capability list --identity "app123"
```

#### **🔍 Filter by Type and Status**

```bash
vault capability list --type "read" --status "active"
```

#### **📄 JSON Output**

```bash
vault capability list --format json
```

#### **📊 Paginated Results**

```bash
vault capability list --limit 10 --offset 20
```

#### **🔍 Advanced Filtering**

```bash
# List all active capabilities for production services
vault capability list \
  --identity "prod-*" \
  --status "active" \
  --limit 100

# List expired capabilities from last hour
vault capability list \
  --status "expired" \
  --format json | jq '.capabilities[] | select((.expires_at | fromdateiso8601) > (now - 3600))'
```

### 📤 **Response Format**

#### **📋 Table Format**

```
Found 25 capabilities:

ID                   Type            Resource                        Identity         Expires
--------------------------------------------------------------------------------------------------------------
cap_1234567890_abc   read            secret:/db/primary             app123           2025-01-13 10:05:00
cap_1234567890_def   write           secret:/api/config             deploy-service   2025-01-13 11:00:00
cap_1234567890_ghi   admin           secret:/system/*               admin-user       2025-01-13 12:00:00
cap_1234567890_jkl   read            secret:/cache/redis           cache-service    2025-01-13 10:30:00
...
```

#### **📄 JSON Format**

```json
{
  "capabilities": [
    {
      "id": "cap_1234567890_abc",
      "type": "read",
      "resource": "secret:/db/primary",
      "identity": "app123",
      "expires_at": "2025-01-13T10:05:00Z",
      "status": "active",
      "used_count": 5,
      "max_uses": 100
    },
    {
      "id": "cap_1234567890_def",
      "type": "write",
      "resource": "secret:/api/config",
      "identity": "deploy-service",
      "expires_at": "2025-01-13T11:00:00Z",
      "status": "active",
      "used_count": 1,
      "max_uses": 50
    }
  ],
  "count": 25,
  "limit": 50,
  "offset": 0,
  "total": 127
}
```

---

## 🗑️ Capability Revoke

### 🎯 **Syntax**

```bash
vault capability revoke [capability-id] [flags]
```

### 📋 **Arguments**

| Argument        | Type   | Description                    |
| --------------- | ------ | ------------------------------ |
| `capability-id` | string | ID of the capability to revoke |

### ⚙️ **Optional Flags**

| Flag       | Type   | Default             | Description           |
| ---------- | ------ | ------------------- | --------------------- |
| `--reason` | string | "Manual revocation" | Reason for revocation |
| `--force`  | bool   | false               | Force revocation      |

### 💡 **Examples**

#### **🗑️ Basic Revocation**

```bash
vault capability revoke cap_1234567890_ghijkl
```

#### **📝 Revocation with Reason**

```bash
vault capability revoke cap_1234567890_ghijkl \
  --reason "Security policy violation - suspicious activity detected"
```

#### **💪 Force Revocation**

```bash
vault capability revoke cap_1234567890_ghijkl --force
```

#### **🚨 Emergency Revocation**

```bash
# Revoke all capabilities for compromised identity
for cap_id in $(vault capability list --identity "compromised-app" --format json | jq -r '.capabilities[].id'); do
  vault capability revoke "$cap_id" --reason "SECURITY INCIDENT - Identity compromised"
done
```

### 📤 **Response Format**

#### **✅ Successful Revocation**

```
Capability cap_1234567890_ghijkl revoked successfully
Reason: Security policy violation - suspicious activity detected
Revoked At: 2025-01-13T10:15:00Z
Audit ID: audit_1234567890_xyz
```

#### **📄 JSON Format**

```json
{
  "capability_id": "cap_1234567890_ghijkl",
  "revoked": true,
  "reason": "Security policy violation - suspicious activity detected",
  "revoked_at": "2025-01-13T10:15:00Z",
  "audit_id": "audit_1234567890_xyz"
}
```

---

## 📊 Capability Status

### 🎯 **Syntax**

```bash
vault capability status [flags]
```

### ⚙️ **Optional Flags**

| Flag        | Type   | Default | Description                       |
| ----------- | ------ | ------- | --------------------------------- |
| `--verbose` | bool   | false   | Show detailed status information  |
| `--format`  | string | table   | Output format (table, json, yaml) |

### 💡 **Examples**

#### **📊 Basic Status**

```bash
vault capability status
```

#### **📋 Verbose Status**

```bash
vault capability status --verbose
```

#### **📄 JSON Status**

```bash
vault capability status --format json
```

### 📤 **Response Format**

#### **📋 Table Format**

```
Aether Vault Agent Status:
  Version: 1.0.0
  Uptime: 2h45m30s
  Connections: 3

Capability Engine:
  Status: Healthy
  Total Capabilities: 127
  Active Capabilities: 45
  Expired Capabilities: 80
  Revoked Capabilities: 2
  Cache Size: 45/10000
  Last Cleanup: 5m ago

Policy Engine:
  Status: Healthy
  Loaded Policies: 5
  Cache Hits: 892
  Cache Misses: 45
  Cache Hit Rate: 95.2%

Performance Metrics:
  Requests/Second: 45.2
  Average Response Time: 12ms
  Validation Cache Hit Rate: 87.3%
```

#### **📄 JSON Format**

```json
{
  "version": "1.0.0",
  "uptime": "2h45m30s",
  "connections": 3,
  "capability_engine": {
    "status": "healthy",
    "total_capabilities": 127,
    "active_capabilities": 45,
    "expired_capabilities": 80,
    "revoked_capabilities": 2,
    "cache_size": "45/10000",
    "last_cleanup": "5m ago"
  },
  "policy_engine": {
    "status": "healthy",
    "loaded_policies": 5,
    "cache_hits": 892,
    "cache_misses": 45,
    "cache_hit_rate": "95.2%"
  },
  "performance_metrics": {
    "requests_per_second": 45.2,
    "average_response_time": "12ms",
    "validation_cache_hit_rate": "87.3%"
  }
}
```

---

## 🎯 Use Cases

### 🗄️ **1. Database Access**

```bash
# Request read access to database
vault capability request \
  --resource "secret:/db/primary" \
  --action read \
  --ttl 300 \
  --identity "web-app" \
  --purpose "Database connection for web application"

# Validate before use
vault capability validate cap_1234567890_abc

# Use capability in application
export DB_CAPABILITY_ID="cap_1234567890_abc"
./web-app

# Revoke when done
vault capability revoke cap_1234567890_abc \
  --reason "Database connection closed"
```

### ⚙️ **2. API Configuration**

```bash
# Request configuration access
vault capability request \
  --resource "secret:/api/production" \
  --action read,write \
  --ttl 600 \
  --identity "config-service" \
  --purpose "API configuration management"

# List all config capabilities
vault capability list --identity "config-service" --type "write"

# Monitor usage
vault capability list --identity "config-service" --status "active" --format json | jq '.capabilities[] | select(.used_count > 10)'
```

### 🐳 **3. Container Deployment**

```bash
# Request capability with container context
vault capability request \
  --resource "secret:/docker/registry" \
  --action read \
  --context '{
    "runtime": {
      "type": "docker",
      "id": "container123",
      "image": "myapp:latest"
    },
    "sourceIP": "10.0.0.100"
  }' \
  --constraints '{
    "environment": {
      "container.namespace": "production"
    }
  }'

# Validate with runtime context
vault capability validate cap_1234567890_xyz \
  --context '{
    "runtime": {
      "type": "docker",
      "id": "container123"
    },
    "sourceIP": "10.0.0.100"
  }'
```

### 🚨 **4. Emergency Revocation**

```bash
# List all active capabilities for a compromised identity
vault capability list --identity "compromised-app" --status "active"

# Revoke all capabilities (script)
for cap_id in $(vault capability list --identity "compromised-app" --format json | jq -r '.capabilities[].id'); do
  vault capability revoke "$cap_id" --reason "Security incident - identity compromised"
done

# Verify revocation
vault capability list --identity "compromised-app" --status "revoked"
```

### 🏢 **5. Enterprise Workflow**

```bash
# Step 1: Request deployment capability
DEPLOY_CAP=$(vault capability request \
  --resource "operation:/deploy/production" \
  --action execute \
  --ttl 1800 \
  --identity "deployment-service" \
  --purpose "Production deployment" \
  --constraints '{
    "timeWindow": {
      "hours": [22,23,0,1,2,3,4,5],
      "days": ["sat","sun"]
    },
    "environment": {
      "type": "production",
      "approval": "required"
    }
  }' --format json | jq -r '.capability.id')

# Step 2: Request database access
DB_CAP=$(vault capability request \
  --resource "secret:/db/primary" \
  --action read,write \
  --ttl 900 \
  --identity "deployment-service" \
  --purpose "Database migration during deployment" \
  --format json | jq -r '.capability.id')

# Step 3: Execute deployment
./deploy.sh --deployment-capability "$DEPLOY_CAP" --database-capability "$DB_CAP"

# Step 4: Cleanup
vault capability revoke "$DEPLOY_CAP" --reason "Deployment completed"
vault capability revoke "$DB_CAP" --reason "Database migration completed"
```

---

## 🔍 Error Handling

### 🚨 **Common Errors**

| Error                      | Cause                          | Solution                                    |
| -------------------------- | ------------------------------ | ------------------------------------------- |
| `resource cannot be empty` | Missing `--resource` flag      | Add `--resource` flag                       |
| `actions cannot be empty`  | Missing `--action` flag        | Add `--action` flag                         |
| `capability not found`     | Invalid capability ID          | Check capability ID with `list`             |
| `connection refused`       | Agent not running              | Start agent with `vault agent start`        |
| `policy denied`            | Request violates policy        | Check policies and adjust request           |
| `constraint violation`     | Request violates constraints   | Adjust constraints or context               |
| `capability expired`       | Capability has expired         | Request new capability                      |
| `max uses exceeded`        | Capability usage limit reached | Request new capability or increase max uses |

### 🔧 **Troubleshooting**

#### **1️⃣ Check Agent Status**

```bash
vault agent status
```

#### **2️⃣ Verify Policy Configuration**

```bash
vault capability request --resource "secret:/test" --action read --verbose
```

#### **3️⃣ Review Audit Logs**

```bash
tail -f ~/.aether-vault/audit.log | grep "capability_request\|capability_validate"
```

#### **4️⃣ Validate Configuration**

```bash
vault agent config --validate
```

#### **5️⃣ Debug Capability Request**

```bash
vault capability request \
  --resource "secret:/test" \
  --action read \
  --verbose \
  --debug
```

---

## 🏆 Best Practices

### ⏰ **1. Use Minimal TTL**

```bash
# ✅ Good: Short TTL for reduced risk
vault capability request --resource "secret:/db" --action read --ttl 300

# ❌ Avoid: Long TTL increases risk
vault capability request --resource "secret:/db" --action read --ttl 3600
```

### 🎯 **2. Request Only Necessary Actions**

```bash
# ✅ Good: Request only read access
vault capability request --resource "secret:/config" --action read

# ❌ Avoid: Requesting unnecessary admin access
vault capability request --resource "secret:/config" --action admin
```

### 📝 **3. Include Purpose and Context**

```bash
# ✅ Good: Include purpose for audit trail
vault capability request \
  --resource "secret:/db" \
  --action read \
  --purpose "Database connection for web-app" \
  --context '{"runtime": {"type": "web-server"}, "sourceIP": "10.0.0.100"}'
```

### 📊 **4. Monitor and Revoke**

```bash
# List active capabilities regularly
vault capability list --status "active"

# Revoke unused capabilities
vault capability revoke cap_1234567890_abc --reason "No longer needed"

# Monitor high-usage capabilities
vault capability list --format json | jq '.capabilities[] | select(.used_count > .max_uses * 0.8)'
```

### 🔒 **5. Use Constraints**

```bash
# ✅ Good: Apply IP and time constraints
vault capability request \
  --resource "secret:/sensitive" \
  --action read \
  --constraints '{
    "ipAddresses": ["10.0.0.0/8"],
    "timeWindow": {"hours": [9,10,11,12,13,14,15,16,17]}
  }'
```

### 🔄 **6. Implement Renewal Strategy**

```bash
# Request capability with renewal in mind
vault capability request \
  --resource "secret:/api" \
  --action read \
  --ttl 1800 \
  --max-uses 1000

# Renew before expiration
if [ $(vault capability validate "$CAP_ID" --format json | jq -r '.valid') = "false" ]; then
  NEW_CAP=$(vault capability request --resource "secret:/api" --action read --format json | jq -r '.capability.id')
  # Update application with new capability
fi
```

---

## 💻 Integration Examples

### 🐹 **Go Application**

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/skygenesisenterprise/aether-vault/package/cli/internal/ipc"
    "github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

func main() {
    // Create IPC client with retry logic
    client, err := ipc.NewClient(&ipc.Config{
        SocketPath:  "~/.aether-vault/agent.sock",
        Timeout:     30 * time.Second,
        RetryCount:  3,
        RetryDelay:  time.Second,
    })
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    defer client.Close()

    // Connect to agent
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := client.ConnectContext(ctx); err != nil {
        log.Fatal("Failed to connect to agent:", err)
    }

    // Request capability with constraints
    request := &types.CapabilityRequest{
        Identity: "production-app",
        Resource: "secret:/db/primary",
        Actions:  []string{"read"},
        TTL:      300,
        Purpose:  "Database connection for production app",
        Constraints: []types.Constraint{
            {
                Type:  "ip",
                Value: "10.0.0.0/8",
            },
            {
                Type:  "time",
                Value: "09:00-17:00",
            },
        },
        Context: map[string]interface{}{
            "runtime": map[string]string{
                "type": "web-server",
                "version": "1.2.3",
            },
            "sourceIP": "10.0.1.100",
        },
    }

    response, err := client.RequestCapabilityContext(ctx, request)
    if err != nil {
        log.Fatal("Failed to request capability:", err)
    }

    if response.Status != "granted" {
        log.Fatal("Capability denied:", response.Message)
    }

    // Use capability
    capabilityID := response.Capability.ID

    // Validate before use
    result, err := client.ValidateCapabilityContext(ctx, capabilityID, nil)
    if err != nil || !result.Valid {
        log.Fatal("Capability validation failed")
    }

    fmt.Printf("✓ Capability %s is valid\n", capabilityID)

    // Simulate database work
    fmt.Println("📊 Connecting to database...")
    time.Sleep(2 * time.Second)
    fmt.Println("✅ Database operation completed")

    // Revoke capability when done
    err = client.RevokeCapabilityContext(ctx, capabilityID, "Database operation completed")
    if err != nil {
        log.Printf("Warning: Failed to revoke capability: %v", err)
    } else {
        fmt.Printf("✓ Capability revoked\n")
    }
}
```

### 🔧 **Shell Script**

```bash
#!/bin/bash

# production-deploy.sh - Production deployment with capability management

set -euo pipefail

# Configuration
readonly SCRIPT_NAME="$(basename "$0")"
readonly LOG_FILE="/var/log/${SCRIPT_NAME%.*}.log"
readonly DEPLOY_RESOURCE="operation:/deploy/production"
readonly DB_RESOURCE="secret:/db/primary"

# Logging function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

# Error handling
error_exit() {
    log "ERROR: $*"
    exit 1
}

# Request capability with retry
request_capability() {
    local resource="$1"
    local actions="$2"
    local ttl="$3"
    local purpose="$4"

    log "📋 Requesting capability for $resource..."

    local response
    response=$(vault capability request \
        --resource "$resource" \
        --action "$actions" \
        --ttl "$ttl" \
        --purpose "$purpose" \
        --constraints '{"ipAddresses": ["10.0.0.0/8"], "timeWindow": {"hours": [22,23,0,1,2,3,4,5]}}' \
        --format json 2>/dev/null) || error_exit "Failed to request capability"

    local cap_id expires_at
    cap_id=$(echo "$response" | jq -r '.capability.id')
    expires_at=$(echo "$response" | jq -r '.capability.expires_at')

    if [[ "$cap_id" == "null" || "$expires_at" == "null" ]]; then
        error_exit "Invalid capability response"
    fi

    log "✓ Capability granted: $cap_id"
    log "✓ Expires at: $expires_at"

    echo "$cap_id"
}

# Validate capability
validate_capability() {
    local cap_id="$1"
    log "🔍 Validating capability..."

    local validation
    validation=$(vault capability validate "$cap_id" --format json 2>/dev/null) || error_exit "Failed to validate capability"

    local is_valid
    is_valid=$(echo "$validation" | jq -r '.valid')

    if [[ "$is_valid" != "true" ]]; then
        error_exit "Capability validation failed"
    fi

    log "✓ Capability is valid"
}

# Revoke capability
revoke_capability() {
    local cap_id="$1"
    local reason="$2"

    log "🗑️  Revoking capability..."
    if vault capability revoke "$cap_id" --reason "$reason" 2>/dev/null; then
        log "✓ Capability revoked"
    else
        log "⚠️  Failed to revoke capability"
    fi
}

# Main deployment
main() {
    log "🚀 Starting production deployment with $SCRIPT_NAME"

    local deploy_cap db_cap=""

    # Setup cleanup trap
    trap 'revoke_capability "$deploy_cap" "Deployment completed"; revoke_capability "$db_cap" "Database access completed"' EXIT INT TERM

    # Request deployment capability
    deploy_cap=$(request_capability "$DEPLOY_RESOURCE" "execute" 1800 "Production deployment")

    # Request database capability
    db_cap=$(request_capability "$DB_RESOURCE" "read,write" 900 "Database migration")

    # Validate capabilities
    validate_capability "$deploy_cap"
    validate_capability "$db_cap"

    # Execute deployment
    log "🚀 Executing deployment..."
    ./deploy.sh --deployment-capability "$deploy_cap" --database-capability "$db_cap"

    log "🎉 Production deployment completed successfully!"
}

# Run main function
main "$@"
```

### 🐍 **Python Integration**

```python
#!/usr/bin/env python3

"""
capability_manager.py - Python capability management
"""

import asyncio
import json
import logging
import subprocess
import sys
from datetime import datetime, timedelta
from typing import Dict, List, Optional

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class CapabilityManager:
    """Manages Aether Vault capabilities"""

    def __init__(self):
        self.capabilities: Dict[str, Dict] = {}

    async def request_capability(
        self,
        resource: str,
        actions: List[str],
        ttl: int = 300,
        purpose: str = "",
        constraints: Optional[Dict] = None,
        context: Optional[Dict] = None
    ) -> Dict:
        """Request a new capability"""

        logger.info(f"Requesting capability for {resource}")

        cmd = [
            "vault", "capability", "request",
            "--resource", resource,
            "--action", ",".join(actions),
            "--ttl", str(ttl)
        ]

        if purpose:
            cmd.extend(["--purpose", purpose])

        if constraints:
            cmd.extend(["--constraints", json.dumps(constraints)])

        if context:
            cmd.extend(["--context", json.dumps(context)])

        cmd.extend(["--format", "json"])

        try:
            result = await asyncio.create_subprocess_exec(
                *cmd,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )

            stdout, stderr = await asyncio.wait_for(
                result.communicate(),
                timeout=30.0
            )

            if result.returncode != 0:
                raise Exception(f"Command failed: {stderr.decode()}")

            response = json.loads(stdout.decode())

            if response.get("status") != "granted":
                raise Exception(f"Capability denied: {response.get('message')}")

            capability = response["capability"]
            self.capabilities[capability["id"]] = capability

            logger.info(f"✓ Capability granted: {capability['id']}")
            return capability

        except Exception as e:
            logger.error(f"Failed to request capability: {e}")
            raise

    async def validate_capability(self, capability_id: str, context: Optional[Dict] = None) -> bool:
        """Validate a capability"""

        logger.info(f"Validating capability {capability_id}")

        cmd = ["vault", "capability", "validate", capability_id, "--format", "json"]

        if context:
            cmd.extend(["--context", json.dumps(context)])

        try:
            result = await asyncio.create_subprocess_exec(
                *cmd,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )

            stdout, stderr = await asyncio.wait_for(
                result.communicate(),
                timeout=10.0
            )

            if result.returncode != 0:
                raise Exception(f"Validation failed: {stderr.decode()}")

            response = json.loads(stdout.decode())
            is_valid = response.get("valid", False)

            logger.info(f"Capability {capability_id} validation: {'✓' if is_valid else '✗'}")
            return is_valid

        except Exception as e:
            logger.error(f"Failed to validate capability: {e}")
            return False

    async def revoke_capability(self, capability_id: str, reason: str = "Manual revocation") -> bool:
        """Revoke a capability"""

        logger.info(f"Revoking capability {capability_id}")

        cmd = [
            "vault", "capability", "revoke", capability_id,
            "--reason", reason
        ]

        try:
            result = await asyncio.create_subprocess_exec(
                *cmd,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )

            stdout, stderr = await asyncio.wait_for(
                result.communicate(),
                timeout=10.0
            )

            if result.returncode != 0:
                raise Exception(f"Revocation failed: {stderr.decode()}")

            if capability_id in self.capabilities:
                del self.capabilities[capability_id]

            logger.info(f"✓ Capability {capability_id} revoked")
            return True

        except Exception as e:
            logger.error(f"Failed to revoke capability: {e}")
            return False

    async def list_capabilities(self, identity: Optional[str] = None, status: Optional[str] = None) -> List[Dict]:
        """List capabilities"""

        logger.info("Listing capabilities")

        cmd = ["vault", "capability", "list", "--format", "json"]

        if identity:
            cmd.extend(["--identity", identity])

        if status:
            cmd.extend(["--status", status])

        try:
            result = await asyncio.create_subprocess_exec(
                *cmd,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )

            stdout, stderr = await asyncio.wait_for(
                result.communicate(),
                timeout=10.0
            )

            if result.returncode != 0:
                raise Exception(f"List failed: {stderr.decode()}")

            response = json.loads(stdout.decode())
            return response.get("capabilities", [])

        except Exception as e:
            logger.error(f"Failed to list capabilities: {e}")
            return []

async def main():
    """Example usage"""

    manager = CapabilityManager()

    try:
        # Request capability
        cap = await manager.request_capability(
            resource="secret:/db/primary",
            actions=["read"],
            ttl=300,
            purpose="Database connection for Python app",
            constraints={
                "ipAddresses": ["10.0.0.0/8"],
                "timeWindow": {"hours": [9,10,11,12,13,14,15,16,17]}
            }
        )

        cap_id = cap["id"]

        # Validate capability
        if await manager.validate_capability(cap_id):
            logger.info("✓ Using capability for database access")
            # Simulate database work
            await asyncio.sleep(2)
            logger.info("✅ Database operation completed")

        # Revoke capability
        await manager.revoke_capability(cap_id, "Python demo completed")

    except Exception as e:
        logger.error(f"Demo failed: {e}")
        sys.exit(1)

if __name__ == "__main__":
    asyncio.run(main())
```

---

<div align="center">

### 🎉 **Master Capability Commands - Complete Control Over Your CBAC System!**

[🚀 Quick Start](QUICK_START.md) • [🔧 Agent Commands](COMMANDS_AGENT.md) • [⚙️ Configuration](CONFIG_OVERVIEW.md) • [🔐 CBAC Overview](CBAC_OVERVIEW.md)

---

**🔐 Enterprise-Grade Capability Management with Comprehensive Security!**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

_Building modern DevOps security infrastructure_

</div>
