package main

import (
	clientset "count/generated/clientset/versioned"
	"count/generated/clientset/versioned/scheme"
	informers "count/generated/informers/externalversions/count/v1"
	listers "count/generated/listers/count/v1"
	countv1 "count/pkg/apis/count/v1"
	"fmt"
	"github.com/golang/glog"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	typeCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"time"
)

const (
	controllerAgentName   = "count-controller"
	SuccessSynced         = "Synced"
	MessageResourceSynced = "Student synced successfully"
)

// Controller is the controller implementation for Student resources
type Controller struct {
	// kubeClientSet is a standard kubernetes clientset
	kubeClientSet kubernetes.Interface
	// countClientSet is a clientset for our own API group
	countClientSet clientset.Interface
	countLister    listers.CountLister

	countSynced cache.InformerSynced
	workQueue   workqueue.RateLimitingInterface

	recorder record.EventRecorder
}

func NewController(kubeClientSet kubernetes.Interface, countClientSet clientset.Interface, countInformer informers.CountInformer) *Controller {
	utilruntime.Must(countv1.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typeCoreV1.EventSinkImpl{
		Interface: kubeClientSet.CoreV1().Events(""),
	})

	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, coreV1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeClientSet:  kubeClientSet,
		countClientSet: countClientSet,
		countLister:    countInformer.Lister(),
		countSynced:    countInformer.Informer().HasSynced,
		workQueue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Counts"),
		recorder:       recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when Count resources change
	countInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueCount,
		UpdateFunc: func(old, new interface{}) {
			oldCount := old.(*countv1.Count)
			newCount := new.(*countv1.Count)
			// 如果版本一致，那么没有实际更新的操作，立即返回
			if oldCount.ResourceVersion == newCount.ResourceVersion {
				return
			}
			controller.enqueueCount(new)
		},
		DeleteFunc: controller.enqueueStudentForDelete,
	})

	return controller
}

// 在此处开始controller的业务
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workQueue.ShutDown()

	glog.Info("开始controller业务，开始一次缓存数据同步")
	if ok := cache.WaitForCacheSync(stopCh, c.countSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("worker启动")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("worker已经启动")
	<-stopCh
	glog.Info("worker已经结束")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// 取数据处理
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workQueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		defer c.workQueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workQueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// 在syncHandler中处理业务
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}

		c.workQueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}
	return true
}

// 处理
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// 从缓存中取对象
	count, err := c.countLister.Counts(namespace).Get(name)
	if err != nil {
		// 如果Count对象被删除了，就会走到这里，所以应该在这里加入执行
		if errors.IsNotFound(err) {
			glog.Infof("Count对象被删除，请在这里执行实际的删除业务: %s/%s ...", namespace, name)
			return nil
		}

		utilruntime.HandleError(fmt.Errorf("failed to list student by: %s/%s", namespace, name))
		return err
	}

	glog.Infof("这里是Count对象的期望状态: %#v ...", count)
	glog.Infof("实际状态是从业务层面得到的，此处应该去的实际状态，与期望状态做对比，并根据差异做出响应(新增或者删除)")

	c.recorder.Event(count, coreV1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// 数据先放入缓存，再入队列
func (c *Controller) enqueueCount(obj interface{}) {
	var key string
	var err error
	fmt.Println("enqueueCount: obj", obj)
	// 将对象放入缓存
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	// 将key放入队列
	c.workQueue.AddRateLimited(key)
}

// 删除操作
func (c *Controller) enqueueStudentForDelete(obj interface{}) {
	var key string
	var err error
	// 从缓存中删除指定对象
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	//再将key放入队列
	c.workQueue.AddRateLimited(key)
}
