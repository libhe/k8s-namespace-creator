package main

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	namespace := "test-namespace"
	kubeconfig := clientcmd.RecommendedConfigPathEnvVar

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Failed to build kubeconfig: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Failed to create clientset: %v\n", err)
		os.Exit(1)
	}

	if err := createNamespace(clientset, namespace); err != nil {
		fmt.Printf("Error creating namespace: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Namespace '%s' created successfully!\n", namespace)
}

func createNamespace(clientset *kubernetes.Clientset, namespace string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return err
}

func deleteNamespace(clientset *kubernetes.Clientset, namespace string) error {
	return clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
}
