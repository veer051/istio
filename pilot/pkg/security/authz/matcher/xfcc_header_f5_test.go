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

package matcher

import (
	"testing"

	routepb "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

// TestXFCCHeaderMatcher tests to make sure the header matcher returns the correct
// filter matcher for Envoy
func TestXFCCHeaderMatcher(t *testing.T) {
	testCases := []struct {
		Name   string
		K      string
		V      string
		Expect *routepb.HeaderMatcher
	}{
		{
			Name: "exact match",
			K:    ":path",
			V:    "/productpage",
			Expect: &routepb.HeaderMatcher{
				Name: ":path",
				HeaderMatchSpecifier: &routepb.HeaderMatcher_ExactMatch{
					ExactMatch: "/productpage",
				},
			},
		},
		{
			Name: "contains match",
			K:    ":path",
			V:    "*/productpage*",
			Expect: &routepb.HeaderMatcher{
				Name: ":path",
				HeaderMatchSpecifier: &routepb.HeaderMatcher_ContainsMatch{
					ContainsMatch: "/productpage",
				},
			},
		},
		{
			Name: "suffix match",
			K:    ":path",
			V:    "*/productpage",
			Expect: &routepb.HeaderMatcher{
				Name: ":path",
				HeaderMatchSpecifier: &routepb.HeaderMatcher_SuffixMatch{
					SuffixMatch: "/productpage",
				},
			},
		},
		{
			Name: "prefix match",
			K:    ":path",
			V:    "/productpage*",
			Expect: &routepb.HeaderMatcher{
				Name: ":path",
				HeaderMatchSpecifier: &routepb.HeaderMatcher_PrefixMatch{
					PrefixMatch: "/productpage",
				},
			},
		},
		{
			Name: "presence match",
			K:    ":path",
			V:    "*",
			Expect: &routepb.HeaderMatcher{
				Name: ":path",
				HeaderMatchSpecifier: &routepb.HeaderMatcher_PresentMatch{
					PresentMatch: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		actual := F5XFCCHeaderMatcherWithContains(tc.K, tc.V)
		if !cmp.Equal(tc.Expect, actual, protocmp.Transform()) {
			t.Errorf("expecting %v, but got %v", tc.Expect, actual)
		}
	}
}
