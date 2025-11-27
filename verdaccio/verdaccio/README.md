# Verdaccio with Reliza Quarantine Filter Helm Chart

This Helm chart deploys [Verdaccio with Reliza Quarantine Filter](https://github.com/relizaio/dockerfile-collection/tree/main/verdaccio), a lightweight private npm proxy registry with built-in package quarantine capabilities, on a Kubernetes cluster.

## What is Reliza Quarantine Filter?

The Reliza Quarantine Filter is a security enhancement for Verdaccio that implements a configurable quarantine period for newly published npm packages. When a package is published to the upstream registry (npmjs.com), it is held in quarantine for a specified number of days before being made available through this proxy. This provides protection against supply chain attacks by giving the community time to identify and report malicious packages before they can be consumed by your projects.

Key benefits:
- **Supply chain security**: Delays availability of new packages, reducing exposure to malicious package attacks
- **Configurable quarantine period**: Set the number of days packages must age before being available (default: 7 days)
- **Transparent operation**: Works as a drop-in replacement for standard Verdaccio

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure (if persistence is enabled)

## Installing the Chart

To install the chart with the release name `my-verdaccio`:

```bash
helm install my-verdaccio ./verdaccio
```

The command deploys Verdaccio on the Kubernetes cluster in the default configuration. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `my-verdaccio` deployment:

```bash
helm delete my-verdaccio
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Parameters

### Global parameters

| Name                      | Description                                     | Value |
| ------------------------- | ----------------------------------------------- | ----- |
| `global.imageRegistry`    | Global Docker image registry                    | `""`  |
| `global.imagePullSecrets` | Global Docker registry secret names as an array| `[]`  |

### Common parameters

| Name                     | Description                                                                             | Value           |
| ------------------------ | --------------------------------------------------------------------------------------- | --------------- |
| `nameOverride`           | String to partially override verdaccio.fullname                                        | `""`            |
| `fullnameOverride`       | String to fully override verdaccio.fullname                                            | `""`            |

### Verdaccio Image parameters

| Name                | Description                                          | Value                                                                                      |
| ------------------- | ---------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| `image`             | Verdaccio with Reliza Quarantine Filter image        | `registry.relizahub.com/library/reliza-verdaccio@sha256:...`                              |
| `imagePullSecrets`  | Docker registry secret names as an array             | `[]`                                                                                       |

### Quarantine Parameters

| Name                | Description                                          | Value                                                                                      |
| ------------------- | ---------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| `quarantineDays`    | Number of days packages must age before being available through the proxy. This is the quarantine period for supply chain security. | `7`                                                                                        |

### Deployment parameters

| Name                                    | Description                                                                               | Value   |
| --------------------------------------- | ----------------------------------------------------------------------------------------- | ------- |
| `replicaCount`                          | Number of Verdaccio replicas to deploy                                                   | `1`     |
| `podAnnotations`                        | Annotations for Verdaccio pods                                                            | `{}`    |
| `podSecurityContext`                    | Set Verdaccio pod's Security Context                                                      | `{}`    |
| `securityContext`                       | Set Verdaccio container's Security Context                                                | `{}`    |

### Traffic Exposure Parameters

| Name                        | Description                                                                                                                      | Value                    |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | ------------------------ |
| `service.type`              | Verdaccio service type                                                                                                           | `ClusterIP`              |
| `service.port`              | Verdaccio service HTTP port                                                                                                      | `4873`                   |
| `ingress.enabled`           | Enable ingress record generation for Verdaccio                                                                                  | `true`                   |
| `ingress.host`              | Default host for the ingress record                                                                                              | `verdaccio.localhost`    |

### Persistence Parameters

| Name                        | Description                                                                                                                      | Value                    |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | ------------------------ |
| `persistence.enabled`       | Enable persistence using Persistent Volume Claims                                                                               | `true`                   |
| `persistence.accessMode`    | Persistent Volume access mode                                                                                                    | `ReadWriteOnce`          |
| `persistence.size`          | Persistent Volume size                                                                                                           | `8Gi`                    |
| `persistence.storageClass`  | Persistent Volume storage class                                                                                                  | `""`                     |

### Traefik Parameters

| Name                        | Description                                                                                                                      | Value                    |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | ------------------------ |
| `useTraefikLe`              | Use Traefik with Let's Encrypt                                                                                                  | `true`                   |
| `traefikBehindLb`           | Use Traefik behind load balancer                                                                                                | `false`                  |
| `ingressHost`               | Ingress host for Traefik                                                                                                        | `verdaccio.localhost`    |

## Configuration and installation details

### Quarantine Period

The `quarantineDays` parameter controls how long newly published packages on npmjs.com must wait before being available through this proxy. For example, with the default value of `7`, a package published to npmjs.com today will only become available through this Verdaccio instance after 7 days.

To adjust the quarantine period:

```bash
helm install my-verdaccio ./verdaccio --set quarantineDays=14
```

### Using npm with Verdaccio

Once Verdaccio is deployed, you can configure npm to use it:

```bash
# Set registry
npm set registry http://verdaccio.localhost

# Add user
npm adduser --registry http://verdaccio.localhost

# Publish package
npm publish --registry http://verdaccio.localhost
```

### Persistence

The chart mounts a Persistent Volume at the `/verdaccio/storage` path. The volume stores the npm packages and metadata.

### Security

The chart runs Verdaccio as a non-root user (UID 10001) for security purposes.

## More Information

For more details about the Reliza Quarantine Filter and the Docker image used by this chart, see:
- [Verdaccio with Reliza Quarantine Filter Dockerfile](https://github.com/relizaio/dockerfile-collection/tree/main/verdaccio)
