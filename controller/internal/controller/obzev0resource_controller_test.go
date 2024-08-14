/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	batchv1 "obzev0/controller/api/v1"
)

var _ = Describe("Obzev0Resource Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		obzev0resource := &batchv1.Obzev0Resource{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind Obzev0Resource")
			err := k8sClient.Get(ctx, typeNamespacedName, obzev0resource)
			if err != nil && errors.IsNotFound(err) {
				resource := &batchv1.Obzev0Resource{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					// TODO(user): Specify other spec details if needed.
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &batchv1.Obzev0Resource{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Obzev0Resource")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		// It("should successfully reconcile the resource", func() {
		// 	By("Reconciling the created resource")
		// 	controllerReconciler := &Obzev0ResourceReconciler{
		// 		Client: k8sClient,
		// 		Scheme: k8sClient.Scheme(),
		// 	}
		//
		// 	_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
		// 		NamespacedName: typeNamespacedName,
		// 	})
		// 	Expect(err).NotTo(HaveOccurred())
		// 	// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
		// 	// Example: If you expect a certain status condition after reconciliation, verify it here.
		// })
	})
})
