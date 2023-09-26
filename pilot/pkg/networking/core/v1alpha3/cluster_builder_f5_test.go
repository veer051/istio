// Copyright Istio Authors. All Rights Reserved.
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

package v1alpha3

import (
	"strings"
	"testing"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"

	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/util"
	"istio.io/istio/pkg/config/protocol"
	"istio.io/istio/pkg/test"
	"istio.io/istio/pkg/test/util/assert"
)

// TestBuildUpstreamClusterTLSContextServiceEntryExternalIstioMutual tests the buildUpstreamClusterTLSContext function
func TestBuildUpstreamClusterTLSContextServiceEntryExternalIstioMutual(t *testing.T) {
	testCases := []struct {
		name         string
		opts         *buildClusterOpts
		tls          *networking.ClientTLSSettings
		meshExternal bool
		result       expectedResult
	}{
		{
			name: "SE location: MESH_EXTERNAL and protocol is HTTP, tls mode ISTIO_MUTUAL with sni having a prefix",
			opts: &buildClusterOpts{
				mutable: &MutableCluster{
					cluster: &cluster.Cluster{
						Name: "test-cluster",
					},
				},
				port: &model.Port{
					Name:     "http",
					Port:     8443,
					Protocol: protocol.HTTP,
				},
			},
			tls: &networking.ClientTLSSettings{
				Mode:            networking.ClientTLSSettings_ISTIO_MUTUAL,
				SubjectAltNames: []string{"SAN"},
				// sni has wildcard so auto_sni should be set to true
				Sni: "*.f5.com",
			},
			result: expectedResult{
				tlsContext: &tls.UpstreamTlsContext{
					CommonTlsContext: &tls.CommonTlsContext{
						TlsParams: &tls.TlsParameters{
							// if not specified, envoy use TLSv1_2 as default for client.
							TlsMaximumProtocolVersion: tls.TlsParameters_TLSv1_3,
							TlsMinimumProtocolVersion: tls.TlsParameters_TLSv1_2,
						},
						TlsCertificateSdsSecretConfigs: []*tls.SdsSecretConfig{
							{
								Name: "default",
								SdsConfig: &core.ConfigSource{
									ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
										ApiConfigSource: &core.ApiConfigSource{
											ApiType:                   core.ApiConfigSource_GRPC,
											SetNodeOnFirstMessageOnly: true,
											TransportApiVersion:       core.ApiVersion_V3,
											GrpcServices: []*core.GrpcService{
												{
													TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
														EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
													},
												},
											},
										},
									},
									InitialFetchTimeout: durationpb.New(time.Second * 0),
									ResourceApiVersion:  core.ApiVersion_V3,
								},
							},
						},
						ValidationContextType: &tls.CommonTlsContext_CombinedValidationContext{
							CombinedValidationContext: &tls.CommonTlsContext_CombinedCertificateValidationContext{
								DefaultValidationContext: &tls.CertificateValidationContext{MatchSubjectAltNames: util.StringToExactMatch([]string{"SAN"})},
								ValidationContextSdsSecretConfig: &tls.SdsSecretConfig{
									Name: "ROOTCA",
									SdsConfig: &core.ConfigSource{
										ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
											ApiConfigSource: &core.ApiConfigSource{
												ApiType:                   core.ApiConfigSource_GRPC,
												SetNodeOnFirstMessageOnly: true,
												TransportApiVersion:       core.ApiVersion_V3,
												GrpcServices: []*core.GrpcService{
													{
														TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
															EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
														},
													},
												},
											},
										},
										InitialFetchTimeout: durationpb.New(time.Second * 0),
										ResourceApiVersion:  core.ApiVersion_V3,
									},
								},
							},
						},
						AlpnProtocols: util.ALPNInMeshWithMxc,
					},
					// sni should not be set, instead auto_sni should be set to true
					// Sni: "some-sni.com",
				},
				err: nil,
			},
			meshExternal: true,
		},
		{
			name: "SE location: MESH_EXTERNAL and protocol is HTTP2, tls mode ISTIO_MUTUAL with sni having a prefix",
			opts: &buildClusterOpts{
				mutable: &MutableCluster{
					cluster: &cluster.Cluster{
						Name: "test-cluster",
					},
				},
				port: &model.Port{
					Name:     "http",
					Port:     8443,
					Protocol: protocol.HTTP2,
				},
			},
			tls: &networking.ClientTLSSettings{
				Mode:            networking.ClientTLSSettings_ISTIO_MUTUAL,
				SubjectAltNames: []string{"SAN"},
				// sni has wildcard so auto_sni should be set to true
				Sni: "*.f5.com",
			},
			result: expectedResult{
				tlsContext: &tls.UpstreamTlsContext{
					CommonTlsContext: &tls.CommonTlsContext{
						TlsParams: &tls.TlsParameters{
							// if not specified, envoy use TLSv1_2 as default for client.
							TlsMaximumProtocolVersion: tls.TlsParameters_TLSv1_3,
							TlsMinimumProtocolVersion: tls.TlsParameters_TLSv1_2,
						},
						TlsCertificateSdsSecretConfigs: []*tls.SdsSecretConfig{
							{
								Name: "default",
								SdsConfig: &core.ConfigSource{
									ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
										ApiConfigSource: &core.ApiConfigSource{
											ApiType:                   core.ApiConfigSource_GRPC,
											SetNodeOnFirstMessageOnly: true,
											TransportApiVersion:       core.ApiVersion_V3,
											GrpcServices: []*core.GrpcService{
												{
													TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
														EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
													},
												},
											},
										},
									},
									InitialFetchTimeout: durationpb.New(time.Second * 0),
									ResourceApiVersion:  core.ApiVersion_V3,
								},
							},
						},
						ValidationContextType: &tls.CommonTlsContext_CombinedValidationContext{
							CombinedValidationContext: &tls.CommonTlsContext_CombinedCertificateValidationContext{
								DefaultValidationContext: &tls.CertificateValidationContext{MatchSubjectAltNames: util.StringToExactMatch([]string{"SAN"})},
								ValidationContextSdsSecretConfig: &tls.SdsSecretConfig{
									Name: "ROOTCA",
									SdsConfig: &core.ConfigSource{
										ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
											ApiConfigSource: &core.ApiConfigSource{
												ApiType:                   core.ApiConfigSource_GRPC,
												SetNodeOnFirstMessageOnly: true,
												TransportApiVersion:       core.ApiVersion_V3,
												GrpcServices: []*core.GrpcService{
													{
														TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
															EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
														},
													},
												},
											},
										},
										InitialFetchTimeout: durationpb.New(time.Second * 0),
										ResourceApiVersion:  core.ApiVersion_V3,
									},
								},
							},
						},
						AlpnProtocols: util.ALPNInMeshWithMxc,
					},
					// sni should not be set, instead auto_sni should be set to true
					// Sni: "some-sni.com",
				},
				err: nil,
			},
			meshExternal: true,
		},
		{
			name: "SE location: MESH_EXTERNAL and protocol is GRPC, tls mode ISTIO_MUTUAL with sni having a prefix",
			opts: &buildClusterOpts{
				mutable: &MutableCluster{
					cluster: &cluster.Cluster{
						Name: "test-cluster",
					},
				},
				port: &model.Port{
					Name:     "http",
					Port:     8443,
					Protocol: protocol.GRPC,
				},
			},
			tls: &networking.ClientTLSSettings{
				Mode:            networking.ClientTLSSettings_ISTIO_MUTUAL,
				SubjectAltNames: []string{"SAN"},
				// sni has wildcard so auto_sni should be set to true
				Sni: "*.f5.com",
			},
			result: expectedResult{
				tlsContext: &tls.UpstreamTlsContext{
					CommonTlsContext: &tls.CommonTlsContext{
						TlsParams: &tls.TlsParameters{
							// if not specified, envoy use TLSv1_2 as default for client.
							TlsMaximumProtocolVersion: tls.TlsParameters_TLSv1_3,
							TlsMinimumProtocolVersion: tls.TlsParameters_TLSv1_2,
						},
						TlsCertificateSdsSecretConfigs: []*tls.SdsSecretConfig{
							{
								Name: "default",
								SdsConfig: &core.ConfigSource{
									ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
										ApiConfigSource: &core.ApiConfigSource{
											ApiType:                   core.ApiConfigSource_GRPC,
											SetNodeOnFirstMessageOnly: true,
											TransportApiVersion:       core.ApiVersion_V3,
											GrpcServices: []*core.GrpcService{
												{
													TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
														EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
													},
												},
											},
										},
									},
									InitialFetchTimeout: durationpb.New(time.Second * 0),
									ResourceApiVersion:  core.ApiVersion_V3,
								},
							},
						},
						ValidationContextType: &tls.CommonTlsContext_CombinedValidationContext{
							CombinedValidationContext: &tls.CommonTlsContext_CombinedCertificateValidationContext{
								DefaultValidationContext: &tls.CertificateValidationContext{MatchSubjectAltNames: util.StringToExactMatch([]string{"SAN"})},
								ValidationContextSdsSecretConfig: &tls.SdsSecretConfig{
									Name: "ROOTCA",
									SdsConfig: &core.ConfigSource{
										ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
											ApiConfigSource: &core.ApiConfigSource{
												ApiType:                   core.ApiConfigSource_GRPC,
												SetNodeOnFirstMessageOnly: true,
												TransportApiVersion:       core.ApiVersion_V3,
												GrpcServices: []*core.GrpcService{
													{
														TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
															EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
														},
													},
												},
											},
										},
										InitialFetchTimeout: durationpb.New(time.Second * 0),
										ResourceApiVersion:  core.ApiVersion_V3,
									},
								},
							},
						},
						AlpnProtocols: util.ALPNInMeshWithMxc,
					},
					// sni should not be set, instead auto_sni should be set to true
					// Sni: "some-sni.com",
				},
				err: nil,
			},
			meshExternal: true,
		},
		{
			name: "SE location: MESH_INTERNAL and protocol is HTTP, tls mode ISTIO_MUTUAL with sni having a prefix",
			opts: &buildClusterOpts{
				mutable: &MutableCluster{
					cluster: &cluster.Cluster{
						Name: "test-cluster",
					},
				},
				port: &model.Port{
					Name:     "http",
					Port:     8443,
					Protocol: protocol.HTTP,
				},
			},
			tls: &networking.ClientTLSSettings{
				Mode:            networking.ClientTLSSettings_ISTIO_MUTUAL,
				SubjectAltNames: []string{"SAN"},
				Sni:             "*.f5.com",
			},
			result: expectedResult{
				tlsContext: &tls.UpstreamTlsContext{
					CommonTlsContext: &tls.CommonTlsContext{
						TlsParams: &tls.TlsParameters{
							// if not specified, envoy use TLSv1_2 as default for client.
							TlsMaximumProtocolVersion: tls.TlsParameters_TLSv1_3,
							TlsMinimumProtocolVersion: tls.TlsParameters_TLSv1_2,
						},
						TlsCertificateSdsSecretConfigs: []*tls.SdsSecretConfig{
							{
								Name: "default",
								SdsConfig: &core.ConfigSource{
									ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
										ApiConfigSource: &core.ApiConfigSource{
											ApiType:                   core.ApiConfigSource_GRPC,
											SetNodeOnFirstMessageOnly: true,
											TransportApiVersion:       core.ApiVersion_V3,
											GrpcServices: []*core.GrpcService{
												{
													TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
														EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
													},
												},
											},
										},
									},
									InitialFetchTimeout: durationpb.New(time.Second * 0),
									ResourceApiVersion:  core.ApiVersion_V3,
								},
							},
						},
						ValidationContextType: &tls.CommonTlsContext_CombinedValidationContext{
							CombinedValidationContext: &tls.CommonTlsContext_CombinedCertificateValidationContext{
								DefaultValidationContext: &tls.CertificateValidationContext{MatchSubjectAltNames: util.StringToExactMatch([]string{"SAN"})},
								ValidationContextSdsSecretConfig: &tls.SdsSecretConfig{
									Name: "ROOTCA",
									SdsConfig: &core.ConfigSource{
										ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
											ApiConfigSource: &core.ApiConfigSource{
												ApiType:                   core.ApiConfigSource_GRPC,
												SetNodeOnFirstMessageOnly: true,
												TransportApiVersion:       core.ApiVersion_V3,
												GrpcServices: []*core.GrpcService{
													{
														TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
															EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
														},
													},
												},
											},
										},
										InitialFetchTimeout: durationpb.New(time.Second * 0),
										ResourceApiVersion:  core.ApiVersion_V3,
									},
								},
							},
						},
						AlpnProtocols: util.ALPNInMeshWithMxc,
					},
					Sni: "*.f5.com",
				},
				err: nil,
			},
			meshExternal: false,
		},
		{
			name: "SE location: MESH_EXTERNAL and protocol is HTTP, tls mode ISTIO_MUTUAL with sni NOT having a prefix",
			opts: &buildClusterOpts{
				mutable: &MutableCluster{
					cluster: &cluster.Cluster{
						Name: "test-cluster",
					},
				},
				port: &model.Port{
					Name:     "http",
					Port:     8443,
					Protocol: protocol.HTTP,
				},
			},
			tls: &networking.ClientTLSSettings{
				Mode:            networking.ClientTLSSettings_ISTIO_MUTUAL,
				SubjectAltNames: []string{"SAN"},
				Sni:             "aspenmesh.f5.com",
			},
			result: expectedResult{
				tlsContext: &tls.UpstreamTlsContext{
					CommonTlsContext: &tls.CommonTlsContext{
						TlsParams: &tls.TlsParameters{
							// if not specified, envoy use TLSv1_2 as default for client.
							TlsMaximumProtocolVersion: tls.TlsParameters_TLSv1_3,
							TlsMinimumProtocolVersion: tls.TlsParameters_TLSv1_2,
						},
						TlsCertificateSdsSecretConfigs: []*tls.SdsSecretConfig{
							{
								Name: "default",
								SdsConfig: &core.ConfigSource{
									ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
										ApiConfigSource: &core.ApiConfigSource{
											ApiType:                   core.ApiConfigSource_GRPC,
											SetNodeOnFirstMessageOnly: true,
											TransportApiVersion:       core.ApiVersion_V3,
											GrpcServices: []*core.GrpcService{
												{
													TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
														EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
													},
												},
											},
										},
									},
									InitialFetchTimeout: durationpb.New(time.Second * 0),
									ResourceApiVersion:  core.ApiVersion_V3,
								},
							},
						},
						ValidationContextType: &tls.CommonTlsContext_CombinedValidationContext{
							CombinedValidationContext: &tls.CommonTlsContext_CombinedCertificateValidationContext{
								DefaultValidationContext: &tls.CertificateValidationContext{MatchSubjectAltNames: util.StringToExactMatch([]string{"SAN"})},
								ValidationContextSdsSecretConfig: &tls.SdsSecretConfig{
									Name: "ROOTCA",
									SdsConfig: &core.ConfigSource{
										ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
											ApiConfigSource: &core.ApiConfigSource{
												ApiType:                   core.ApiConfigSource_GRPC,
												SetNodeOnFirstMessageOnly: true,
												TransportApiVersion:       core.ApiVersion_V3,
												GrpcServices: []*core.GrpcService{
													{
														TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
															EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
														},
													},
												},
											},
										},
										InitialFetchTimeout: durationpb.New(time.Second * 0),
										ResourceApiVersion:  core.ApiVersion_V3,
									},
								},
							},
						},
						AlpnProtocols: util.ALPNInMeshWithMxc,
					},
					Sni: "aspenmesh.f5.com",
				},
				err: nil,
			},
			meshExternal: true,
		},
		{
			name: "SE location: MESH_EXTERNAL and protocol is TCP, tls mode ISTIO_MUTUAL with sni having a wildcard prefix",
			opts: &buildClusterOpts{
				mutable: &MutableCluster{
					cluster: &cluster.Cluster{
						Name: "test-cluster",
					},
				},
				port: &model.Port{
					Name:     "tcp",
					Port:     8443,
					Protocol: protocol.TCP,
				},
			},
			tls: &networking.ClientTLSSettings{
				Mode:            networking.ClientTLSSettings_ISTIO_MUTUAL,
				SubjectAltNames: []string{"SAN"},
				Sni:             "*.f5.com",
			},
			result: expectedResult{
				tlsContext: &tls.UpstreamTlsContext{
					CommonTlsContext: &tls.CommonTlsContext{
						TlsParams: &tls.TlsParameters{
							// if not specified, envoy use TLSv1_2 as default for client.
							TlsMaximumProtocolVersion: tls.TlsParameters_TLSv1_3,
							TlsMinimumProtocolVersion: tls.TlsParameters_TLSv1_2,
						},
						TlsCertificateSdsSecretConfigs: []*tls.SdsSecretConfig{
							{
								Name: "default",
								SdsConfig: &core.ConfigSource{
									ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
										ApiConfigSource: &core.ApiConfigSource{
											ApiType:                   core.ApiConfigSource_GRPC,
											SetNodeOnFirstMessageOnly: true,
											TransportApiVersion:       core.ApiVersion_V3,
											GrpcServices: []*core.GrpcService{
												{
													TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
														EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
													},
												},
											},
										},
									},
									InitialFetchTimeout: durationpb.New(time.Second * 0),
									ResourceApiVersion:  core.ApiVersion_V3,
								},
							},
						},
						ValidationContextType: &tls.CommonTlsContext_CombinedValidationContext{
							CombinedValidationContext: &tls.CommonTlsContext_CombinedCertificateValidationContext{
								DefaultValidationContext: &tls.CertificateValidationContext{MatchSubjectAltNames: util.StringToExactMatch([]string{"SAN"})},
								ValidationContextSdsSecretConfig: &tls.SdsSecretConfig{
									Name: "ROOTCA",
									SdsConfig: &core.ConfigSource{
										ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
											ApiConfigSource: &core.ApiConfigSource{
												ApiType:                   core.ApiConfigSource_GRPC,
												SetNodeOnFirstMessageOnly: true,
												TransportApiVersion:       core.ApiVersion_V3,
												GrpcServices: []*core.GrpcService{
													{
														TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
															EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "sds-grpc"},
														},
													},
												},
											},
										},
										InitialFetchTimeout: durationpb.New(time.Second * 0),
										ResourceApiVersion:  core.ApiVersion_V3,
									},
								},
							},
						},
						AlpnProtocols: util.ALPNInMeshWithMxc,
					},
					Sni: "*.f5.com",
				},
				err: nil,
			},
			meshExternal: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			test.SetForTest(t, &features.CarrierGradeServiceEntryIstioMutual, true)
			proxy := newSidecarProxy()
			cb := NewClusterBuilder(proxy, nil, model.DisabledCache{})
			tc.opts.meshExternal = tc.meshExternal
			ret, err := cb.buildUpstreamClusterTLSContext(tc.opts, tc.tls)
			if err != nil && tc.result.err == nil || err == nil && tc.result.err != nil {
				t.Errorf("expecting:\n err=%v but got err=%v", tc.result.err, err)
			} else if diff := cmp.Diff(tc.result.tlsContext, ret, protocmp.Transform()); diff != "" {
				t.Errorf("got diff: `%v", diff)
			}
			if tc.meshExternal && strings.HasPrefix(tc.tls.Sni, "*.") && (tc.opts.port.Protocol.IsHTTP() || tc.opts.port.Protocol.IsHTTP2()) {
				assert.Equal(t, tc.opts.mutable.httpProtocolOptions.UpstreamHttpProtocolOptions.AutoSni, true)
			} else if tc.opts.mutable.httpProtocolOptions != nil {
				t.Errorf("%v: expected httpProtocolOptions to be nil", tc.name)
			}
		})
	}
}
