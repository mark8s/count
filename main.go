package main

import (
	clientset "count/generated/clientset/versioned"
	"count/generated/informers/externalversions"
	"count/pkg/signals"
	"flag"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

func main() {
	flag.Parse()

	// 处理信号量
	stopCh := signals.SetupSignalHandler()

	// 处理入参
	cfg, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	countClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	countInformerFactory := externalversions.NewSharedInformerFactory(countClient, time.Second*30)
	//得到controller
	controller := NewController(kubeClient, countClient,
		countInformerFactory.Mark8s().V1().Counts())

	//启动informer
	go countInformerFactory.Start(stopCh)

	//controller开始处理消息
	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}

}
