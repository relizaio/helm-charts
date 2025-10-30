# Harbor Automated Chart

Automated Harbor Helm chart customization with Reliza CI/CD integration.

**Harbor:** 1.18.0 (v2.14.0) | **Go:** 1.25.3

## Quick Start

```bash
# Local
make setup && make lint

# CI/CD
git add modifications/ && git commit -m "feat: change" && git push
# CI builds and publishes automatically (3-5 min)
```

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
├── harbor-helm/           # Generated (by CI)
└── Makefile               # Build automation
```

## CI/CD

**Two stages:**
1. Build chart from modifications → commit to git
2. Publish with `reliza-helm-action` → registry

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
echo "myConfig: {enabled: true}" > modifications/values/my-config.yaml
make clean && make setup && make lint
git add modifications/ && git commit -m "feat: add config" && git push
```

### Upgrade Harbor
```bash
make clean && make setup HARBOR_VERSION=1.19.0
git add && git commit -m "chore: upgrade" && git push
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

## License

Apache 2.0
