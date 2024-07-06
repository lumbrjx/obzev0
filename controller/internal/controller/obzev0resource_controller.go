package controller

import (
	"context"
	v1 "github.com/lumbrjx/obzev0/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"log"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
)

var setupLog = ctrl.Log.WithName("setup")

func SetupInformers(mgr ctrl.Manager) {
	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		log.Fatal(err, "unable to create clientset")
		os.Exit(1)
	}
	setupLog.Info("Setting up informers")
	ctx := context.Background()
	crInformer, err := mgr.GetCache().GetInformer(ctx, &v1.Obzev0Resource{})
	if err != nil {
		setupLog.Error(err, "unable to create CR informer")
		os.Exit(1)
	}
	setupLog.Info("CR informer created")

	crInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			handleAdd(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			handleUpdate(oldObj, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			handleDelete(obj)
		},
	})
	setupLog.Info("Event handlers added to CR informer")
	listNodes(clientset) // Start the informer
	// go func() {
	// 	setupLog.Info("Starting informer")
	// 	if err := mgr.GetCache().Start(ctx); err != nil {
	// 		setupLog.Error(err, "error starting informer")
	// 		os.Exit(1)
	// 	}
	// }()
}

func handleAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.Errorf("Error getting key for object: %v", err)
		return
	}
	klog.Infof("Custom Resource added: %s", key)
}

func handleUpdate(oldObj, newObj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err != nil {
		klog.Errorf("Error getting key for object: %v", err)
		return
	}
	klog.Infof("Custom Resource updated: %s", key)
}

func handleDelete(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.Errorf("Error getting key for object: %v", err)
		return
	}
	klog.Infof("Custom Resource deleted: %s", key)
}
func listNodes(clientset *kubernetes.Clientset) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Fatalf("Error listing nodes: %v", err)
	}

	klog.Info("Listing all nodes in the cluster:")
	for _, node := range nodes.Items {
		klog.Infof("Node Name: %s", node.Name)
	}
}
