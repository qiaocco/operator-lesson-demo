package main

import (
	"context"

	"k8s.io/client-go/tools/cache"

	"github.com/operator-crd/pkg/generated/informers/externalversions"

	clientset "github.com/operator-crd/pkg/generated/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	client, err := clientset.NewForConfig(config)

	crdv1 := client.CrdV1()
	list, err := crdv1.Foos("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, item := range list.Items {
		println(item.Name)
	}

	factory := externalversions.NewSharedInformerFactory(client, 0)
	factory.Crd().V1().Foos().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// todo
		},
	})

}
