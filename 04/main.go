package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/rest"
)

func main() {

	//RESTClient
	//config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	config.APIPath = "/api"

	// client
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	// get data
	pod := v1.Pod{}
	err = restClient.Get().Namespace("default").Resource("pods").Name("wx-notify").Do(context.TODO()).Into(&pod)
	if err != nil {
		fmt.Printf("%+v", err)
	} else {
		fmt.Printf(pod.Name)
	}
	// config
	//config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	//if err != nil {
	//	panic(err)
	//}
	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	panic(err)
	//}
	//coreV1 := clientset.CoreV1()
	//pod, err := coreV1.Pods("default").Get(context.TODO(), "test", v1.GetOptions{})
	//if err != nil {
	//	println(err)
	//} else {
	//	println(pod.Name)
	//}
}
