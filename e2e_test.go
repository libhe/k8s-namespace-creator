package main

import (
	"context"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientset *kubernetes.Clientset
	namespace = "test-namespace"
)

func TestE2E(t *testing.T) {
	ginkgo.RunSpecs(t, "E2E Suite")
}

var _ = ginkgo.Describe("Namespace E2E Tests", func() {
	ginkgo.BeforeSuite(func() {
		// Load kubeconfig from the default location
		kubeconfig := clientcmd.RecommendedHomeFile // Directly use the default kubeconfig path
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		clientset, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.Context("when creating and checking the namespace", func() {
		ginkgo.It("should exist after creation", func() {
			// Create the namespace using the function from main package
			err := createNamespace(clientset, namespace)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Verify the namespace exists
			_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, v1.GetOptions{})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.AfterSuite(func() {
		// Cleanup: Delete the namespace after tests
		err := deleteNamespace(clientset, namespace)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})
})
