<div align="center">

# 🚀 Aether Vault CLI - Quick Start

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![Go](https://img.shields.io/badge/Go-1.25+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![Cobra](https://img.shields.io/badge/Cobra-1.8+-lightgrey?style=for-the-badge&logo=go)](https://github.com/spf13/cobra) [![Viper](https://img.shields.io/badge/Viper-1.16+-green?style=for-the-badge&logo=go)](https://github.com/spf13/viper) [![DevOps](https://img.shields.io/badge/DevOps-Ready-orange?style=for-the-badge&logo=devops)](https://www.devops.com/)

**🔐 Get Started in Minutes - Enterprise-Grade Secret Management with CBAC System**

A comprehensive quick start guide to get you up and running with **Aether Vault CLI** in minutes. Learn installation, initialization, capability management, and integration patterns for modern DevOps workflows.

[📋 Prerequisites](#-prerequisites) • [🚀 Installation](#-installation) • [⚙️ Initial Setup](#️-initial-setup) • [🔐 First Capability](#-first-capability) • [💻 Integration Examples](#-integration-examples) • [🛠️ Configuration](#️-configuration) • [🔍 Troubleshooting](#-troubleshooting)

[![GitHub stars](https://img.shields.io/github/stars/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/stargazers) [![GitHub forks](https://img.shields.io/github/forks/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/network)

</div>

---

## 📋 Prerequisites

### 🎯 **System Requirements**

- **Operating System**: Linux, macOS, or Windows (WSL2)
- **Go**: Version 1.25 or later (for building from source)
- **Unix Socket Support**: Required for IPC communication
- **Git**: For cloning and version control
- **Make**: Build automation tool (included with most systems)

### 🔧 **Optional Dependencies**

- **Docker**: For containerized deployment
- **PostgreSQL**: For audit log storage (optional)
- **jq**: JSON processing for shell scripts
- **Python 3.8+**: For Python integration examples

---

## 🚀 Installation

### 🎯 **Option 1: Download Binary (Recommended)**

```bash
# Download the latest binary for your platform
curl -L https://github.com/skygenesisenterprise/aether-vault/releases/latest/download/vault-linux-amd64 -o vault

# Make it executable
chmod +x vault

# Move to system PATH
sudo mv vault /usr/local/bin/

# Verify installation
vault version
```

**Platform-specific downloads**:

```bash
# macOS (Apple Silicon)
curl -L https://github.com/skygenesisenterprise/aether-vault/releases/latest/download/vault-darwin-arm64 -o vault

# macOS (Intel)
curl -L https://github.com/skygenesisenterprise/aether-vault/releases/latest/download/vault-darwin-amd64 -o vault

# Windows
curl -L https://github.com/skygenesisenterprise/aether-vault/releases/latest/download/vault-windows-amd64.exe -o vault.exe
```

### ⚙️ **Option 2: Build from Source**

```bash
# Clone the repository
git clone https://github.com/skygenesisenterprise/aether-vault.git
cd aether-vault/package/cli

# Build the CLI
make build

# Install to system PATH
make install

# Verify installation
vault version
```

### 📦 **Option 3: Package Manager (Coming Soon)**

```bash
# macOS with Homebrew
brew install aether-vault/tap/vault

# Linux with apt
sudo apt update
sudo apt install aether-vault-cli

# Linux with yum
sudo yum install aether-vault-cli

# Windows with Chocolatey
choco install aether-vault-cli
```

---

## ⚙️ Initial Setup

### 🎯 **1. Initialize Local Environment**

```bash
# Initialize your local vault environment
vault init

# Output:
# ✓ Created configuration directory: /home/user/.aether-vault
# ✓ Generated default configuration file
# ✓ Created policy directory: /home/user/.aether-vault/policies
# ✓ Created audit log file: /home/user/.aether-vault/audit.log
# ✓ Local environment initialized successfully
```

**Custom initialization options**:

```bash
# Custom path
vault init --path /custom/path/to/vault

# Force reinitialization
vault init --force

# Verbose output
vault init --verbose
```

### 🚀 **2. Start the Agent**

```bash
# Start the Aether Vault Agent
vault agent start

# Output:
# Starting Aether Vault Agent...
# Configuration loaded from /home/user/.aether-vault/agent.yaml
# Capability engine initialized with Ed25519 signing
# Policy engine loaded 1 policies from /home/user/.aether-vault/policies
# Audit system started with file logging
# IPC server listening on /home/user/.aether-vault/agent.sock
# Agent started successfully (PID: 12345)
```

**Advanced agent options**:

```bash
# Start in development mode
vault agent start --mode development

# Custom configuration
vault agent start --config /custom/path/config.yaml

# Debug mode
vault agent start --log-level debug

# Custom socket path
vault agent start --socket-path /tmp/vault.sock
```

### 🔍 **3. Verify Setup**

```bash
# Check agent status
vault agent status

# Output:
# Aether Vault Agent Status:
#   Running: true
#   PID: 12345
#   Uptime: 2m30s
#   Version: 1.0.0
#
# IPC Server:
#   Socket Path: /home/user/.aether-vault/agent.sock
#   Active Connections: 0
#
# Capability Engine:
#   Status: Healthy
#   Total Capabilities: 0
#
# Policy Engine:
#   Active Policies: 1
#   Last Evaluation: 2025-01-13T10:05:00Z
```

**Comprehensive status check**:

```bash
# Full status with diagnostics
vault agent status --verbose

# JSON format for automation
vault agent status --format json

# Check specific components
vault agent status --component capability-engine
vault agent status --component policy-engine
vault agent status --component ipc-server
```

---

## 🔐 First Capability

### 🎯 **1. Request a Read Capability**

```bash
# Request a capability to read a database secret
vault capability request \
  --resource "secret:/db/primary" \
  --action read \
  --ttl 300 \
  --purpose "Database connection for my app"

# Output:
# Capability Request Result:
#   Status: granted
#   Request ID: req_1234567890_abcdef
#   Processing Time: 45ms
#
# Capability Details:
#   ID: cap_1234567890_ghijkl
#   Type: read
#   Resource: secret:/db/primary
#   Actions: read
#   Identity: user
#   Issuer: aether-vault-agent
#   TTL: 300 seconds
#   Max Uses: 100
#   Issued At: 2025-01-13T10:00:00Z
#   Expires At: 2025-01-13T10:05:00Z
```

**Advanced capability requests**:

```bash
# Multiple actions
vault capability request \
  --resource "secret:/app/config" \
  --action read,write \
  --ttl 600 \
  --max-uses 50 \
  --constraint "ip:192.168.1.0/24"

# Time-based constraints
vault capability request \
  --resource "operation:/deploy/production" \
  --action execute \
  --ttl 1800 \
  --constraint "time:09:00-17:00" \
  --constraint "day:mon-fri"

# Environment constraints
vault capability request \
  --resource "secret:/db/backup" \
  --action read,write \
  --ttl 3600 \
  --constraint "env:production"
```

### 🔍 **2. Validate the Capability**

```bash
# Validate the capability (using the ID from above)
vault capability validate cap_1234567890_ghijkl

# Output:
# Capability Validation Result:
#   Valid: true
#   Validation Time: 12ms
#   Checked Constraints: 0
#
# Capability Info:
#   ID: cap_1234567890_ghijkl
#   Resource: secret:/db/primary
#   Actions: read
#   Uses Remaining: 100
#   Expires In: 4m48s
```

### 📋 **3. List Active Capabilities**

```bash
# List all active capabilities
vault capability list --status active

# Output:
# Found 1 capabilities:
#
# ID                   Type            Resource                        Identity         Expires
# --------------------------------------------------------------------------------------------------------------
# cap_1234567890_ghi   read            secret:/db/primary             user             2025-01-13 10:05:00
```

**Capability management**:

```bash
# List all capabilities
vault capability list --all

# List by resource
vault capability list --resource "secret:/db/*"

# List by identity
vault capability list --identity "app:*"

# Show detailed capability info
vault capability show cap_1234567890_ghijkl

# Revoke capability
vault capability revoke cap_1234567890_ghijkl --reason "Testing completed"
```

---

## 💻 Integration Examples

### 🐹 **Go Application Example**

Create a production-ready Go application with capability management:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/ipc"
	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Create IPC client with retry logic
	client, err := ipc.NewClient(&ipc.Config{
		SocketPath:  os.Getenv("VAULT_AGENT_SOCKET_PATH"),
		Timeout:     30 * time.Second,
		RetryCount:  3,
		RetryDelay:  time.Second,
	})
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer client.Close()

	// Connect to agent with timeout
	ctx, connectCancel := context.WithTimeout(ctx, 10*time.Second)
	defer connectCancel()

	if err := client.ConnectContext(ctx); err != nil {
		log.Fatal("Failed to connect to agent:", err)
	}

	// Request capability with constraints
	request := &types.CapabilityRequest{
		Identity: "production-app",
		Resource: "secret:/db/primary",
		Actions:  []string{"read"},
		TTL:      300,
		Purpose:  "Production application database access",
		Constraints: []types.Constraint{
			{
				Type:  "ip",
				Value: "10.0.0.0/8",
			},
			{
				Type:  "env",
				Value: "production",
			},
		},
	}

	response, err := client.RequestCapabilityContext(ctx, request)
	if err != nil {
		log.Fatal("Failed to request capability:", err)
	}

	if response.Status != "granted" {
		log.Fatal("Capability denied:", response.Message)
	}

	log.Printf("✓ Capability %s granted", response.Capability.ID)

	// Validate capability
	result, err := client.ValidateCapabilityContext(ctx, response.Capability.ID, nil)
	if err != nil || !result.Valid {
		log.Fatal("Capability validation failed")
	}

	// Setup capability renewal
	renewalTicker := time.NewTicker(time.Duration(response.Capability.TTL/2) * time.Second)
	defer renewalTicker.Stop()

	// Main application loop
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
			// Revoke capability
			err := client.RevokeCapabilityContext(ctx, response.Capability.ID, "Application shutdown")
			if err != nil {
				log.Printf("Warning: Failed to revoke capability: %v", err)
			}
			return
		case <-renewalTicker.C:
			// Renew capability
			renewal, err := client.RenewCapabilityContext(ctx, response.Capability.ID, response.Capability.TTL)
			if err != nil {
				log.Printf("Failed to renew capability: %v", err)
				continue
			}
			log.Printf("✓ Capability renewed until %s", renewal.ExpiresAt.Format(time.RFC3339))
		default:
			// Simulate database work
			log.Printf("✓ Accessing database with capability %s", response.Capability.ID)
			time.Sleep(5 * time.Second)
		}
	}
}
```

Build and run:

```bash
go mod init production-app
go mod tidy
go build -o production-app
./production-app
```

### 🔧 **Shell Script Example**

Advanced shell script with error handling and logging:

```bash
#!/bin/bash

# production-deploy.sh - Deployment script with capability management

set -euo pipefail

# Configuration
readonly SCRIPT_NAME="$(basename "$0")"
readonly LOG_FILE="/var/log/${SCRIPT_NAME%.*}.log"
readonly VAULT_SOCKET="${VAULT_AGENT_SOCKET_PATH:-$HOME/.aether-vault/agent.sock}"
readonly DEPLOY_RESOURCE="operation:/deploy/production"
readonly DEPLOY_TTL=900

# Logging function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

# Error handling
error_exit() {
    log "ERROR: $*"
    exit 1
}

# Cleanup function
cleanup() {
    local cap_id="$1"
    if [[ -n "$cap_id" ]]; then
        log "🗑️  Revoking deployment capability..."
        if vault capability revoke "$cap_id" --reason "Deployment completed" 2>/dev/null; then
            log "✓ Capability revoked"
        else
            log "⚠️  Failed to revoke capability"
        fi
    fi
}

# Check prerequisites
check_prerequisites() {
    log "🔍 Checking prerequisites..."

    command -v vault >/dev/null 2>&1 || error_exit "vault CLI not found"
    command -v jq >/dev/null 2>&1 || error_exit "jq not found"

    [[ -S "$VAULT_SOCKET" ]] || error_exit "Vault agent socket not found at $VAULT_SOCKET"

    # Check agent status
    if ! vault agent status --format json | jq -r '.running' | grep -q true; then
        error_exit "Vault agent is not running"
    fi

    log "✓ Prerequisites check passed"
}

# Request deployment capability
request_capability() {
    log "📋 Requesting deployment capability..."

    local response
    response=$(vault capability request \
        --resource "$DEPLOY_RESOURCE" \
        --action execute \
        --ttl "$DEPLOY_TTL" \
        --purpose "Production deployment via $SCRIPT_NAME" \
        --constraint "env:production" \
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

# Execute deployment
execute_deployment() {
    local cap_id="$1"
    log "🚀 Executing deployment..."

    # Simulate deployment steps
    log "📦 Building application..."
    sleep 2

    log "🔄 Running tests..."
    sleep 1

    log "🚢 Deploying to production..."
    sleep 3

    log "🔍 Verifying deployment..."
    sleep 1

    log "✅ Deployment completed successfully"
}

# Main execution
main() {
    log "🚀 Starting production deployment with $SCRIPT_NAME"

    local cap_id=""

    # Setup cleanup trap
    trap 'cleanup "$cap_id"' EXIT INT TERM

    check_prerequisites
    cap_id=$(request_capability)
    validate_capability "$cap_id"
    execute_deployment "$cap_id"

    log "🎉 Production deployment completed successfully!"
}

# Run main function
main "$@"
```

Make executable and run:

```bash
chmod +x production-deploy.sh
./production-deploy.sh
```

### 🐍 **Python Example**

Production-ready Python integration with async support:

```python
#!/usr/bin/env python3

"""
production_app.py - Production application with capability management
"""

import asyncio
import json
import logging
import os
import signal
import subprocess
import sys
from contextlib import asynccontextmanager
from dataclasses import dataclass
from datetime import datetime, timedelta
from pathlib import Path
from typing import Optional

import aiofiles

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='[%(asctime)s] %(levelname)s: %(message)s',
    handlers=[
        logging.FileHandler('/var/log/production_app.log'),
        logging.StreamHandler(sys.stdout)
    ]
)
logger = logging.getLogger(__name__)

@dataclass
class VaultConfig:
    socket_path: str = os.getenv('VAULT_AGENT_SOCKET_PATH', f'{Path.home()}/.aether-vault/agent.sock')
    timeout: int = 30
    retry_count: int = 3
    retry_delay: float = 1.0

class VaultClient:
    """Async Vault client for capability management"""

    def __init__(self, config: VaultConfig):
        self.config = config
        self.capability_id: Optional[str] = None

    async def _run_command(self, args: list) -> dict:
        """Run vault command and return JSON result"""
        cmd = ['vault'] + args + ['--format', 'json']

        for attempt in range(self.config.retry_count):
            try:
                process = await asyncio.create_subprocess_exec(
                    *cmd,
                    stdout=asyncio.subprocess.PIPE,
                    stderr=asyncio.subprocess.PIPE
                )

                stdout, stderr = await asyncio.wait_for(
                    process.communicate(),
                    timeout=self.config.timeout
                )

                if process.returncode != 0:
                    raise subprocess.CalledProcessError(
                        process.returncode, cmd, stderr.decode()
                    )

                return json.loads(stdout.decode())

            except (asyncio.TimeoutError, subprocess.CalledProcessError, json.JSONDecodeError) as e:
                logger.warning(f"Command attempt {attempt + 1} failed: {e}")
                if attempt < self.config.retry_count - 1:
                    await asyncio.sleep(self.config.retry_delay)
                else:
                    raise

    async def check_agent_status(self) -> bool:
        """Check if agent is running"""
        try:
            result = await self._run_command(['agent', 'status'])
            return result.get('running', False)
        except Exception as e:
            logger.error(f"Failed to check agent status: {e}")
            return False

    async def request_capability(
        self,
        resource: str,
        actions: list,
        ttl: int,
        purpose: str,
        constraints: Optional[list] = None
    ) -> dict:
        """Request a capability"""
        args = [
            'capability', 'request',
            '--resource', resource,
            '--action', ','.join(actions),
            '--ttl', str(ttl),
            '--purpose', purpose
        ]

        if constraints:
            for constraint in constraints:
                args.extend(['--constraint', constraint])

        result = await self._run_command(args)

        if result.get('status') != 'granted':
            raise Exception(f"Capability denied: {result.get('message', 'Unknown error')}")

        self.capability_id = result['capability']['id']
        return result['capability']

    async def validate_capability(self, capability_id: str) -> bool:
        """Validate a capability"""
        result = await self._run_command(['capability', 'validate', capability_id])
        return result.get('valid', False)

    async def revoke_capability(self, capability_id: str, reason: str) -> bool:
        """Revoke a capability"""
        try:
            await self._run_command([
                'capability', 'revoke', capability_id,
                '--reason', reason
            ])
            return True
        except Exception as e:
            logger.warning(f"Failed to revoke capability: {e}")
            return False

    async def renew_capability(self, capability_id: str, ttl: int) -> dict:
        """Renew a capability"""
        result = await self._run_command([
            'capability', 'renew', capability_id,
            '--ttl', str(ttl)
        ])
        return result

class ProductionApp:
    """Production application with capability management"""

    def __init__(self, vault_client: VaultClient):
        self.vault_client = vault_client
        self.running = False
        self.capability: Optional[dict] = None
        self.renewal_task: Optional[asyncio.Task] = None

    async def setup_signal_handlers(self):
        """Setup graceful shutdown handlers"""
        def signal_handler(signum, frame):
            logger.info(f"Received signal {signum}")
            self.running = False

        signal.signal(signal.SIGINT, signal_handler)
        signal.signal(signal.SIGTERM, signal_handler)

    async def acquire_capability(self):
        """Acquire initial capability"""
        logger.info("📋 Requesting database capability...")

        self.capability = await self.vault_client.request_capability(
            resource="secret:/db/primary",
            actions=["read", "write"],
            ttl=300,
            purpose="Production application database access",
            constraints=[
                "env:production",
                "ip:10.0.0.0/8"
            ]
        )

        logger.info(f"✓ Capability granted: {self.capability['id']}")
        logger.info(f"✓ Expires at: {self.capability['expires_at']}")

        # Validate capability
        if not await self.vault_client.validate_capability(self.capability['id']):
            raise Exception("Capability validation failed")

        logger.info("✓ Capability validated")

    async def setup_renewal(self):
        """Setup capability renewal task"""
        if not self.capability:
            return

        renewal_interval = self.capability['ttl'] // 2

        async def renewal_loop():
            while self.running:
                await asyncio.sleep(renewal_interval)

                if not self.running or not self.capability:
                    break

                try:
                    logger.info("🔄 Renewing capability...")
                    renewed = await self.vault_client.renew_capability(
                        self.capability['id'],
                        self.capability['ttl']
                    )
                    self.capability = renewed['capability']
                    logger.info(f"✓ Capability renewed until {self.capability['expires_at']}")
                except Exception as e:
                    logger.error(f"Failed to renew capability: {e}")
                    break

        self.renewal_task = asyncio.create_task(renewal_loop())

    async def do_work(self):
        """Simulate application work"""
        while self.running:
            logger.info(f"✓ Processing data with capability {self.capability['id']}")
            await asyncio.sleep(5)

    async def cleanup(self):
        """Cleanup resources"""
        logger.info("🧹 Cleaning up...")

        self.running = False

        if self.renewal_task:
            self.renewal_task.cancel()
            try:
                await self.renewal_task
            except asyncio.CancelledError:
                pass

        if self.capability:
            await self.vault_client.revoke_capability(
                self.capability['id'],
                "Application shutdown"
            )
            logger.info("✓ Capability revoked")

    async def run(self):
        """Run the application"""
        await self.setup_signal_handlers()

        # Check agent status
        if not await self.vault_client.check_agent_status():
            raise Exception("Vault agent is not running")

        await self.acquire_capability()
        await self.setup_renewal()

        self.running = True
        logger.info("🚀 Application started")

        try:
            await self.do_work()
        finally:
            await self.cleanup()

async def main():
    """Main entry point"""
    config = VaultConfig()
    client = VaultClient(config)
    app = ProductionApp(client)

    try:
        await app.run()
    except Exception as e:
        logger.error(f"Application failed: {e}")
        sys.exit(1)

if __name__ == "__main__":
    asyncio.run(main())
```

Install dependencies and run:

```bash
pip install aiofiles
python3 production_app.py
```

---

## 🛠️ Configuration

### 🔍 **View Current Configuration**

```bash
# Show current configuration
vault config show

# Output:
# Agent Configuration:
#   Mode: standard
#   Log Level: info
#   Socket Path: /home/user/.aether-vault/agent.sock
#
# Capability Engine:
#   Enable: true
#   Default TTL: 300
#   Max TTL: 3600
#   Max Uses: 100
#
# Policy Engine:
#   Policy Directory: /home/user/.aether-vault/policies
#   Cache Size: 1000
#   Evaluation Timeout: 5s
```

### ⚙️ **Generate Default Configuration**

```bash
# Generate default configuration to file
vault agent config --generate --output ~/.aether-vault/agent.yaml

# View the generated file
cat ~/.aether-vault/agent.yaml
```

**Example configuration file**:

```yaml
# Aether Vault Agent Configuration
agent:
  mode: standard
  log_level: info
  socket_path: ~/.aether-vault/agent.sock
  pid_file: ~/.aether-vault/agent.pid

capability_engine:
  enable: true
  default_ttl: 300
  max_ttl: 3600
  max_uses: 100
  signing_key_file: ~/.aether-vault/signing.key

policy_engine:
  policy_directory: ~/.aether-vault/policies
  cache_size: 1000
  evaluation_timeout: 5s
  reload_interval: 30s

audit:
  enable: true
  log_file: ~/.aether-vault/audit.log
  rotation_size: 100MB
  retention_days: 30

security:
  enable_tls: false
  cert_file: ""
  key_file: ""

performance:
  max_connections: 100
  connection_timeout: 30s
  read_timeout: 30s
  write_timeout: 30s
```

### 🌍 **Environment Variables**

```bash
# Set common environment variables
export VAULT_LOG_LEVEL=debug
export VAULT_AGENT_SOCKET_PATH=/tmp/vault.sock
export VAULT_CONFIG_FILE=/custom/path/config.yaml
export VAULT_POLICY_DIR=/custom/path/policies

# Use with commands
vault capability status --verbose
```

**Available environment variables**:

- `VAULT_LOG_LEVEL`: Set logging level (debug, info, warn, error)
- `VAULT_AGENT_SOCKET_PATH`: Custom socket path
- `VAULT_CONFIG_FILE`: Custom configuration file
- `VAULT_POLICY_DIR`: Custom policy directory
- `VAULT_AUDIT_LOG_FILE`: Custom audit log file
- `VAULT_MODE`: Agent mode (standard, development, production)

---

## 🔍 Troubleshooting

### 🚨 **Common Issues**

#### **Agent Not Running**

```bash
# Check if agent is running
vault agent status

# Start agent if not running
vault agent start

# Check for socket file
ls -la ~/.aether-vault/agent.sock

# Check process
ps aux | grep vault-agent

# Check system logs
journalctl -u vault-agent -f
```

#### **Permission Denied**

```bash
# Check file permissions
ls -la ~/.aether-vault/

# Fix permissions
chmod 700 ~/.aether-vault/
chmod 600 ~/.aether-vault/config.yaml
chmod 600 ~/.aether-vault/agent.sock

# Check ownership
ls -la ~/.aether-vault/ | grep "$(whoami)"
```

#### **Capability Denied**

```bash
# Request with verbose output
vault capability request \
  --resource "secret:/test" \
  --action read \
  --verbose

# Check policies
ls ~/.aether-vault/policies/
cat ~/.aether-vault/policies/default.json

# Check audit logs for denial reason
grep "denied" ~/.aether-vault/audit.log | tail -5

# Test with different identity
vault capability request \
  --resource "secret:/test" \
  --action read \
  --identity "test-app" \
  --verbose
```

#### **Connection Issues**

```bash
# Test agent connection
vault capability status

# Check socket path
echo $VAULT_AGENT_SOCKET_PATH
ls -la "$VAULT_AGENT_SOCKET_PATH"

# Test with custom socket path
vault capability status --socket-path /tmp/vault.sock

# Test socket connectivity
nc -U "$VAULT_AGENT_SOCKET_PATH" <<< '{"method":"ping"}'
```

### 🐛 **Debug Mode**

```bash
# Enable debug logging
export VAULT_LOG_LEVEL=debug

# Start agent in debug mode
vault agent start --log-level debug --mode development

# Run commands with verbose output
vault capability request --resource "secret:/test" --action read --verbose

# Enable debug for specific components
vault agent start --debug capability-engine
vault agent start --debug policy-engine
vault agent start --debug ipc-server
```

### 🔄 **Reset Environment**

```bash
# Stop agent gracefully
vault agent stop

# Force stop if needed
pkill -f vault-agent

# Remove all data (WARNING: This deletes everything)
rm -rf ~/.aether-vault/

# Reinitialize from scratch
vault init
vault agent start

# Verify setup
vault agent status
vault capability status
```

### 📊 **Performance Issues**

```bash
# Check agent performance
vault agent status --verbose

# Monitor capability cache
watch -n 5 'vault capability status --format json | jq .cache_stats'

# Check policy evaluation performance
vault agent status --component policy-engine --verbose

# Monitor IPC connections
vault agent status --component ipc-server --verbose

# Enable performance profiling
vault agent start --profile cpu
vault agent start --profile memory
vault agent start --profile trace
```

---

## 🎯 Next Steps

### 📚 **Learn More**

1. **📖 Architecture Deep Dive**: [../ARCHITECTURE_DEEP_DIVE.md](../ARCHITECTURE_DEEP_DIVE.md)
2. **🔐 CBAC System**: [CBAC_OVERVIEW.md](CBAC_OVERVIEW.md)
3. **📋 Policy Management**: [CBAC_POLICIES.md](CBAC_POLICIES.md)
4. **🔗 Integration Guides**: [INTEGRATION_OVERVIEW.md](INTEGRATION_OVERVIEW.md)
5. **🛠️ Configuration**: [CONFIG_OVERVIEW.md](CONFIG_OVERVIEW.md)

### 🚀 **Advanced Topics**

1. **🎛️ Custom Policies**: Create sophisticated access control policies
2. **⚡ Constraints**: Use IP, time, and environment constraints
3. **📊 Audit Integration**: Set up SIEM integration
4. **🔄 High Availability**: Deploy multiple agents for redundancy
5. **⚙️ Performance Tuning**: Optimize for high-throughput scenarios
6. **🐳 Container Deployment**: Docker and Kubernetes deployment
7. **🔐 Security Hardening**: Production security configuration

### 🏭 **Production Deployment**

1. **🔒 Security Hardening**: Configure for production security
2. **📈 Monitoring Setup**: Set up metrics and alerting
3. **💾 Backup Strategy**: Implement proper backup procedures
4. **📋 Compliance**: Configure for regulatory compliance
5. **🆘 Disaster Recovery**: Plan for outage scenarios
6. **🔄 Updates & Patching**: Maintain and update the system

### 🤝 **Community & Support**

- **📖 Documentation**: [https://docs.aethervault.com](https://docs.aethervault.com)
- **🐛 GitHub Issues**: [Report Issues](https://github.com/skygenesisenterprise/aether-vault/issues)
- **💡 GitHub Discussions**: [Community Discussions](https://github.com/skygenesisenterprise/aether-vault/discussions)
- **📧 Email Support**: support@skygenesisenterprise.com
- **💬 Discord**: [Join our Discord](https://discord.gg/aethervault)

---

<div align="center">

### 🎉 **Congratulations! You're Ready to Use Aether Vault CLI!**

[⭐ Star This Repo](https://github.com/skygenesisenterprise/aether-vault) • [📚 Read Full Docs](../README.md) • [🐛 Report Issues](https://github.com/skygenesisenterprise/aether-vault/issues)

---

**🔐 Enterprise-Grade Secret Management with CBAC System!**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

_Building modern DevOps security infrastructure_

</div>
