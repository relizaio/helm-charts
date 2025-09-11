# Reliza Common Helm Chart

A Library Helm Chart for grouping common logic between Reliza charts. This chart is not deployable by itself.

## Description

This chart provides a collection of useful Helm template helpers that can be used across multiple Reliza Helm charts. It includes functions for:

- **Names and Labels**: Standardized naming conventions and label generation
- **Images**: Image reference and pull policy helpers  
- **Resources**: Resource preset and limit helpers
- **Storage**: Storage class and volume helpers
- **Networking**: Service and ingress helpers
- **Security**: Security context and capability helpers
- **Validation**: Input validation and error handling
- **Utilities**: Common utility functions

## Usage

### As a Chart Dependency

Add this chart as a dependency in your `Chart.yaml`:

```yaml
dependencies:
- name: common
  repository: oci://registry.relizahub.com/charts
  version: "^1.0.0"
```

### Template Functions

The chart provides numerous template functions under the `common.*` namespace:

```yaml
# Names and namespaces
{{ include "common.names.name" . }}
{{ include "common.names.fullname" . }}
{{ include "common.names.namespace" . }}

# Labels
{{ include "common.labels.standard" . }}
{{ include "common.labels.matchLabels" . }}

# Images
{{ include "common.images.image" . }}
{{ include "common.images.pullSecrets" . }}

# Resources
{{ include "common.resources.preset" . }}

# And many more...
```

## Template Categories

### Core Templates
- `_names.tpl` - Name generation functions
- `_labels.tpl` - Label standardization
- `_images.tpl` - Image reference helpers
- `_tplvalues.tpl` - Template value rendering

### Advanced Templates  
- `_affinities.tpl` - Pod affinity/anti-affinity
- `_capabilities.tpl` - Kubernetes API capabilities
- `_resources.tpl` - Resource presets and limits
- `_storage.tpl` - Storage class helpers
- `_secrets.tpl` - Secret management
- `_ingress.tpl` - Ingress configuration

### Validation Templates
- `validations/` - Input validation functions for various services

## License

Copyright Reliza Incorporated. All Rights Reserved.
Licensed under the MIT License.

## Maintainers

- [Reliza Incorporated](https://github.com/relizaio)

## Sources

- [GitHub Repository](https://github.com/relizaio/helm-charts)
- [Documentation](https://github.com/relizaio/dockerfile-collection/tree/main/helm-charts)
