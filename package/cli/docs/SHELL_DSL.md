# Aether Vault Shell DSL

The `vault shell` command provides an interactive shell and scripting environment with a Domain Specific Language (DSL) for automating Vault operations.

## Features

- **Interactive REPL Mode**: Type commands directly in an interactive shell
- **Script Execution**: Run scripts with `.vault` extension
- **Variable Management**: Set and use variables throughout your session
- **Control Flow**: Conditional statements and loops
- **Built-in Commands**: DSL commands for common automation tasks
- **Vault Integration**: Full access to all Vault commands

## Usage

### Interactive Shell

```bash
# Start interactive shell
vault shell

# Start with specific context
vault shell --interactive
```

### Script Execution

```bash
# Execute a script file
vault shell script.vault

# Execute with verbose output
vault shell --verbose script.vault
```

### Single Command

```bash
# Execute single command
vault shell -c "set env production; read secret/myapp/${env}"
```

## DSL Commands

### Variable Management

```bash
# Set variables
set env production
set app myapp
set version 1.2.3

# Use variables with ${} syntax
echo "Deploying ${app} v${version} to ${env}"

# Get variable value
get env
```

### Control Flow

#### If Statements

```bash
# Conditional execution
if "${env}" == "production" then
    echo "Running in production mode"
    read secret/production/db
fi
```

#### For Loops

```bash
# Loop through values
for env in "staging,production,development" do
    echo "Processing ${env}"
    read secret/${env}/config
done
```

#### Wait Command

```bash
# Wait for condition (with timeout)
wait "secret/ready == true" 30s
```

### Built-in Commands

```bash
# Print messages
echo "Hello World"

# Pause execution
sleep 5s
sleep 1m
sleep 30s

# Show help
help

# Clear screen
clear

# Show command history
history

# Show all variables
vars

# Exit shell
exit
quit
```

## Vault Command Integration

All standard Vault commands are available directly:

```bash
# Read secrets
read secret/myapp/database

# Write secrets
write secret/myapp/database host="db.example.com" port=5432

# List secrets
list secret/

# Delete secrets
delete secret/myapp/old_config

# Check status
status

# Login
login --method token
```

## Script Examples

### Basic Deployment Script

```vault
# deploy.vault - Basic deployment automation

# Configuration
set app myapp
set env production
set version 1.2.3

echo "Starting deployment of ${app} v${version} to ${env}"

# Check current status
status

# Read database configuration
read secret/${env}/database

# Update application config
write secret/${app}/${env}/config version="${version}" deployed_at=$(date)

# Verify deployment
if "${version}" == "1.2.3" then
    echo "Deployment successful"
else
    echo "Deployment failed"
fi
```

### Multi-Environment Script

```vault
# multi-env.vault - Process multiple environments

set app myapp
set version 2.0.0

for env in "staging,production" do
    echo "Processing ${env} environment"

    # Read current config
    read secret/${app}/${env}/config

    # Update version
    write secret/${app}/${env}/config version="${version}"

    # Wait for propagation
    sleep 10s

    # Verify update
    read secret/${app}/${env}/config
done

echo "All environments updated successfully"
```

### Health Check Script

```vault
# health-check.vault - System health monitoring

set timeout 30s
set retries 3

echo "Starting health check..."

# Check Vault status
status

# Check critical secrets
for secret in "secret/database,secret/redis,secret/api" do
    echo "Checking ${secret}"

    # Wait for secret to be available
    wait "${secret} == accessible" ${timeout}

    # Read secret
    read ${secret}

    sleep 5s
done

echo "Health check completed"
```

## Variable Substitution

Variables can be used anywhere with the `${variable}` syntax:

```bash
set base_path secret/myapp
set env production

# Variables in paths
read ${base_path}/${env}/database

# Variables in values
write ${base_path}/config env="${env}" status="active"

# Variables in conditions
if "${env}" == "production" then
    echo "Production mode"
fi
```

## Error Handling

Scripts will stop on the first error by default. Use conditional statements for error handling:

```vault
# Safe secret reading
if "secret/myapp/config" == "exists" then
    read secret/myapp/config
else
    echo "Config not found, using defaults"
    write secret/myapp/config mode="default"
fi
```

## History and Navigation

In interactive mode:

- `↑/↓` arrows: Navigate command history
- `history`: Show all commands
- `clear`: Clear screen
- `vars`: Show all variables
- `help`: Show available commands

## File Extensions

Use `.vault` extension for script files to enable syntax highlighting in editors that support it.

## Best Practices

1. **Use descriptive variable names**
2. **Add comments with `#`**
3. **Handle errors with conditions**
4. **Use appropriate timeouts in wait commands**
5. **Test scripts in non-production environments first**

## Integration with CI/CD

The shell DSL can be integrated into CI/CD pipelines:

```bash
# In CI/CD script
vault shell deploy.vault --verbose
vault shell health-check.vault
```
