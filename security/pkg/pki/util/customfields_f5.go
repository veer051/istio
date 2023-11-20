// Copyright Aspen Mesh Authors
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

package util

import (
	"encoding/json"
	"fmt"
	"sort"
)

const uriPrefix = "uri://"

// AspenMeshCertCustomFields is used to deserialize annotations from
// the ServiceAccount to create the proper Secret name when x509 extensions
// with SAN DNS is enabled
type SAN struct {
	DNS []string `json:"DNS,omitempty"`
	URI []string `json:"URI,omitempty"`
}

type CertificateAspenMeshIOCustomFields struct {
	SAN *SAN `json:"SAN,omitempty"`
}

type CustomFields struct {
	cF             *CertificateAspenMeshIOCustomFields
	namespace      string
	serviceAccount string
}

func NewCustomFields(saName, podNamespace, customFields string) (*CustomFields, error) {
	if saName == "" {
		return nil, fmt.Errorf("SERVICE_ACCOUNT env not defined")
	}
	if podNamespace == "" {
		return nil, fmt.Errorf("POD_NAMESPACE env not defined")
	}

	aspenMeshCertCustomFields := &CertificateAspenMeshIOCustomFields{}
	if customFields != "" {
		err := json.Unmarshal([]byte(customFields), aspenMeshCertCustomFields)
		if err != nil {
			return nil, fmt.Errorf("CERTIFICATE_CUSTOM_FIELDS env is not valid json: %v", err)
		}
	}

	return &CustomFields{
		serviceAccount: saName,
		namespace:      podNamespace,
		cF:             aspenMeshCertCustomFields,
	}, nil
}

func (c *CustomFields) GetSANDNSNames() []string {
	var fqdns []string

	if c.cF != nil && c.cF.SAN != nil {
		fqdns = append(fqdns, c.cF.SAN.DNS...)
	}

	if len(fqdns) == 0 {
		fqdns = []string{fmt.Sprintf("%s.%s.svc.cluster.local", c.serviceAccount, c.namespace)}
	}

	sort.Strings(fqdns)
	return fqdns
}

func (c *CustomFields) GetSANURINames() ([]string, error) {
	urisSlice := []string{}
	uris := []string{}

	if c.cF != nil && c.cF.SAN != nil {
		uris = append(uris, c.cF.SAN.URI...)
	}

	for _, uri := range uris {
		uriString := uriPrefix + uri
		urisSlice = append(urisSlice, uriString)
	}

	sort.Strings(urisSlice)
	return urisSlice, nil
}
