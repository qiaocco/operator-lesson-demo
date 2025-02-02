package main

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	coreV1 := clientset.CoreV1()
	pod, err := coreV1.Pods("default").Get(context.TODO(), "wx-notify", v1.GetOptions{})
	if err != nil {
		println(err)
	} else {
		println(pod.Name)
	}
}
