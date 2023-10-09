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
	"fmt"
	"strings"

	rbacpb "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"

	"istio.io/istio/pilot/pkg/security/authz/matcher"
)

// requestF5XFCCGenerator is a F5 request header generator for XFCC headers only
type requestF5XFCCGenerator struct{}

func (requestF5XFCCGenerator) permission(_, _ string, _ bool) (*rbacpb.Permission, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (requestF5XFCCGenerator) principal(key, value string, forTCP, _ bool) (*rbacpb.Principal, error) {
	if forTCP {
		return nil, fmt.Errorf("%q is HTTP only", key)
	}

	header, err := extractNameInBrackets(strings.TrimPrefix(key, attrRequestHeader))
	if err != nil {
		return nil, err
	}
	m := matcher.F5XFCCHeaderMatcherWithContains(header, value)
	return principalHeader(m), nil
}
