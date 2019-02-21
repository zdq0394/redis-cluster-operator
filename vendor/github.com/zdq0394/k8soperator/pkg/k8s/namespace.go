package k8s

import (
	"fmt"

	"github.com/zdq0394/k8soperator/pkg/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Namespace the namespace service that knows how to interact with k8s to manage them
type Namespace interface {
	CreateNamespace(namespace string) error
	DeleteNamespace(namespace string) error
	NamespaceExists(namespace string) (bool, error)
}

// NamespaceService is the namespace service implement Namespace interface.
type NamespaceService struct {
	kubeClient kubernetes.Interface
	logger     log.Logger
}

// NewNamespaceService return a NamespaceService instance.
func NewNamespaceService(c kubernetes.Interface, l log.Logger) *NamespaceService {
	return &NamespaceService{
		kubeClient: c,
		logger:     l,
	}
}

//CreateNamespace create a namespace with name `namespace`
func (ns *NamespaceService) CreateNamespace(namespace string) error {
	kns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	if _, err := ns.kubeClient.CoreV1().Namespaces().Create(kns); err != nil {
		return err
	}

	return nil
}

// DeleteNamespace delete a namespace with name `namespace`
func (ns *NamespaceService) DeleteNamespace(namespace string) error {
	if err := ns.kubeClient.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

// NamespaceExists check wether the namespace with name `namespace` exists
func (ns *NamespaceService) NamespaceExists(namespace string) (bool, error) {
	if _, err := ns.kubeClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{}); err != nil {
		return false, fmt.Errorf("namespace %s not exists", namespace)
	}
	return true, nil
}
