// package controller
//
// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
//
// 	"obzev0/common/proto/latency"
// 	pb "obzev0/common/proto/latency"
//
// 	v1 "obzev0/controller/api/v1"
//
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/tools/cache"
// 	"k8s.io/klog/v2"
// 	ctrl "sigs.k8s.io/controller-runtime"
// )
//
// var setupLog = ctrl.Log.WithName("setup")
//
// func SetupInformers(mgr ctrl.Manager) {
// 	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
// 	if err != nil {
// 		log.Fatal(err, "unable to create clientset")
// 		os.Exit(1)
// 	}
// 	setupLog.Info("Setting up informers")
// 	ctx := context.Background()
// 	crInformer, err := mgr.GetCache().GetInformer(ctx, &v1.Obzev0Resource{})
// 	if err != nil {
// 		setupLog.Error(err, "unable to create CR informer")
// 		os.Exit(1)
// 	}
// 	setupLog.Info("CR informer created")
//
// 	crInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
// 		AddFunc: func(obj interface{}) {
// 			handleAdd(obj)
// 		},
// 		UpdateFunc: func(oldObj, newObj interface{}) {
// 			handleUpdate(oldObj, newObj)
// 		},
// 		DeleteFunc: func(obj interface{}) {
// 			handleDelete(obj)
// 		},
// 	})
// 	setupLog.Info("Event handlers added to CR informer")
//
// 	// Get all DaemonSets
// 	daemonSets, err := clientset.AppsV1().
// 		DaemonSets("").
// 		List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}
//
// 	fmt.Printf(
// 		"There are %d daemonsets in the cluster\n",
// 		len(daemonSets.Items),
// 	)
//
// 	for _, ds := range daemonSets.Items {
// 		fmt.Printf("DaemonSet Name: %s, Namespace: %s\n", ds.Name, ds.Namespace)
//
// 		// List Pods for each DaemonSet
// 		pods, err := clientset.CoreV1().
// 			Pods(ds.Namespace).
// 			List(context.TODO(), metav1.ListOptions{
// 				LabelSelector: fmt.Sprintf(
// 					"app=%s",
// 					ds.Spec.Template.Labels["app"],
// 				),
// 			})
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		for _, pod := range pods.Items {
// 			fmt.Println(pod.Status.PodIP)
// 			if pod.Status.Phase == corev1.PodRunning {
//
// 				ip := pod.Status.PodIP
// 				port := "50051" // Replace with your gRPC server port
// 				address := fmt.Sprintf("%s:%s", ip, port)
// 				fmt.Printf("Connecting to gRPC server at %s\n", address)
//
// 				// Create a gRPC connection
// 				conn, err := grpc.NewClient(
// 					address,
// 					grpc.WithTransportCredentials(insecure.NewCredentials()),
// 				)
// 				if err != nil {
// 					log.Printf("Failed to connect to %s: %v\n", address, err)
// 					continue
// 				}
// 				defer conn.Close()
//
// 				// Create gRPC client and make a request
// 				client := pb.NewLatencyServiceClient(conn)
// 				response, err := client.StartTcpServer(
// 					context.Background(),
// 					&pb.RequestForTcp{Config: &latency.TcpConfig{
// 						ReqDelay: 1,
// 						ResDelay: 0,
// 						Server:   "7090",
// 						Client:   "8080",
// 					}},
// 				)
// 				if err != nil {
// 					log.Printf("Error calling gRPC method: %v\n", err)
// 				} else {
// 					fmt.Printf("Response from gRPC server: %s\n", response.Message)
// 				}
//
// 				fmt.Printf(
// 					"Successfully connected to gRPC server at %s\n",
// 					address,
// 				)
// 			}
// 		}
// 	}
// }
//
// func handleAdd(obj interface{}) {
// 	key, err := cache.MetaNamespaceKeyFunc(obj)
// 	if err != nil {
// 		klog.Errorf("Error getting key for object: %v", err)
// 		return
// 	}
// 	klog.Infof("Custom Resource added: %s", key)
// }
//
// func handleUpdate(oldObj, newObj interface{}) {
// 	key, err := cache.MetaNamespaceKeyFunc(newObj)
// 	if err != nil {
// 		klog.Errorf("Error getting key for object: %v", err)
// 		return
// 	}
// 	klog.Infof("Custom Resource updated: %s", key)
// }
//
// func handleDelete(obj interface{}) {
// 	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
// 	if err != nil {
// 		klog.Errorf("Error getting key for object: %v", err)
// 		return
// 	}
// 	klog.Infof("Custom Resource deleted: %s", key)
// }
//
// func listNodes(clientset *kubernetes.Clientset) {
// 	nodes, err := clientset.CoreV1().
// 		Nodes().
// 		List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		klog.Fatalf("Error listing nodes: %v", err)
// 	}
//
// 	klog.Info("Listing all nodes in the cluster:")
// 	for _, node := range nodes.Items {
// 		klog.Infof("Node Name: %s", node.Name)
// 	}
// }
//

package controller

import (
	"context"
	"fmt"
	"log"
	"os"

	"obzev0/common/proto/latency"
	pb "obzev0/common/proto/latency"

	v1 "obzev0/controller/api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
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
			defer conn.Close()

			// Create gRPC client and make a request (temp)
			client := pb.NewLatencyServiceClient(conn)
			response, err := client.StartTcpServer(
				context.Background(),
				&pb.RequestForTcp{Config: &latency.TcpConfig{

					ReqDelay: 1,
					ResDelay: 1,
					Server:   "7090",
					Client:   "8080",
				}},
			)
			if err != nil {
				log.Printf("Error calling gRPC method: %v\n", err)
			} else {
				fmt.Printf("Response from gRPC server: %s\n", response.Message)
			}

			fmt.Printf(
				"Successfully connected to gRPC server at %s\n",
				address,
			)
		}
	}
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
	nodes, err := clientset.CoreV1().
		Nodes().
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Fatalf("Error listing nodes: %v", err)
	}

	klog.Info("Listing all nodes in the cluster:")
	for _, node := range nodes.Items {
		klog.Infof("Node Name: %s", node.Name)
	}
}
