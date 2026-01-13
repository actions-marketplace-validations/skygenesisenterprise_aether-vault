<div align="center">

# 🚀 Aether Vault CLI - Integration Overview

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![Go](https://img.shields.io/badge/Go-1.25+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![Cobra](https://img.shields.io/badge/Cobra-1.8+-lightgrey?style=for-the-badge&logo=go)](https://github.com/spf13/cobra) [![Viper](https://img.shields.io/badge/Viper-1.16+-green?style=for-the-badge&logo=go)](https://github.com/spf13/viper) [![DevOps](https://img.shields.io/badge/DevOps-Ready-orange?style=for-the-badge&logo=devops)](https://www.devops.com/)

**🔗 Complete Integration Guide - Connect Your Infrastructure**

Comprehensive integration guide covering IPC client integration, container deployment patterns, CI/CD pipelines, application frameworks, and best practices for incorporating Aether Vault's capability-based access control into your existing infrastructure.

[🔗 Integration Approaches](#-integration-approaches) • [📋 Integration Patterns](#-integration-patterns) • [🐳 Container Integration](#-container-integration) • [☸️ Kubernetes Integration](#️-kubernetes-integration) • [🚀 CI/CD Integration](#-ci-cd-integration) • [💻 Application Integration](#-application-integration) • [🔍 Best Practices](#-best-practices) • [🧪 Testing Integration](#-testing-integration) • [🔧 Troubleshooting](#-troubleshooting)

[![GitHub stars](https://img.shields.io/github/stars/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/stargazers) [![GitHub forks](https://img.shields.io/github/forks/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/network)

</div>

---

## 🔗 Integration Approaches

### 1️⃣ **IPC Client Integration**

The primary integration method is through IPC (Inter-Process Communication) client, which communicates with Aether Vault Agent via Unix sockets.

**When to Use:**

- 🖥️ Applications running on the same host as the agent
- 📦 Local services and daemons
- 🐳 Containerized applications with sidecar agents
- 🧪 Development and testing environments

**Benefits:**

- ⚡ Low latency communication
- 🌐 No network dependencies
- 📦 Simple client library
- 🔄 Local fallback operation

---

### 2️⃣ **Sidecar Pattern**

Deploy Aether Vault Agent as a sidecar container alongside your application containers.

**When to Use:**

- ☸️ Kubernetes deployments
- 🐳 Docker Compose applications
- 🏗️ Microservices architectures
- 🔄 Container orchestration platforms

**Benefits:**

- 🏠 Isolated agent lifecycle
- 🤝 Shared access within pod
- 📈 Easy scaling and updates
- 🏢 Standard container patterns

---

### 3️⃣ **Daemon Service Pattern**

Run Aether Vault Agent as a system-wide daemon service.

**When to Use:**

- 🖥️ Host-level services
- 🏢 Multiple applications per host
- 🏭️ Traditional VM deployments
- 🦿 Bare metal installations

**Benefits:**

- 🤝 Shared resource across applications
- 🔄 Centralized management
- ⚡ System-level integration
- 🚀 Boot-time startup

---

### 4️⃣ **Library Integration**

Embed Aether Vault capabilities directly into your application using Go libraries.

**When to Use:**

- 🐹 Go applications with high-performance requirements
- 🔧 Custom integration needs
- 📦 Embedded systems
- 🚀 Maximum performance scenarios

**Benefits:**

- 🚀 Maximum performance
- 🔧 Custom behavior
- 📦 No external dependencies
- 💪 Tight integration

---

## 📋 Integration Patterns

### 🔄 **Request-Response Pattern**

Applications request capabilities when needed and use them immediately.

```go
// Request capability when needed
capability, err := client.RequestCapability(&types.CapabilityRequest{
    Resource: "secret:/db/primary",
    Actions:  []string{"read"},
    TTL:      300,
})

// Use capability immediately
if err == nil && capability.Status == "granted" {
    secret := getSecret(capability.Capability.ID)
}
```

**Use Cases:**

- 📦 On-demand secret access
- 🔄 Sporadic resource usage
- 🎯 Event-driven applications
- 🌐 Microservices calls

### 💾 **Capability Caching Pattern**

Cache capabilities for short-term reuse to reduce request overhead.

```go
// Check cache first
if capability, exists := cache.Get("db_read"); exists && !isExpired(capability) {
    return useCapability(capability)
}

// Request new capability
capability, err := client.RequestCapability(request)
if err == nil {
    // Cache for future use
    cache.Set("db_read", capability, 5*time.Minute)
    return useCapability(capability)
}
```

**Use Cases:**

- ⚡ High-frequency operations
- 📦 Batch processing
- 🚀 Long-running applications
- 📊 Performance-critical services

### 🔐 **Pre-Authorization Pattern**

Request capabilities during application startup for known operations.

```go
func initializeApp() error {
    // Pre-authorize common operations
    capabilities := map[string]*types.Capability{
        "db_read": requestCapability("secret:/db/primary", "read", 300),
        "config_read": requestCapability("secret:/config/app", "read", 300),
        "log_write": requestCapability("log:/app", "write", 300),
    }

    // Store for application lifetime
    app.capabilities = capabilities
    return nil
}
```

**Use Cases:**

- 🚀 Service initialization
- 🔧 Known resource requirements
- 📈 Startup optimization
- 🔄 Background services

### ⚡ **Just-In-Time Pattern**

Request capabilities immediately before each use for maximum security.

```go
func getDatabaseConnection() (*sql.DB, error) {
    // Request capability just before use
    capability, err := client.RequestCapability(&types.CapabilityRequest{
        Resource: "secret:/db/primary",
        Actions:  []string{"read"},
        TTL:      60, // Very short TTL
    })

    if err != nil || capability.Status != "granted" {
        return nil, fmt.Errorf("failed to get database capability")
    }

    // Use capability immediately
    return connectToDatabase(capability.Capability.ID)
}
```

**Use Cases:**

- 🔒 High-security environments
- 🛡️ Sensitive operations
- 📊 Audit-critical applications
- 🚨 Compliance requirements

---

## 🐳 Container Integration

### 🏠 **Method 1: Host Socket Mount**

Mount the host's Unix socket into the container:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o vault-cli ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/vault-cli /usr/local/bin/vault
COPY --from=builder /app/cmd/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Create vault user
RUN adduser -D -s /bin/sh vault

USER vault
ENTRYPOINT ["/entrypoint.sh"]
```

```bash
#!/bin/bash
# entrypoint.sh

# Wait for agent socket
while [ ! -S /var/run/vault/agent.sock ]; do
    echo "Waiting for Vault agent socket..."
    sleep 1
done

echo "Vault agent socket found, starting application..."
exec "$@"
```

```yaml
# docker-compose.yml
version: "3.8"

services:
  app:
    build: .
    volumes:
      - /var/run/vault:/var/run/vault:ro
    environment:
      - VAULT_AGENT_SOCKET_PATH=/var/run/vault/agent.sock
    depends_on:
      - vault-agent

  vault-agent:
    image: aether-vault/agent:latest
    volumes:
      - /var/run/vault:/var/run/vault
      - ./config:/etc/aether-vault
    command: ["start", "--config", "/etc/aether-vault/agent.yaml"]
```

### 🚀 **Method 2: Sidecar Container**

Deploy agent as a sidecar:

```yaml
# kubernetes.yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app
spec:
  containers:
    - name: app
      image: my-app:latest
      volumeMounts:
        - name: vault-socket
          mountPath: /var/run/vault
      env:
        - name: VAULT_AGENT_SOCKET_PATH
          value: /var/run/vault/agent.sock

    - name: vault-agent
      image: aether-vault/agent:latest
      args:
        - "start"
        - "--mode"
        - "standard"
      volumeMounts:
        - name: vault-socket
          mountPath: /var/run/vault
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000

  volumes:
    - name: vault-socket
      emptyDir: {}
```

### 🔧 **Method 3: Init Container**

Use init container to set up capabilities:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app
spec:
  initContainers:
    - name: vault-init
      image: aether-vault/cli:latest
      command: ["sh", "-c"]
      args:
        - |
          vault capability request \
            --resource "secret:/db/primary" \
            --action read \
            --ttl 3600 \
            --output /etc/capabilities/db.json
      volumeMounts:
        - name: capabilities
          mountPath: /etc/capabilities

  containers:
    - name: app
      image: my-app:latest
      volumeMounts:
        - name: capabilities
          mountPath: /etc/capabilities
          readOnly: true
      env:
        - name: CAPABILITY_FILE
          value: /etc/capabilities/db.json
```

---

## ☸️ Kubernetes Integration

### 🏛️ **Method 1: DaemonSet**

Deploy agent as a DaemonSet for node-level access:

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: vault-agent
  labels:
    app: vault-agent
spec:
  selector:
    matchLabels:
      app: vault-agent
  template:
    metadata:
      labels:
        app: vault-agent
    spec:
      serviceAccountName: vault-agent
      containers:
        - name: vault-agent
          image: aether-vault/agent:latest
          args:
            - "start"
            - "--mode"
            - "standard"
          volumeMounts:
            - name: socket-dir
              mountPath: /var/run/vault
            - name: config-dir
              mountPath: /etc/aether-vault
          securityContext:
            runAsUser: 1000
            runAsGroup: 1000
      volumes:
        - name: socket-dir
          hostPath: /var/run
        - name: config-dir
          hostPath: /etc/aether-vault

---
apiVersion: v1
kind: Service
metadata:
  name: vault-agent
  labels:
    app: vault-agent
spec:
  selector:
    app: vault-agent
  ports:
    - port: 9090
      name: metrics
    - port: 8080
      name: http
```

### 🚀 **Method 2: Service with Headless Service**

Use headless service for agent communication:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: vault-agent
spec:
  selector:
    app: vault-agent
  ports:
    - port: 8080
      name: http
---
apiVersion: v1
kind: Endpoints
metadata:
  name: vault-agent
subsets:
  - addresses:
      - port: 8080
```

### 🔧 **Method 3: ConfigMap with Volumes**

Use ConfigMaps to share agent configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: vault-agent-config
data:
  agent.yaml: |
    version: "1.0"
    mode: "standard"
    capabilities:
      enable: true
      default_ttl: 300
      max_ttl: 3600

---
apiVersion: v1
kind: Pod
metadata:
  name: my-app
spec:
  containers:
    - name: app
      image: my-app:latest
      volumeMounts:
        - name: vault-socket
          mountPath: /var/run/vault
      env:
        - name: VAULT_AGENT_SOCKET_PATH
          value: /var/run/vault/agent.sock

    - name: vault-agent
      image: aether-vault/agent:latest
      args:
        - "start"
        - "--config"
        "/etc/aether-vault/agent.yaml"
      volumeMounts:
        - name: vault-socket
          mountPath: /var/run/vault
        - name: vault-config
          mountPath: /etc/aether-vault
          readOnly: true

  volumes:
    - name: vault-socket
      emptyDir: {}
    - name: vault-config
      configMap:
        name: vault-agent-config
```

---

## 🚀 CI/CD Integration

### 🐙 **GitHub Actions**

```yaml
name: Deploy with Vault Capabilities

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  setup-vault:
    runs-on: ubuntu-latest
    steps:
      - name: Setup Vault CLI
        run: |
          curl -L https://github.com/skygenesisenterprise/aether-vault/releases/latest/download/vault-linux-amd64 -o vault
          chmod +x vault
          sudo mv vault /usr/local/bin/

      - name: Start Vault Agent
        run: |
          vault agent start --mode development

      - name: Request Deployment Capability
        run: |
          vault capability request \
            --resource "operation:/deploy/production" \
            --action execute \
            --ttl 600 \
            --identity "github-actions" \
            --purpose "Production deployment" \
            --output capability.json

      - name: Deploy Application
        run: |
          ./deploy.sh --capability capability.json

      - name: Cleanup
        if: always()
        run: |
          vault agent stop

  deploy:
    needs: setup-vault
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Use Vault Capability
        run: |
          vault capability request \
            --resource "operation:/deploy/production" \
            --action execute \
            --ttl 600 \
            --identity "github-actions" \
            --purpose "Production deployment" \
            --output capability.json
          ./deploy.sh --capability capability.json
```

### 🔄 **Jenkins Pipeline**

```groovy
pipeline {
    agent any
    stages {
        stage('Setup Vault') {
            steps {
                sh 'curl -L https://github.com/skygenesisenterprise/aether-vault/releases/latest/download/vault-linux-amd64 -o vault'
                sh 'chmod +x vault'
                sh 'sudo mv vault /usr/local/bin/'

                sh 'vault agent start --mode development &'
                sh 'sleep 5'
            }
        }
        stage('Request Capability') {
            steps {
                sh '''
                    vault capability request \\
                        --resource "secret:/build/${env.BRANCH_NAME}" \\
                        --action read,write \\
                        --ttl 300 \\
                        --identity "jenkins" \\
                        --purpose "Build pipeline" \\
                        --output capability.json
                '''
            }
        }
        stage('Build and Deploy') {
            steps {
                sh './build.sh --capability capability.json'
                sh './deploy.sh'
            }
        }
        stage('Cleanup') {
            steps {
                sh 'vault agent stop || true'
            }
        }
    }
}
```

### 🐳 **Docker CI/CD**

```dockerfile
# Multi-stage build with Vault CLI
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Install Vault CLI
RUN curl -L https://github.com/skygenesisenterprise/aether-vault/releases/latest/download/vault-linux-amd64 -o /usr/local/bin/vault
RUN chmod +x /usr/local/bin/vault

WORKDIR /app
COPY --from=builder /app /app
COPY --from=builder /app/capability-handler.sh /capability-handler.sh
RUN chmod +x /capability-handler.sh

# Start agent and app
CMD ["/capability-handler.sh"]
```

```bash
#!/bin/bash
# capability-handler.sh

# Start vault agent in background
/usr/local/bin/vault agent start --mode development &

# Wait for agent
sleep 2

# Request build capability
/usr/local/bin/vault capability request \
    --resource "secret:/build/${TARGET}" \
    --action read,write \
    --ttl 300 \
    --identity "docker-ci" \
    --purpose "Docker build" \
    --output /tmp/build-capability.json

# Run build
./build.sh --capability /tmp/build-capability.json

# Cleanup
/usr/local/bin/vault agent stop || true
```

---

## 💻 Application Integration

### 🐹 **Go Application**

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/skygenesisenterprise/aether-vault/package/cli/internal/ipc"
    "github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

type VaultClient struct {
    client  *ipc.Client
    cache   map[string]*types.Capability
    mutex   sync.RWMutex
}

func NewVaultClient() (*VaultClient, error) {
    client, err := ipc.NewClient(nil)
    if err != nil {
        return nil, err
    }

    return &VaultClient{
        client: client,
        cache:  make(map[string]*types.Capability),
    }
}

func (v *VaultClient) GetCapability(resource, action string, ttl int64) (*types.Capability, error) {
    v.mutex.RLock()
    defer v.mutex.RUnlock()

    cacheKey := fmt.Sprintf("%s:%s", resource, action)

    // Check cache first
    if cap, exists := v.cache[cacheKey]; exists && !isExpired(cap) {
        return cap, nil
    }

    // Request new capability
    request := &types.CapabilityRequest{
        Resource: resource,
        Actions:  []string{action},
        TTL:      ttl,
        Identity: os.Getenv("APP_ID"),
        Purpose:  fmt.Sprintf("Access to %s for %s", resource, os.Getenv("APP_NAME")),
    }

    response, err := v.client.RequestCapability(request)
    if err != nil {
        return nil, err
    }

    if response.Status != "granted" {
        return nil, fmt.Errorf("capability denied: %s", response.Message)
    }

    // Cache capability
    v.cache[cacheKey] = response.Capability
    return response.Capability, nil
}

func (v *VaultClient) UseCapability(capabilityID string, usage func() error) error {
    // Validate capability before use
    result, err := v.client.ValidateCapability(capabilityID, nil)
    if err != nil {
        return err
    }
    if !result.Valid {
        return fmt.Errorf("capability validation failed")
    }

    // Use capability
    if err := usage(); err != nil {
        // Revoke capability on error
        v.client.RevokeCapability(capabilityID, "Usage failed")
        return err
    }

    return nil
}

func main() {
    vault, err := NewVaultClient()
    if err != nil {
        log.Fatal(err)
    }
    defer vault.client.Close()

    // Example: Get database capability
    dbCap, err := vault.GetCapability("secret:/db/primary", "read", 300)
    if err != nil {
        log.Fatal(err)
    }

    // Use capability
    if err := vault.UseCapability(dbCap.ID, func() error {
        // Database connection logic here
        fmt.Println("Connecting to database...")
        time.Sleep(2 * time.Second)
        fmt.Println("Database operation completed")
        return nil
    }); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Successfully used capability %s\n", dbCap.ID)
}
```

### 🐍 **Python Application**

```python
#!/usr/bin/env python3

import json
import subprocess
import os
import time
from typing import Dict, Optional

class VaultClient:
    def __init__(self, socket_path: Optional[str] = None):
        self.socket_path = socket_path or os.getenv("VAULT_AGENT_SOCKET_PATH")
        self.cache = {}

    def _run_command(self, cmd: list) -> Dict:
        """Run vault CLI command and return JSON result"""
        full_cmd = ["vault"] + cmd
        if self.socket_path:
            full_cmd.extend(["--socket-path", self.socket_path])

        try:
            result = subprocess.run(
                full_cmd,
                capture_output=True,
                text=True,
                check=True
            )
            return json.loads(result.stdout)
        except subprocess.CalledProcessError as e:
            raise Exception(f"Vault command failed: {e}")

    def get_capability(self, resource: str, action: str, ttl: int = 300) -> Dict:
        """Request a capability"""
        cache_key = f"{resource}:{action}"

        # Check cache first
        if cache_key in self.cache:
            cap = self.cache[cache_key]
            expires_at = time.strptime(cap["expires_at"], "%Y-%m-%dT%H:%M:%SZ")
            if time.time() < expires_at:
                return cap

        # Request new capability
        response = self._run_command([
            "capability", "request",
            "--resource", resource,
            "--action", action,
            "--ttl", str(ttl),
            "--identity", os.getenv("APP_ID", ""),
            "--purpose", f"Access to {resource} for {os.getenv('APP_NAME', '')}",
            "--format", "json"
        ])

        if response.get("status") != "granted":
            raise Exception(f"Capability denied: {response.get('message', 'Unknown error')}")

        # Cache capability
        self.cache[cache_key] = response["capability"]
        return response["capability"]

    def use_capability(self, capability_id: str, usage_func) -> bool:
        """Use a capability"""
        # Validate capability
        result = self._run_command(["capability", "validate", capability_id, "--format", "json"])
        if not result.get("valid", False):
            print(f"Capability {capability_id} validation failed")
            return False

        # Use capability
        try:
            success = usage_func()
            return success
        except Exception as e:
            print(f"Error using capability: {e}")
            return False

    def revoke_capability(self, capability_id: str, reason: str = "Usage completed") -> bool:
        """Revoke a capability"""
        try:
            subprocess.run([
                "capability", "revoke", capability_id,
                "--reason", reason
            ], check=True)
            return True
        except Exception as e:
            print(f"Failed to revoke capability: {e}")
            return False

# Example usage
def main():
    vault = VaultClient()

    try:
        # Get database capability
        db_cap = vault.get_capability("secret:/db/primary", "read", 300)

        # Use capability for database access
        if vault.use_capability(db_cap["id"], lambda: access_database(db_cap["id"])):
            print("Database operation completed successfully")
        else:
            print("Database operation failed")

        # Revoke capability when done
        vault.revoke_capability(db_cap["id"], "Database operation completed")

    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main()
```

### 🟨 **Node.js Application**

```javascript
const { spawn } = require("child_process");
const fs = require("fs");
const path = require("path");

class VaultClient {
  constructor(socketPath) {
    this.socketPath = socketPath || process.env.VAULT_AGENT_SOCKET_PATH;
    this.cache = new Map();
  }

  async _runCommand(args) {
    const cmd = ["vault", ...args];
    if (this.socketPath) {
      cmd.push("--socket-path", this.socketPath);
    }

    return new Promise((resolve, reject) => {
      const child = spawn(cmd[0], cmd.slice(1), {
        stdio: ["pipe", "pipe"],
        env: {
          ...process.env,
        },
      });

      let stdout = "";
      let stderr = "";

      child.stdout.on("data", (data) => {
        stdout += data.toString();
      });

      child.stderr.on("data", (data) => {
        stderr += data.toString();
      });

      child.on("close", (code) => {
        if (code !== 0) {
          reject(new Error(`Vault command failed: ${stderr}`));
          return;
        }

        try {
          const result = JSON.parse(stdout);
          resolve(result);
        } catch (e) {
          resolve(stdout.trim());
        }
      });
    });
  }

  async getCapability(resource, action, ttl = 300) {
    const cacheKey = `${resource}:${action}`;

    // Check cache first
    if (this.cache.has(cacheKey)) {
      const cap = this.cache.get(cacheKey);
      const expiresAt = new Date(cap.expires_at);
      if (Date.now() < expiresAt) {
        return cap;
      }
      this.cache.delete(cacheKey);
    }

    // Request new capability
    const response = await this._runCommand([
      "capability",
      "request",
      "--resource",
      resource,
      "--action",
      action,
      "--ttl",
      ttl.toString(),
      "--identity",
      process.env.APP_ID || "",
      "--purpose",
      `Access to ${resource} for ${process.env.APP_NAME || ""}`,
      "--format",
      "json",
    ]);

    if (response.status !== "granted") {
      throw new Error(
        `Capability denied: ${response.message || "Unknown error"}`,
      );
    }

    // Cache capability
    this.cache.set(cacheKey, response.capability);
    return response.capability;
  }

  async validateCapability(capabilityId) {
    const result = await this._runCommand([
      "capability",
      "validate",
      capabilityId,
      "--format",
      "json",
    ]);

    return result.valid || false;
  }

  async useCapability(capabilityId, usageFunc) {
    // Validate capability
    const isValid = await this.validateCapability(capabilityId);
    if (!isValid) {
      throw new Error(`Capability ${capabilityId} validation failed`);
    }

    // Use capability
    try {
      await usageFunc();
      return true;
    } catch (error) {
      console.error(`Error using capability: ${error}`);
      return false;
    }
  }

  async revokeCapability(capabilityId, reason = "Usage completed") {
    await this._runCommand([
      "capability",
      "revoke",
      capabilityId,
      "--reason",
      reason,
    ]);
  }
}

// Example usage
async function main() {
  const vault = new VaultClient();

  try {
    // Get database capability
    const dbCap = await vault.getCapability("secret:/db/primary", "read", 300);

    // Use capability for database access
    const success = await vault.useCapability(dbCap.id, async () => {
      console.log(`Using capability ${dbCap.id} for database access`);
      // Database connection logic here
      return true;
    });

    if (success) {
      console.log("Database operation completed successfully");
    }

    // Revoke capability when done
    await vault.revokeCapability(dbCap.id, "Database operation completed");
  } catch (error) {
    console.error(`Error: ${error.message}`);
  }
}

module.exports = { VaultClient };
```

---

## 🔍 Best Practices

### 1️⃣ **Error Handling**

```go
// Robust error handling with retries
func (v *VaultClient) GetCapabilityWithRetry(resource, action string, ttl int64, maxRetries int) (*types.Capability, error) {
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        cap, err := v.GetCapability(resource, action, ttl)
        if err == nil {
            return cap, nil
        }

        lastErr = err

        // Check if error is retryable
        if !isRetryableError(err) {
            break
        }

        // Exponential backoff
        backoff := time.Duration(1<<uint(i)) * time.Second
        time.Sleep(backoff)
    }

    return nil, lastErr
}

func isRetryableError(err error) bool {
    return strings.Contains(err.Error(), "connection refused") ||
           strings.Contains(err.Error(), "timeout")
}
```

### 2️⃣ **Capability Lifecycle Management**

```go
type CapabilityManager struct {
    client   *VaultClient
    caps     map[string]*types.Capability
    mutex    sync.RWMutex
    cleanup  chan string
}

func (cm *CapabilityManager) RequestCapability(resource, action string, ttl int64) (*types.Capability, error) {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()

    // Check existing capability
    key := fmt.Sprintf("%s:%s", resource, action)
    if cap, exists := cm.caps[key]; exists {
        if time.Now().Before(cap.ExpiresAt) {
            return cap, nil
        }
        // Capability expired, remove it
        delete(cm.caps, key)
    }

    // Request new capability
    cap, err := cm.client.RequestCapability(&types.CapabilityRequest{
        Resource: resource,
        Actions:  []string{action},
        TTL:      ttl,
    })

    if err != nil {
        return nil, err
    }

    // Store capability
    cm.caps[key] = cap

    // Schedule cleanup
    go func() {
        time.Sleep(time.Until(cap.ExpiresAt))
        cm.cleanup <- key
    }()

    return cap, nil
}

func (cm *CapabilityManager) StartCleanupWorker() {
    go func() {
        for key := range cm.cleanup {
            delete(cm.caps, key)
        }
    }
}
```

### 3️⚡ **Performance Optimization**

```go
// Batch capability requests
func (v *VaultClient) BatchRequestCapabilities(requests []CapabilityRequest) ([]CapabilityResponse, error) {
    responses := make([]CapabilityResponse, len(requests))

    // Send requests in parallel
    var wg sync.WaitGroup
    for i, req := range requests {
        wg.Add(1)
        go func(idx int, r CapabilityRequest) {
            defer wg.Done()

            response, err := v.client.RequestCapability(&r)
            if err != nil {
                responses[idx] = CapabilityResponse{Error: err.Error()}
                return
            }
            responses[idx] = *response
        }(i, req)
    }

    wg.Wait()
    return responses, nil
}

// Connection pooling
type ConnectionPool struct {
    connections chan *ipc.Client
    available    chan *ipc.Client
    maxConnections int
}

func (p *ConnectionPool) Get() (*ipc.Client, error) {
    select {
    case conn := <-p.available:
        return conn, nil
    default:
        // Create new connection
        conn, err := ipc.NewClient(nil)
        if err != nil {
            return nil, err
        }
        return conn, nil
    }
}

func (p *ConnectionPool) Put(conn *ipc.Client) {
    select {
    case p.available <- conn:
    default:
        // Pool full, discard connection
        conn.Close()
    }
}
```

---

## 🧪 Testing Integration

### 🧪 **Unit Testing with Mock Client**

```go
// Mock Vault client for testing
type MockVaultClient struct {
    capabilities map[string]*types.Capability
    responses    map[string]*types.CapabilityResponse
    mutex        sync.RWMutex
}

func (m *MockVaultClient) RequestCapability(req *types.CapabilityRequest) (*types.CapabilityResponse, error) {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    key := fmt.Sprintf("%s:%s", req.Resource, strings.Join(req.Actions, ","))

    // Check if response exists
    if resp, exists := m.responses[key]; exists {
        return resp, nil
    }

    // Return mock response
    resp := &types.CapabilityResponse{
        Status: "granted",
        Capability: &types.Capability{
            ID:      fmt.Sprintf("mock-cap-%s", key),
            Resource: req.Resource,
            Actions: req.Actions,
            ExpiresAt: time.Now().Add(time.Duration(req.TTL) * time.Second),
        },
    }

    m.responses[key] = resp
    return resp, nil
}

func TestDatabaseAccess(t *testing.T) {
    mock := &MockVaultClient{}
    app := &Application{vault: mock}

    // Test capability request
    cap, err := mock.RequestCapability(&types.CapabilityRequest{
        Resource: "secret:/db/primary",
        Actions:  []string{"read"},
        TTL:      300,
    })

    require.NoError(t, err)
    require.NotNil(t, cap)

    // Test application logic
    err = app.AccessDatabase(cap.ID)
    require.NoError(t, err)
}
```

### 🔧 **Integration Testing**

```bash
#!/bin/bash
# integration-test.sh

set -e

echo "🧪 Starting integration tests..."

# Test 1: IPC Connection
echo "📡 Testing IPC connection..."
vault capability status

# Test 2: Capability Request
echo "🔐 Testing capability request..."
CAP_ID=$(vault capability request \
  --resource "secret:/test" \
  --action read \
  --ttl 60 \
  --format json | jq -r '.capability.id')

echo "✅ Capability granted: $CAP_ID"

# Test 3: Capability Validation
echo "🔍 Testing capability validation..."
VALID=$(vault capability validate "$CAP_ID" --format json | jq -r '.valid')

if [ "$VALID" = "true" ]; then
    echo "✅ Capability validation passed"
else
    echo "❌ Capability validation failed"
    exit 1
fi

# Test 4: Capability Revocation
echo "🗑️ Testing capability revocation..."
REVOKE_RESULT=$(vault capability revoke "$CAP_ID" --reason "Test completed")

if [ "$REVOKE_RESULT" = "Capability revoked successfully" ]; then
    echo "✅ Capability revocation passed"
else
    echo "❌ Capability revocation failed"
    exit 1
fi

echo "🎉 All integration tests passed!"
```

---

## 🔧 Troubleshooting

### 🚨 **Common Integration Issues**

#### **Connection Problems**

```bash
# Check if agent is running
vault agent status

# Check socket file permissions
ls -la /var/run/vault/agent.sock

# Test IPC connection
vault capability status

# Check environment variables
echo "VAULT_AGENT_SOCKET_PATH: $VAULT_AGENT_SOCKET_PATH"

# Test with custom socket path
vault capability status --socket-path /tmp/test-vault.sock
```

#### **Capability Issues**

```bash
# Test capability request
vault capability request \
  --resource "secret:/test" \
  --action read \
  --verbose

# Check policies
ls ~/.aether-vault/policies/

# Review audit logs
tail -f ~/.aether-vault/audit.log | grep "capability_request"

# Validate configuration
vault agent config --validate
```

#### **Performance Issues**

```bash
# Monitor agent performance
vault agent status --verbose

# Check capability metrics
vault capability status --format json | jq '.capability_engine'

# Monitor IPC performance
time vault capability status
```

---

## 🔍 Advanced Patterns

### 🔒 **Multi-Region Deployment**

```yaml
# kubernetes-deployment.yaml
apiVersion: v1
kind: Deployment
metadata:
  name: vault-agents
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: vault-agent
        region: us-west-2
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: topology.kubernetes.io/zone
                    operator: In
                    values:
                      - us-west-2-a
                      - us-west-2-b
      containers:
        - name: vault-agent
          image: aether-vault/agent:latest
          args:
            - "start"
            - "--mode"
            - "standard"
            - "--region"
            - "us-west-2"
          env:
            - name: VAULT_AGENT_REGION
              value: "us-west-2"
            - name: VAULT_AGENT_MODE
              value: "standard"
          volumeMounts:
            - name: socket-dir
              mountPath: /var/run/vault
            - name: config-dir
              mountPath: /etc/aether-vault
          resources:
            limits:
              cpu: "100m"
              memory: "128Mi"
          securityContext:
            runAsUser: 1000
            runAsGroup: 1000
      volumes:
        - name: socket-dir
          hostPath: /var/run
        - name: config-dir
          hostPath: /etc/aether-vault

---
apiVersion: v1
kind: Service
metadata:
  name: vault-agent-service
spec:
  selector:
    app: vault-agent
  ports:
    - port: 8080
      name: http
    - port: 9090
      name: metrics
```

### 🔄 **High Availability Deployment**

```yaml
# vault-ha.yaml
apiVersion: v1
kind: Deployment
metadata:
  name: vault-agent-ha
  annotations:
    service.beta.kubernetes.io/load-balancer: "internal"
spec:
  replicas: 2
  selector:
    matchLabels:
      app: vault-agent-ha
  template:
    metadata:
      labels:
        app: vault-agent-ha
    spec:
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - podAffinityTerm:
                - labelSelector:
                  matchExpressions:
                    - key: app
                      operator: In
                      values:
                        - vault-agent-ha
      containers:
        - name: vault-agent
          image: aether-vault/agent:latest
          args:
            - "start"
            - "--mode"
            - "standard"
            - "--enable-ha"
          env:
            - name: VAULT_AGENT_HA
              value: "true"
            - name: VAULT_AGENT_PEER_ADDRESSES
              value: "vault-agent-ha-0.vault-agent-ha-1.vault-agent-ha-2"
          volumeMounts:
            - name: socket-dir
              mountPath: /var/run/vault
            - name: config-dir
              mountPath: /etc/aether-vault
          resources:
            limits:
              cpu: "200m"
              memory: "256Mi"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
              periodSeconds: 5
```

---

<div align="center">

### 🎉 **Master Integration Guide - Connect Your Infrastructure!**

[🚀 Quick Start](QUICK_START.md) • [🔧 Agent Commands](COMMANDS_AGENT.md) • [🔐 Capability Commands](COMMANDS_CAPABILITY.md) • [⚙️ Configuration](CONFIG_OVERVIEW.md)

---

**🔗 Enterprise-Grade Integration with Multiple Deployment Patterns!**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

_Building next-generation DevOps security infrastructure_

</div>
