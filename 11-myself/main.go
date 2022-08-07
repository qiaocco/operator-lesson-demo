package main

import (
	"log"

	"github.com/qiaocc/client-go-demo/11/pkg"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. config
	// 2. client
	// 3. informer
	// 4. add event handler
	// 5. informer.Start

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		// 集群内部
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}
	factory := informers.NewSharedInformerFactory(clientset, 0)
	// 创建我们关注的资源对象
	serviceInformer := factory.Core().V1().Services()
	ingresseInformer := factory.Networking().V1().Ingresses()
	controller := pkg.NewController(clientset, serviceInformer, ingresseInformer)

	stopChan := make(chan struct{})
	factory.Start(stopChan)
	factory.WaitForCacheSync(stopChan)

	controller.Run(stopChan)
}
