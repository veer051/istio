// Copyright 2019 Istio Authors
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

package ca

import (
	"fmt"
	"sort"

	"istio.io/istio/security/pkg/pki/util"
)

func getCSRHosts(csrPem []byte) ([]string, error) {
	csr, err := util.ParsePemEncodedCSR(csrPem)
	if err != nil {
		return []string{}, fmt.Errorf("CSR is not valid (%v)", err)
	}

	hosts, err := util.ExtractIDs(csr.Extensions)
	if err != nil {
		return []string{}, fmt.Errorf("could not extract extension from CSR (%v)", err)
	}

	sort.Strings(hosts)
	return hosts, nil
}

func isIdentityInHosts(identities, hosts []string) bool {
	for _, id := range identities {
		for _, h := range hosts {
			if id == h {
				return true
			}
		}
	}
	return false
}

func addCSRHostsToIds(ids, hosts []string) []string {
	identities := append(ids, hosts...)

	deDupedIDs := make(map[string]interface{})
	for _, id := range identities {
		var present interface{}
		deDupedIDs[id] = present
	}

	res := []string{}
	for id := range deDupedIDs {
		res = append(res, id)
	}

	sort.Strings(res)
	return res
}
