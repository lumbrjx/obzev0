package controller

import (
	"context"
	"os"

	v1 "github.com/lumbrjx/obzev0/api/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
)

var setupLog = ctrl.Log.WithName("setup")

func SetupInformers(mgr ctrl.Manager) {
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

	// Start the informer
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
