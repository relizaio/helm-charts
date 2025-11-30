# Harbor Automated Chart

Automated Harbor Helm chart customization with Reliza CI/CD integration.

**Harbor:** 1.18.0 (v2.14.0) | **Go:** 1.25.3

## Quick Start

```bash
# Build locally (required before committing)
./build-local.sh

# Or with specific version
./build-local.sh 1.19.0

# Then commit and push
git add harbor-helm/
git commit -m "chore: update Harbor chart"
git push
```

## Local Build Process

**Important:** Chart generation must be done locally before committing. CI only validates and packages existing files.

### Using the build script (recommended)

```bash
# Default build (Harbor 1.18.0, iteration 1)
./build-local.sh

# Specific Harbor version
./build-local.sh 1.19.0

# Specific version and iteration
./build-local.sh 1.19.0 2
```

### Using Make

```bash
make clean          # Remove previous artifacts
make setup          # Build tool and generate chart
make lint           # Validate chart
```

### What the build does

1. Builds the `harbor-modifier` Go tool
2. Pulls official Harbor chart from helm.goharbor.io
3. Applies Reliza modifications from `modifications/`
4. Builds chart dependencies
5. Sets chart version to `{HARBOR_VERSION}-reliza.{ITERATION}`
6. Validates the generated chart

## What It Does

Customizes official Harbor chart with:
- **Reliza-CD compatibility** - Smart image references work with tag replacement
- **Label standardization** - Fixes subchart conflicts
- **Image digest support** - Pin images by digest
- **Reliza PostgreSQL** - Alternative database
- **Traefik IngressRoute** - Native Traefik support

## Structure

```
harbor-automated/
├── cmd/harbor-modifier/    # Go tool
├── modifications/          # All customizations
├── examples/               # Example values
├── harbor-helm/           # Generated locally, committed to git
├── build-local.sh         # Local build script
└── Makefile               # Build automation
```

## CI/CD

**Workflow:**
1. **Local:** Run `./build-local.sh` to generate/update `harbor-helm/`
2. **Commit:** Push changes including `harbor-helm/` directory
3. **CI:** Validates chart and publishes to registry

CI does NOT regenerate the chart - it only validates what you committed.

**Secrets (already exist):**
- `RH_LIBRARY_HELM_LOGIN`
- `RH_LIBRARY_HELM_PASS`
- `RELIZA_HARBOR_HELM_API_ID`
- `RELIZA_HARBOR_HELM_API_KEY`

## Usage

### Install from Registry
```bash
helm install harbor oci://registry.relizahub.com/library/harbor \
  --version 1.18.0-reliza.1 -n harbor --create-namespace
```

### Add Modification
```bash
# 1. Add your modification
echo "myConfig: {enabled: true}" > modifications/values/my-config.yaml

# 2. Rebuild locally
./build-local.sh

# 3. Commit everything
git add modifications/ harbor-helm/
git commit -m "feat: add config"
git push
```

### Upgrade Harbor Version
```bash
# 1. Clean and rebuild with new version
make clean
./build-local.sh 1.19.0

# 2. Commit
git add harbor-helm/
git commit -m "chore: upgrade Harbor to 1.19.0"
git push
```

## Documentation

- **README.md** (this file) - Quick reference
- **examples/README.md** - Configuration examples and patterns
- **modifications/README.md** - Modifications structure

## Commands

```bash
make setup    # Build and generate
make clean    # Clean artifacts
make lint     # Validate
make help     # Show all
```
