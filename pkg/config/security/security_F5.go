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

package security

import (
	"unicode"

	"istio.io/istio/pkg/util/sets"
)

// ValidCarrierGradeCipherSuites contains a list of all ciphers supported in Gateway.server.tls.cipherSuites
// Extracted from: `bssl ciphers -openssl-name ALL | rg -v PSK`
var ValidCarrierGradeCipherSuites = sets.New(
	"ECDHE-ECDSA-AES256-GCM-SHA384", // TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
	"ECDHE-ECDSA-CHACHA20-POLY1305", // TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256
	"ECDHE-RSA-AES256-GCM-SHA384",   // TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
	"ECDHE-RSA-CHACHA20-POLY1305",   // TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
	// TLS_DHE_RSA_WITH_AES_256_GCM_SHA384 not supported
	// TLS_DHE_RSA_WITH_AES_256_CCM not supported
	"ECDHE-ECDSA-AES128-GCM-SHA256", // TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
	"ECDHE-RSA-AES128-GCM-SHA256",   // TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
)

func IsValidCarrierGradeCipherSuite(cs string) bool {
	if cs == "" || cs == "ALL" {
		return true
	}
	if !unicode.IsNumber(rune(cs[0])) && !unicode.IsLetter(rune(cs[0])) {
		// Not all of these are correct, but this is needed to support advanced cases like - and + operators
		// without needing to parse the full expression
		return true
	}
	return ValidCarrierGradeCipherSuites.Contains(cs)
}
