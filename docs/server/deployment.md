<div align="center">

# 🚀 Aether Vault Deployment Guide

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-vault/blob/main/LICENSE) [![Docker](https://img.shields.io/badge/Docker-Ready-blue?style=for-the-badge&logo=docker)](https://www.docker.com/) [![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326ce5?style=for-the-badge&logo=kubernetes)](https://kubernetes.io/) [![Linux](https://img.shields.io/badge/Linux-Deployments-FCC624?style=for-the-badge&logo=linux)](https://www.linux.org/) [![Cloud](https://img.shields.io/badge/Cloud-Ready-00A8E8?style=for-the-badge&logo=aws)](https://aws.amazon.com/)

**🔥 Complete Deployment Guide for Enterprise Authentication & Vault Server**

Comprehensive deployment documentation for Aether Vault Server, covering Docker, Kubernetes, cloud deployments, and production best practices.

[🐳 Docker Deployment](#-docker-deployment) • [☸️ Kubernetes Deployment](#️-kubernetes-deployment) • [☁️ Cloud Deployment](#️-cloud-deployment) • [🖥️ Traditional Deployment](#️-traditional-deployment) • [🔧 Production Setup](#-production-setup) • [📊 Monitoring & Logging](#-monitoring--logging) • [🛠️ Troubleshooting](#️-troubleshooting)

</div>

---

## 🐳 Docker Deployment

### 🎯 **Quick Start with Docker**

```bash
# 1. Clone the repository
git clone https://github.com/skygenesisenterprise/aether-vault.git
cd aether-vault/server

# 2. Build and run with Docker Compose
docker-compose up -d

# 3. Access the application
curl http://localhost:8080/api/v1/system/health
```

### 📄 **Dockerfile**

```dockerfile
# Multi-stage Dockerfile for production deployment
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aether-vault .

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S vault && \
    adduser -u 1001 -S vault -G vault

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/aether-vault .

# Copy configuration files
COPY .env.example .env

# Create directories
RUN mkdir -p /app/logs /app/data && \
    chown -R vault:vault /app

# Switch to non-root user
USER vault

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/system/health || exit 1

# Start the application
CMD ["./aether-vault"]
```

### 📋 **Docker Compose**

```yaml
# docker-compose.yml
version: "3.8"

services:
  # Aether Vault Application
  vault:
    build: .
    container_name: aether-vault
    ports:
      - "8080:8080"
    environment:
      - VAULT_SERVER_HOST=0.0.0.0
      - VAULT_SERVER_PORT=8080
      - VAULT_SERVER_ENVIRONMENT=production
      - VAULT_DATABASE_HOST=postgres
      - VAULT_DATABASE_PORT=5432
      - VAULT_DATABASE_USER=vault
      - VAULT_DATABASE_PASSWORD=vault_password
      - VAULT_DATABASE_DBNAME=vault
      - VAULT_DATABASE_SSLMODE=disable
      - VAULT_JWT_SECRET=your-super-secret-jwt-key-here-must-be-very-long
      - VAULT_SECURITY_ENCRYPTION_KEY=your-32-character-encryption-key-here
      - VAULT_AUDIT_ENABLED=true
      - VAULT_AUDIT_LOG_LEVEL=info
    depends_on:
      - postgres
    volumes:
      - ./logs:/app/logs
    restart: unless-stopped
    networks:
      - vault-network

  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: vault-postgres
    environment:
      - POSTGRES_DB=vault
      - POSTGRES_USER=vault
      - POSTGRES_PASSWORD=vault_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    restart: unless-stopped
    networks:
      - vault-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U vault"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis Cache (Optional)
  redis:
    image: redis:7-alpine
    container_name: vault-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped
    networks:
      - vault-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3

  # Nginx Reverse Proxy (Optional)
  nginx:
    image: nginx:alpine
    container_name: vault-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - vault
    restart: unless-stopped
    networks:
      - vault-network

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local

networks:
  vault-network:
    driver: bridge
```

### 🔧 **Production Docker Compose**

```yaml
# docker-compose.prod.yml
version: "3.8"

services:
  vault:
    build:
      context: .
      dockerfile: Dockerfile
    image: aether-vault:latest
    container_name: aether-vault-prod
    restart: always
    environment:
      - VAULT_SERVER_ENVIRONMENT=production
      - VAULT_DATABASE_HOST=postgres
      - VAULT_DATABASE_SSLMODE=require
    env_file:
      - .env.prod
    volumes:
      - /opt/aether-vault/logs:/app/logs
      - /opt/aether-vault/data:/app/data
    deploy:
      resources:
        limits:
          cpus: "2.0"
          memory: 2G
        reservations:
          cpus: "1.0"
          memory: 1G
    networks:
      - vault-network
    healthcheck:
      test:
        [
          "CMD",
          "wget",
          "--no-verbose",
          "--tries=1",
          "--spider",
          "http://localhost:8080/api/v1/system/health",
        ]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:15-alpine
    container_name: vault-postgres-prod
    restart: always
    environment:
      - POSTGRES_DB=vault
      - POSTGRES_USER=vault
      - POSTGRES_PASSWORD_FILE=/run/secrets/db_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./postgresql.conf:/etc/postgresql/postgresql.conf
    networks:
      - vault-network
    secrets:
      - db_password
    deploy:
      resources:
        limits:
          cpus: "1.0"
          memory: 2G
        reservations:
          cpus: "0.5"
          memory: 1G

secrets:
  db_password:
    file: ./secrets/db_password.txt

volumes:
  postgres_data:
    driver: local

networks:
  vault-network:
    driver: bridge
```

---

## ☸️ Kubernetes Deployment

### 📋 **Kubernetes Manifests**

#### **Namespace**

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: aether-vault
  labels:
    name: aether-vault
    app: aether-vault
```

#### **ConfigMap**

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: vault-config
  namespace: aether-vault
data:
  server.env: |
    VAULT_SERVER_HOST=0.0.0.0
    VAULT_SERVER_PORT=8080
    VAULT_SERVER_ENVIRONMENT=production
    VAULT_DATABASE_HOST=postgres-service
    VAULT_DATABASE_PORT=5432
    VAULT_DATABASE_USER=vault
    VAULT_DATABASE_DBNAME=vault
    VAULT_DATABASE_SSLMODE=require
    VAULT_SECURITY_KDF_ITERATIONS=100000
    VAULT_SECURITY_SALT_LENGTH=32
    VAULT_JWT_EXPIRATION=3600
    VAULT_AUDIT_ENABLED=true
    VAULT_AUDIT_LOG_LEVEL=info
    VAULT_AUDIT_LOG_FORMAT=json
```

#### **Secret**

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: vault-secrets
  namespace: aether-vault
type: Opaque
data:
  # Base64 encoded values
  db-password: dmF1bHRfcGFzc3dvcmQ= # vault_password
  jwt-secret: eW91ci1zdXBlci1zZWNyZXQtand0LWtleS1oZXJlLW11c3QtYmUtdmVyeS1sb25n
  encryption-key: eW91ci0zMi1jaGFyYWN0ZXItZW5jcnlwdGlvbi1rZXk=
```

#### **Deployment**

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: aether-vault
  namespace: aether-vault
  labels:
    app: aether-vault
spec:
  replicas: 3
  selector:
    matchLabels:
      app: aether-vault
  template:
    metadata:
      labels:
        app: aether-vault
    spec:
      containers:
        - name: vault
          image: aether-vault:latest
          ports:
            - containerPort: 8080
          envFrom:
            - configMapRef:
                name: vault-config
            - secretRef:
                name: vault-secrets
          env:
            - name: VAULT_DATABASE_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: vault-secrets
                  key: db-password
            - name: VAULT_JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: vault-secrets
                  key: jwt-secret
            - name: VAULT_SECURITY_ENCRYPTION_KEY
              valueFrom:
                secretKeyRef:
                  name: vault-secrets
                  key: encryption-key
          resources:
            requests:
              cpu: 100m
              memory: 256Mi
            limits:
              cpu: 500m
              memory: 512Mi
          livenessProbe:
            httpGet:
              path: /api/v1/system/health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /api/v1/system/health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
          volumeMounts:
            - name: logs
              mountPath: /app/logs
      volumes:
        - name: logs
          emptyDir: {}
      imagePullSecrets:
        - name: registry-secret
```

#### **Service**

```yaml
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: vault-service
  namespace: aether-vault
  labels:
    app: aether-vault
spec:
  selector:
    app: aether-vault
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
```

#### **Ingress**

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: vault-ingress
  namespace: aether-vault
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
spec:
  tls:
    - hosts:
        - vault.example.com
      secretName: vault-tls
  rules:
    - host: vault.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: vault-service
                port:
                  number: 80
```

#### **HorizontalPodAutoscaler**

```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: vault-hpa
  namespace: aether-vault
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: aether-vault
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

### 🚀 **Deploy to Kubernetes**

```bash
# 1. Create namespace
kubectl apply -f namespace.yaml

# 2. Apply configuration
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml

# 3. Deploy database (if needed)
kubectl apply -f postgres-deployment.yaml

# 4. Deploy application
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
kubectl apply -f hpa.yaml

# 5. Check deployment
kubectl get pods -n aether-vault
kubectl get services -n aether-vault
kubectl get ingress -n aether-vault

# 6. Check logs
kubectl logs -f deployment/aether-vault -n aether-vault
```

---

## ☁️ Cloud Deployment

### 🟠 **AWS Deployment**

#### **AWS ECS Task Definition**

```json
{
  "family": "aether-vault",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",
  "memory": "1024",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::account:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "vault",
      "image": "your-account.dkr.ecr.region.amazonaws.com/aether-vault:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "VAULT_SERVER_HOST",
          "value": "0.0.0.0"
        },
        {
          "name": "VAULT_SERVER_PORT",
          "value": "8080"
        },
        {
          "name": "VAULT_SERVER_ENVIRONMENT",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "VAULT_JWT_SECRET",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:vault/jwt-secret"
        },
        {
          "name": "VAULT_SECURITY_ENCRYPTION_KEY",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:vault/encryption-key"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/aether-vault",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "healthCheck": {
        "command": [
          "CMD-SHELL",
          "wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/system/health || exit 1"
        ],
        "interval": 30,
        "timeout": 5,
        "retries": 3
      }
    }
  ]
}
```

#### **AWS CloudFormation Template**

```yaml
# cloudformation.yaml
AWSTemplateFormatVersion: "2010-09-09"
Description: "Aether Vault Deployment on AWS"

Parameters:
  Environment:
    Type: String
    Default: production
    AllowedValues: [development, staging, production]

Resources:
  # VPC
  VPC:
    Type: AWS::EC2::VPC
    Properties:
      CidrBlock: 10.0.0.0/16
      EnableDnsHostnames: true
      EnableDnsSupport: true
      Tags:
        - Key: Name
          Value: !Sub "${Environment}-vault-vpc"

  # ECS Cluster
  ECSCluster:
    Type: AWS::ECS::Cluster
    Properties:
      ClusterName: !Sub "${Environment}-vault-cluster"
      CapacityProviders:
        - FARGATE
        - FARGATE_SPOT

  # Application Load Balancer
  LoadBalancer:
    Type: AWS::ElasticLoadBalancingV2::LoadBalancer
    Properties:
      Name: !Sub "${Environment}-vault-alb"
      Scheme: internet-facing
      Type: application
      Subnets:
        - !Ref PublicSubnet1
        - !Ref PublicSubnet2
      SecurityGroups:
        - !Ref LoadBalancerSecurityGroup

  # Secrets Manager
  JWTSecret:
    Type: AWS::SecretsManager::Secret
    Properties:
      Name: !Sub "${Environment}/vault/jwt-secret"
      Description: "JWT signing secret for Aether Vault"
      GenerateSecretString:
        SecretStringTemplate: "{}"
        GenerateStringKey: "jwt-secret"
        PasswordLength: 64
        ExcludeCharacters: '"@/\\'

  EncryptionKeySecret:
    Type: AWS::SecretsManager::Secret
    Properties:
      Name: !Sub "${Environment}/vault/encryption-key"
      Description: "Encryption key for Aether Vault"
      GenerateSecretString:
        SecretStringTemplate: "{}"
        GenerateStringKey: "encryption-key"
        PasswordLength: 32
        ExcludeCharacters: '"@/\\'

Outputs:
  ClusterName:
    Description: "ECS Cluster Name"
    Value: !Ref ECSCluster
    Export:
      Name: !Sub "${Environment}-vault-cluster-name"

  LoadBalancerDNS:
    Description: "Load Balancer DNS Name"
    Value: !GetAtt LoadBalancer.DNSName
    Export:
      Name: !Sub "${Environment}-vault-alb-dns"
```

### 🔵 **Azure Deployment**

#### **Azure Container Instance**

```yaml
# azure-container.yaml
apiVersion: 2021-09-01
location: eastus
name: aether-vault-aci
properties:
  containers:
    - name: vault
      properties:
        image: aether-vault:latest
        ports:
          - port: 8080
        resources:
          requests:
            cpu: 1.0
            memoryInGb: 2.0
        environmentVariables:
          - name: VAULT_SERVER_HOST
            value: 0.0.0.0
          - name: VAULT_SERVER_PORT
            value: 8080
          - name: VAULT_SERVER_ENVIRONMENT
            value: production
        secureEnvironmentVariables:
          - name: VAULT_JWT_SECRET
            value: ${JWT_SECRET}
          - name: VAULT_SECURITY_ENCRYPTION_KEY
            value: ${ENCRYPTION_KEY}
  osType: Linux
  ipAddress:
    type: Public
    ports:
      - port: 8080
    dnsNameLabel: vault-aci-demo
tags: null
type: Microsoft.ContainerInstance/containerGroups
```

### 🟡 **Google Cloud Deployment**

#### **Cloud Run Service**

```yaml
# cloud-run.yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: aether-vault
  annotations:
    run.googleapis.com/ingress: all
spec:
  template:
    metadata:
      annotations:
        run.googleapis.com/cpu-throttling: "false"
        run.googleapis.com/memory: "2Gi"
    spec:
      containerConcurrency: 80
      containers:
        - image: gcr.io/your-project/aether-vault:latest
          ports:
            - containerPort: 8080
          env:
            - name: VAULT_SERVER_HOST
              value: 0.0.0.0
            - name: VAULT_SERVER_PORT
              value: 8080
            - name: VAULT_SERVER_ENVIRONMENT
              value: production
          resources:
            limits:
              cpu: "2"
              memory: 2Gi
```

---

## 🖥️ Traditional Deployment

### 🐧 **Linux System Deployment**

#### **System Requirements**

- **CPU**: 2+ cores recommended
- **Memory**: 4GB+ RAM recommended
- **Storage**: 20GB+ SSD recommended
- **OS**: Ubuntu 20.04+, CentOS 8+, RHEL 8+

#### **Installation Script**

```bash
#!/bin/bash
# install.sh - Aether Vault Installation Script

set -e

# Configuration
VAULT_USER="vault"
VAULT_HOME="/opt/aether-vault"
VAULT_VERSION="latest"
SERVICE_NAME="aether-vault"

echo "🚀 Installing Aether Vault Server..."

# Create system user
sudo useradd -r -s /bin/false -d $VAULT_HOME $VAULT_USER || true

# Create directories
sudo mkdir -p $VAULT_HOME/{bin,config,logs,data}
sudo chown -R $VAULT_USER:$VAULT_USER $VAULT_HOME

# Download binary (replace with actual download URL)
cd /tmp
wget https://github.com/skygenesisenterprise/aether-vault/releases/latest/download/aether-vault-linux-amd64.tar.gz
tar -xzf aether-vault-linux-amd64.tar.gz
sudo mv aether-vault $VAULT_HOME/bin/
sudo chmod +x $VAULT_HOME/bin/aether-vault

# Copy configuration
sudo cp .env.example $VAULT_HOME/config/.env
sudo chown $VAULT_USER:$VAULT_USER $VAULT_HOME/config/.env

# Create systemd service
sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null <<EOF
[Unit]
Description=Aether Vault Server
After=network.target postgresql.service

[Service]
Type=simple
User=$VAULT_USER
Group=$VAULT_USER
WorkingDirectory=$VAULT_HOME
ExecStart=$VAULT_HOME/bin/aether-vault
Restart=always
RestartSec=10
Environment=VAULT_SERVER_ENVIRONMENT=production

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$VAULT_HOME/logs $VAULT_HOME/data

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and enable service
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME

echo "✅ Installation completed!"
echo "🔧 Configure: $VAULT_HOME/config/.env"
echo "🚀 Start with: sudo systemctl start $SERVICE_NAME"
echo "📊 Status: sudo systemctl status $SERVICE_NAME"
```

#### **Configuration Management**

```bash
# Configuration validation script
#!/bin/bash
# validate-config.sh

CONFIG_FILE="/opt/aether-vault/config/.env"

echo "🔍 Validating configuration..."

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Configuration file not found: $CONFIG_FILE"
    exit 1
fi

# Source environment variables
source $CONFIG_FILE

# Validate required variables
if [ -z "$VAULT_JWT_SECRET" ]; then
    echo "❌ VAULT_JWT_SECRET is required"
    exit 1
fi

if [ -z "$VAULT_SECURITY_ENCRYPTION_KEY" ]; then
    echo "❌ VAULT_SECURITY_ENCRYPTION_KEY is required"
    exit 1
fi

# Check secret lengths
if [ ${#VAULT_JWT_SECRET} -lt 32 ]; then
    echo "❌ JWT_SECRET must be at least 32 characters"
    exit 1
fi

if [ ${#VAULT_SECURITY_ENCRYPTION_KEY} -ne 32 ]; then
    echo "❌ ENCRYPTION_KEY must be exactly 32 characters"
    exit 1
fi

echo "✅ Configuration validation passed!"
```

---

## 🔧 Production Setup

### 🛡️ **Security Hardening**

#### **System Hardening**

```bash
# Security hardening script
#!/bin/bash
# hardening.sh

echo "🔒 Hardening production environment..."

# Update system
sudo apt update && sudo apt upgrade -y

# Install security tools
sudo apt install -y fail2ban ufw logwatch rkhunter

# Configure firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw --force enable

# Configure fail2ban
sudo tee /etc/fail2ban/jail.local > /dev/null <<EOF
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 3

[sshd]
enabled = true
port = ssh
logpath = /var/log/auth.log
EOF

sudo systemctl restart fail2ban

# Set file permissions
sudo chmod 600 /opt/aether-vault/config/.env
sudo chmod 700 /opt/aether-vault/data
sudo chown -R vault:vault /opt/aether-vault

echo "✅ Security hardening completed!"
```

#### **SSL/TLS Configuration**

```nginx
# SSL Configuration for Nginx
server {
    listen 443 ssl http2;
    server_name vault.example.com;

    # SSL certificates
    ssl_certificate /etc/ssl/certs/vault.crt;
    ssl_certificate_key /etc/ssl/private/vault.key;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 📊 **Monitoring Setup**

#### **Prometheus Configuration**

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "aether-vault"
    static_configs:
      - targets: ["localhost:8080"]
    metrics_path: "/metrics"
    scrape_interval: 10s

  - job_name: "postgres"
    static_configs:
      - targets: ["localhost:9187"]

  - job_name: "node"
    static_configs:
      - targets: ["localhost:9100"]
```

#### **Grafana Dashboard**

```json
{
  "dashboard": {
    "title": "Aether Vault Monitoring",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{status}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      },
      {
        "title": "Database Connections",
        "type": "graph",
        "targets": [
          {
            "expr": "pg_stat_database_numbackends",
            "legendFormat": "Active Connections"
          }
        ]
      }
    ]
  }
}
```

---

## 📊 Monitoring & Logging

### 📈 **Application Metrics**

#### **Prometheus Metrics**

```go
// Add metrics to the application
var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "path"},
    )

    activeConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )
)

// Metrics middleware
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())

        requestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
        requestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
    }
}
```

### 📝 **Structured Logging**

#### **Log Configuration**

```go
// Enhanced logging configuration
func setupLogging(cfg *config.Config) {
    // Create logger
    log := logrus.New()

    // Set log level
    level, err := logrus.ParseLevel(cfg.Audit.LogLevel)
    if err != nil {
        level = logrus.InfoLevel
    }
    log.SetLevel(level)

    // Set formatter
    if cfg.Audit.LogFormat == "json" {
        log.SetFormatter(&logrus.JSONFormatter{
            TimestampFormat: time.RFC3339,
        })
    } else {
        log.SetFormatter(&logrus.TextFormatter{
            TimestampFormat: time.RFC3339,
            FullTimestamp:   true,
        })
    }

    // Add hooks for structured logging
    log.AddHook(&ContextHook{})

    // Log to file
    if file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
        log.SetOutput(io.MultiWriter(os.Stdout, file))
    }
}
```

### 🚨 **Alerting Rules**

#### **Prometheus Alert Rules**

```yaml
# alerts.yml
groups:
  - name: aether-vault
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is above 5% for 5 minutes"

      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          description: "95th percentile response time is above 1 second"

      - alert: DatabaseConnectionFailure
        expr: up{job="postgres"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database connection failure"
          description: "PostgreSQL is down"
```

---

## 🛠️ Troubleshooting

### 🔍 **Common Issues**

#### **Database Connection Issues**

```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Test connection
psql -h localhost -U vault -d vault -c "SELECT 1;"

# Check logs
sudo tail -f /var/log/postgresql/postgresql-15-main.log

# Common solutions
# 1. Check if PostgreSQL is running
# 2. Verify connection parameters
# 3. Check network connectivity
# 4. Validate SSL settings
```

#### **Authentication Issues**

```bash
# Check JWT secret
echo $VAULT_JWT_SECRET | wc -c

# Test JWT generation
go run -c 'package main; import ("fmt"; "time"; "github.com/golang-jwt/jwt/v5"); func main() { token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"exp": time.Now().Add(time.Hour).Unix()}); ss, _ := token.SignedString([]byte("your-secret")); fmt.Println(ss) }'

# Check encryption key
echo -n $VAULT_SECURITY_ENCRYPTION_KEY | wc -c
```

#### **Performance Issues**

```bash
# Check system resources
top
htop
iotop
df -h

# Check application logs
tail -f /opt/aether-vault/logs/app.log

# Check database performance
psql -U vault -d vault -c "SELECT * FROM pg_stat_activity;"

# Profile application
go tool pprof http://localhost:8080/debug/pprof/profile
```

### 📊 **Health Checks**

#### **Application Health Check**

```bash
# Basic health check
curl -f http://localhost:8080/api/v1/system/health || exit 1

# Detailed health check
curl -s http://localhost:8080/api/v1/system/health | jq .

# Database health check
psql -U vault -d vault -c "SELECT 1;" > /dev/null 2>&1 && echo "DB OK" || echo "DB FAIL"

# Memory check
free -m | grep -q "^Mem:" && echo "Memory OK" || echo "Memory FAIL"

# Disk space check
df / | grep -q "/$" && echo "Disk OK" || echo "Disk FAIL"
```

#### **Monitoring Script**

```bash
#!/bin/bash
# monitor.sh

VAULT_URL="http://localhost:8080"
LOG_FILE="/var/log/vault-monitor.log"

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a $LOG_FILE
}

check_vault() {
    if curl -f -s $VAULT_URL/api/v1/system/health > /dev/null; then
        log "✅ Vault is healthy"
        return 0
    else
        log "❌ Vault is unhealthy"
        return 1
    fi
}

check_database() {
    if pg_isready -h localhost -p 5432 -U vault > /dev/null 2>&1; then
        log "✅ Database is ready"
        return 0
    else
        log "❌ Database is not ready"
        return 1
    fi
}

check_system() {
    # Check memory
    MEMORY_USAGE=$(free | awk '/Mem/{printf("%.1f", $3/$2 * 100.0)}')
    if (( $(echo "$MEMORY_USAGE > 80" | bc -l) )); then
        log "⚠️ High memory usage: ${MEMORY_USAGE}%"
    fi

    # Check disk
    DISK_USAGE=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
    if [ $DISK_USAGE -gt 80 ]; then
        log "⚠️ High disk usage: ${DISK_USAGE}%"
    fi
}

# Run checks
check_vault
check_database
check_system
```

---

## 📚 Deployment Scripts

### 🚀 **Complete Deployment Script**

```bash
#!/bin/bash
# deploy.sh - Complete deployment script

set -e

# Configuration
ENVIRONMENT=${1:-production}
VERSION=${2:-latest}
VAULT_HOME="/opt/aether-vault"
SERVICE_NAME="aether-vault"

echo "🚀 Deploying Aether Vault Server ($ENVIRONMENT)..."

# Create backup
if [ -d "$VAULT_HOME" ]; then
    echo "📦 Creating backup..."
    sudo tar -czf /tmp/vault-backup-$(date +%Y%m%d-%H%M%S).tar.gz $VAULT_HOME
fi

# Download new version
echo "📥 Downloading version $VERSION..."
cd /tmp
wget https://github.com/skygenesisenterprise/aether-vault/releases/download/$VERSION/aether-vault-linux-amd64.tar.gz
tar -xzf aether-vault-linux-amd64.tar.gz

# Stop service
echo "⏹️ Stopping service..."
sudo systemctl stop $SERVICE_NAME || true

# Update binary
echo "🔄 Updating binary..."
sudo mv aether-vault $VAULT_HOME/bin/
sudo chmod +x $VAULT_HOME/bin/aether-vault

# Run database migrations
echo "🗄️ Running migrations..."
sudo -u vault $VAULT_HOME/bin/aether-vault migrate

# Start service
echo "▶️ Starting service..."
sudo systemctl start $SERVICE_NAME

# Wait for startup
echo "⏳ Waiting for startup..."
sleep 10

# Health check
echo "🔍 Health check..."
if curl -f http://localhost:8080/api/v1/system/health > /dev/null; then
    echo "✅ Deployment successful!"
else
    echo "❌ Deployment failed!"
    sudo systemctl status $SERVICE_NAME
    exit 1
fi

echo "🎉 Deployment completed successfully!"
```

---

## 🔗 Related Documentation

- [📖 Server Documentation](./README.md)
- [🏗️ Architecture Guide](./architecture.md)
- [⚙️ Configuration Guide](./configuration.md)
- [📚 API Documentation](./api.md)
- [🔒 Security Guide](./security.md)

---

<div align="center">

### 🚀 **Deploy Your Aether Vault Server with Confidence!**

[📖 Full Documentation](../../README.md) • [🐛 Report Issues](https://github.com/skygenesisenterprise/aether-vault/issues) • [💡 Deployment Help](https://github.com/skygenesisenterprise/aether-vault/discussions)

---

**🔧 Production-Ready Deployment for Enterprise Authentication & Vault Management**

**Made with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

</div>
