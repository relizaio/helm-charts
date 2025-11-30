#!/bin/bash
#
# Local build script for Harbor Helm chart
# Run this BEFORE committing changes to generate the harbor-helm/ directory
#
# Usage:
#   ./build-local.sh                    # Build with default version (1.18.0)
#   ./build-local.sh 1.19.0             # Build with specific Harbor version
#   ./build-local.sh 1.19.0 2           # Build with version and iteration
#

set -e

HARBOR_VERSION="${1:-1.18.0}"
RELIZA_ITERATION="${2:-1}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "$SCRIPT_DIR"

echo "=============================================="
echo "Harbor Chart Local Build"
echo "=============================================="
echo "Harbor Version: $HARBOR_VERSION"
echo "Reliza Iteration: $RELIZA_ITERATION"
echo "Chart Version: ${HARBOR_VERSION}-reliza.${RELIZA_ITERATION}"
echo "=============================================="
echo ""

# Step 1: Add Harbor Helm repo (if not exists)
echo "Step 1/4: Adding Harbor Helm repo..."
helm repo add harbor https://helm.goharbor.io 2>/dev/null || true
helm repo update
echo "✅ Helm repo ready"
echo ""

# Step 2: Build and generate chart using Make
echo "Step 2/4: Building and generating chart (make setup)..."
make setup HARBOR_VERSION="$HARBOR_VERSION"
echo ""

# Step 3: Build chart dependencies
echo "Step 3/4: Building chart dependencies..."
cd harbor-helm
helm dependency build
cd ..
echo "✅ Dependencies built"
echo ""

# Step 4: Update chart version
echo "Step 4/4: Setting chart version..."
CHART_VERSION="${HARBOR_VERSION}-reliza.${RELIZA_ITERATION}"
sed -i "s/^version:.*/version: $CHART_VERSION/" harbor-helm/Chart.yaml
echo "✅ Chart version set to: $CHART_VERSION"
echo ""

# Validate
echo "=============================================="
echo "Validating chart..."
echo "=============================================="
helm lint harbor-helm

echo ""
echo "Testing template rendering..."
helm template test harbor-helm \
  --set expose.type=clusterIP \
  --set expose.tls.auto.commonName=harbor.local > /dev/null
echo "✅ Template renders successfully"

echo ""
echo "Verifying Reliza modifications..."
grep -q "Reliza customization" harbor-helm/templates/_helpers.tpl && echo "✅ _helpers.tpl modified"
grep -q "relizapostgresql" harbor-helm/values.yaml && echo "✅ values.yaml modified"

echo ""
echo "=============================================="
echo "✅ BUILD COMPLETE"
echo "=============================================="
echo ""
echo "Chart ready at: harbor-helm/"
echo "Version: $CHART_VERSION"
echo ""
echo "Next steps:"
echo "  1. Review changes: git diff harbor-helm/"
echo "  2. Stage files:    git add harbor-helm/"
echo "  3. Commit:         git commit -m 'chore: update Harbor chart to $CHART_VERSION'"
echo "  4. Push:           git push"
echo ""
