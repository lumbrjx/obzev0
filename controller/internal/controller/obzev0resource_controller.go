package controller

import (
	"context"
	"fmt"
	"log"
	"os"

	v1 "obzev0/controller/api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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

	labelSelector := "app=grpc-server"
	daemonSetName := "grpc-server-daemonset"

	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	if daemonSetName != "" {
		ds, err := clientset.AppsV1().
			DaemonSets("default").
			Get(context.TODO(), daemonSetName, metav1.GetOptions{})
		if err != nil {
			setupLog.Error(err, "unable to get specific DaemonSet")
			os.Exit(1)
		}
		fmt.Printf("DaemonSet Name: %s, Namespace: %s\n", ds.Name, ds.Namespace)
		listOptions.LabelSelector = fmt.Sprintf(
			"app=%s",
			ds.Spec.Template.Labels["app"],
		)
	} else {
		// If no specific DaemonSet name, list DaemonSets based on label selector
		fmt.Println("Listing DaemonSets based on label selector:", labelSelector)
		daemonSets, err := clientset.AppsV1().
			DaemonSets("").
			List(context.TODO(), listOptions)
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf(
			"There are %d daemonsets in the cluster\n",
			len(daemonSets.Items),
		)

		for _, ds := range daemonSets.Items {
			fmt.Printf("DaemonSet Name: %s, Namespace: %s\n", ds.Name, ds.Namespace)
		}
	}

	pods, err := clientset.CoreV1().
		Pods("default").
		List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Status.PodIP)
		if pod.Status.Phase == corev1.PodRunning {
			ip := pod.Status.PodIP
			port := "50051"
			address := fmt.Sprintf("%s:%s", ip, port)
			fmt.Printf("Connecting to gRPC server at %s\n", address)

			conn, err := grpc.NewClient(
				address,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				log.Printf("Failed to connect to %s: %v\n", address, err)
				continue
			}

			crInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					handleAdd(obj, conn)
				},
				UpdateFunc: func(oldObj, newObj interface{}) {
					handleUpdate(newObj, conn)
				},
				DeleteFunc: func(obj interface{}) {
					handleDelete(obj)
				},
			})

			fmt.Printf(
				"Successfully connected to gRPC server at %s\n",
				address,
			)
		}
	}
}
