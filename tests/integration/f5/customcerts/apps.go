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
	"encoding/json"

	"istio.io/istio/pkg/config/protocol"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/echo/deployment"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/security/pkg/pki/util"
)

type EchoDeployments struct {
	// Namespace echo apps will be deployed
	Namespace namespace.Instance
	// Standard echo app to be used by tests
	EchoPod echo.Instance
}

var svcAnnotationConfig = map[string]*util.CertificateAspenMeshIOCustomFields{
	"san-none": {
		SAN: &util.SAN{},
	},
	"san-uri": {
		SAN: &util.SAN{
			URI: []string{"http://test.example.com/get"},
		},
	},
	"san-dns": {
		SAN: &util.SAN{
			DNS: []string{"http.echo.com"},
		},
	},
	"san-multi": {
		SAN: &util.SAN{
			DNS: []string{"http.echo.com", "test.echo.com"},
			URI: []string{"http://test.example.com/get", "http://test.example.com/post", "com.example.uuid:cc1bf660-9e35-4169-9bbf-2a893ea2fa5d"},
		},
	},
}

var EchoPorts = []echo.Port{
	{Name: "http", Protocol: protocol.HTTP, ServicePort: 80, WorkloadPort: 18080},
}

var CustomCerts = certAnnotation("certificate.aspenmesh.io/customFields", "")

func certAnnotation(name string, value string) echo.Annotation {
	return echo.Annotation{
		Name: name,
		Type: "customcert",
		Default: echo.AnnotationValue{
			Value: value,
		},
	}
}

func SetupApps(ctx resource.Context, apps map[string]*EchoDeployments) error {
	var err error
	for name, app := range apps {
		app.Namespace, err = namespace.New(ctx, namespace.Config{
			Prefix: name,
			Inject: true,
		})
		if err != nil {
			return err
		}
		cfg, err := json.Marshal(svcAnnotationConfig[name])
		if err != nil {
			return err
		}
		annotations := echo.NewAnnotations()
		annotations.Set(CustomCerts, string(cfg))

		builder := deployment.New(ctx).
			WithClusters(ctx.Clusters()...).
			With(&apps[name].EchoPod, echo.Config{
				Service:                   name,
				Namespace:                 app.Namespace,
				Ports:                     EchoPorts,
				ServiceAccount:            true,
				ServiceAccountAnnotations: annotations,
			})

		_, err = builder.Build()
		if err != nil {
			return err
		}

	}
	return nil
}
