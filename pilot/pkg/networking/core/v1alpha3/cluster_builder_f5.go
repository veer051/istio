// Copyright F5 Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha3

import (
	"strings"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	http "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"

	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/features"
)

func isCarrierGradeExternalIstioMutualServiceEntriesForceAutoSNI(opts *buildClusterOpts, tls *networking.ClientTLSSettings) bool {
	return features.CarrierGradeServiceEntryIstioMutual && opts.meshExternal && strings.HasPrefix(tls.GetSni(), "*.") &&
		(opts.port.Protocol.IsHTTP() || opts.port.Protocol.IsHTTP2())
}

func (*ClusterBuilder) changeAutoSniAndAutoSanValidation(mc *MutableCluster, setAutoSni, setAutoSanValidation bool) {
	if mc.httpProtocolOptions == nil {
		mc.httpProtocolOptions = &http.HttpProtocolOptions{}
	}
	if mc.httpProtocolOptions.UpstreamHttpProtocolOptions == nil {
		mc.httpProtocolOptions.UpstreamHttpProtocolOptions = &core.UpstreamHttpProtocolOptions{}
	}
	if setAutoSni {
		mc.httpProtocolOptions.UpstreamHttpProtocolOptions.AutoSni = true
	}
	if setAutoSanValidation {
		mc.httpProtocolOptions.UpstreamHttpProtocolOptions.AutoSanValidation = true
	}
}
