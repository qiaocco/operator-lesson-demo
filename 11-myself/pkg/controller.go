package pkg

import (
	"context"
	"reflect"
	"time"

	v14 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/util/wait"

	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v12 "k8s.io/api/networking/v1"

	"k8s.io/apimachinery/pkg/util/runtime"

	informer "k8s.io/client-go/informers/core/v1"
	netInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	coreLister "k8s.io/client-go/listers/core/v1"
	v1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	workNum  = 5
	maxRetry = 10
)

type controller struct {
	client        kubernetes.Interface
	ingressLister v1.IngressLister
	serviceLister coreLister.ServiceLister
	queue         workqueue.RateLimitingInterface
}

func (c *controller) addService(obj interface{}) {
	c.enqueue(obj)
}

func (c *controller) updateService(oldObj interface{}, newObj interface{}) {
	// todo 比较annotation
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}
	c.enqueue(newObj)
}

func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	c.queue.Add(key)
}

func (c *controller) deleteIngress(obj interface{}) {
	ingress := obj.(*v12.Ingress)
	// ingress对应的service
	ownerReference := v13.GetControllerOf(ingress)
	if ownerReference == nil {
		return
	}

	// 判断是否是service
	if ownerReference.Kind != "Service" {
		return
	}

	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

func (c *controller) Run(stopChan <-chan struct{}) {
	for i := 0; i < workNum; i++ {
		go wait.Until(c.worker, time.Minute, stopChan) // 始终有5个在执行
	}
	<-stopChan
}

func (c *controller) worker() {
	for c.processNextWorkItem() {

	}
}

func (c *controller) processNextWorkItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(item)

	key := item.(string)
	err := c.syncService(key)
	if err != nil {
		c.handleError(key, err)
	}
	return true
}

func (c *controller) syncService(item string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(item)
	if err != nil {
		return err
	}

	// 删除
	service, err := c.serviceLister.Services(namespace).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	// 新增和删除
	_, ok := service.GetAnnotations()["ingress/http"]
	ingress, err := c.ingressLister.Ingresses(namespace).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if ok && errors.IsNotFound(err) {
		// create ingress
		var ig = c.constructIngress(service)
		_, err = c.client.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ig, v13.CreateOptions{})
		if err != nil {
			return err
		}

	} else if !ok && ingress != nil {
		// delete ingress
		err = c.client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), name, v13.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil

}

func (c *controller) handleError(key string, err error) {
	if c.queue.NumRequeues(key) <= maxRetry {
		// 出现错误后，等一段时间重试
		c.queue.AddRateLimited(key)
		return
	}

	runtime.HandleError(err)
	c.queue.Forget(key)
}

func (c *controller) constructIngress(service *v14.Service) *v12.Ingress {
	pathType := v12.PathTypePrefix

	icn := "nginx"
	return &v12.Ingress{
		ObjectMeta: v13.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
			Annotations: map[string]string{
				"ingress/http": "true",
			},
			OwnerReferences: []v13.OwnerReference{
				*v13.NewControllerRef(service, v14.SchemeGroupVersion.WithKind("Service"))},
		},
		Spec: v12.IngressSpec{
			IngressClassName: &icn,
			Rules: []v12.IngressRule{
				{
					Host: "www.example.com",
					IngressRuleValue: v12.IngressRuleValue{
						HTTP: &v12.HTTPIngressRuleValue{
							Paths: []v12.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: v12.IngressBackend{
										Service: &v12.IngressServiceBackend{
											Name: service.Name,
											Port: v12.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func NewController(client kubernetes.Interface, serviceInformer informer.ServiceInformer, ingressInformer netInformer.IngressInformer) controller {
	c := controller{
		client:        client,
		ingressLister: ingressInformer.Lister(), // 其实就是indexer 获取资源对象状态，避免与api server交互
		serviceLister: serviceInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})
	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})
	return c
}
