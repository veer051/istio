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
	"strings"

	routepb "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
)

// F5XFCCHeaderMatcherWithContains adds the contains match option for XFCC headers
func F5XFCCHeaderMatcherWithContains(k, v string) *routepb.HeaderMatcher {
	// We must check "*" first to make sure we'll generate a non empty value in the prefix/suffix case.
	// Empty prefix/suffix value is invalid in HeaderMatcher.
	if v == "*" {
		return &routepb.HeaderMatcher{
			Name: k,
			HeaderMatchSpecifier: &routepb.HeaderMatcher_PresentMatch{
				PresentMatch: true,
			},
		}
	} else if strings.HasPrefix(v, "*") && strings.HasSuffix(v, "*") {
		return &routepb.HeaderMatcher{
			Name: k,
			HeaderMatchSpecifier: &routepb.HeaderMatcher_ContainsMatch{
				ContainsMatch: v[1 : len(v)-1],
			},
		}
	} else if strings.HasPrefix(v, "*") {
		return &routepb.HeaderMatcher{
			Name: k,
			HeaderMatchSpecifier: &routepb.HeaderMatcher_SuffixMatch{
				SuffixMatch: v[1:],
			},
		}
	} else if strings.HasSuffix(v, "*") {
		return &routepb.HeaderMatcher{
			Name: k,
			HeaderMatchSpecifier: &routepb.HeaderMatcher_PrefixMatch{
				PrefixMatch: v[:len(v)-1],
			},
		}
	}
	return &routepb.HeaderMatcher{
		Name: k,
		HeaderMatchSpecifier: &routepb.HeaderMatcher_ExactMatch{
			ExactMatch: v,
		},
	}
}
