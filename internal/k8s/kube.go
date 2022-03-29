package k8s

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	k8serrs "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

const DefaultKubeConfFmt = "%s/.kube/config"

// **************************************************************
// kube client to access k8s resource
// **************************************************************
type KClient struct {
	*kubernetes.Clientset
}

func NewKClient() (*KClient, error) {
	kubeConfPath := fmt.Sprintf(DefaultKubeConfFmt, os.Getenv("HOME"))
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfPath)
	if err != nil {
		return nil, err
	}

	kc, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &KClient{kc}, nil
}

func (c *KClient) GetKubeNodes(ctx context.Context) ([]string, error) {
	nodeList, err := c.CoreV1().Nodes().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var nodeSlice []string
	for _, node := range nodeList.Items {
		nodeSlice = append(nodeSlice, node.Name)
	}
	return nodeSlice, nil
}

func (c *KClient) CreateNamespace(ctx context.Context, namespace string) error {
	ns := &corev1.Namespace{}
	ns.Name = namespace
	_, err := c.CoreV1().Namespaces().Get(ctx, namespace, v1.GetOptions{})
	if err != nil {
		if k8serrs.IsNotFound(err) {
			_, err = c.CoreV1().Namespaces().Create(ctx, ns, v1.CreateOptions{})
		}
	}
	return err
}
