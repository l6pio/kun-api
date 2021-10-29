package service

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"l6p.io/kun/api/pkg/core"
)

func FindNamespaceGarbage(conf *core.Config) error {
	namespaces, err := conf.KubeClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, namespace := range namespaces.Items {
		println(namespace.Name)
	}
	return nil
}
