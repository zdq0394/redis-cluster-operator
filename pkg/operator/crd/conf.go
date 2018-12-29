package crd

import (
	"fmt"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// Conf is the configuration required to create a CRD
type Conf struct {
	Kind       string
	NamePlural string
	Group      string
	Version    string
	Scope      apiextensionsv1beta1.ResourceScope
	Categories []string
}

func (c *Conf) getName() string {
	return fmt.Sprintf("%s.%s", c.NamePlural, c.Group)
}
