package client_go

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func Connect() {
	config := &rest.Config{
		// 指定API服务器的地址
		Host: "https://192.168.122.127:6443",
		// 指定用于身份验证的Bearer Token
		BearerToken: "a1232123131231231231231231233112",
		// 忽略SSL证书验证，仅在测试环境下使用
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}

	// 使用配置创建Kubernetes客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 使用clientset进行操作，例如列出所有命名空间
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}
}
