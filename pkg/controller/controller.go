package controller

import (
	"context"
	"github.com/acorn-io/baaah"
	"github.com/acorn-io/baaah/pkg/restconfig"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/ebauman/crder"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/crd"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/namespace"
	v1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
)

type Controller struct {
	Router     *router.Router
	Scheme     *runtime.Scheme
	restconfig *rest.Config
}

func New() (*Controller, error) {
	scheme := runtime.NewScheme()

	cfg, err := restconfig.New(scheme)

	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))

	r, err := baaah.NewRouter("hf-provisioner-digitalocean", scheme, &baaah.Options{
		RESTConfig: cfg,
		Namespace:  namespace.Resolve(),
	})

	if err != nil {
		return nil, err
	}

	routes(r)

	return &Controller{
		Router:     r,
		Scheme:     scheme,
		restconfig: cfg,
	}, nil
}

func (c *Controller) Start(ctx context.Context) error {
	crds := crd.Setup()

	if err := crder.InstallUpdateCRDs(c.restconfig, crds...); err != nil {
		return err
	}

	return c.Router.Start(ctx)
}
