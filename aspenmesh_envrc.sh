#!/bin/bash

# Copyright 2023 F5 Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


cat <<EOF > ./.envrc
export TOOLS_REGISTRY_PROVIDER=gcr.io
export PROJECT_ID=f5-gcs-7056-ptg-aspenmesh-pub/tw-istio-testing
export TOOLS_REGISTRY_REPO=build-tools
export BUILD_TOOLS_ORG=F5-External
export AUTH_HEADER="Authorization: Bearer $(gcloud auth print-access-token)"
export ISTIO_ENVOY_BASE_URL=https://storage.googleapis.com/tw-istio-private-build/proxy
export ISTIO_ZTUNNEL_BASE_URL=https://storage.googleapis.com/tw-istio-private-build/ztunnel
export ISTIO_BASE_REGISTRY=gcr.io/istio-release
export BASE_VERSION=1.18-2023-07-21T13-00-12
EOF

