# Verdaccio Helm Chart

This Helm chart deploys Verdaccio, a lightweight private npm proxy registry, on a Kubernetes cluster.

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
| `image.repository`  | Verdaccio image repository                           | `registry.relizahub.com/library/reliza-verdaccio`                                         |
| `image.tag`         | Verdaccio image tag (immutable tags are recommended)| `sha256:5d7ea34f17ce57cb0e0c8bbb6a1adfdecc19e49f85d2a059c81f61dc606f8724`              |
| `image.pullPolicy`  | Verdaccio image pull policy                          | `IfNotPresent`                                                                             |

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
