package digitalocean

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/digitalocean/godo"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/config"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/namespace"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetGodoClient(vmName string, req router.Request) (*godo.Client, error) {
	// lookup the do token
	token := config.ResolveConfigItemName(vmName, req, "token")
	if token == "" {
		// check for token secret?
		tokenSecret := config.ResolveConfigItemName(vmName, req, "token-secret")
		if tokenSecret == "" {
			return nil, fmt.Errorf("unable to resolve token/token secret for digitalocean api")
		}

		secret := &v1.Secret{}
		err := req.Client.Get(req.Ctx, client.ObjectKey{
			Namespace: namespace.Resolve(),
			Name:      tokenSecret,
		}, secret)
		if err != nil {
			return nil, fmt.Errorf("error retrieving token secret: %s", err.Error())
		}

		token = string(secret.Data["token"])
	}

	return godo.NewFromToken(token), nil
}
