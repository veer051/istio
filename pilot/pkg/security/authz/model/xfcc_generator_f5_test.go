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
	"testing"

	rbacpb "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"istio.io/istio/pkg/util/protomarshal"
)

// TestXFCCRequestHeaderGenerator tests that the XFCC request header generator
// returns the correct header for XFCC headers
func TestXFCCRequestHeaderGenerator(t *testing.T) {
	cases := []struct {
		name   string
		g      generator
		key    string
		value  string
		forTCP bool
		want   interface{}
	}{
		{
			name:  "requestF5XFCCGenerator",
			g:     requestF5XFCCGenerator{},
			key:   "request.headers[x-forwarded-client-cert]",
			value: "*foo*",
			want: yamlPrincipal(t, `
         header:
          containsMatch: foo
          name: x-forwarded-client-cert`),
		},
		{
			name:  "requestF5XFCCGenerator",
			g:     requestF5XFCCGenerator{},
			key:   "request.headers[x-forwarded-client-cert]",
			value: "*foo",
			want: yamlPrincipal(t, `
         header:
          suffixMatch: foo
          name: x-forwarded-client-cert`),
		},
		{
			name:  "requestF5XFCCGenerator",
			g:     requestF5XFCCGenerator{},
			key:   "request.headers[x-forwarded-client-cert]",
			value: "foo*",
			want: yamlPrincipal(t, `
         header:
          prefixMatch: foo
          name: x-forwarded-client-cert`),
		},
		{
			name:  "requestF5XFCCGenerator",
			g:     requestF5XFCCGenerator{},
			key:   "request.headers[x-forwarded-client-cert]",
			value: "*",
			want: yamlPrincipal(t, `
         header:
          presentMatch: true
          name: x-forwarded-client-cert`),
		},
		{
			name:  "requestF5XFCCGenerator",
			g:     requestF5XFCCGenerator{},
			key:   "request.headers[x-forwarded-client-cert]",
			value: "foo",
			want: yamlPrincipal(t, `
         header:
          exactMatch: foo
          name: x-forwarded-client-cert`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got interface{}
			var err error
			// nolint: gocritic
			if _, ok := tc.want.(*rbacpb.Permission); ok {
				got, err = tc.g.permission(tc.key, tc.value, tc.forTCP)
				if err != nil {
					t.Errorf("permission returned error")
				}
			} else if _, ok := tc.want.(*rbacpb.Principal); ok {
				got, err = tc.g.principal(tc.key, tc.value, false, tc.forTCP)
				if err != nil {
					t.Errorf("principal returned error")
				}
			} else {
				_, err1 := tc.g.principal(tc.key, tc.value, false, tc.forTCP)
				_, err2 := tc.g.permission(tc.key, tc.value, tc.forTCP)
				if err1 == nil || err2 == nil {
					t.Fatalf("wanted error")
				}
				return
			}
			if diff := cmp.Diff(got, tc.want, protocmp.Transform()); diff != "" {
				var gotYaml string
				gotProto, ok := got.(proto.Message)
				if !ok {
					t.Fatal("failed to extract proto")
				}
				if gotYaml, err = protomarshal.ToYAML(gotProto); err != nil {
					t.Fatalf("%s: failed to parse yaml: %s", tc.name, err)
				}
				t.Errorf("got:\n %v\n but want:\n %v", gotYaml, tc.want)
			}
		})
	}
}
