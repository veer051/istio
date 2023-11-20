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
	"reflect"
	"sort"
	"testing"
)

func TestNewCustomFields(t *testing.T) {
	testCases := map[string]struct {
		saName       string
		podNamespace string
		customFields string
		expectErr    bool
	}{
		"Missing Service Account env": {
			saName:    "",
			expectErr: true,
		},
		"Missing pod Namespace env": {
			podNamespace: "",
			expectErr:    true,
		},
		"Has Service Account and pod Namespace but missing Certificate Custom Fields Annotation is empty": {
			saName:       "sleep",
			podNamespace: "default",
			customFields: "",
			expectErr:    false,
		},
		"Has Service Account and pod Namespace but missing Certificate Custom Fields Annotation env": {
			saName:       "sleep",
			podNamespace: "default",
			expectErr:    false,
		},
		"Has Service Account and pod Namespace but Certificate Custom Fields Annotation is not valid json": {
			saName:       "sleep",
			podNamespace: "default",
			customFields: "aspenmesh stuff",
			expectErr:    true,
		},
		"Certificate Custom Fields Annotation is valid format": {
			saName:       "sleep",
			podNamespace: "default",
			customFields: `{ "SAN": { "DNS": [ "foo.com" ] } }`,
			expectErr:    false,
		},
		"Certificate Custom Fields Annotation is valid format with uuid": {
			saName:       "sleep",
			podNamespace: "default",
			customFields: `{ "SAN": { "DNS": [ "foo.com" ], "URI": [ "1eb1f8fa-5607-4783-9a73-3e7630140833" ] } }`,
			expectErr:    false,
		},
	}

	for k, tc := range testCases {
		t.Run(k, func(t *testing.T) {
			_, err := NewCustomFields(tc.saName, tc.podNamespace, tc.customFields)
			if tc.expectErr {
				if err == nil {
					t.Errorf("%s expected an error, but did not receive one", k)
				}
			}
		})
	}
}

