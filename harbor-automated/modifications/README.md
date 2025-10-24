# Modifications

Harbor chart customizations as separate files.

## Structure

```
modifications/
├── helpers/       # Template helpers (.tpl)
├── templates/     # Custom templates (.yaml)
├── values/        # Values additions (.yaml)
└── chart/         # Chart.yaml modifications (.yaml)
```

## How It Works

`harbor-modifier` applies these files:
- `helpers/` → Appended to `_helpers.tpl`
- `templates/` → Copied to `templates/`
- `values/` → Merged into `values.yaml`
- `chart/` → Merged into `Chart.yaml`

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

**helpers/** - Template helpers
- `chart.tpl` - Chart label
- `labels.tpl` - Standard labels
- `image-reference.tpl` - Image digest support

**templates/** - Custom resources
- `traefik-ingressroute.yaml` - Traefik routing
- `traefik-middleware.yaml` - HTTPS redirect

**values/** - Configuration
- `labels.yaml` - Label customization
- `image-digests.yaml` - Image digests
- `postgresql.yaml` - Reliza PostgreSQL

**chart/** - Chart metadata
- `dependencies.yaml` - Reliza PostgreSQL dependency
