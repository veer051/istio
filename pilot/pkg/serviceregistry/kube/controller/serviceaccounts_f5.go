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

package controller

import (
	"sync"

	v1 "k8s.io/api/core/v1"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/serviceregistry/kube"
	"istio.io/istio/pkg/kube/kclient"
)

// ServiceAccountCache is an eventually consistent ServiceAccount cache
type ServiceAccountCache struct {
	sas kclient.Client[*v1.ServiceAccount]

	sync.RWMutex

	// ServiceAccount namespace/name and their annotations
	serviceAccounts map[string]map[string]string
}

func newServiceAccountCache(sas kclient.Client[*v1.ServiceAccount]) *ServiceAccountCache {
	out := &ServiceAccountCache{
		sas:             sas,
		serviceAccounts: make(map[string]map[string]string),
	}

	return out
}

// onEvent updates the serviceAccounts map when a Service Account has triggered
// a create/update/delete operation
// func (sa *ServiceAccountCache) onEvent(curr interface{}, ev model.Event) error {
func (sa *ServiceAccountCache) onEvent(_, sac *v1.ServiceAccount, ev model.Event) error {
	sa.Lock()
	defer sa.Unlock()

	saNameByNamespace := kube.KeyFunc(sac.Name, sac.Namespace)
	switch ev {
	case model.EventAdd, model.EventUpdate:
		sa.serviceAccounts[saNameByNamespace] = sac.GetAnnotations()
	case model.EventDelete:
		delete(sa.serviceAccounts, saNameByNamespace)
	}
	return nil
}

func (sa *ServiceAccountCache) getServiceAccountAnnotations(name, namespace string) map[string]string {
	sa.RLock()
	defer sa.RUnlock()

	saNameByNamespace := kube.KeyFunc(name, namespace)
	if annotations, ok := sa.serviceAccounts[saNameByNamespace]; ok {
		return annotations
	}

	return nil
}
