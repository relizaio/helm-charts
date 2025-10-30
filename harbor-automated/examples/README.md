# Harbor Configuration Examples

This directory contains example configurations for deploying Harbor with Reliza customizations.

### How It Works

1. **Configure credentials once** in `relizapostgresql.auth.*`
2. **Harbor templates automatically use** these values for:
   - Database host: `{release-name}-relizapostgresql`
   - Database username: `relizapostgresql.auth.username`
   - Database password: `relizapostgresql.auth.password`
   - Database name: `relizapostgresql.auth.database`

### External Database

For external databases, use `database.external.*`:

```yaml
database:
  type: external
  external:
    host: my-postgres.example.com
    port: "5432"
    username: harbor
    password: "secret"
    coreDatabase: registry
    sslmode: disable
```

## Example Files

### 1. `values-reliza-postgresql.yaml`
**Purpose:** Basic Harbor with Reliza PostgreSQL  
**Use Case:** Simple deployment with internal database

**Features:**
- Reliza PostgreSQL with 20Gi storage
- ClusterIP exposure (access via port-forward)
- Minimal configuration

**Deploy:**
```bash
helm install harbor ./harbor-helm \
  -f examples/values-reliza-postgresql.yaml \
  -n harbor --create-namespace
```

### 2. `values-reliza-full.yaml`
**Purpose:** Complete production-ready configuration  
**Use Case:** Full-featured Harbor deployment

**Features:**
- Reliza PostgreSQL with 20Gi storage
- Traefik IngressRoute with Let's Encrypt
- S3 storage example (commented)
- Resource limits configured
- Metrics enabled
- All components configured

**Deploy:**
```bash
helm install harbor ./harbor-helm \
  -f examples/values-reliza-full.yaml \
  -n harbor --create-namespace
```

### 3. `values-traefik.yaml`
**Purpose:** Harbor with Traefik ingress  
**Use Case:** Traefik-based Kubernetes clusters

**Features:**
- Reliza PostgreSQL
- Traefik IngressRoute
- Let's Encrypt TLS
- HTTPS redirect
- Custom middlewares support

**Deploy:**
```bash
helm install harbor ./harbor-helm \
  -f examples/values-traefik.yaml \
  -n harbor --create-namespace
```

### 4. `values-k3s-simple.yaml`
**Purpose:** Minimal k3s testing  
**Use Case:** Quick local testing

**Features:**
- Reliza PostgreSQL with 5Gi storage
- ClusterIP exposure
- Minimal resource requests
- Small PVCs for testing

**Deploy:**
```bash
helm install harbor ./harbor-helm \
  -f examples/values-k3s-simple.yaml \
  -n harbor --create-namespace

# Access via port-forward
kubectl port-forward -n harbor svc/harbor 8080:80
```

### 5. `values-k3s-test.yaml`
**Purpose:** k3s with Traefik  
**Use Case:** Local k3s cluster testing

**Features:**
- Reliza PostgreSQL with 5Gi storage
- Traefik IngressRoute (k3s default)
- Local-path storage class
- Reduced resource limits
- No TLS (easier testing)

**Deploy:**
```bash
# Add harbor.local to /etc/hosts
echo "127.0.0.1 harbor.local" | sudo tee -a /etc/hosts

helm install harbor ./harbor-helm \
  -f examples/values-k3s-test.yaml \
  -n harbor --create-namespace

# Access at http://harbor.local
```

## Common Customizations

### Change Database Credentials

```yaml
relizapostgresql:
  auth:
    username: myuser
    password: "MySecurePassword123!"
    database: mydb
```

### Increase Database Storage

```yaml
relizapostgresql:
  primary:
    persistence:
      size: 100Gi
```

### Configure Database Resources

```yaml
relizapostgresql:
  primary:
    resources:
      requests:
        memory: 1Gi
        cpu: 500m
      limits:
        memory: 4Gi
        cpu: 2000m
```

### Enable Database Metrics

```yaml
relizapostgresql:
  metrics:
    enabled: true
    serviceMonitor:
      enabled: true  # If using Prometheus Operator
```

## Migration from Old Configuration

If you have old configurations with `database.internal.*`:

### Before (Old)
```yaml
database:
  type: internal
  internal:
    password: "changeit"
```

### After (New)
```yaml
relizapostgresql:
  enabled: true
  auth:
    username: harbor
    password: "changeit"
    database: registry

database:
  type: internal
```

## Troubleshooting

### Check Database Connection

```bash
# Get PostgreSQL pod
kubectl get pods -n harbor | grep relizapostgresql

# Check logs
kubectl logs -n harbor <release>-relizapostgresql-0

# Test connection from Harbor core
kubectl exec -n harbor <release>-harbor-core-xxx -- \
  psql -h <release>-relizapostgresql -U harbor -d registry -c "SELECT 1"
```

### Verify Configuration

```bash
# Check what values Harbor is using
kubectl get cm -n harbor <release>-harbor-core -o yaml | grep POSTGRESQL
```

### Common Issues

1. **Connection refused**: Check if PostgreSQL pod is running
2. **Authentication failed**: Verify `relizapostgresql.auth.password` matches
3. **Database not found**: Verify `relizapostgresql.auth.database` is correct

## Additional Resources

- [Harbor Documentation](https://goharbor.io/docs/)
- [Reliza PostgreSQL Chart](https://registry.relizahub.com/library/postgresql)
- [Traefik IngressRoute](https://doc.traefik.io/traefik/routing/providers/kubernetes-crd/)
