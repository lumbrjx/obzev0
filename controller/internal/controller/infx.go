package controller

import (
	"context"
	"log"
	v1 "obzev0/controller/api/v1"

	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

func processCustomResource(obz *v1.Obzev0Resource, conn *grpc.ClientConn) {
	name := obz.GetName()
	namespace := obz.GetNamespace()
	latencyConfig := obz.Spec.LatencyServiceConfig
	tcAConfig := obz.Spec.TcAnalyserServiceConfig

	klog.Infof("Custom Resource processed: %s/%s", namespace, name)
	klog.Infof("TCP Server Configuration: %+v", latencyConfig)
	klog.Infof("Tc Analyser Configuration: %+v", tcAConfig)

	svcConfig := GrpcServiceConfig{
		LatencyConfig: latencyConfig,
		TcAConfig:     tcAConfig,
	}

	err := callGrpcServices(conn, svcConfig)
	if err != nil {
		log.Printf("Error calling gRPC services: %v\n", err)
	}

	defer conn.Close()
}

func handleAdd(obj interface{}, conn *grpc.ClientConn) {
	obz, ok := obj.(*v1.Obzev0Resource)
	if !ok {
		klog.Errorf("Error converting object to Obzev0Resource: %v", obj)
		return
	}

	processCustomResource(obz, conn)
}

func handleUpdate(newObj interface{}, conn *grpc.ClientConn) {
	obz, ok := newObj.(*v1.Obzev0Resource)
	if !ok {
		klog.Errorf("Error converting object to Obzev0Resource: %v", newObj)
		return
	}

	processCustomResource(obz, conn)
}

func handleDelete(obj interface{}) {
	obz, ok := obj.(*v1.Obzev0Resource)
	if !ok {
		klog.Errorf("Error converting object to Obzev0Resource: %v", obj)
		return
	}

	name := obz.GetName()
	namespace := obz.GetNamespace()

	klog.Infof("Custom Resource deleted: %s/%s", namespace, name)
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
