//go:build integ
// +build integ

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

package security

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"istio.io/istio/pkg/http/headers"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/echo/check"
	"istio.io/istio/pkg/test/framework/components/echo/common/ports"
	"istio.io/istio/pkg/test/framework/components/echo/config"
	"istio.io/istio/pkg/test/framework/components/echo/config/param"
	"istio.io/istio/pkg/test/framework/components/echo/match"
)

// TestXFCCAuthorization_Conditions modifies TestAuthorization_Conditions test and
// tests contains authorization with conditions using XFCC headers.
func TestXFCCAuthorization_Conditions(t *testing.T) {
	// nolint: staticcheck
	framework.NewTest(t).
		Features("security.authorization.conditions").
		RequiresSingleCluster().
		Run(func(t framework.TestContext) {
			allowed := apps.Ns1.A
			denied := apps.Ns2.A

			from := allowed.Append(denied)
			fromMatch := match.AnyServiceName(from.NamespacedNames())
			toMatch := match.Not(fromMatch)
			to := toMatch.GetServiceMatches(apps.Ns1.All)
			fromAndTo := to.Instances().Append(from)

			config.New(t).
				Source(config.File("testdata/authz/mtls.yaml.tmpl")).
				Source(config.File("testdata/authz/xfcc-conditions.yaml.tmpl").WithParams(param.Params{
					"Allowed": allowed,
					"Denied":  denied,
				})).
				BuildAll(nil, to).
				Apply()

			newTrafficTest(t, fromAndTo).
				FromMatch(fromMatch).
				ToMatch(toMatch).
				Run(func(t framework.TestContext, from echo.Instance, to echo.Target) {
					allow := allowValue(from.NamespacedName() == allowed.Config().NamespacedName())
					cases := []struct {
						path    string
						headers http.Header
						allow   allowValue
					}{
						// Test headers.
						{
							path:  "/request-xfcc-headers",
							allow: allow,
						},
						// test suffix match with baz* (envoy sanitizes xfcc header, so these should never work)
						{
							path:  "/request-xfcc-headers-baz",
							allow: false,
						},
					}

					for _, c := range cases {
						c := c
						xHeader := ""
						if c.headers != nil {
							xHeader = "?x-forwarded-client-cert=" + c.headers.Get("x-forwarded-client-cert")
						}
						testName := fmt.Sprintf("%s%s(%s)/http", c.path, xHeader, c.allow)
						t.NewSubTest(testName).Run(func(t framework.TestContext) {
							newAuthzTest().
								From(from).
								To(to).
								PortName(ports.HTTP.Name).
								Path(c.path).
								Allow(c.allow).
								Headers(c.headers).
								BuildAndRun(t)
						})
					}
					t.NewSubTest("headerstoolarge(deny)").Run(func(t framework.TestContext) {
						tooLargeTest := newAuthzTest().
							From(from).
							To(to).
							PortName(ports.HTTP.Name).
							Path("/request-xfcc-headers").
							Allow(false).
							Headers(headers.New().With("x-forwarded-client-cert", generateLargeHeader()).Build()).
							Build(t)
						tooLargeTest.opts.Check = check.ErrorOrStatus(431)
						tooLargeTest.Run(t)
					})
				})
		})
}

func generateLargeHeader() string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return strings.Repeat(letters, 3000)
}
