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

package ca

import (
	"reflect"
	"sort"
	"testing"
)

const (
	csrWithoutSanExtension = `
-----BEGIN CERTIFICATE REQUEST-----
MIICWDCCAUACAQAwEzERMA8GA1UEChMISnVqdSBvcmcwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQCv+wnhAbUu+2iYYYmehVquT1BZePJSVA7mHWzFih1a
eGS+fHnlq1TDIgMfh3GhQSWyNGgompCM6+OCeZiNIk77wiIK9Wmm0Wckk9X5qHyB
w8eAtgrmVUzxxVBKKJGYAO8cu12X2TK+5cXJq0WZd5zqqh6nX/LoopFTeXPiCu7x
FSIzoDfnYrKFp5TmjYFaUR6B/0qhoBd4/cA+VcHxuR8iq1tfZ1IeV7O/g6UUP7AE
Nz4RDxYwWgPeIcVp1npD4WBIEP18yh1MtJ9UnG0o/lNJVqiOZLL6j4dNUWJ/KScM
H7SLEIwPcwRjYcaFHg1u12KSGcbm57cEOH6+3EpmuCzlAgMBAAGgADANBgkqhkiG
9w0BAQsFAAOCAQEACVDTH1SYp+t0m0a3s6BaUAhcamdpp8pNpwoQGTnJRg6Ojf7C
nJ4XKMEPYx99JjnKmV4OcjkpNw0m6a6T/Ka8k5Zjel1Pi1Ou7h9mOCZUXFoUuh6N
BOhD7Bs6i6kVKPvUCgqjmZxEp2HUM936IYBQnYh5nhfP/YL5MUDmiy1D/DFtTwww
T1JrRI5Up6gm7n0BHQLK64/+jAd69ih9JMCnGrl1AekNZSDYyfZuMjXYyhLJMi1c
gWSiPRrze8PA1iibmwU/Z4EGkR2o8Gtu2VyalZsOnbu+EUWgr4hz20JE8IPt7v2o
kIm0u1vpISHhYz9pgWSj6y1oN30gjXN047TkEg==
-----END CERTIFICATE REQUEST-----`

	// spiffe://cluster.local/ns/default/sa/httpbin
	csrWithSpiffeExtensionOnly = `
-----BEGIN CERTIFICATE REQUEST-----
MIICpTCCAY0CAQAwEzERMA8GA1UEChMISnVqdSBvcmcwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQCz+I4Qfy/npskwqlBT7T1l4nbUCLmCkVdt3jBSOsyC
Ptw+0aDnb9HsM1lfkVfdHwOLtCofQYwBM6xv/lxrMLyJJKz19A8UFdkGsn6cueQ/
nA/6TeKHITtJSrH31v2yhCx+cs1K7zwfbaV06BfBkJlmEpYpSAe4Eygg62J1djEp
/ECkkwZTh6zBVIpROYpn2aSa/AWuD8bV0wqW+9NrTHv0aCCfTE7T0llyfsTuN6Uh
KJRV2v5WxBwc2vQ+mXZHfkFJie6RR78XUB67d68C6taR9xXJY2C9dC1Ukl4pZw85
HyK8zHaCI0MNc6nvJ5SazTXHZFcvQOZD1IOSL/48UWf5AgMBAAGgTTBLBgkqhkiG
9w0BCQ4xPjA8MDoGA1UdEQEB/wQwMC6GLHNwaWZmZTovL2NsdXN0ZXIubG9jYWwv
bnMvZGVmYXVsdC9zYS9odHRwYmluMA0GCSqGSIb3DQEBCwUAA4IBAQByJoI8jgdm
NtBezXeUE1h4xMjAnSCfOuHEG0pvsPN+Mp7IcjSFaKUebktNvK1ESNmodU3pBRtk
VyDWeBW0HVcUR+fRyrok9+zG1Y2c54xxI419kx5kTfH9NtNpU+iwNTkVNiApX6A7
Ccmx7cbu/HA0BbbrWWBqUkxyeC7t9uRvCxODcgZEiXwDVdPKfNwbVgN7r6aD4+SK
JVUT4brKpYd50SeEm9ANmjNxTf9LdeNvdU83oO/gtAyyaOIfRcm16EcvOb5WpfN7
YJxZHz2teZajcQYS0uQNi63MMtLrseGly6y0ENuNzL/VcGrUp9OI/JNDwYV7UWms
h8nEs6j6Oluu
-----END CERTIFICATE REQUEST-----`

	// spiffe://cluster.local/ns/default/sa/httpbin,httpbin.org
	csrWithSpiffeAndSingleDNSExtensions = `
-----BEGIN CERTIFICATE REQUEST-----
MIICsjCCAZoCAQAwEzERMA8GA1UEChMISnVqdSBvcmcwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQCYqKC36iSl4/e/FOPB0XgZIs4i+XAUmptMsDd2Ziqd
lDHcEzybBAN2XpZkv28PjdVOxWKpBXkMAMLtvC+BBrDuYmJ2Hc1tclePG7kGVl63
vMVMjJ3Fb8hbYir3ek0w2qSh/Y1iaJ4Lj2Xvinl4HRPTMHjCxymDfGRerx6c7w/s
RMi+u/jU8R5e7wszSoAJqPevwkzqsAGvqvoLe/Xfx/tYGJBAdtha0aVL8IVncSp9
56XIqMD51qjTCtLe8z9uxOxhsEI1sM7J3JHVT5MBbAP6Yi/XvBZ64GPpbyfNm4fg
FDvJUwpDjRKMoJ91MArHqTajME74f+hqyXzN1yOhiW9BAgMBAAGgWjBYBgkqhkiG
9w0BCQ4xSzBJMEcGA1UdEQEB/wQ9MDuGLHNwaWZmZTovL2NsdXN0ZXIubG9jYWwv
bnMvZGVmYXVsdC9zYS9odHRwYmluggtodHRwYmluLm9yZzANBgkqhkiG9w0BAQsF
AAOCAQEAYAhP1F6UXElAn3MuduiTW9N2775MRcw12rE1GEwKGQ/ceUz1XnL3ZNle
Yt8vPGnpegbzCEa6NtEtb7sTZ/HrirVEThJz8YtvukzLLivY9xb2IVJtCJYuBqo9
NgGH5zEfEKsHkwDLY+WFJpQxBzgQiHWz9HRN0SG9u5H6jnWZf9iaz/jeIH/mYbCv
COr6MnG+qR3QBDmfHSB3eGZhZ95bWjPLw0FfqU/YE2HJxVkTnGO0nt4V2WUR2lIc
gvliaNjbI1AGsfdz0BkasDGee99j9OhGollmOvFO0G7/tpgFyke5cifFDKftcq+J
yfW6OcBAoUTzk2+xBTg70mA1cMKD5g==
-----END CERTIFICATE REQUEST-----`

	csrWithSpiffeAndMutleipleDNSExtensions = `
-----BEGIN CERTIFICATE REQUEST-----
MIICwzCCAasCAQAwEzERMA8GA1UEChMISnVqdSBvcmcwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDjIlSIy4w7xahlUXNZaQIkrAt95kzLpPvpuYAUxQnI
cS3EULl6rp1H4qyMyTjqWE4xcn/HBkKOLmkDl7L+0l0trEPjMhS1/PGlLlUTA1iu
QqCIk1A00bXVgZuw+L1+Wgwon60y/xEAMWW39kPW4piNVGdq+Omdja6+4rNzBNf0
WJRLT52cwAUkXqjdF8yVImSCPGHr+4sJne+3rs/GQ/CafLv6hjMDjU003SGUkLFI
8oPlpVtw+ifxcg04sm6v8HrHeZZ/KCHWgM+KTCVD+or/5DGZkXIYaZLwtkgSaz+3
9PF1m5KeLhdUnQ3eUOh/kdkj+5kEBsGwoHKR5b5LEty7AgMBAAGgazBpBgkqhkiG
9w0BCQ4xXDBaMFgGA1UdEQEB/wROMEyGLHNwaWZmZTovL2NsdXN0ZXIubG9jYWwv
bnMvZGVmYXVsdC9zYS9odHRwYmluggtodHRwYmluLm9yZ4IPbXkuYXNwZW5tZXNo
LmlvMA0GCSqGSIb3DQEBCwUAA4IBAQAetQ7+74AIkD6aKD7ygnvhhZ001LtiAQ4q
4rF83msATmNdsd8G2hx2NAWxsweUiNtB4dVce/XUQH5EHBe+wRP9cdrdPYbGMTsH
H2vniEtCieMNgg1tELVFD9B4o6/L1Srl8PxAKYxq6FMY8FdXy9etF/o22tpiTjBb
tFh5BZ4KnXSxlhk1/0HlAJFCUAFM83gYQ8wkoQuqDEbaww5RhWPI+0U6Z183r4HG
wR5FdxJo5CNBNvjMMC+tKQ7uOtej0vdGAvxAC6c4tssvPmt+2rqN7BEesznn1qA8
T10ns4TXr8/EbIVc6cSf/EQpIQmo/i68sM247pE8VX7a4ZnA7N0B
-----END CERTIFICATE REQUEST-----`

	csrWithSpiffeAndUUIDDNSExtension = `
-----BEGIN CERTIFICATE REQUEST-----
MIIC0jCCAboCAQAwEzERMA8GA1UEChMISnVqdSBvcmcwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDO8w80InxtAxB4hcDngqAwfwfa8y76jnE7Kn2nvPIh
8ObY3eJPqQRd1IanbvHwTKItainijcsVMHrSnXA9uozDUR6AVF0sAxUIVxLUrtG9
GPaeBcFX5rv/dh8jC2/uqi0gy5SKz/9os33PkbfsYajhiak9A+gxANs198Yd7JHR
KFswrpGk0VAuw4LKrLkNB6NmX9xrbY8eWKovqpww2wKGcSgxxnY+KPP7fupfLNyG
pqUjlSk1HBDnKqEHBSOux7OY2k09eiNMz6/sSIB7OFrAqe0ZOdkqIakka9cFNddw
IYiMIe5IHCpZ1cUvyvAFSR51nYSSi3AIJ3YmC81hH7pfAgMBAAGgejB4BgkqhkiG
9w0BCQ4xazBpMGcGA1UdEQEB/wRdMFuGK3V1aWQ6Ly9mNTYxNjU4OS1jMTg3LTQ0
MDEtOTUwMC0zY2E5ODY1NmQ0MGOGLHNwaWZmZTovL2NsdXN0ZXIubG9jYWwvbnMv
ZGVmYXVsdC9zYS9odHRwYmluMA0GCSqGSIb3DQEBCwUAA4IBAQAIfP4Pr+rjK6wh
WZTpG3Hp0/iX6IebdOHTaLzQB3qAwKDmWuNlGoRfdJhZuQLrG/M6fxU10JG0Bu69
1x0pj7Ec8h7NKpwfXFw+CsUM6THnB2s4hbGlnf7AZNaqjLs/wDR/+riGO9oQqKTW
BgdhAE0HJ39iUcTD4dls/FIaMotxJLGlwabpx+E8ca24xg3KDvXBfoi3u0rJMn1D
LK3I4oICiXbdAEDNmqdpVmQrenX2m49v9StJz2G0Raixyj6SwvM7/a3NYjP96H37
a9MyUC5Du/5sRzK564UQzXby+iw5/HL8d796qI9Ndy1myo9ihLi37pV9YgD3lCKe
bhlXM9pK
-----END CERTIFICATE REQUEST-----`

	csrWithSpiffeAndMutleipleUUIDExtensions = `
-----BEGIN CERTIFICATE REQUEST-----
MIIDBjCCAe4CAQAwEzERMA8GA1UEChMISnVqdSBvcmcwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDDUoBpvpMeo3dficyi9HWYF4D01l11Y5htk8k5l8Xt
0GEx0uvyA8sDxSPMViHwrXZwSgd/e8P7Vqc6uaPwmX89dsWK6kgniqhlzCP99BuH
E+iPDpS/T6XCroOQjstWuWC58CWRL6FfGioGEcxfLROqjA63ECCoP5dzLg6LxKYf
hmilqumJMMBtTTGkWF6mRYY2tFUW/8pr+Y+bDKttvZtFv+RZWaLAWI/2CdEtpg6J
bTCoJ1vT1D4eTC9KQbRizoWZ4QoUJJoTolGxTDEQPP8YKKDo5ADhKo6l7a/v+6XK
GlzrQGaBkZFwBuM6B9SiW0UrQS7HRsgyg8C/h23ysKZ5AgMBAAGgga0wgaoGCSqG
SIb3DQEJDjGBnDCBmTCBlgYDVR0RAQH/BIGLMIGIhit1dWlkOi8vZjU2MTY1ODkt
YzE4Ny00NDAxLTk1MDAtM2NhOTg2NTZkNDBjhixzcGlmZmU6Ly9jbHVzdGVyLmxv
Y2FsL25zL2RlZmF1bHQvc2EvaHR0cGJpboYrdXVpZDovLzYwMTdjYzg1LTczYzct
NGZiYy05YjZmLWJlZGU2MmQwNjMwMDANBgkqhkiG9w0BAQsFAAOCAQEAWrdihKJm
5aWOeti84nRa5Tca06QQE9OQ3NFAmhb1vC5AIzRbilZsnXlMDMLrRaAblLrQqcOY
qi9oRyElmRhEdXYqAy8EvKjmPLU49UjDYUdx8UCfUAYX41iHOfKiOxcd7hzJsSZT
lzhMBe1IUWSeVwINGOciHhd+HJ0fu1noUVQqR0+2njOY/y54pvE+WpBN+mYgy9yW
2aSj0M9HC+/GUeOwQsoZLV/WAeRD88ymSma0plnnlKiRl07gY+pqs2uOGkH6b8ob
jTwFdZcFwbnFbVXfIHJoRql6GEM5lahAEYYZRx+KynCzkNSU/yfPiQ1I8ZDiZHnh
H0MzLtGLR0QO8A==
-----END CERTIFICATE REQUEST-----`

	csrWithSpiffeAndUUIDANDSingleDNSExtensions = `
-----BEGIN CERTIFICATE REQUEST-----
MIIC5TCCAc0CAQAwEzERMA8GA1UEChMISnVqdSBvcmcwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQC40L4vUsltqzzrvtQ1QtwjQwaN3AQOoWWH0yUc/F1Y
dgEr0BnWwWxqtIh5L5Q9Z91jkE0NYZ3ARil+KmxACzozIZPNZOf1amg24EXqbNv5
8uy5C1VIcaHd2TU/X/lCnYnfQtKNyXt5a0J926iQWOKP1y50rTtw1XPvw7o29orB
p9ItV84MGTnYwmhqIfGBJ+MVejxo0NBHlR2LLVv2S8SMk4MYbuHj1szl7JhPFc36
RBLadJBAbutgNOLwSrmqpZ5oTdu/eiBt3mdmQQtF6rMRaNVtVXlnN+xQgoxfFx0H
KK5jnIs88BbThC+DtmONHO/uL8ygkYz14UiPMLv0nS7zAgMBAAGggYwwgYkGCSqG
SIb3DQEJDjF8MHoweAYDVR0RAQH/BG4wbIYrdXVpZDovL2Y1NjE2NTg5LWMxODct
NDQwMS05NTAwLTNjYTk4NjU2ZDQwY4Ysc3BpZmZlOi8vY2x1c3Rlci5sb2NhbC9u
cy9kZWZhdWx0L3NhL2h0dHBiaW6CD215LmFzcGVubWVzaC5pbzANBgkqhkiG9w0B
AQsFAAOCAQEAQfZWhVj6cbeyMK4cJaHbNFIHIjyUymuEV34mbTVWPNNhXLfB+Lol
fodLn7nYtsD8jcLeVXsq35saIido+2P7PUzUp4Lk56pDQo7YSLx7BXIssyVxbOIf
zChkSwIvoDtgqHudF7zQ3sKETrmeSvwpba1a/BhyisJWnccouowgmtUF3DptPVYo
Eyt7cLtT9cHrJ/WWJMDA5SE+oMg+KQL602U0kiu1iywEjVJO42xdISnlSZje0+cT
0/0juCK3kAMqvnVsr8gZ0xU125Vni6u91700jK0xgB8svHJ1LSxhQOLdnShIrSjZ
r7rQ561uYGatocWTwUlisUjoNC8NvSwqMg==
-----END CERTIFICATE REQUEST-----`
)

