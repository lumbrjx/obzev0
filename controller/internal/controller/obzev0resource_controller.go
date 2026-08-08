package controller

import (
	"context"
	"fmt"
	v1 "obzev0/controller/api/v1"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var setupLog = ctrl.Log.WithName("setup")

// Package-level singletons set once in SetupInformers and used by infx.go.
var (
	eventRecorder record.EventRecorder
	ctrlClient    client.Client
)

// PodConnection pairs a gRPC connection with the metadata needed for blast-radius filtering.
type PodConnection struct {
	Conn      *grpc.ClientConn
	PodName   string
	Namespace string
	Labels    map[string]string
	NodeName  string
}

func SetupInformers(mgr ctrl.Manager, recorder record.EventRecorder, c client.Client) {
	eventRecorder = recorder
	ctrlClient = c

	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "unable to create clientset")
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

	daemonSetInformer, err := mgr.GetCache().GetInformer(ctx, &appsv1.DaemonSet{})
	if err != nil {
		setupLog.Error(err, "unable to create DaemonSet informer")
		os.Exit(1)
	}
	setupLog.Info("DaemonSet informer created")

	podInformer, err := mgr.GetCache().GetInformer(ctx, &corev1.Pod{})
	if err != nil {
		setupLog.Error(err, "unable to create Pod informer")
		os.Exit(1)
	}
	setupLog.Info("Pod informer created")

	gRPCConnections := make(map[string]*PodConnection)

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
				cleanupConnections(ds, gRPCConnections)
			}
		},
	})

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
				deleteConnection(pod.Status.PodIP, gRPCConnections)
			}
		},
	})

	crInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { handleCREvent(obj, gRPCConnections) },
		UpdateFunc: func(oldObj, newObj interface{}) { handleCREvent(newObj, gRPCConnections) },
		DeleteFunc: func(obj interface{}) { handleCRDelete(obj, gRPCConnections) },
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
	connections map[string]*PodConnection,
) {
	for {
		pods, err := clientset.CoreV1().
			Pods(ds.Namespace).
			List(context.TODO(), metav1.ListOptions{
				LabelSelector: labels.SelectorFromSet(ds.Spec.Selector.MatchLabels).String(),
			})
		if err != nil {
			setupLog.Error(err, "Failed to list pods for DaemonSet", "name", ds.Name)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase == corev1.PodRunning {
				go retryConnectToPod(&pod, connections)
			}
		}

		time.Sleep(30 * time.Second)
	}
}

func retryConnectToPod(pod *corev1.Pod, connections map[string]*PodConnection) {
	ip := pod.Status.PodIP
	address := fmt.Sprintf("%s:50051", ip)

	for {
		if _, exists := connections[address]; exists {
			return
		}

		setupLog.Info("Attempting to connect to gRPC server", "address", address)
		conn, err := grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(LoggingInterceptor),
		)
		if err != nil {
			setupLog.Error(err, "Failed to connect to gRPC server", "address", address)
			time.Sleep(5 * time.Second)
			continue
		}

		connections[address] = &PodConnection{
			Conn:      conn,
			PodName:   pod.Name,
			Namespace: pod.Namespace,
			Labels:    pod.Labels,
			NodeName:  pod.Spec.NodeName,
		}
		setupLog.Info("Successfully connected to gRPC server", "address", address)
		return
	}
}

func cleanupConnections(ds *appsv1.DaemonSet, connections map[string]*PodConnection) {
	for address, pc := range connections {
		if pc.Namespace == ds.Namespace {
			pc.Conn.Close()
			delete(connections, address)
			setupLog.Info("Cleaned up connection for deleted DaemonSet", "address", address)
		}
	}
}

func deleteConnection(podIP string, connections map[string]*PodConnection) {
	address := fmt.Sprintf("%s:50051", podIP)
	if pc, exists := connections[address]; exists {
		pc.Conn.Close()
		delete(connections, address)
		setupLog.Info("Cleaned up connection for deleted Pod", "address", address)
	}
}

func handleCREvent(obj interface{}, connections map[string]*PodConnection) {
	cr, ok := obj.(*v1.Obzev0Resource)
	if !ok {
		setupLog.Error(nil, "Failed to cast object to Obzev0Resource")
		return
	}

	targets := filterConnections(connections, cr.Spec.BlastRadius)
	setupLog.Info("Dispatching chaos experiment",
		"cr", cr.Name,
		"totalPods", len(connections),
		"targetedPods", len(targets),
	)

	go processCustomResource(cr, targets)
}

func handleCRDelete(obj interface{}, connections map[string]*PodConnection) {
	cr, ok := obj.(*v1.Obzev0Resource)
	if !ok {
		setupLog.Error(nil, "Failed to cast object to Obzev0Resource on delete")
		return
	}
	stopScheduler(cr.Namespace + "/" + cr.Name)
	setupLog.Info("Custom Resource deleted", "name", cr.Name, "namespace", cr.Namespace)
}
