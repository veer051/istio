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

package inject

import (
	"fmt"
	"os"
	"testing"

	meshapi "istio.io/api/mesh/v1alpha1"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/serviceregistry/provider"
	"istio.io/istio/pkg/cluster"
	"istio.io/istio/pkg/config/mesh"
	"istio.io/istio/pkg/util/sets"
)

// TestInjection_CustomCertFields is a modified version of TestInjection that uses data in
// testdata/customFields to test proper injection of customFields annotation
func TestInjection_CustomCertFields(t *testing.T) {
	type testCase struct {
		in            string
		want          string
		setFlags      []string
		inFilePath    string
		mesh          func(m *meshapi.MeshConfig)
		skipWebhook   bool
		expectedError string
		setup         func()
		teardown      func()

		// svcAnnotations added field to test for injected annotations from the associated
		// ServiceAccount
		svcAnnotations map[string]string
	}

	cases := []testCase{
		{
			in:   "customcertfields.yaml",
			want: "customcertfields.yaml.injected",
			svcAnnotations: map[string]string{
				"certificate.aspenmesh.io/customFields": "{ 'SAN': { 'DNS': [ 'foo.com' ] } }",
			},
		},
	}

	// Keep track of tests we add options above
	// We will search for all test files and skip these ones
	alreadyTested := sets.New[string]()

	for _, t := range cases {
		if t.want != "" {
			alreadyTested.Insert(t.want)
		} else {
			alreadyTested.Insert(t.in + ".injected")
		}
	}

	// Automatically add any other test files in the folder. This ensures we don't
	// forget to add to this list, that we don't have duplicates, etc
	// Keep track of all golden files so we can ensure we don't have unused ones later
	allOutputFiles := sets.New[string]()

	// Preload default settings. Computation here is expensive, so this speeds the tests up substantially
	writeInjectionSettings(t, "default", nil, "")
	defaultTemplate, defaultValues, defaultMesh := readInjectionSettings(t, "default")

	for i, c := range cases {
		c := c
		testName := fmt.Sprintf("[%02d] %s", i, c.want)
		if c.expectedError != "" {
			testName = fmt.Sprintf("[%02d] %s", i, c.in)
		}
		t.Run(testName, func(t *testing.T) {
			if c.setup != nil {
				c.setup()
			} else {
				// Tests with custom setup modify global state and cannot run in parallel
				t.Parallel()
			}
			if c.teardown != nil {
				t.Cleanup(c.teardown)
			}

			mc, err := mesh.DeepCopyMeshConfig(defaultMesh)
			if err != nil {
				t.Fatal(err)
			}
			sidecarTemplate, valuesConfig := defaultTemplate, defaultValues
			if c.setFlags != nil || c.inFilePath != "" {
				writeInjectionSettings(t, fmt.Sprintf("%s.%d", c.in, i), c.setFlags, c.inFilePath)
				sidecarTemplate, valuesConfig, mc = readInjectionSettings(t, fmt.Sprintf("%s.%d", c.in, i))
			}
			if c.mesh != nil {
				c.mesh(mc)
			}

			inputFilePath := "testdata/customfields/" + c.in
			wantFilePath := "testdata/customfields/" + c.want
			in, err := os.Open(inputFilePath)
			if err != nil {
				t.Fatalf("Failed to open %q: %v", inputFilePath, err)
			}
			t.Cleanup(func() {
				_ = in.Close()
			})

			// Exit early if we don't need to test webhook. We can skip errors since its redundant
			// and painful to test here.
			if c.expectedError != "" || c.skipWebhook {
				return
			}
			// Next run the webhook test. This one is a bit trickier as the webhook operates
			// on Pods, but the inputs are Deployments/StatefulSets/etc. As a result, we need
			// to convert these to pods, then run the injection This test will *not*
			// overwrite golden files, as we do not have identical textual output as
			// kube-inject. Instead, we just compare the desired/actual pod specs.
			t.Run("webhook", func(t *testing.T) {
				webhook := &Webhook{
					Config:     sidecarTemplate,
					meshConfig: mc,
					env: &model.Environment{
						PushContext: &model.PushContext{
							ProxyConfigs: &model.ProxyConfigs{},
						},
					},
					valuesConfig: valuesConfig,
					revision:     "default",
					kubeRegistry: MockServiceRegistry{
						ServiceAnnotations: c.svcAnnotations,
					},
				}
				// Split multi-part yaml documents. Input and output will have the same number of parts.
				inputYAMLs := splitYamlFile(inputFilePath, t)
				wantYAMLs := splitYamlFile(wantFilePath, t)
				for i := 0; i < len(inputYAMLs); i++ {
					t.Run(fmt.Sprintf("yamlPart[%d]", i), func(t *testing.T) {
						runWebhook(t, webhook, inputYAMLs[i], wantYAMLs[i], false)
						t.Run("idempotency", func(t *testing.T) {
							runWebhook(t, webhook, wantYAMLs[i], wantYAMLs[i], true)
						})
					})
				}
			})
		})
	}

	// Make sure we don't have any stale test data leftover, as it can cause confusion.

	for _, c := range cases {
		delete(allOutputFiles, c.want)
	}

	if len(allOutputFiles) != 0 {
		t.Fatalf("stale golden files found: %v", allOutputFiles.UnsortedList())
	}
}

// MockServiceRegistry Instance implementation, where fields are set individually.
type MockServiceRegistry struct {
	ProviderID         provider.ID
	ClusterID          cluster.ID
	ServiceAnnotations map[string]string

	model.Controller
	model.ServiceDiscovery
}

func (r MockServiceRegistry) Provider() provider.ID {
	return r.ProviderID
}

func (r MockServiceRegistry) Cluster() cluster.ID {
	return r.ClusterID
}

func (r MockServiceRegistry) GetServiceAccountAnnotations(name, namespace string) map[string]string {
	return r.ServiceAnnotations
}
