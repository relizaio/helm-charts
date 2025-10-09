# Relizaio Harbor Helm Chart

Harbor container registry with PostgreSQL 17 and Traefik ingress support.

## Features

- Harbor 2.14.0
- PostgreSQL 17 (Relizaio)
- Traefik IngressRoute or standard Ingress
- S3/GCS/Azure storage support
- External database support

## Quick Start

```bash
# Add dependencies
helm dependency build

# Install with defaults
helm install my-harbor . -n harbor --create-namespace

# Install with S3 storage
helm install my-harbor . -n harbor --create-namespace -f examples/values-s3.yaml

# Install with external database
helm install my-harbor . -n harbor --create-namespace -f examples/values-external-db.yaml
```

## Configuration

### Storage Options

**Filesystem** (default):
```yaml
harbor:
  persistence:
    imageChartStorage:
      type: filesystem
```

**S3**:
```yaml
harbor:
  persistence:
    imageChartStorage:
      type: s3
      s3:
        region: us-east-1
        bucket: my-harbor-registry
```

### Database Options

**Bundled PostgreSQL** (default):
```yaml
externalDatabase:
  enabled: false
postgresql:
  enabled: true
```

**External PostgreSQL**:
```yaml
externalDatabase:
  enabled: true
  host: postgres.example.com
  password: "SecurePassword"
postgresql:
  enabled: false

harbor:
  database:
    external:
      host: "postgres.example.com"
      password: "SecurePassword"
```

### Ingress Options

**Standard Ingress** (with cert-manager):
```yaml
ingress:
  enabled: true
  type: ingress
  host: harbor.example.com
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  tls:
    enabled: true
    secretName: harbor-tls
```

**Traefik IngressRoute** (with Let's Encrypt):
```yaml
ingress:
  enabled: true
  type: ingressroute
  host: harbor.example.com
  createMiddlewares: true
  httpsRedirect:
    enabled: true
  tls:
    enabled: true
    certResolver: le
```

## Examples

See `examples/` directory:
- `values-ingress.yaml` - Standard Ingress with cert-manager
- `values-traefik.yaml` - Traefik IngressRoute with middlewares
- `values-s3.yaml` - S3 storage configuration
- `values-external-db.yaml` - External PostgreSQL configuration
- `values-s3-external-db.yaml` - Both S3 and external PostgreSQL

## Important Notes

### Storage

- When using S3/GCS/Azure, registry PVC is not created
- JobService and Trivy always need PVCs
- Set `registry.size: 0` when using object storage for clarity

### Database

- Harbor subchart receives values under `harbor:` key
- External database requires configuration in both parent and Harbor subchart
- Bundled PostgreSQL uses `<release-name>-postgresql` as hostname

## Values

Key configuration options:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `externalDatabase.enabled` | Use external PostgreSQL | `false` |
| `postgresql.enabled` | Deploy bundled PostgreSQL | `true` |
| `harbor.persistence.imageChartStorage.type` | Storage backend | `filesystem` |
| `ingress.enabled` | Enable ingress | `false` |
| `ingress.type` | Ingress type (`ingress` or `ingressroute`) | `ingress` |

See `values.yaml` for all options.

## Upgrading

```bash
helm upgrade my-harbor . -n harbor -f your-values.yaml
```

## Uninstalling

```bash
helm uninstall my-harbor -n harbor
```

PVCs are kept by default (`resourcePolicy: keep`).
