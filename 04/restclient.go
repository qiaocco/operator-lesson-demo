package main

import (
	"context"
	"fmt"

	v12 "k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	config.GroupVersion = &v12.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	config.APIPath = "/apis"
	// deployments /apis/apps/v1/namespaces/{ns}/deployments

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	deploys := v12.DeploymentList{}
	err = restClient.Get().
		Namespace("kube-system").
		Resource("deployments").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec). // 参数及参数的序列化方式
		Do(context.TODO()).
		Into(&deploys)
	if err != nil {
		panic(err)
	}

	for _, deploy := range deploys.Items {
		fmt.Printf("Namespace=%v, name=%v", deploy.Namespace, deploy.Name)
	}
}
