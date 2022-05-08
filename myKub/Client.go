package myKub

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

func KubApi() kubernetes.Interface {
	// Bootstrap k8s configuration from local Kubernetes config file
	kubeconf := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	//log.Println("Using kubeconf file: ", kubeconf)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconf)
	if err != nil {
		log.Fatal("ERROR : ", err)
	}

	//Create a rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("ERROR : ", err)
	}
	return clientset
}
