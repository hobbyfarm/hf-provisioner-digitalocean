package v1alpha1

import (
	provisioning_hobbyfarm_io "github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const Version = "v1alpha1"

var SchemeGroupVersion = schema.GroupVersion{
	Group:   provisioning_hobbyfarm_io.Group,
	Version: Version,
}

func AddToScheme(scheme *runtime.Scheme) error {
	return AddToSchemeWithGV(scheme, SchemeGroupVersion)
}

func AddToSchemeWithGV(scheme *runtime.Scheme, schemeGroupVersion schema.GroupVersion) error {
	scheme.AddKnownTypes(schemeGroupVersion,
		&Droplet{},
		&DropletList{},
		&Key{},
		&KeyList{})

	scheme.AddKnownTypes(schemeGroupVersion, &metav1.Status{})

	if schemeGroupVersion == SchemeGroupVersion {
		metav1.AddToGroupVersion(scheme, schemeGroupVersion)
	}

	return nil
}
