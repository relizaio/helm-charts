# Modifications

Harbor chart customizations as separate files.

## Structure

```
modifications/
├── template-overlays/ # Complete template file replacements
├── helpers/           # Template helpers (.tpl)
├── templates/         # Custom templates (.yaml)
├── values/            # Values additions (.yaml)
└── chart/             # Chart.yaml modifications (.yaml)
```

## How It Works

`harbor-modifier` applies these files:
- `template-overlays/` → Copied to `templates/` (overwrites originals)
- `helpers/` → Appended to `_helpers.tpl`
- `templates/` → Copied to `templates/` (new files)
- `values/` → Merged into `values.yaml`
- `chart/` → Merged into `Chart.yaml`

### Reliza-CD Compatibility

Template overlays use the `harbor.imageRef` helper for smart image references:
- **Problem**: Harbor templates use `repository:tag`, but reliza-cd puts full image references (with digests) into the `repository` field
- **Solution**: `harbor.imageRef` helper checks if repository contains `:` - if yes, uses as-is; otherwise appends tag
- **Result**: Works for both manual deployments (appends tag) and reliza-cd (uses full reference as-is)
- **Bonus**: Supports optional digest parameter for manual digest pinning

## Adding Modifications

```bash
# Add file
cat > values/my-config.yaml << EOF
myConfig:
  enabled: true
EOF

# Test
cd .. && make clean && make setup

# Commit
git add modifications/ && git commit -m "feat: add config" && git push
```

## Files

**template-overlays/** - Complete template replacements (use harbor.imageRef helper)
- `core/core-dpl.yaml` - Core service
- `core/core-pre-upgrade-job.yaml` - Pre-upgrade job with inline conditional
- `exporter/exporter-dpl.yaml` - Metrics exporter
- `jobservice/jobservice-dpl.yaml` - Job service
- `nginx/deployment.yaml` - Nginx reverse proxy
- `portal/deployment.yaml` - Web portal
- `redis/statefulset.yaml` - Redis cache
- `registry/registry-dpl.yaml` - Registry and registryctl
- `trivy/trivy-sts.yaml` - Trivy scanner

Note: Harbor's internal database templates (database-ss.yaml, database-svc.yaml, database-secret.yaml) 
are NOT included - they've been completely removed in favor of relizapostgresql subchart.

**helpers/** - Template helpers
- `chart.tpl` - Chart label
- `labels.tpl` - Standard labels
- `image-ref.tpl` - Smart image reference (reliza-cd compatible)

**Template Patches (applied by main.go):**
- `harbor.autoGenCertForNginx` - Patched to exclude Traefik type (TLS handled by Traefik, not nginx)
- `registry-cm.yaml` - Patched to use token auth when TLS enabled (fixes robot account authentication)
- `registry-dpl.yaml` - Patched to mount token certificate when TLS enabled (required for token auth)

**templates/** - Custom resources
- `traefik-ingressroute.yaml` - Traefik routing with priorities (API, chartrepo, registry, service, UI)
- `traefik-middleware.yaml` - HTTPS redirect

**values/** - Configuration
- `labels.yaml` - Label customization
- `image-digests.yaml` - Image digests
- `postgresql.yaml` - Reliza PostgreSQL

**chart/** - Chart metadata
- `dependencies.yaml` - Reliza PostgreSQL dependency
- `name.yaml` - Chart name override (harbor-helm)

**Root files**
- `.helmignore` - Files to exclude from packaging (CI secrets)
