// Copyright Istio Authors
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
	"reflect"
	"sort"
	"testing"

	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pkg/test"
)

func TestDefaultCarrierGradeCipherSuite(t *testing.T) {
	test.SetForTest(t, &features.CarrierGradeCipherSuite, true)

	testCases := []struct {
		name            string
		server          *networking.Server
		expectedCiphers []string
	}{
		{
			name: "destination rule does not set any",
			server: &networking.Server{
				Tls: &networking.ServerTLSSettings{},
			},
			expectedCiphers: []string{
				"ECDHE-ECDSA-AES256-GCM-SHA384",
				"ECDHE-ECDSA-CHACHA20-POLY1305",
				"ECDHE-RSA-AES256-GCM-SHA384",
				"ECDHE-RSA-CHACHA20-POLY1305",
				"ECDHE-ECDSA-AES128-GCM-SHA256",
				"ECDHE-RSA-AES128-GCM-SHA256",
			},
		},
		{
			name: "destination rule sets an unsupported carrier grade cipher",
			server: &networking.Server{
				Tls: &networking.ServerTLSSettings{
					CipherSuites: []string{
						"ECDHE-ECDSA-AES256-GCM-SHA384",
						"AES128-SHA",
					},
				},
			},
			expectedCiphers: []string{
				"ECDHE-ECDSA-AES256-GCM-SHA384",
			},
		},
	}
	// for Envoy v1.22 the default ciphers used when not specified are
	// https://www.envoyproxy.io/docs/envoy/v1.22.11/api-v3/extensions/transport_sockets/tls/v3/common.proto.html?highlight=cipher
	// this has changed in https://www.envoyproxy.io/docs/envoy/v1.26.0/api-v3/extensions/transport_sockets/tls/v3/common.proto.html
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			suites := filteredGatewayCipherSuites(tc.server)
			sort.Strings(suites)
			sort.Strings(tc.expectedCiphers)
			if !reflect.DeepEqual(suites, tc.expectedCiphers) {
				t.Errorf("got cipher suites (%v) but expected (%v)", suites, tc.expectedCiphers)
			}
		})
	}
}

func TestDefaultCipherSuite(t *testing.T) {
	testCases := []struct {
		name            string
		server          *networking.Server
		expectedCiphers []string
	}{
		{
			name: "destination rule does not set any",
			server: &networking.Server{
				Tls: &networking.ServerTLSSettings{},
			},
			expectedCiphers: []string{},
		},
	}

	// the default ciphers for Istio 1.18 can be found at https://www.envoyproxy.io/docs/envoy/v1.26.0/api-v3/extensions/transport_sockets/tls/v3/common.proto.html
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			suites := filteredGatewayCipherSuites(tc.server)
			if !reflect.DeepEqual(suites, tc.expectedCiphers) {
				t.Errorf("got cipher suites (%v) but expected (%v)", suites, tc.expectedCiphers)
			}
		})
	}
}
