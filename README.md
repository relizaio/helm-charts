# This repository contains a set of public helm charts by Reliza


For all charts below use the following to add Reliza chart repository on your system:

```
helm repo add reliza https://registry.relizahub.com/chartrepo/library
helm repo update
```

## 1. ECR-Regcred Helm Chart

This chart is built to be used stand-alone or as a dependency which creates a regcred secret for AWS ECR or a regcred file for a regular container registry.

When working with ECR, AWS IAM API ID and Key are used to obtain a token which in turn is used to authenticate to the ECR registry. However, the token itself is valid for only 12 hours. That becomes problematic for CD process on Kubernetes. This chart automatically resolves the registry token based on IAM credentials an does so on startup and then refreshes every 11 hours.

When working with a standard container registry, we sometimes need a way to pass registry login and token directly to the chart. If this is the case, set *secret_type* value to *regular* (it is set to *ecr* by default).

Currently, this chart supports 3 modes configurable via create_secret_in_chart property in values.
- "none" (default) - You would need to provide `reliza-ecr-regcred` secret yourself, the secret needs to contain 2 data entries: AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY.
- "regular" - would create a secret based on IAM credentials supplied in plain text via values `aws_id` and `aws_key`. IMPORTANT: this is only recommended for local testing.
- "sealed" - leverages [Bitnami Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets). Expects `aws_id` and `aws_key` to be populated with sealed secrets. Also allows sealed_secrets_scope configuration to be either "namespace-wide" (default) or "cluster-wide" or "strict" (see [here](https://github.com/bitnami-labs/sealed-secrets#scopes) for more details).

Most basic standalone installation for ECR:

```
helm install ecr-regcred --set create_secret_in_chart=regular --set aws_id=YOUR_AWS_IAM_ID --set aws_key=YOUR_AWS_IAM_KEY -n default reliza/ecr-regcred
```

Most basic standalone installation for a Regular Registry:

```
helm install ecr-regcred --set secret_type=regular --set create_secret_in_chart=regular --set aws_id=YOUR_AWS_IAM_ID --set aws_key=YOUR_AWS_IAM_KEY -n default reliza/ecr-regcred
```


To add to your chart as dependency:

Add this chart to dependencies section of your Chart.yaml as following:

```
dependencies:
  - name: ecr-regcred
    repository: "https://registry.relizahub.com/chartrepo/library"
    version: ">=0.0.3"
```

After that run `helm dependency update`

Then use same properties as above to configure but under `ecr-regcred` umbrella yaml key. Refer to values.yaml in this chart for full list of settings.


## 2. Reliza Watcher Helm Chart
Create your Instance on [Reliza Hub](https://relizahub.com) and obtain Instance API ID and API Key.

Basic installation to monitor all namespaces (creates secret in chart):

```
kubectl create ns reliza-watcher
helm install reliza-watcher -n reliza-watcher --set create_secret_in_chart=regular --set relizaApiId=actual_reliza_api_id --set relizaApiKey=actual_reliza_api_key reliza/reliza-watcher
```

Currently, this chart supports 3 modes for the `reliza-watcher` secret containing Reliza Hub credentials configurable via `create_secret_in_chart` property in values.
- "none" (default) - You would need to provide `reliza-watcher` secret yourself in the namespace where reliza-watcher chart is deployed, the secret needs to contain 2 data entries: reliza-api-id and reliza-api-key.
- "regular" - would create a secret supplied in plain text via values `relizaApiId` and `relizaApiKey` (as in the Basic Installation example above). IMPORTANT: this is only recommended for local testing.
- "sealed" - leverages [Bitnami Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets). Expects `relizaApiId` and `relizaApiKey` to be populated with sealed secrets. Also allows sealed_secrets_scope configuration to be either "namespace-wide" (default) or "cluster-wide" or "strict" (see [here](https://github.com/bitnami-labs/sealed-secrets#scopes) for more details).

**Notes:**

1. If you would like to watch only specific namespaces, say *default* and *myappnamespace*, leverage the `namespace` property when installing helmchart (provide comma-separated namespaces as a value for namespace key) - as shown below:

```
helm install reliza-watcher -n reliza-watcher --set namespace="default\,myappnamespace" reliza/reliza-watcher
```

2. Sender id can be set to any string via *sender* property. Data from different senders will be combined on the Reliza Hub in the instance view.

### Install Using Helm In A Multi-Namespace Kubernetes Clusters

When you wish to watch different instances deployed in different namespaces of a kubernetes cluster, multiple instances of reliza-watcher are required and can be deployed as follows:

Assume you have two instances *instance-A* and *instance-B* on [Reliza Hub](https://relizahub.com) deployed in namespaces *ns-A*  and *ns-B* respectively. Then:

1. Obtain Instance API ID and API Key for the instances *instance-A* and *instance-B* from [Reliza Hub](https://relizahub.com).
2. Issue following commands replacing <RELIZA_API_ID_FOR_INSTANCE_A> and <RELIZA_API_KEY_INSTANCE_A> with values obtained from Reliza Hub:
```
kubectl create secret generic reliza-watcher -n <ns-A> --from-literal=reliza-api-id=<RELIZA_API_ID_FOR_INSTANCE_A> --from-literal=reliza-api-key=<RELIZA_API_KEY_INSTANCE_A>
helm install reliza-watcher -n <ns-A> --set namespace="ns-A" reliza/reliza-watcher
```
3. Issue following commands replacing <RELIZA_API_ID_FOR_INSTANCE_B> and <RELIZA_API_KEY_INSTANCE_B> with values obtained from Reliza Hub:
```
kubectl create secret generic reliza-watcher -n <ns-B> --from-literal=reliza-api-id=<RELIZA_API_ID_FOR_INSTANCE_B> --from-literal=reliza-api-key=<RELIZA_API_KEY_INSTANCE_B>
helm install reliza-watcher -n <ns-B> --set namespace="ns-B" reliza/reliza-watcher
```

Note that this last example also shows reliza-watcher secret created outside of the helm chart.

### Use Reliza Watcher as a Dependency

Note that this approach is convenient for testing, but we do not recommend it for major persistent environments, especially for Production. That is because it binds watcher to the application helm chart, so if such chart is getting undeployed, the Watcher is also getting undeployed - and observability is lost.

To add Reliza Watcher to your chart as a dependency:

Add this chart to dependencies section of your Chart.yaml as following:

```
dependencies:
  - name: reliza-watcher
    repository: "https://registry.relizahub.com/chartrepo/library"
    version: ">=0.0.0"
```

After that run `helm dependency update`

Then use same properties as above to configure but under `reliza-watcher` umbrella yaml key. Refer to values.yaml in this chart for full list of settings.

Sample:

```
reliza-watcher:
  enabled: true
  namespace: relizahub
  create_secret_in_chart: regular
  relizaApiId: placeholder_id
  relizaApiKey: placeholder_key
```

Also we recommend to initialize it as disabled in root values.yaml, and then customize per environments, like the following:

```
reliza-watcher:
  enabled: true
```