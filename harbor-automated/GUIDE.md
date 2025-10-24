# Harbor Automated Chart - Complete Guide

## Quick Start

### Local Development
```bash
make setup    # Build and generate chart
make lint     # Validate
```

### CI/CD (Automatic)
```bash
git add modifications/
git commit -m "feat: add feature"
git push
# CI builds, commits chart, and publishes (3-5 min)
```

---

## CI/CD Integration

### Two-Stage Pipeline

**Stage 1:** Build chart (`.github/workflows/harbor-build-chart.yml`)
- Triggers on source changes
- Builds Go tool → Generates chart → Commits with `[skip ci]`

**Stage 2:** Publish chart (`.github/workflows/helmbuild.yml`)
- Triggers on `harbor-helm/` changes
- Uses `reliza-helm-action` → Publishes to registry

### Workflow
```
You commit → Stage 1 builds → Commits harbor-helm/ → Stage 2 publishes
```

### Monitoring
- **Stage 1:** `GitHub → Actions → Build Harbor Chart`
- **Stage 2:** `GitHub → Actions → Push Helm Charts → build-harbor-automated`
- **RelizaHub:** `Projects → Harbor Helm`

---

## Adding Modifications

### 1. Add File
```bash
cat > modifications/values/my-config.yaml << EOF
myConfig:
  enabled: true
EOF
```

### 2. Test
```bash
make clean && make setup
make lint
```

### 3. Commit
```bash
git add modifications/
git commit -m "feat: add my config"
git push
```

---

## Upgrading Harbor

```bash
# 1. Update version
vim Makefile  # HARBOR_VERSION?=1.19.0

# 2. Test
make clean && make setup HARBOR_VERSION=1.19.0

# 3. Commit
git add Makefile
git commit -m "chore: upgrade to Harbor 1.19.0"
git push
```

---

## Using as Subchart

```yaml
# Chart.yaml
dependencies:
  - name: harbor
    version: 1.18.0-reliza.1
    repository: oci://registry.relizahub.com/library
    alias: registry

# values.yaml
registry:
  externalURL: https://registry.example.com
  relizapostgresql:
    enabled: true
  postgresql:
    enabled: false
```

---

## Modifications

### Label Standardization
Fixes subchart selector conflicts:
```yaml
selector:
  matchLabels:
    app.kubernetes.io/name: harbor
    app.kubernetes.io/instance: {{ .Release.Name }}
    # No version - stable across upgrades
```

### Image Digest Support
```yaml
imageDigests:
  core:
    digest: "sha256:abc123..."
```

### Reliza PostgreSQL
```yaml
relizapostgresql:
  enabled: true
postgresql:
  enabled: false
```

### Traefik IngressRoute
```yaml
expose:
  type: traefik
  traefik:
    enabled: true
    host: harbor.example.com
    tls:
      certResolver: le
```

---

## Troubleshooting

### Chart not building
```bash
# Check logs
GitHub → Actions → Build Harbor Chart

# Test locally
make clean && make setup
```

### Chart not publishing
```bash
# Check if committed
git log harbor-helm/

# Check Stage 2
GitHub → Actions → Push Helm Charts
```

### `[skip ci]` behavior
- Stage 1 commits with `[skip ci]` to prevent loop
- Stage 2 (`helmbuild.yml`) ignores `[skip ci]` and runs anyway
- This is correct - Stage 2 should publish the chart

---

## FAQ

**Q: What do I commit?**  
A: Only source files (`cmd/`, `modifications/`, `examples/`). CI commits `harbor-helm/`.

**Q: How long does CI take?**  
A: 3-5 minutes total (2-3 min build + 1-2 min publish).

**Q: Can I test before pushing?**  
A: Yes, run `make clean && make setup && make lint` locally.

**Q: How do I rollback?**  
A: Revert the commit or install previous version from registry.

**Q: Where is the chart published?**  
A: `registry.relizahub.com/library/harbor`

---

## Commands

```bash
make setup    # Build tool and generate chart
make clean    # Remove generated files
make lint     # Validate chart
make build    # Build Go tool only
make help     # Show all commands
```

---

## Files to Commit

✅ Commit: `cmd/`, `modifications/`, `examples/`, `*.md`, `Makefile`, `go.mod`  
❌ Don't commit: `bin/`, `packages/`, `*.tgz`  
⚙️ CI commits: `harbor-helm/`
