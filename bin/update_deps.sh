#!/bin/bash

# Copyright 2019 Istio Authors
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

set -exo pipefail

# shellcheck disable=SC1001
LATEST_DEB11_DISTROLESS_SHA256=$(crane digest gcr.io/distroless/static-debian11 | awk -F\: '{print $2}')
sed -i -E "s/sha256:[a-z0-9]+/sha256:${LATEST_DEB11_DISTROLESS_SHA256}/g" docker/Dockerfile.distroless

# shellcheck disable=SC1001
LATEST_IPTABLES_DISTROLESS_SHA256=$(crane digest gcr.io/istio-release/iptables | awk -F\: '{print $2}')
sed -i -E "s/sha256:[a-z0-9]+/sha256:${LATEST_IPTABLES_DISTROLESS_SHA256}/g" pilot/docker/Dockerfile.proxyv2
sed -i -E "s/sha256:[a-z0-9]+/sha256:${LATEST_IPTABLES_DISTROLESS_SHA256}/g" pilot/docker/Dockerfile.ztunnel
