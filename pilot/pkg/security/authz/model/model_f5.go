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

package model

import "strconv"

const (
	// F5 Aspenmesh XFCC Annotation
	xfccAnnotation          = "authz.contains.aspenmesh.io/xfcc"
	f5attrXfccRequestHeader = "x-forwarded-client-cert"
)

func useXfccHeader(annotations map[string]string) bool {
	f5useXfcc := false
	if xfcc, found := annotations[xfccAnnotation]; found {
		// strconv.Parsebool returns false on error
		f5useXfcc, _ = strconv.ParseBool(xfcc)
	}
	return f5useXfcc
}
