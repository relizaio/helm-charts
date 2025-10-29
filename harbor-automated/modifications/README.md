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

**template-overlays/** - Complete template replacements (reliza-cd compatible)
- `core/core-dpl.yaml` - Core deployment with inline image conditional
- `core/core-pre-upgrade-job.yaml` - Pre-upgrade job with inline conditional
- `database/database-ss.yaml` - Database statefulset (2 images)
- `exporter/exporter-dpl.yaml` - Metrics exporter
- `jobservice/jobservice-dpl.yaml` - Job service
- `nginx/deployment.yaml` - Nginx proxy
- `portal/deployment.yaml` - Web portal
- `redis/statefulset.yaml` - Redis cache
- `registry/registry-dpl.yaml` - Registry and controller (2 images)
- `trivy/trivy-sts.yaml` - Trivy scanner

**helpers/** - Template helpers
- `chart.tpl` - Chart label
- `labels.tpl` - Standard labels
- `image-ref.tpl` - Smart image reference (reliza-cd compatible)
- `image-reference.tpl` - Legacy image digest support (unused)

**templates/** - Custom resources
- `traefik-ingressroute.yaml` - Traefik routing
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
