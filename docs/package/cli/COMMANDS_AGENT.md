<div align="center">

# 🚀 Aether Vault CLI - Agent Commands Reference

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![Go](https://img.shields.io/badge/Go-1.25+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![Cobra](https://img.shields.io/badge/Cobra-1.8+-lightgrey?style=for-the-badge&logo=go)](https://github.com/spf13/cobra) [![Viper](https://img.shields.io/badge/Viper-1.16+-green?style=for-the-badge&logo=go)](https://github.com/spf13/viper) [![DevOps](https://img.shields.io/badge/DevOps-Ready-orange?style=for-the-badge&logo=devops)](https://www.devops.com/)

**🔧 Complete Agent Command Reference - Manage Your Security Daemon**

Comprehensive reference for all `vault agent` commands, including daemon management, configuration, monitoring, and troubleshooting. The Aether Vault Agent is the core security daemon providing CBAC, policy evaluation, and IPC communication.

[📋 Command Overview](#-command-overview) • [🚀 Agent Start](#-agent-start) • [🛑 Agent Stop](#-agent-stop) • [📊 Agent Status](#-agent-status) • [🔄 Agent Reload](#-agent-reload) • [⚙️ Agent Config](#️-agent-config) • [🎯 Agent Modes](#-agent-modes) • [📝 Configuration](#-configuration) • [🏥 Health Monitoring](#-health-monitoring) • [🔍 Troubleshooting](#-troubleshooting)

[![GitHub stars](https://img.shields.io/github/stars/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/stargazers) [![GitHub forks](https://img.shields.io/github/forks/skygenesisenterprise/aether-vault?style=social)](https://github.com/skygenesisenterprise/aether-vault/network)

</div>

---

## 📋 Command Overview

### 🎯 **Command Structure**

```bash
vault agent [subcommand] [flags]
```

### 📚 **Available Subcommands**

| Subcommand | Description            | Purpose                                     |
| ---------- | ---------------------- | ------------------------------------------- |
| `start`    | Start the agent daemon | Initialize and run the security daemon      |
| `stop`     | Stop the agent daemon  | Graceful shutdown of running agent          |
| `status`   | Show agent status      | Display comprehensive status information    |
| `reload`   | Reload configuration   | Refresh config and policies without restart |
| `config`   | Manage configuration   | Show, generate, or validate configuration   |

### 🔄 **Command Flow**

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   vault agent   │    │  Agent Daemon    │    │  Configuration  │
│   Commands      │◄──►│  (Running)       │◄──►│  Management     │
│  (CLI Interface)│    │  IPC Server      │    │  (YAML Files)   │
└─────────────────┘    │  CBAC Engine     │    └─────────────────┘
         │              │  Policy Engine   │              │
         │              │  Audit System    │              │
         ▼              └──────────────────┘              ▼
   User Interaction        Component Status        Settings Storage
   Process Control         Health Monitoring        Environment Overrides
```

---

## 🚀 Agent Start

### 🎯 **Syntax**

```bash
vault agent start [flags]
```

### 📋 **Optional Flags**

| Flag            | Type   | Default                    | Description                                 |
| --------------- | ------ | -------------------------- | ------------------------------------------- |
| `--config`      | string | -                          | Path to agent configuration file            |
| `--mode`        | string | standard                   | Agent mode: standard, hardened, development |
| `--log-level`   | string | info                       | Log level: debug, info, warn, error         |
| `--enable-cbac` | bool   | true                       | Enable Capability-Based Access Control      |
| `--policy-dir`  | string | -                          | Directory for policy files                  |
| `--socket-path` | string | ~/.aether-vault/agent.sock | Unix socket path                            |

### 💡 **Examples**

#### **🎯 Start Agent with Default Settings**

```bash
vault agent start
```

#### **🔒 Start Agent in Hardened Mode**

```bash
vault agent start --mode hardened --log-level warn
```

#### **⚙️ Start Agent with Custom Configuration**

```bash
vault agent start --config /etc/aether-vault/agent.yaml --enable-cbac
```

#### **🐛 Start Agent with Debug Logging**

```bash
vault agent start --log-level debug --mode development
```

#### **📁 Start Agent with Custom Policy Directory**

```bash
vault agent start --policy-dir /etc/aether-vault/policies --socket-path /tmp/vault.sock
```

### 🔄 **Startup Process**

When the agent starts, it performs these sequential steps:

1. **📁 Configuration Loading**
   - Load configuration from file or use defaults
   - Override with environment variables
   - Validate configuration syntax

2. **🔧 Component Initialization**
   - Initialize Ed25519 signing keys
   - Setup capability engine
   - Initialize policy engine
   - Start audit system

3. **📋 Policy Loading**
   - Scan policy directory
   - Load and validate policies
   - Build policy cache
   - Set up policy reloading

4. **🌐 IPC Server Start**
   - Create Unix socket
   - Start listening for connections
   - Setup authentication
   - Configure connection limits

5. **🏥 Health Monitoring**
   - Start health check routines
   - Initialize metrics collection
   - Setup monitoring endpoints

6. **🧹 Cleanup Routine**
   - Start background cleanup
   - Setup expired capability removal
   - Initialize log rotation

### 📤 **Output Examples**

#### **✅ Successful Start**

```
Starting Aether Vault Agent...
Configuration loaded from /home/user/.aether-vault/agent.yaml
Capability engine initialized with Ed25519 signing
Policy engine loaded 5 policies from /home/user/.aether-vault/policies
Audit system started with file logging
IPC server listening on /home/user/.aether-vault/agent.sock
Agent started successfully (PID: 12345)
Health checks enabled (interval: 30s)
Metrics server started on port 9090
```

#### **⚠️ Start with Issues**

```
Starting Aether Vault Agent...
Warning: Failed to load policy directory /custom/policies: No such file or directory
Configuration loaded from /home/user/.aether-vault/agent.yaml
Capability engine initialized with Ed25519 signing
Policy engine loaded 0 policies
Audit system started with file logging
IPC server listening on /home/user/.aether-vault/agent.sock
Agent started successfully (PID: 12345)
Warning: No policies loaded - all requests will be denied
```

#### **❌ Failed Start**

```
Starting Aether Vault Agent...
Error: Failed to create socket: permission denied
Please check socket directory permissions and try again
```

---

## 🛑 Agent Stop

### 🎯 **Syntax**

```bash
vault agent stop [flags]
```

### 📋 **Optional Flags**

| Flag        | Type | Default | Description                 |
| ----------- | ---- | ------- | --------------------------- |
| `--timeout` | int  | 30      | Shutdown timeout in seconds |
| `--force`   | bool | false   | Force immediate shutdown    |

### 💡 **Examples**

#### **🛑 Stop Agent Gracefully**

```bash
vault agent stop
```

#### **⏱️ Stop with Custom Timeout**

```bash
vault agent stop --timeout 60
```

#### **💪 Force Immediate Shutdown**

```bash
vault agent stop --force
```

### 🔄 **Shutdown Process**

The agent performs a graceful shutdown sequence:

1. **🚫 Stop Accepting Connections**
   - IPC server stops accepting new connections
   - Signal clients about shutdown
   - Wait for active connections to finish

2. **⏳ Complete In-Flight Operations**
   - Finish processing ongoing requests
   - Complete capability validations
   - Finalize audit event logging

3. **📝 Flush Audit Logs**
   - Ensure all audit events are written
   - Close log files properly
   - Flush buffers to disk

4. **🧹 Cleanup Resources**
   - Release temporary files
   - Close database connections
   - Clean up socket files

5. **🔌 Close Connections**
   - Close all active client connections
   - Shutdown IPC server
   - Release network resources

### 📤 **Output Examples**

#### **✅ Successful Stop**

```
Stopping Aether Vault Agent...
IPC server stopped accepting new connections
Waiting for in-flight operations to complete... (3 active)
Flushed 15 audit events to disk
Cleaned up temporary resources
Agent stopped successfully
```

#### **⚠️ Stop with Timeout**

```
Stopping Aether Vault Agent...
IPC server stopped accepting new connections
Waiting for in-flight operations to complete... (5 active)
Timeout reached, forcing shutdown...
Flushed 12 audit events to disk
Agent stopped with warnings
```

#### **❌ Agent Not Running**

```
Error: Failed to connect to agent: connection refused
Agent is not running or not accessible
```

---

## 📊 Agent Status

### 🎯 **Syntax**

```bash
vault agent status [flags]
```

### 📋 **Optional Flags**

| Flag          | Type   | Default | Description                      |
| ------------- | ------ | ------- | -------------------------------- |
| `--verbose`   | bool   | false   | Show detailed status information |
| `--format`    | string | table   | Output format: table, json, yaml |
| `--component` | string | -       | Show specific component status   |

### 💡 **Examples**

#### **📊 Basic Status**

```bash
vault agent status
```

#### **📋 Verbose Status**

```bash
vault agent status --verbose
```

#### **📄 JSON Status**

```bash
vault agent status --format json
```

#### **🔍 Component-Specific Status**

```bash
vault agent status --component capability-engine
vault agent status --component policy-engine
vault agent status --component ipc-server
```

### 📤 **Response Format**

#### **📋 Table Format**

```
Aether Vault Agent Status:
  Running: true
  PID: 12345
  Uptime: 2h45m30s
  Version: 1.0.0
  Mode: standard

IPC Server:
  Socket Path: /home/user/.aether-vault/agent.sock
  Active Connections: 3
  Max Connections: 100
  Server Uptime: 2h45m30s
  Total Requests: 1,247

Capability Engine:
  Status: Healthy
  Total Capabilities: 127
  Active Capabilities: 45
  Expired Capabilities: 80
  Revoked Capabilities: 2
  Cache Size: 45/10000
  Last Cleanup: 5m ago
  Signing Algorithm: Ed25519

Policy Engine:
  Status: Healthy
  Loaded Policies: 5
  Cache Hits: 892
  Cache Misses: 45
  Cache Hit Rate: 95.2%
  Last Policy Reload: 1h ago
  Default Decision: deny

Audit System:
  Status: Healthy
  Total Events: 1,247
  Buffer Size: 234/1000
  Last Flush: 2m ago
  Log File: /home/user/.aether-vault/audit.log
  Log Size: 15.2MB
  Rotation Enabled: true

System Resources:
  Memory Usage: 45.2MB
  CPU Usage: 2.1%
  File Descriptors: 12/1024
  Goroutines: 8
  Uptime: 2h45m30s
```

#### **📄 JSON Format**

```json
{
  "running": true,
  "pid": 12345,
  "uptime": "2h45m30s",
  "version": "1.0.0",
  "mode": "standard",
  "ipc_server": {
    "socket_path": "/home/user/.aether-vault/agent.sock",
    "active_connections": 3,
    "max_connections": 100,
    "uptime": "2h45m30s",
    "total_requests": 1247
  },
  "capability_engine": {
    "status": "healthy",
    "total_capabilities": 127,
    "active_capabilities": 45,
    "expired_capabilities": 80,
    "revoked_capabilities": 2,
    "cache_size": "45/10000",
    "last_cleanup": "5m ago",
    "signing_algorithm": "Ed25519"
  },
  "policy_engine": {
    "status": "healthy",
    "loaded_policies": 5,
    "cache_hits": 892,
    "cache_misses": 45,
    "cache_hit_rate": "95.2%",
    "last_policy_reload": "1h ago",
    "default_decision": "deny"
  },
  "audit_system": {
    "status": "healthy",
    "total_events": 1247,
    "buffer_size": "234/1000",
    "last_flush": "2m ago",
    "log_file": "/home/user/.aether-vault/audit.log",
    "log_size": "15.2MB",
    "rotation_enabled": true
  },
  "system_resources": {
    "memory_usage": "45.2MB",
    "cpu_usage": "2.1%",
    "file_descriptors": "12/1024",
    "goroutines": 8,
    "uptime": "2h45m30s"
  }
}
```

#### **📋 Component-Specific Status**

```
Capability Engine Status:
  Health: Healthy
  Total Capabilities: 127
  Active Capabilities: 45
  Expired Capabilities: 80
  Revoked Capabilities: 2
  Cache Size: 45/10000
  Last Cleanup: 5m ago
  Signing Algorithm: Ed25519
  Key ID: key_1234567890_abcdef
  Performance: 45.2 requests/second
```

---

## 🔄 Agent Reload

### 🎯 **Syntax**

```bash
vault agent reload [flags]
```

### 📋 **Optional Flags**

| Flag         | Type   | Default | Description                     |
| ------------ | ------ | ------- | ------------------------------- |
| `--config`   | string | -       | Reload specific configuration   |
| `--policies` | bool   | true    | Reload policies                 |
| `--force`    | bool   | false   | Force reload without validation |

### 💡 **Examples**

#### **🔄 Reload Agent**

```bash
vault agent reload
```

#### **⚙️ Reload Configuration Only**

```bash
vault agent reload --config
```

#### **📋 Reload Policies Only**

```bash
vault agent reload --policies
```

#### **💪 Force Reload**

```bash
vault agent reload --force
```

### 🔄 **Reload Process**

1. **📁 Configuration Reload**
   - Reload configuration from file
   - Validate new configuration
   - Apply changes to running components

2. **📋 Policy Cache Refresh**
   - Refresh cached policies from disk
   - Validate policy syntax
   - Update policy engine cache

3. **🔧 Policy Engine Restart**
   - Restart policy engine with new policies
   - Rebuild evaluation cache
   - Maintain existing connections

4. **🌐 Maintain Connections**
   - Keep existing client connections active
   - Apply new policies to new requests
   - Graceful transition to new configuration

5. **📝 Log Reload Events**
   - Log all reload events for audit
   - Track configuration changes
   - Monitor reload success/failure

### 📤 **Output Examples**

#### **✅ Successful Reload**

```
Reloading Aether Vault Agent...
Configuration reloaded from /home/user/.aether-vault/agent.yaml
Policy cache refreshed
Policy engine restarted with 5 policies
Existing connections maintained
Agent reloaded successfully
```

#### **⚠️ Reload with Issues**

```
Reloading Aether Vault Agent...
Warning: Failed to reload configuration: file not found, using current config
Policy cache refreshed
Policy engine restarted with 5 policies
Existing connections maintained
Agent reloaded successfully with warnings
```

#### **❌ Failed Reload**

```
Reloading Aether Vault Agent...
Error: Failed to validate new policies: syntax error in policy file
Reload aborted - agent continues with current configuration
```

---

## ⚙️ Agent Config

### 🎯 **Syntax**

```bash
vault agent config [flags]
```

### 📋 **Optional Flags**

| Flag         | Type   | Default | Description                    |
| ------------ | ------ | ------- | ------------------------------ |
| `--output`   | string | -       | Output configuration to file   |
| `--validate` | bool   | false   | Validate configuration only    |
| `--generate` | bool   | false   | Generate default configuration |
| `--config`   | string | -       | Path to configuration file     |

### 💡 **Examples**

#### **📋 Show Current Configuration**

```bash
vault agent config
```

#### **🔧 Generate Default Configuration**

```bash
vault agent config --generate
```

#### **📁 Generate Configuration to File**

```bash
vault agent config --generate --output /etc/aether-vault/agent.yaml
```

#### **✅ Validate Configuration**

```bash
vault agent config --validate --config /etc/aether-vault/agent.yaml
```

### 📤 **Response Format**

#### **📋 Current Configuration (Table Format)**

```
Agent Configuration:
  Mode: standard
  Log Level: info
  Socket Path: /home/user/.aether-vault/agent.sock
  PID File: /home/user/.aether-vault/agent.pid

Capability Engine:
  Enable: true
  Default TTL: 300
  Max TTL: 3600
  Max Uses: 100
  Signing Algorithm: ed25519
  Enable Usage Tracking: true
  Cleanup Interval: 60

Policy Engine:
  Enable: true
  Directory: /home/user/.aether-vault/policies
  Cache Enable: true
  Cache TTL: 300
  Cache Size: 1000
  Enable Reloading: true
  Reload Interval: 60
  Default Decision: deny

Audit System:
  Enable: true
  Log File: /home/user/.aether-vault/audit.log
  Buffer Size: 1000
  Flush Interval: 60
  Enable Rotation: true
  Max File Size: 104857600
  Max Backup Files: 10

IPC Server:
  Timeout: 30
  Max Connections: 100
  Enable Auth: true
  Enable TLS: false

Health Monitoring:
  Enable Checks: true
  Check Interval: 30
  Enable Metrics: true
  Metrics Port: 9090
```

#### **📄 Default Configuration (YAML Format)**

```yaml
# Aether Vault Agent Configuration
version: "1.0"

# Agent Settings
mode: "standard"
log_level: "info"
socket_path: "/home/user/.aether-vault/agent.sock"
pid_file: "/home/user/.aether-vault/agent.pid"

# Capability Engine
capability_engine:
  enable: true
  default_ttl: 300
  max_ttl: 3600
  max_uses: 100
  signing_algorithm: "ed25519"
  enable_usage_tracking: true
  cleanup_interval: 60

  # Key Management
  keys:
    private_key_file: "/home/user/.aether-vault/private.key"
    public_key_file: "/home/user/.aether-vault/public.key"
    auto_generate: true

# Policy Engine
policy_engine:
  enable: true
  directory: "/home/user/.aether-vault/policies"
  cache:
    enable: true
    ttl: 300
    size: 1000
  enable_reloading: true
  reload_interval: 60
  default_decision: "deny"
  enable_validation: true

# Audit System
audit:
  enable: true
  log_file: "/home/user/.aether-vault/audit.log"
  log_level: "info"
  enable_buffer: true
  buffer_size: 1000
  flush_interval: 60
  enable_rotation: true
  max_file_size: 104857600
  max_backup_files: 10
  enable_compression: false

# IPC Server
ipc:
  timeout: 30
  max_connections: 100
  enable_auth: true
  enable_tls: false
  tls_cert_file: "/home/user/.aether-vault/server.crt"
  tls_key_file: "/home/user/.aether-vault/server.key"

# Health Monitoring
health:
  enable_checks: true
  check_interval: 30
  enable_metrics: true
  metrics_port: 9090
```

#### **✅ Validation Results**

```
Validating configuration: /etc/aether-vault/agent.yaml
✓ Configuration file syntax is valid
✓ All required fields are present
✓ Capability engine configuration is valid
✓ Policy engine configuration is valid
✓ Audit system configuration is valid
✓ IPC server configuration is valid
✓ Health monitoring configuration is valid

Configuration is valid
```

---

## 🎯 Agent Modes

### 🎯 **Standard Mode**

Default mode for normal operation with balanced security and performance.

```bash
vault agent start --mode standard
```

**Characteristics:**

- ✅ Full CBAC functionality
- ✅ Standard security policies
- ✅ Normal performance optimization
- ✅ Comprehensive audit logging
- ✅ Default capability TTLs (300s)
- ✅ Standard connection limits (100)

### 🔒 **Hardened Mode**

Enhanced security mode for high-security environments.

```bash
vault agent start --mode hardened
```

**Characteristics:**

- ✅ Stricter security policies
- ✅ Reduced capability TTLs (60s)
- ✅ Enhanced audit logging
- ✅ Additional validation checks
- ✅ Limited connection rates (50)
- ✅ Stricter constraint enforcement
- ✅ Extended audit retention

### 🐛 **Development Mode**

Relaxed security mode for development and testing.

```bash
vault agent start --mode development
```

**Characteristics:**

- ✅ Longer capability TTLs (3600s)
- ✅ Debug-level logging
- ✅ Relaxed security policies
- ✅ Additional debugging information
- ✅ Performance monitoring enabled
- ✅ Higher connection limits (200)
- ✅ Verbose error messages

### 📊 **Mode Comparison**

| Feature         | Standard | Hardened | Development |
| --------------- | -------- | -------- | ----------- |
| Default TTL     | 300s     | 60s      | 3600s       |
| Max TTL         | 3600s    | 300s     | 7200s       |
| Max Connections | 100      | 50       | 200         |
| Log Level       | info     | warn     | debug       |
| Audit Detail    | standard | enhanced | verbose     |
| Validation      | standard | strict   | relaxed     |
| Performance     | balanced | secure   | optimized   |

---

## 📝 Configuration

### 📁 **Configuration File Structure**

The agent uses YAML configuration files with comprehensive customization options:

```yaml
# Aether Vault Agent Configuration
version: "1.0"

# Basic Settings
mode: "standard"
log_level: "info"
socket_path: "/home/user/.aether-vault/agent.sock"
pid_file: "/home/user/.aether-vault/agent.pid"

# Capability Engine Configuration
capability_engine:
  enable: true
  default_ttl: 300
  max_ttl: 3600
  max_uses: 100
  issuer: "aether-vault-agent"
  signing_algorithm: "ed25519"
  enable_usage_tracking: true
  cleanup_interval: 60

  # Key Management
  keys:
    private_key_file: "/home/user/.aether-vault/private.key"
    public_key_file: "/home/user/.aether-vault/public.key"
    auto_generate: true

# Policy Engine Configuration
policy_engine:
  enable: true
  directory: "/home/user/.aether-vault/policies"
  default_decision: "deny"

  # Cache Configuration
  cache:
    enable: true
    ttl: 300
    size: 1000

  # Reloading
  enable_reloading: true
  reload_interval: 60

  # Validation
  enable_validation: true

# Audit System Configuration
audit:
  enable: true
  log_file: "/home/user/.aether-vault/audit.log"
  log_level: "info"

  # Buffer Configuration
  enable_buffer: true
  buffer_size: 1000
  flush_interval: 60

  # Log Rotation
  enable_rotation: true
  max_file_size: 104857600 # 100MB
  max_backup_files: 10
  enable_compression: false

  # Security
  enable_signature: false
  signature_key_file: "/home/user/.aether-vault/audit-signature.key"

  # SIEM Integration
  enable_siem: false
  siem_endpoint: "https://siem.company.com/events"
  siem_format: "json"

# IPC Server Configuration
ipc:
  timeout: 30
  max_connections: 100
  enable_auth: true

  # TLS Configuration
  enable_tls: false
  tls_cert_file: "/home/user/.aether-vault/server.crt"
  tls_key_file: "/home/user/.aether-vault/server.key"

  # Authentication
  auth_timeout: 30
  conn_timeout: 60

# Storage Configuration
storage:
  enable_persistence: true
  storage_file: "/home/user/.aether-vault/capabilities.json"
  enable_compression: false
  enable_encryption: false
  encryption_key_file: "/home/user/.aether-vault/storage.key"

# Health Monitoring
health:
  enable_checks: true
  check_interval: 30
  enable_metrics: true
  metrics_port: 9090
```

### 🌍 **Environment Variables**

Configuration can be overridden with environment variables:

```bash
export VAULT_AGENT_MODE="hardened"
export VAULT_AGENT_LOG_LEVEL="debug"
export VAULT_AGENT_SOCKET_PATH="/tmp/vault.sock"
export VAULT_AGENT_POLICY_DIR="/etc/aether-vault/policies"
export VAULT_AGENT_AUDIT_FILE="/var/log/vault/audit.log"
export VAULT_AGENT_CONFIG_FILE="/etc/aether-vault/agent.yaml"
```

**Available Environment Variables:**

| Variable                  | Description        | Default                    |
| ------------------------- | ------------------ | -------------------------- |
| `VAULT_AGENT_MODE`        | Agent mode         | standard                   |
| `VAULT_AGENT_LOG_LEVEL`   | Log level          | info                       |
| `VAULT_AGENT_SOCKET_PATH` | Unix socket path   | ~/.aether-vault/agent.sock |
| `VAULT_AGENT_POLICY_DIR`  | Policy directory   | ~/.aether-vault/policies   |
| `VAULT_AGENT_AUDIT_FILE`  | Audit log file     | ~/.aether-vault/audit.log  |
| `VAULT_AGENT_CONFIG_FILE` | Configuration file | ~/.aether-vault/agent.yaml |

---

## 🏥 Health Monitoring

### 🔍 **Health Checks**

The agent performs regular health checks on all components:

```bash
# Check agent health
vault agent status --verbose
```

**Health Check Components:**

1. **🌐 IPC Server**: Socket accessibility and connection handling
2. **🔐 Capability Engine**: Token generation and validation
3. **📋 Policy Engine**: Policy loading and evaluation
4. **📝 Audit System**: Log writing and rotation
5. **💾 Storage**: Persistence and cleanup
6. **🖥️ System Resources**: Memory, CPU, file descriptors

### 📊 **Metrics**

The agent can expose metrics for monitoring:

```yaml
health:
  enable_metrics: true
  metrics_port: 9090
```

**Available Metrics:**

- `vault_agent_capabilities_total`: Total capabilities created
- `vault_agent_capabilities_active`: Currently active capabilities
- `vault_agent_policy_evaluations_total`: Policy evaluations performed
- `vault_agent_audit_events_total`: Audit events logged
- `vault_agent_ipc_connections_active`: Active IPC connections
- `vault_agent_memory_usage_bytes`: Memory usage in bytes
- `vault_agent_cpu_usage_percent`: CPU usage percentage
- `vault_agent_uptime_seconds`: Agent uptime in seconds

### 📈 **Monitoring Examples**

#### **📊 Check Metrics**

```bash
# Get metrics in Prometheus format
curl http://localhost:9090/metrics

# Monitor specific metrics
curl -s http://localhost:9090/metrics | grep vault_agent_capabilities
```

#### **🏥 Health Check Script**

```bash
#!/bin/bash
# health-check.sh

# Check if agent is running
if ! vault agent status > /dev/null 2>&1; then
    echo "❌ Agent is not running"
    exit 1
fi

# Check component health
status=$(vault agent status --format json)
capability_status=$(echo "$status" | jq -r '.capability_engine.status')
policy_status=$(echo "$status" | jq -r '.policy_engine.status')

if [[ "$capability_status" != "healthy" || "$policy_status" != "healthy" ]]; then
    echo "❌ Agent components are not healthy"
    echo "Capability Engine: $capability_status"
    echo "Policy Engine: $policy_status"
    exit 1
fi

echo "✅ Agent is healthy"
```

---

## 🔍 Troubleshooting

### 🚨 **Common Issues**

#### **🚫 Agent Won't Start**

**Symptoms:**

```
Error: Failed to start agent: address already in use
```

**Solutions:**

```bash
# Check if agent is already running
vault agent status

# Kill existing agent
pkill -f "vault agent"

# Remove stale socket file
rm -f ~/.aether-vault/agent.sock

# Start agent
vault agent start
```

#### **🔐 Permission Issues**

**Symptoms:**

```
Error: Failed to create socket: permission denied
```

**Solutions:**

```bash
# Check socket directory permissions
ls -la ~/.aether-vault/

# Fix permissions
chmod 755 ~/.aether-vault/
chmod 644 ~/.aether-vault/agent.sock

# Start with different socket path
vault agent start --socket-path /tmp/vault.sock
```

#### **⚙️ Configuration Errors**

**Symptoms:**

```
Error: Failed to load configuration: invalid YAML
```

**Solutions:**

```bash
# Validate configuration
vault agent config --validate --config ~/.aether-vault/agent.yaml

# Generate new default config
vault agent config --generate --output ~/.aether-vault/agent.yaml

# Check YAML syntax
python -c "import yaml; yaml.safe_load(open('~/.aether-vault/agent.yaml'))"
```

#### **📋 Policy Loading Issues**

**Symptoms:**

```
Warning: Failed to load policy directory: No such file or directory
```

**Solutions:**

```bash
# Create policy directory
mkdir -p ~/.aether-vault/policies

# Add example policy
cat > ~/.aether-vault/policies/default.json << EOF
{
  "id": "default",
  "name": "Default Policy",
  "version": "1.0",
  "status": "active",
  "rules": [
    {
      "id": "allow-local",
      "effect": "allow",
      "resources": ["secret:*"],
      "actions": ["*"],
      "identities": ["*"],
      "priority": 100
    }
  ]
}
EOF

# Reload agent
vault agent reload
```

### 🐛 **Debug Mode**

For detailed troubleshooting, start the agent in debug mode:

```bash
vault agent start --log-level debug --mode development
```

This enables:

- 📝 Detailed logging of all operations
- 📚 Stack traces for errors
- 📊 Performance metrics
- ✅ Additional validation checks

### 📋 **Log Analysis**

#### **📝 Agent Logs**

```bash
# View agent logs
tail -f ~/.aether-vault/agent.log

# Search for errors
grep -i error ~/.aether-vault/agent.log

# View recent capability requests
grep "capability_request" ~/.aether-vault/audit.log | tail -10
```

#### **🖥️ System Logs**

```bash
# Check system logs for agent issues
journalctl -u vault-agent -f

# Check for socket issues
ss -xl | grep vault
```

### 🔧 **Advanced Troubleshooting**

#### **🔍 Component Diagnostics**

```bash
# Check specific component status
vault agent status --component capability-engine --verbose
vault agent status --component policy-engine --verbose
vault agent status --component ipc-server --verbose
```

#### **📊 Performance Analysis**

```bash
# Monitor resource usage
watch -n 5 'vault agent status --verbose | grep -E "(Memory|CPU|Goroutines)"'

# Check capability cache performance
vault agent status --verbose | grep -A 10 "Capability Engine"
```

#### **🔄 Recovery Procedures**

```bash
# Full agent reset
vault agent stop
rm -rf ~/.aether-vault/
vault init
vault agent start

# Configuration reset
vault agent stop
rm ~/.aether-vault/agent.yaml
vault agent config --generate --output ~/.aether-vault/agent.yaml
vault agent start
```

---

## 🏆 Best Practices

### 🎯 **1. Production Deployment**

```bash
# Use hardened mode
vault agent start --mode hardened --log-level warn

# Use systemd service
sudo systemctl enable vault-agent
sudo systemctl start vault-agent

# Regular health checks
vault agent status --verbose
```

### ⚙️ **2. Configuration Management**

```bash
# Use configuration management
vault agent config --generate --output /etc/aether-vault/agent.yaml

# Validate before deployment
vault agent config --validate --config /etc/aether-vault/agent.yaml

# Use environment variables for secrets
export VAULT_AGENT_SIGNING_KEY_FILE="/etc/vault/secrets/private.key"
```

### 📊 **3. Monitoring**

```bash
# Regular health checks
vault agent status --verbose

# Monitor logs
tail -f ~/.aether-vault/audit.log | grep "ERROR\|WARN"

# Check capability usage
vault capability list --status "active" | wc -l

# Monitor metrics
curl http://localhost:9090/metrics
```

### 🔐 **4. Security**

```bash
# Use appropriate file permissions
chmod 600 ~/.aether-vault/agent.yaml
chmod 700 ~/.aether-vault/

# Regular cleanup
vault agent status --verbose | grep "Last Cleanup"

# Use hardened mode in production
vault agent start --mode hardened
```

### 🔄 **5. Maintenance**

```bash
# Regular configuration reload
vault agent reload

# Policy updates
vault agent reload --policies

# Log rotation monitoring
ls -la ~/.aether-vault/audit.log*

# Backup configuration
cp ~/.aether-vault/agent.yaml ~/.aether-vault/agent.yaml.backup
```

---

## 🚀 Integration Examples

### 🖥️ **Systemd Service**

```ini
[Unit]
Description=Aether Vault Agent
After=network.target

[Service]
Type=forking
User=vault
Group=vault
ExecStart=/usr/local/bin/vault agent start --config /etc/aether-vault/agent.yaml
ExecReload=/usr/local/bin/vault agent reload
ExecStop=/usr/local/bin/vault agent stop
PIDFile=/var/run/vault/agent.pid
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### 🐳 **Docker Container**

```dockerfile
FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY vault /usr/local/bin/vault
RUN chmod +x /usr/local/bin/vault

RUN adduser -D -s /bin/sh vault

USER vault
EXPOSE 9090

CMD ["vault", "agent", "start", "--mode", "standard"]
```

### ☸️ **Kubernetes Deployment**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vault-agent
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vault-agent
  template:
    metadata:
      labels:
        app: vault-agent
    spec:
      containers:
        - name: vault-agent
          image: aether-vault/agent:latest
          args:
            - "start"
            - "--mode"
            - "standard"
            - "--log-level"
            - "info"
          volumeMounts:
            - name: config
              mountPath: /etc/aether-vault
            - name: data
              mountPath: /home/vault/.aether-vault
          ports:
            - containerPort: 9090
              name: metrics
      volumes:
        - name: config
          configMap:
            name: vault-agent-config
        - name: data
          emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: vault-agent-metrics
spec:
  selector:
    app: vault-agent
  ports:
    - port: 9090
      targetPort: 9090
      name: metrics
```

---

<div align="center">

### 🎉 **Master Agent Commands - Complete Control Over Your Security Daemon!**

[🚀 Quick Start](QUICK_START.md) • [📋 Capability Commands](COMMANDS_CAPABILITY.md) • [⚙️ Configuration](CONFIG_OVERVIEW.md) • [🔐 CBAC Overview](CBAC_OVERVIEW.md)

---

**🔧 Enterprise-Grade Agent Management with Comprehensive Monitoring!**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

_Building modern DevOps security infrastructure_

</div>