func TestGetCSRHosts(t *testing.T) {
	testCases := map[string]struct {
		csrPem        []byte
		expectErr     bool
		expectedHosts []string
	}{
		"CSR contains no SAN extension": {
			csrPem:        []byte(csrWithoutSanExtension),
			expectErr:     true,
			expectedHosts: []string{},
		},
		"CSR contains spiffee SAN URI extension": {
			csrPem:        []byte(csrWithSpiffeExtensionOnly),
			expectErr:     false,
			expectedHosts: []string{"spiffe://cluster.local/ns/default/sa/httpbin"},
		},
		"CSR contains spiffee SAN URI extension and a single DNS entry": {
			csrPem:        []byte(csrWithSpiffeAndSingleDNSExtensions),
			expectErr:     false,
			expectedHosts: []string{"spiffe://cluster.local/ns/default/sa/httpbin", "httpbin.org"},
		},
		"CSR contains spiffee SAN URI extension and multiple DNS entries": {
			csrPem:        []byte(csrWithSpiffeAndMutleipleDNSExtensions),
			expectErr:     false,
			expectedHosts: []string{"my.aspenmesh.io", "spiffe://cluster.local/ns/default/sa/httpbin", "httpbin.org"},
		},
		"CSR contains spiffee SAN URI extension and UUID DNS entries": {
			csrPem:        []byte(csrWithSpiffeAndUUIDDNSExtension),
			expectErr:     false,
			expectedHosts: []string{"uri://uuid://f5616589-c187-4401-9500-3ca98656d40c", "spiffe://cluster.local/ns/default/sa/httpbin"},
		},
		"CSR contains spiffee SAN URI extension and multiple UUID DNS entries": {
			csrPem:    []byte(csrWithSpiffeAndMutleipleUUIDExtensions),
			expectErr: false,
			expectedHosts: []string{
				"uri://uuid://f5616589-c187-4401-9500-3ca98656d40c", "spiffe://cluster.local/ns/default/sa/httpbin",
				"uri://uuid://6017cc85-73c7-4fbc-9b6f-bede62d06300",
			},
		},
		"CSR contains spiffee SAN URI extension, DNS name and UUID DNS entries": {
			csrPem:        []byte(csrWithSpiffeAndUUIDANDSingleDNSExtensions),
			expectErr:     false,
			expectedHosts: []string{"my.aspenmesh.io", "uri://uuid://f5616589-c187-4401-9500-3ca98656d40c", "spiffe://cluster.local/ns/default/sa/httpbin"},
		},
	}

	for k, tc := range testCases {
		t.Run(k, func(t *testing.T) {
			hosts, err := getCSRHosts(tc.csrPem)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("%s expected an error, but did not received one", k)
				}
			} else {
				if err != nil {
					t.Fatalf("%s did not expect an error, but received one (%v)", k, err)
				}
			}

			sort.Strings(tc.expectedHosts)
			if !reflect.DeepEqual(tc.expectedHosts, hosts) {
				t.Fatalf("%s did not receive (%v) expected hosts (%v)", k, hosts, tc.expectedHosts)
			}
		})
	}
}