func TestGetSANDNSNames(t *testing.T) {
	testCases := map[string]struct {
		saName         string
		podNamespace   string
		customFields   string
		expectErr      bool
		expectedSANDNS []string
		expectedSANURI []string
	}{
		"Certificate Custom Fields Annotation is not set": {
			saName:         "httpbin",
			podNamespace:   "httpbin-ns",
			customFields:   "",
			expectErr:      false,
			expectedSANDNS: []string{"httpbin.httpbin-ns.svc.cluster.local"},
			expectedSANURI: []string{},
		},
		"Certificate Custom Fields Annotation is empty": {
			saName:         "sleep",
			podNamespace:   "sleep-ns",
			customFields:   `{ "foo": "bar" }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.sleep-ns.svc.cluster.local"},
			expectedSANURI: []string{},
		},
		"Certificate Custom Fields Annotation has no SAN field": {
			saName:         "sleep",
			podNamespace:   "sleep-ns",
			customFields:   "",
			expectErr:      false,
			expectedSANDNS: []string{"sleep.sleep-ns.svc.cluster.local"},
			expectedSANURI: []string{},
		},
		"Certificate Custom Fields Annotation has no DNS field": {
			saName:         "sleep",
			podNamespace:   "sleep-ns",
			customFields:   `{ "SAN": { "foo": [ "bar" ] } }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.sleep-ns.svc.cluster.local"},
			expectedSANURI: []string{},
		},
		"Certificate Custom Fields Annotation has a single DNS value": {
			saName:         "sleep",
			podNamespace:   "default",
			customFields:   `{ "SAN": { "DNS": [ "foo.com" ] } }`,
			expectErr:      false,
			expectedSANDNS: []string{"foo.com"},
			expectedSANURI: []string{},
		},
		"Certificate Custom Fields Annotation has multiple DNS values": {
			saName:       "sleep",
			podNamespace: "default",
			customFields: `{
						"SAN": {
							"DNS": [
								"foo.com",
								"my.aspenmesh.io",
								"zoo.gz"
							]
						}
			}`,
			expectErr:      false,
			expectedSANDNS: []string{"my.aspenmesh.io", "foo.com", "zoo.gz"},
			expectedSANURI: []string{},
		},
		"Certificate Custom Fields Annotation has a single URI value": {
			saName:         "sleep",
			podNamespace:   "default",
			customFields:   `{ "SAN": { "URI": [ "a94907d5-42b7-477e-96f6-81036e0bf989" ] } }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.default.svc.cluster.local"},
			expectedSANURI: []string{"uri://a94907d5-42b7-477e-96f6-81036e0bf989"},
		},
		"Certificate Custom Fields Annotation has no URI value": {
			saName:         "sleep",
			podNamespace:   "default",
			customFields:   `{ "SAN": { } }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.default.svc.cluster.local"},
			expectedSANURI: []string{},
		},
		"Certificate Custom Fields Annotation has a schema specified URI value": {
			saName:         "sleep",
			podNamespace:   "default",
			customFields:   `{ "SAN": { "URI": [ "uuid://a94907d5-42b7-477e-96f6-81036e0bf989" ] } }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.default.svc.cluster.local"},
			expectedSANURI: []string{"uri://uuid://a94907d5-42b7-477e-96f6-81036e0bf989"},
		},
		"Certificate Custom Fields Annotation has a bad URI value": {
			saName:         "sleep",
			podNamespace:   "default",
			customFields:   `{ "SAN": { "URI": [ "a94907d5-42b7-477e-96f6-81036e0bf989" ] } }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.default.svc.cluster.local"},
			expectedSANURI: []string{"uri://a94907d5-42b7-477e-96f6-81036e0bf989"},
		},
		"Certificate Custom Fields Annotation has a double spaced URI value": {
			saName:         "sleep",
			podNamespace:   "default",
			customFields:   `{ "SAN": { "URI": [ "uuid:  a94907d5-42b7-477e-96f6-81036e0bf989" ] } }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.default.svc.cluster.local"},
			expectedSANURI: []string{"uri://uuid:  a94907d5-42b7-477e-96f6-81036e0bf989"},
		},
		"Certificate Custom Fields Annotation has a one / in URI value": {
			saName:         "sleep",
			podNamespace:   "default",
			customFields:   `{ "SAN": { "URI": [ "uuid:/a94907d5-42b7-477e-96f6-81036e0bf989" ] } }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.default.svc.cluster.local"},
			expectedSANURI: []string{"uri://uuid:/a94907d5-42b7-477e-96f6-81036e0bf989"},
		},
		"Certificate Custom Fields Annotation has an improper proper UUID value": {
			saName:       "sleep",
			podNamespace: "default",
			// Removed '-' character from uuid used above between 477e and 96f6
			customFields:   `{ "SAN": { "URI": [ "uuid://a94907d5-42b7-477e96f6-81036e0bf989" ] } }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.default.svc.cluster.local"},
			expectedSANURI: []string{"uri://uuid://a94907d5-42b7-477e96f6-81036e0bf989"},
		},
		"Certificate Custom Fields Annotation has one improper proper and one proper UUID value": {
			saName:       "sleep",
			podNamespace: "default",
			// Removed '-' character from uuid used above between 477e and 96f6
			customFields:   `{ "SAN": { "URI": [ "uuid: 6017cc85-73c7-4fbc-9b6f-bede62d06300", "uuid://a94907d5-42b7-477e96f6-81036e0bf989" ] } }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.default.svc.cluster.local"},
			expectedSANURI: []string{"uri://uuid: 6017cc85-73c7-4fbc-9b6f-bede62d06300", "uri://uuid://a94907d5-42b7-477e96f6-81036e0bf989"},
		},
		"Certificate Custom Fields Annotation has multiple valid UUID values": {
			saName:       "sleep",
			podNamespace: "default",
			// Removed '-' character from uuid used above between 477e and 96f6
			customFields:   `{ "SAN": { "URI": [ "uuid: 6017cc85-73c7-4fbc-9b6f-bede62d06300", "uuid://a94907d5-42b7-477e-96f6-81036e0bf989" ] } }`,
			expectErr:      false,
			expectedSANDNS: []string{"sleep.default.svc.cluster.local"},
			expectedSANURI: []string{"uri://uuid: 6017cc85-73c7-4fbc-9b6f-bede62d06300", "uri://uuid://a94907d5-42b7-477e-96f6-81036e0bf989"},
		},
	}

	for k, tc := range testCases {
		t.Run(k, func(t *testing.T) {
			cF, err := NewCustomFields(tc.saName, tc.podNamespace, tc.customFields)
			if err != nil {
				t.Fatalf("%s had an unexpected error when CustomFields was constructed", k)
			}

			sort.Strings(tc.expectedSANDNS)

			sanDNS := cF.GetSANDNSNames()
			if !reflect.DeepEqual(sanDNS, tc.expectedSANDNS) {
				t.Errorf("%s expected SAN DNS names (%v) did not match returned SAN DNS names (%v)",
					k, tc.expectedSANDNS, sanDNS)
			}

			sanURI, err := cF.GetSANURINames()
			if err != nil {
				if !tc.expectErr {
					t.Errorf("%s expected no error but error was thrown: %s", k, err)
				}
			} else {
				if tc.expectErr {
					t.Errorf("%s did not error but error was expected", k)
				}
			}
			if !reflect.DeepEqual(sanURI, tc.expectedSANURI) {
				t.Errorf("%s expected SAN URI names (%v) did not match returned SAN URI names (%v)",
					k, tc.expectedSANURI, sanURI)
			}
		})
	}
}
