# Istio

This repository is a "pseudo fork" of the open source [Istio](https://github.com/istio/istio) repository. As a pseudo fork, changes from the open source repository are regularly pulled here but changes here are not directly pushed back to the upstream open source repository. This allows F5 to maintain our own feature set privately for the benefit of our customers.

<!-- markdownlint-disable-next-line MD026 -->
## README.md is customer facing!

Do not make changes to [./README.md](./README.md) that are not customer facing. README.md is delivered to customers via release and extras tarballs. Internal notes should be kept in this file instead.

## Environment Variables for Building Locally

Since we now build our own build-tools, the following environment variables need to be exported to pull the correct image.

```script
export TOOLS_REGISTRY_PROVIDER=gcr.io
export PROJECT_ID=f5-gcs-7056-ptg-aspenmesh-pub/tw-istio-testing
export TOOLS_REGISTRY_REPO=build-tools
export BUILD_TOOLS_ORG=F5-External
```
