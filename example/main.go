package main

import (
	"os"

	mkitlogger "github.com/layer5io/meshkit/logger"
	mkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"k8s.io/client-go/kubernetes"
)

func main() {
	// nginx contains the deployment manifest for nginx.
	nginx := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2 # tells deployment to run 2 pods matching the template
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
`

	// Create an instance of the meshkit logger handler, providing standardized output format.
	log, err := mkitlogger.New("ExampleApp", mkitlogger.Options{Format: mkitlogger.JsonLogFormat, DebugLevel: false})
	if err != nil {
		os.Exit(1)
	}
	log.Info("Successfully instantiated meshkit logger")

	// Detect the kubeconfig on the local system.
	config, err := mkitkube.DetectKubeConfig()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Info(config.Host)

	// Create Kubernetes client set for the detected kubeconfig. 'kubernetes' is from the Kubernetes Go client.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Create an instance of the meshkit Kubernetes client ...
	client, err := mkitkube.New(clientset, *config)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// ... and use it to deploy nginx to the cluster.
	err2 := client.ApplyManifest([]byte(nginx), mkitkube.ApplyOptions{
		Namespace: "default",
		Update:    true,
		Delete:    false,
	})

	if err2 != nil {
		log.Error(err2)
		os.Exit(1)
	}
}
