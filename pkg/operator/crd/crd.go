package crd

import (
	"fmt"
	"time"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionscli "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeversion "k8s.io/kubernetes/pkg/util/version"
)

const (
	checkCRDInterval = 2 * time.Second
	crdReadyTimeout  = 3 * time.Minute
)

var (
	clusterMinVersion = kubeversion.MustParseGeneric("v1.7.0")
	defCategories     = []string{"all"}
)

// Client is the CRD client that knows how to interact with k8s to manage them.
type Client interface {
	// EnsurePresent will ensure the the CRD is present, this also means that
	// apart from creating the CRD if is not present it will wait until is
	// ready, this is a blocking operation and will return an error if timesout
	// waiting.
	EnsurePresent(conf Conf) error
	// WaitToBePresent will wait until the CRD is present, it will check if
	// is present at regular intervals until it timesout, in case of timeout
	// will return an error.
	WaitToBePresent(name string, timeout time.Duration) error
	// Delete will delete the CRD.
	Delete(name string) error
}

type client struct {
	aeClient apiextensionscli.Interface
}

// NewClient returns a new CRD client.
func NewClient(aeClient apiextensionscli.Interface) Client {
	return &client{
		aeClient: aeClient,
	}
}

func (c *client) EnsurePresent(conf Conf) error {
	if err := c.validClusterForCRDs(); err != nil {
		return err
	}

	crd := &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: conf.getName(),
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   conf.Group,
			Version: conf.Version,
			Scope:   conf.Scope,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural:     conf.NamePlural,
				Kind:       conf.Kind,
				Categories: c.addDefaultCaregories(conf.Categories),
			},
		},
	}
	_, err := c.aeClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return fmt.Errorf("error creating crd %s:%s", conf.getName(), err)
		}
		return nil
	}
	if err := c.WaitToBePresent(conf.getName(), crdReadyTimeout); err != nil {
		return err
	}
	return nil
}

func (c *client) WaitToBePresent(name string, timeout time.Duration) error {
	if err := c.validClusterForCRDs(); err != nil {
		return err
	}
	timeOut := time.After(timeout)
	t := time.NewTicker(checkCRDInterval)
	for {
		select {
		case <-t.C:
			_, err := c.aeClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(name, metav1.GetOptions{})
			if err == nil {
				return nil
			}
		case <-timeOut:
			return fmt.Errorf("Timeout waiting for CRD to be created")
		}
	}
}

func (c *client) Delete(name string) error {
	if err := c.validClusterForCRDs(); err != nil {
		return err
	}
	return c.aeClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(name, &metav1.DeleteOptions{})
}

func (c *client) validClusterForCRDs() error {
	v, err := c.aeClient.Discovery().ServerVersion()
	if err != nil {
		return err
	}
	parsedV, err := kubeversion.ParseGeneric(v.GitVersion)
	if err != nil {
		return err
	}
	if parsedV.LessThan(clusterMinVersion) {
		return fmt.Errorf("Not a valid cluster version for CRDs (required >=1.7)")
	}
	return nil
}

func (c *client) addDefaultCaregories(categories []string) []string {
	currentCats := make(map[string]bool)
	for _, ca := range categories {
		currentCats[ca] = true
	}

	// Add default categories if required.
	for _, ca := range defCategories {
		if _, ok := currentCats[ca]; !ok {
			categories = append(categories, ca)
		}
	}
	return categories
}
