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

package aggregate

import (
	"istio.io/istio/pilot/pkg/serviceregistry"
	"istio.io/istio/pilot/pkg/serviceregistry/provider"
	"istio.io/istio/pkg/cluster"
	"istio.io/pkg/log"
)

// GetRegistry returns a copy of the registry associated with clusterID
func (c *Controller) GetRegistry(clusterID cluster.ID, provider provider.ID) serviceregistry.Instance {
	c.storeLock.RLock()
	defer c.storeLock.RUnlock()

	idx, ok := c.getRegistryIndex(clusterID, provider)
	if !ok {
		log.Warnf("failed to retrieve registry for clusterID:[%v], provider:[%v], returning default registry", clusterID, provider)
	}

	return c.registries[idx]
}
