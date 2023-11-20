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
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/resource"
)

var (
	istioInst istio.Instance
	apps      map[string]*EchoDeployments
)

var Names = [4]string{"san-none", "san-uri", "san-dns", "san-multi"}

// TestMain defines the entrypoint for dualstack tests using a dualstack Istio installation.
// If a test requires a custom install it should go into its own package, otherwise it should go
// here to reuse a single install across tests.
func TestMain(m *testing.M) {
	apps = make(map[string]*EchoDeployments)
	for _, name := range Names {
		apps[name] = &EchoDeployments{}
	}
	framework.
		NewSuite(m).
		Setup(istio.Setup(&istioInst, setupConfig)).
		Setup(func(ctx resource.Context) error {
			return SetupApps(ctx, apps)
		}).
		Run()
}

func setupConfig(_ resource.Context, cfg *istio.Config) {
	if cfg == nil {
		return
	}
	cfg.ControlPlaneValues = `
values:
  meshConfig:
    defaultConfig:
      proxyMetadata:
        CERTIFICATE_CUSTOM_FIELDS: "true"
`
}
