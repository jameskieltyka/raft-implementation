package discovery

import (
	"fmt"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Kubernetes struct{}

func (k *Kubernetes) GetNodes(name, namespace string) []string {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//Grab the raft deployment details
	deploy, err := clientset.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	//Get the pod ips from the deployment selector labels
	selectors := deploy.Spec.Selector.MatchLabels
	selectorString := joinSelectors(selectors)
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: selectorString})
	result := []string{}

	currentIP := os.Getenv("POD_IP")
	for _, pod := range pods.Items {
		if pod.Status.PodIP != currentIP {
			result = append(result, pod.Status.PodIP)
		}
	}

	return result
}

func joinSelectors(selectors map[string]string) string {
	var result []string
	for k, v := range selectors {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(result, ", ")
}
