//go:build integ
// +build integ

//  Copyright Aspen Mesh Authors
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

package customcerts

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"testing"

	envoyAdmin "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	auth "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/istioctl"
	"istio.io/istio/pkg/util/protomarshal"
)

func getSecretsConfigDumpFromInstance(t *testing.T, ctx framework.TestContext, inst echo.Instance) *envoyAdmin.SecretsConfigDump {
	istioCtl := istioctl.NewOrFail(t, ctx, istioctl.Config{})

	workload := inst.WorkloadsOrFail(t)[0]
	args := []string{
		"pc", "secret", fmt.Sprintf("%s.%s", workload.PodName(), inst.NamespacedName().Namespace.Name()), "-o", "json",
	}
	output, _ := istioCtl.InvokeOrFail(t, args)

	var jsonOutput interface{}
	if err := json.Unmarshal([]byte(output), &jsonOutput); err != nil {
		t.Fatalf("Could not unmarshal response %s", output)
	}

	dump := &envoyAdmin.SecretsConfigDump{}
	if err := protomarshal.Unmarshal([]byte(output), dump); err != nil {
		t.Fatal(err)
	}
	return dump
}

func getDefaultX509Cert(t *testing.T, ctx framework.TestContext, inst echo.Instance) *x509.Certificate {
	secretDump := getSecretsConfigDumpFromInstance(t, ctx, inst)

	secretTyped := &auth.Secret{}
	err := secretDump.DynamicActiveSecrets[0].GetSecret().UnmarshalTo(secretTyped)
	if err != nil {
		t.Fatalf("Failed to unmarshal default secret for instance %s", inst.NamespacedName().Name)
	}
	block, _ := pem.Decode(secretTyped.GetTlsCertificate().GetCertificateChain().GetInlineBytes())
	if block == nil {
		t.Fatalf("Failed to decode default secret PEM for instance %s", inst.NamespacedName().Name)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse x509 cert for instance %s", inst.NamespacedName().Name)
	}
	return cert
}
