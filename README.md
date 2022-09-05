# This repository contains a set of public helm charts by Reliza

1. ecr-regcred

This chart is built to be used stand-alone or as a dependency which creates a regcred secret for AWS ECR.

When working with ECR, AWS IAM API ID and Key are used to obtain a token which in turn is used to authenticate to the ECR registry. However, the token itself is valid for only 12 hours. That becomes problematic for CD process on Kubernetes. This chart automatically resolves the registry token based on IAM credentials an does so on startup and then refreshes every 11 hours.

Currently, this chart supports 3 modes configurable via create_secret_in_chart property in values.
- "none" (default) - You would need to provide `reliza-ecr-regcred` secret yourself, the secret needs to contain 2 data entries: AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY.
- "regular" - would create a secret based on IAM credentials supplied in plain text and base64-d via values `aws_id` and `aws_key`. IMPORTANT: this is only recommended for local testing.
- "sealed" - leverages [Bitnami Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets). Expects `aws_id` and `aws_key` to be populated with sealed secrets. Also allows sealed_secrets_scope configuration to be either "namespace-wide" (default) or "cluster-wide" or "strict" (see [here](https://github.com/bitnami-labs/sealed-secrets#scopes) for more details).