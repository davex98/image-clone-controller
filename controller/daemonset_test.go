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

var _ = Describe("Daemonset controller", func() {

	const (
		DaemonName      = "busybox"
		DaemonNamespace = "default"
		DaemonImage     = "busybox"

		timeout  = time.Second * 10
		interval = time.Second * 1
	)

	When("creating daemonset with invalid image", func() {
		It("should change the image to use one from the private repository", func() {
			By("editing daemonset container images")
			ctx := context.Background()
			daemonset := &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      DaemonName,
					Namespace: DaemonNamespace,
				},
				Spec: appsv1.DaemonSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app":"busybox"},
					},
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app":"busybox"}},
						Spec:       v1.PodSpec{
							Containers: []v1.Container{{Image: DaemonImage, Name: "busybox"}},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, daemonset)).Should(Succeed())


			daemonLookupKey := types.NamespacedName{Name: DaemonName, Namespace: DaemonNamespace}
			createdDaemon := &appsv1.DaemonSet{}

			Eventually(func() string {
				err := k8sClient.Get(ctx, daemonLookupKey, createdDaemon)
				if err != nil {
					return ""
				}
				image := createdDaemon.Spec.Template.Spec.Containers[0].Image
				return image
			}, timeout, interval).Should(Equal("burghardtkubermatic/busybox"))
		})
	})

	When("daemonset has a proper image", func() {
		It("should do nothing", func() {
			ctx := context.Background()
			daemonset := &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kubermatic",
					Namespace: DaemonNamespace,
				},
				Spec: appsv1.DaemonSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app":"busybox"},
					},
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app":"busybox"}},
						Spec:       v1.PodSpec{
							Containers: []v1.Container{{Image: "burghardtkubermatic/busybox", Name: "busybox"}},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, daemonset)).Should(Succeed())


			daemonLookupKey := types.NamespacedName{Name: "kubermatic", Namespace: DaemonNamespace}
			createdDaemon := &appsv1.DaemonSet{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, daemonLookupKey, createdDaemon)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Eventually(func() string {
				err := k8sClient.Get(ctx, daemonLookupKey, createdDaemon)
				if err != nil {
					return ""
				}

				image := createdDaemon.Spec.Template.Spec.Containers[0].Image
				return image
			}, timeout, interval).Should(Equal("burghardtkubermatic/busybox"))

		})
	})
})

