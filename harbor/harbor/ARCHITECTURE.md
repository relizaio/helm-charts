# Architecture

## Chart Structure

```
harbor/
├── Chart.yaml              # Chart metadata and dependencies
├── values.yaml             # Default configuration
├── templates/              # Kubernetes manifests
│   ├── _helpers.tpl       # Template helpers
│   ├── ingress.yaml       # Standard Ingress
│   ├── ingressroute.yaml  # Traefik IngressRoute
│   ├── middlewares.yaml   # Traefik Middlewares
│   └── NOTES.txt          # Post-install notes
└── examples/              # Configuration examples
    ├── values-s3.yaml
    ├── values-external-db.yaml
    └── values-s3-external-db.yaml
```

## Dependencies

1. **PostgreSQL 17** (Relizaio) - Optional bundled database
2. **Harbor** (GoHarbor) - Container registry

## How It Works

### Value Passing

Values under `harbor:` in parent chart are passed to Harbor subchart:

```yaml
# Parent (harbor/values.yaml)
harbor:
  externalURL: http://harbor.local
  persistence:
    imageChartStorage:
      type: s3

# Subchart (harbor) receives:
externalURL: http://harbor.local
persistence:
  imageChartStorage:
    type: s3
```

### Storage

Harbor uses different storage for different components:

| Component | Storage | Configurable |
|-----------|---------|--------------|
| Container Images | PVC or S3/GCS/Azure | Yes (`imageChartStorage.type`) |
| JobService Logs | PVC | No (always PVC) |
| Trivy Database | PVC | No (always PVC) |
| PostgreSQL | PVC or External | Yes (`externalDatabase.enabled`) |

**Key Point**: When `imageChartStorage.type=s3`, registry PVC is not created (Harbor chart has conditional logic).

### Database

Two modes:

**Bundled** (default):
- PostgreSQL deployed as `<release-name>-postgresql`
- Harbor connects automatically

**External**:
- Set `externalDatabase.enabled=true`
- Set `postgresql.enabled=false`
- Configure `harbor.database.external`

### Ingress

Two types supported:

**Standard Ingress**:
```yaml
ingress:
  type: ingress
```

**Traefik IngressRoute**:
```yaml
ingress:
  type: ingressroute
```

## Configuration Flow

```
User values.yaml
    ↓
Parent chart merges with defaults
    ↓
Values under harbor: → Harbor subchart
Values under postgresql: → PostgreSQL subchart
    ↓
Templates rendered
    ↓
Kubernetes resources created
```

## Key Design Decisions

1. **Minimal templates** - Only ingress/ingressroute, no validation or monitoring
2. **Explicit configuration** - Database settings in both parent and Harbor subchart
3. **Examples over docs** - Show don't tell with working examples
4. **No magic** - Clear, straightforward value passing
