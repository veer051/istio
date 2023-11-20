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
	"context"
	"reflect"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/serviceregistry/kube"
	"istio.io/istio/pilot/pkg/serviceregistry/util/xdsfake"
)

var ctx = context.Background()

func TestServiceAccountCache(t *testing.T) {
	t.Run("fakeApiserver", func(t *testing.T) {
		t.Parallel()
		c, fx := NewFakeControllerWithOptions(t, FakeControllerOptions{Mode: EndpointsOnly})
		go c.Run(c.stop)
		testServiceAccountCache(t, c, fx)
	})
}

// Checks that events from the watcher create the proper internal structures
func TestServiceAccountsCacheEvents(t *testing.T) {
	t.Parallel()

	NewFakeControllerWithOptions(t, FakeControllerOptions{Mode: EndpointsOnly})
	svcAcctCache := newServiceAccountCache(nil)

	f := svcAcctCache.onEvent
	ns := "default"
	svcAcctName := "sa-1"
	svcAcct1 := metav1.ObjectMeta{Name: svcAcctName, Namespace: ns}

	annotations := map[string]string{"certificate.aspenmesh.io/customFields": `{ "SAN": { "DNS": [ "foo.com" ] , "URI": [ "http://test.example.com/get" ] } }`}
	svcAcct1Annotated := metav1.ObjectMeta{Name: svcAcctName, Namespace: ns, Annotations: annotations}

	if err := f(nil, &v1.ServiceAccount{ObjectMeta: svcAcct1}, model.EventAdd); err != nil {
		t.Error(err)
	}

	mdAnnotations := svcAcctCache.getServiceAccountAnnotations(svcAcctName, ns)
	if mdAnnotations != nil {
		t.Error("ServiceAccount create failed")
	}

	if err := f(nil, &v1.ServiceAccount{ObjectMeta: svcAcct1Annotated}, model.EventUpdate); err != nil {
		t.Error(err)
	}

	mdAnnotations = svcAcctCache.getServiceAccountAnnotations(svcAcctName, ns)
	if !reflect.DeepEqual(mdAnnotations, annotations) {
		t.Error("ServiceAccount update failed")
	}

	if err := f(nil, &v1.ServiceAccount{ObjectMeta: svcAcct1Annotated}, model.EventDelete); err != nil {
		t.Error(err)
	}

	mdAnnotations = svcAcctCache.getServiceAccountAnnotations(svcAcctName, ns)
	if mdAnnotations != nil {
		t.Error("ServiceAccount delete failed")
	}
}

func testServiceAccountCache(t *testing.T, c *FakeController, fx *xdsfake.Updater) {
	initServiceAcctTestEnv(t, c.client.Kube(), fx)

	serviceAcctName := "sa-1"
	serviceAcctNamespace := "default"
	annotations := map[string]string{"certificate.aspenmesh.io/customFields": `{ "SAN": { "DNS": [ "foo.com" ] , "URI": [ "http://test.example.com/get" ] } }`}

	svcAccounts := []*v1.ServiceAccount{
		generateServiceAccount(serviceAcctName, serviceAcctNamespace, annotations),
	}

	for _, svcAccount := range svcAccounts {
		svcAccount := svcAccount
		addServiceAccounts(t, c, svcAccount)
		// wait for service accounts events
		err := waitForServiceAccount(c, svcAccount.Name, svcAccount.Namespace)
		if err != nil {
			t.Errorf("waitForServiceAccount error: %v", err)
		}
	}

	ann := c.sa.getServiceAccountAnnotations(serviceAcctName, serviceAcctNamespace)
	if !reflect.DeepEqual(annotations, ann) {
		t.Errorf("expected annotations don't match received: %v, %v", annotations, ann)
	}
}

// Prepare k8s. This can be used in multiple tests, to
// avoid duplicating creation, which can be tricky. It can be used with the fake or
// standalone apiserver.
func initServiceAcctTestEnv(t *testing.T, ki kubernetes.Interface, fx *xdsfake.Updater) {
	cleanupSvcAcct(ki)

	for _, n := range []string{"default"} {
		_, err := ki.CoreV1().Namespaces().Create(context.TODO(), &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: n,
				Labels: map[string]string{
					"istio-injection": "enabled",
				},
			},
		}, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("failed creating test namespace: %v", err)
		}

		// K8S 1.10 also checks if service account exists
		_, err = ki.CoreV1().ServiceAccounts(n).Create(context.TODO(), &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
				Annotations: map[string]string{
					"kubernetes.io/enforce-mountable-secrets": "false",
				},
			},
			Secrets: []v1.ObjectReference{
				{
					Name: "default-token-2",
					UID:  "1",
				},
			},
		}, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("failed creating test service account: %v", err)
		}

		_, err = ki.CoreV1().Secrets(n).Create(context.TODO(), &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default-token-2",
				Annotations: map[string]string{
					"kubernetes.io/service-account.name": "default",
					"kubernetes.io/service-account.uid":  "1",
				},
			},
			Type: v1.SecretTypeServiceAccountToken,
			Data: map[string][]byte{
				"token": []byte("1"),
			},
		}, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("failed creating test secret: %v", err)
		}
	}
	fx.Clear()
}

func cleanupSvcAcct(ki kubernetes.Interface) {
	for _, n := range []string{"default", "ns-1"} {
		n := n
		svcAccts, err := ki.CoreV1().ServiceAccounts(n).List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			// Make sure the svcAccts don't exist
			for _, svcAcct := range svcAccts.Items {
				_ = ki.CoreV1().ServiceAccounts(svcAcct.Namespace).Delete(context.TODO(), svcAcct.Name, metav1.DeleteOptions{})
			}
		}
	}
}

func generateServiceAccount(name, namespace string, annotations map[string]string) *v1.ServiceAccount {
	return &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: annotations,
			Namespace:   namespace,
		},
	}
}

func addServiceAccounts(t *testing.T, controller *FakeController, svcAccts ...*v1.ServiceAccount) {
	for _, svcAcct := range svcAccts {
		p, _ := controller.client.Kube().CoreV1().ServiceAccounts(svcAcct.Namespace).Get(context.TODO(), svcAcct.Name, metav1.GetOptions{})
		var _ *v1.ServiceAccount
		var err error
		if p == nil {
			_, err = controller.client.Kube().CoreV1().ServiceAccounts(svcAcct.Namespace).Create(context.TODO(), svcAcct, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Cannot create %s in namespace %s (error: %v)", svcAcct.GetName(), svcAcct.GetNamespace(), err)
			}
		} else {
			_, err = controller.client.Kube().CoreV1().ServiceAccounts(svcAcct.Namespace).Update(context.TODO(), svcAcct, metav1.UpdateOptions{})
			if err != nil {
				t.Fatalf("Cannot update %s in namespace %s (error: %v)", svcAcct.GetName(), svcAcct.GetNamespace(), err)
			}
		}
	}
}

func waitForServiceAccount(c *FakeController, name, namespace string) error {
	saNameByNamespace := kube.KeyFunc(name, namespace)
	return wait.PollUntilContextTimeout(ctx, 10*time.Millisecond, 5*time.Second, true, func(context.Context) (bool, error) {
		c.sa.RLock()
		defer c.sa.RUnlock()

		if _, ok := c.sa.serviceAccounts[saNameByNamespace]; ok {
			return true, nil
		}
		return false, nil
	})
}
