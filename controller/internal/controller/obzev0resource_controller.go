package controller

import (
	"context"
	"fmt"
	v1 "obzev0/controller/api/v1"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
)

var setupLog = ctrl.Log.WithName("setup")

func SetupInformers(mgr ctrl.Manager) {
	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "unable to create clientset")
		os.Exit(1)
	}

	setupLog.Info("Setting up informers")
	ctx := context.Background()

	// Set up CR informer
	crInformer, err := mgr.GetCache().GetInformer(ctx, &v1.Obzev0Resource{})
	if err != nil {
		setupLog.Error(err, "unable to create CR informer")
		os.Exit(1)
	}
	setupLog.Info("CR informer created")

	// Set up DaemonSet informer
	daemonSetInformer, err := mgr.GetCache().GetInformer(ctx, &appsv1.DaemonSet{})
	if err != nil {
		setupLog.Error(err, "unable to create DaemonSet informer")
		os.Exit(1)
	}
	setupLog.Info("DaemonSet informer created")

	// Set up Pod informer
	podInformer, err := mgr.GetCache().GetInformer(ctx, &corev1.Pod{})
	if err != nil {
		setupLog.Error(err, "unable to create Pod informer")
		os.Exit(1)
	}
	setupLog.Info("Pod informer created")

	// Create a map to store gRPC connections
	gRPCConnections := make(map[string]*grpc.ClientConn)

	// Set up DaemonSet event handler
	daemonSetInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ds := obj.(*appsv1.DaemonSet)
			if isTargetDaemonSet(ds) {
				setupLog.Info("New target DaemonSet added", "name", ds.Name)
				go retryConnectToPods(clientset, ds, gRPCConnections)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newDS := newObj.(*appsv1.DaemonSet)
			if isTargetDaemonSet(newDS) {
				setupLog.Info("Target DaemonSet updated", "name", newDS.Name)
				go retryConnectToPods(clientset, newDS, gRPCConnections)
			}
		},
		DeleteFunc: func(obj interface{}) {
			ds := obj.(*appsv1.DaemonSet)
			if isTargetDaemonSet(ds) {
				setupLog.Info("Target DaemonSet deleted", "name", ds.Name)
				// Clean up connections for deleted DaemonSet
				cleanupConnections(ds, gRPCConnections)
			}
		},
	})

	// Set up Pod event handler
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			if isTargetPod(pod) {
				setupLog.Info("New target Pod added", "name", pod.Name)
				go retryConnectToPod(pod, gRPCConnections)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newPod := newObj.(*corev1.Pod)
			if isTargetPod(newPod) && newPod.Status.Phase == corev1.PodRunning {
				setupLog.Info("Target Pod became ready", "name", newPod.Name)
				go retryConnectToPod(newPod, gRPCConnections)
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			if isTargetPod(pod) {
				setupLog.Info("Target Pod deleted", "name", pod.Name)
				// Clean up connection for deleted Pod
				deleteConnection(pod.Status.PodIP, gRPCConnections)
			}
		},
	})

	// Set up CR event handler
	crInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { handleCREvent(obj, gRPCConnections) },
		UpdateFunc: func(oldObj, newObj interface{}) { handleCREvent(newObj, gRPCConnections) },
		DeleteFunc: func(obj interface{}) { handleCREvent(obj, gRPCConnections) },
	})
}

func isTargetDaemonSet(ds *appsv1.DaemonSet) bool {
	return ds.Labels["app"] == "grpc-server"
}

func isTargetPod(pod *corev1.Pod) bool {
	return pod.Labels["app"] == "grpc-server"
}

func retryConnectToPods(
	clientset *kubernetes.Clientset,
	ds *appsv1.DaemonSet,
	connections map[string]*grpc.ClientConn,
) {
	for {
		pods, err := clientset.CoreV1().
			Pods(ds.Namespace).
			List(context.TODO(), metav1.ListOptions{
				LabelSelector: labels.SelectorFromSet(ds.Spec.Selector.MatchLabels).
					String(),
			})
		if err != nil {
			setupLog.Error(
				err,
				"Failed to list pods for DaemonSet",
				"name",
				ds.Name,
			)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase == corev1.PodRunning {
				go retryConnectToPod(&pod, connections)
			}
		}

		time.Sleep(30 * time.Second) // Wait before checking again
	}
}

func retryConnectToPod(pod *corev1.Pod, connections map[string]*grpc.ClientConn) {
	ip := pod.Status.PodIP
	port := "50051"
	address := fmt.Sprintf("%s:%s", ip, port)

	for {
		if _, exists := connections[address]; exists {
			return // Connection already exists
		}

		setupLog.Info("Attempting to connect to gRPC server", "address", address)
		conn, err := grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(LoggingInterceptor),
		)
		if err != nil {
			setupLog.Error(
				err,
				"Failed to connect to gRPC server",
				"address",
				address,
			)
			time.Sleep(5 * time.Second)
			continue
		}

		connections[address] = conn
		setupLog.Info("Successfully connected to gRPC server", "address", address)
		return
	}
}

func cleanupConnections(
	ds *appsv1.DaemonSet,
	connections map[string]*grpc.ClientConn,
) {
	for address, conn := range connections {
		// Check if this connection belongs to the deleted DaemonSet
		// This is a simplification; you might need a more robust way to associate connections with DaemonSets
		if strings.HasPrefix(address, ds.Name) {
			conn.Close()
			delete(connections, address)
			setupLog.Info(
				"Cleaned up connection for deleted DaemonSet",
				"address",
				address,
			)
		}
	}
}

func deleteConnection(podIP string, connections map[string]*grpc.ClientConn) {
	address := fmt.Sprintf("%s:50051", podIP)
	if conn, exists := connections[address]; exists {
		conn.Close()
		delete(connections, address)
		setupLog.Info("Cleaned up connection for deleted Pod", "address", address)
	}
}

func handleCREvent(obj interface{}, connections map[string]*grpc.ClientConn) {
	cr, ok := obj.(*v1.Obzev0Resource)
	if !ok {
		setupLog.Error(nil, "Failed to cast object to Obzev0Resource")
		return
	}

	for _, conn := range connections {
		CheckConnection(conn)
		processCustomResource(cr, conn)
	}
}
