package controller_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Deployment controller", func() {

	const (
		DeploymentName      = "nginx"
		DeploymentNamespace = "default"
		DeploymentImage = "nginx"

		timeout  = time.Second * 10
		interval = time.Second * 1
	)

	When("creating deployment with invalid images", func() {
		It("should change the images to use ones from the private repository", func() {
			By("editing deployment containers images")
			ctx := context.Background()
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      DeploymentName,
					Namespace: DeploymentNamespace,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app":"nginx"},
					},
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app":"nginx"}},
						Spec:       v1.PodSpec{
							Containers: []v1.Container{{Image: DeploymentImage, Name: "nginx"}},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, deployment)).Should(Succeed())
			time.Sleep(time.Second * 5)

			deploymentLookupKey := types.NamespacedName{Name: DeploymentName, Namespace: DeploymentNamespace}
			createdDeployment := &appsv1.Deployment{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, deploymentLookupKey, createdDeployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(createdDeployment.Spec.Template.Spec.Containers[0].Image).Should(Equal("burghardtkubermatic/nginx"))
		})
	})

	When("deployment has a proper image", func() {
		It("should do nothing", func() {
			ctx := context.Background()
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kubermatic",
					Namespace: DeploymentNamespace,
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app":"nginx"},
					},
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app":"nginx"}},
						Spec:       v1.PodSpec{
							Containers: []v1.Container{{Image: "burghardtkubermatic/nginx", Name: "nginx"}},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, deployment)).Should(Succeed())

			deploymentLookupKey := types.NamespacedName{Name: "kubermatic", Namespace: DeploymentNamespace}
			createdDeployment := &appsv1.Deployment{}

			Eventually(func() string {
				err := k8sClient.Get(ctx, deploymentLookupKey, createdDeployment)
				if err != nil {
					return ""
				}

				image := createdDeployment.Spec.Template.Spec.Containers[0].Image
				return image
			}, timeout, interval).Should(Equal("burghardtkubermatic/nginx"))

		})
	})


})

