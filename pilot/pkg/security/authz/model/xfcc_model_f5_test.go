// Copyright © 2022, F5 Networks, Inc. All rights reserved.
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

package model

import (
	"strings"
	"testing"

	rbacpb "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"

	authzpb "istio.io/api/security/v1beta1"
	"istio.io/istio/pkg/util/protomarshal"
)

// TestXFCCModel_Generate tests that the Envoy filters are being generated properly
// for XFCC headers.
func TestXFCCModel_Generate(t *testing.T) {
	annotations := make(map[string]string)
	annotations["authz.contains.aspenmesh.io/xfcc"] = "true"
	rule := yamlRule(t, `
from:
- source:
    requestPrincipals: ["td-1/ns/foo/sa/sleep-1"]
    notRequestPrincipals: ["td-1/ns/foo/sa/sleep-2"]
- source:
    requestPrincipals: ["td-1/ns/foo/sa/sleep-3"]
    notRequestPrincipals: ["td-1/ns/foo/sa/sleep-4"]
to:
- operation:
    ports: ["8001"]
    notPorts: ["8002"]
- operation:
    ports: ["8003"]
    notPorts: ["8004"]
when:
- key: request.headers["X-Forwarded-Client-Cert"]
  values: ["*URI=com.example.nfType:smf*"]
  notValues: ["*DNS=sleep.example.com*"]
- key: request.headers["X-FORWARDED-CLIENT-CERT"]
  values: ["*DNS=httpbin.example.com*"]
- key: request.headers["x-forwarded-client-cert"]
  notValues: ["*URI=com.example.productpage*"]
`)

	cases := []struct {
		name    string
		forTCP  bool
		action  rbacpb.RBAC_Action
		rule    *authzpb.Rule
		want    []string
		notWant []string
	}{
		{
			name:   "allow-http",
			action: rbacpb.RBAC_ALLOW,
			rule:   rule,
			want: []string{
				"URI=com.example.nfType:smf",
				"DNS=sleep.example.com",
				"DNS=httpbin.example.com",
				"URI=com.example.productpage",
				"containsMatch",
			},
			notWant: []string{
				"URI=com.example.nfType:smf*",
				"DNS=sleep.example.com*",
				"DNS=httpbin.example.com*",
				"URI=com.example.productpage*",
				"suffixMatch",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := New(tc.rule, annotations)
			if err != nil {
				t.Fatal(err)
			}
			p, _ := m.Generate(tc.forTCP, false, tc.action)
			var gotYaml string
			if p != nil {
				if gotYaml, err = protomarshal.ToYAML(p); err != nil {
					t.Fatalf("%s: failed to parse yaml: %s", tc.name, err)
				}
			}

			for _, want := range tc.want {
				if !strings.Contains(gotYaml, want) {
					t.Errorf("got:\n%s but not found %s", gotYaml, want)
				}
			}
			for _, notWant := range tc.notWant {
				if strings.Contains(gotYaml, notWant) {
					t.Errorf("got:\n%s but not want %s", gotYaml, notWant)
				}
			}
		})
	}
}
