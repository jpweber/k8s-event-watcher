package main

import (
	"log"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "pods", v1.NamespaceAll,
		fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				log.Printf("Pod Added: %s, In Namespace: %s, Phase: %s\n",
					pod.ObjectMeta.Name,
					pod.ObjectMeta.Namespace,
					pod.Status.Phase,
				)

				// log container status information from added pod
				for _, cStatus := range pod.Status.ContainerStatuses {
					state := "Running"
					reason := ""

					if cStatus.State.Terminated != nil {
						state = "Terminated"
						reason = cStatus.State.Terminated.Reason
					}
					if cStatus.State.Waiting != nil {
						state = "Waiting"
						reason = cStatus.State.Waiting.Reason
					}

					log.Printf("Pod Added[Container Status] Name: %s, State: %s, Reason: %s, Ready %t, RestartCount %d",
						cStatus.Name,
						state,
						reason,
						cStatus.Ready,
						cStatus.RestartCount,
					)
				}

			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				log.Printf("Pod Deleted: %s, In Namespace: %s, Generation: %d, Phase: %s\n",
					pod.ObjectMeta.Name,
					pod.ObjectMeta.Namespace,
					pod.ObjectMeta.Generation,
					pod.Status.Phase,
					pod.Status.Message,
					pod.Status.Reason,
				)
				// log container status information from deleted pod
				for _, cStatus := range pod.Status.ContainerStatuses {
					state := "Running"
					reason := ""

					if cStatus.State.Terminated != nil {
						state = "Terminated"
						reason = cStatus.State.Terminated.Reason
					}
					if cStatus.State.Waiting != nil {
						state = "Waiting"
						reason = cStatus.State.Waiting.Reason
					}

					log.Printf("Pod Deleted[Container Status] Name: %s, State: %s, Reason: %s, Ready %t, RestartCount %d",
						cStatus.Name,
						state,
						reason,
						cStatus.Ready,
						cStatus.RestartCount,
					)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldPod := oldObj.(*v1.Pod)
				newPod := newObj.(*v1.Pod)
				log.Printf("Pod Modified old: %s, In Namespace: %s, Generation: %d, Phase: %s",
					oldPod.ObjectMeta.Name,
					oldPod.ObjectMeta.Namespace,
					oldPod.ObjectMeta.Generation,
					oldPod.Status.Phase,
				)
				// log container status information from old modified pod
				for _, cStatus := range oldPod.Status.ContainerStatuses {
					state := "Running"
					reason := ""

					if cStatus.State.Terminated != nil {
						state = "Terminated"
						reason = cStatus.State.Terminated.Reason
					}
					if cStatus.State.Waiting != nil {
						state = "Waiting"
						reason = cStatus.State.Waiting.Reason
					}

					log.Printf("Pod Modified old[Container Status] Name: %s, State: %s, Reason: %s, Ready %t, RestartCount %d",
						cStatus.Name,
						state,
						reason,
						cStatus.Ready,
						cStatus.RestartCount,
					)
				}

				log.Printf("Pod Modified New: %s, In Namespace: %s, Generation: %d, Phase: %s",
					newPod.ObjectMeta.Name,
					newPod.ObjectMeta.Namespace,
					newPod.ObjectMeta.Generation,
					newPod.Status.Phase,
				)
				// log container status information from old modified pod
				for _, cStatus := range newPod.Status.ContainerStatuses {
					state := "Running"
					reason := ""

					if cStatus.State.Terminated != nil {
						state = "Terminated"
						reason = cStatus.State.Terminated.Reason
					}
					if cStatus.State.Waiting != nil {
						state = "Waiting"
						reason = cStatus.State.Waiting.Reason
					}

					log.Printf("Pod Modified New[Container Status] Name: %s, State: %s, Reason: %s, Ready %t, RestartCount %d",
						cStatus.Name,
						state,
						reason,
						cStatus.Ready,
						cStatus.RestartCount,
					)
				}
			},
		},
	)
	stop := make(chan struct{})
	done := make(chan bool)
	go controller.Run(stop)
	log.Println("K8s Event Watcher Started")
	<-done
}