func TestIsIdentityInHosts(t *testing.T) {
	testCases := map[string]struct {
		identities []string
		hosts      []string
		expected   bool
	}{
		"identities are empty": {
			hosts: []string{"foo"},
		},
		"hosts are empty": {
			identities: []string{"foo"},
		},
		"no hosts in identity list": {
			identities: []string{"spiffe://cluster.local/ns/default/sa/httpbin", "bar"},
			hosts:      []string{"httpbin.org", "my.aspenmesh.io"},
			expected:   false,
		},
		"spiffe hosts in identity list": {
			identities: []string{"spiffe://cluster.local/ns/default/sa/httpbin", "bar"},
			hosts:      []string{"httpbin.org", "spiffe://cluster.local/ns/default/sa/httpbin"},
			expected:   true,
		},
	}

	for k, tc := range testCases {
		t.Run(k, func(t *testing.T) {
			isFound := isIdentityInHosts(tc.identities, tc.hosts)

			if isFound != tc.expected {
				t.Errorf("%s expected to be found (%v) did not match returned value (%v)", k, tc.expected, isFound)
			}
		})
	}
}

func TestAddCSRHostsToIds(t *testing.T) {
	testCases := map[string]struct {
		identities    []string
		hosts         []string
		expectedHosts []string
	}{
		"identities are empty": {
			hosts:         []string{"foo.com"},
			expectedHosts: []string{"foo.com"},
		},
		"hosts are empty": {
			identities:    []string{"bar.gz"},
			expectedHosts: []string{"bar.gz"},
		},
		"spiffe host in identity and hosts": {
			identities:    []string{"spiffe://cluster.local/ns/default/sa/httpbin"},
			hosts:         []string{"spiffe://cluster.local/ns/default/sa/httpbin"},
			expectedHosts: []string{"spiffe://cluster.local/ns/default/sa/httpbin"},
		},
		"spiffe host in identity and hosts has DNS name": {
			identities:    []string{"spiffe://cluster.local/ns/default/sa/httpbin"},
			hosts:         []string{"spiffe://cluster.local/ns/default/sa/httpbin", "aspenmesh.io"},
			expectedHosts: []string{"spiffe://cluster.local/ns/default/sa/httpbin", "aspenmesh.io"},
		},
		"both identity and hosts contain multiple entries": {
			identities:    []string{"my.aspenmesh.io", "aspenmesh.io", "spiffe://cluster.local/ns/default/sa/httpbin"},
			hosts:         []string{"spiffe://cluster.local/ns/default/sa/httpbin", "aspenmesh.io"},
			expectedHosts: []string{"spiffe://cluster.local/ns/default/sa/httpbin", "aspenmesh.io", "my.aspenmesh.io"},
		},
	}

	for k, tc := range testCases {
		t.Run(k, func(t *testing.T) {
			hosts := addCSRHostsToIds(tc.identities, tc.hosts)

			sort.Strings(tc.expectedHosts)
			if !reflect.DeepEqual(tc.expectedHosts, hosts) {
				t.Errorf("%s expected hosts (%v) did not match returned value (%v)", k, tc.expectedHosts, hosts)
			}
		})
	}
}
