package commands

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type CloneConfig struct {
	ApplicatioName              string
	SourceNamespace             string
	TargetNamespace             string
	ComponentSourceURLOverrides string
	OutputFile                  string
}

/*
modified variables to be inside a struct instead, and using references to instances of struct

	to access them in order to avoid using global variables
*/
type EnvironmentVariables struct {
	Host        string
	Namespace   string
	Bearertoken string
}

func NewOpenShiftClient() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return kubernetes.NewForConfig(config)
}
