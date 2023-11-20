//go:build integ
// +build integ

//  Copyright Aspen Mesh Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package customcerts

import (
	"crypto/x509"
	"fmt"
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/echo"
)

func TestInjectedCerts(t *testing.T) {
	framework.NewTest(t).
		Features("security.peer.trust-domain-validation").
		Run(func(ctx framework.TestContext) {
			for k, val := range apps {
				ctx.NewSubTestf("TestInjectedCerts/%s", k).
					Run(func(ctx framework.TestContext) {
						validateWorkloadCerts(t, ctx, val.EchoPod)
					})
			}
		})
}

func validateWorkloadCerts(t *testing.T, ctx framework.TestContext, echoInstance echo.Instance) {
	cert := getDefaultX509Cert(t, ctx, echoInstance)

	expectedDNS := svcAnnotationConfig[echoInstance.NamespacedName().Name].SAN.DNS
	if expectedDNS == nil {
		// If we don't have DNS specified, we expect the default svcacctname.namesspace.svc.cluster.local
		// our service account name is same as echo name
		expectedDNS = []string{fmt.Sprintf("%s.%s.svc.cluster.local",
			echoInstance.NamespacedName().Name, echoInstance.NamespacedName().Namespace.Name())}
	}
	foundDNS := containsAllDNSEntries(cert, expectedDNS)
	if !foundDNS {
		t.Fatalf("Did not find DNS entries in cert for %s", echoInstance.NamespacedName().Name)
	}

	expectedURI := svcAnnotationConfig[echoInstance.NamespacedName().Name].SAN.URI
	if expectedURI != nil {
		foundURI := containsAllURIEntries(cert, expectedURI)
		if !foundURI {
			t.Fatalf("Did not find URI entries in cert for %s", echoInstance.NamespacedName().Name)
		}
	}
}

func containsAllDNSEntries(certificate *x509.Certificate, expectedDNSNames []string) bool {
	expectedDNSEntries := len(expectedDNSNames)
	dnsEntries := 0

	for _, expectedDNSName := range expectedDNSNames {
		for _, dnsName := range certificate.DNSNames {
			if dnsName == expectedDNSName {
				dnsEntries++
			}
		}
	}

	return dnsEntries == expectedDNSEntries
}

func containsAllURIEntries(certificate *x509.Certificate, expectedURINames []string) bool {
	expectedURIEntries := len(expectedURINames)
	uriEntries := 0

	for _, expectedURIName := range expectedURINames {
		for _, uriName := range certificate.URIs {
			if uriName.String() == expectedURIName {
				uriEntries++
			}
		}
	}

	return uriEntries == expectedURIEntries
}
