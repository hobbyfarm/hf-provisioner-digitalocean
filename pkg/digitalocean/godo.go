package digitalocean

import (
	"context"
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/digitalocean/godo"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/config"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/namespace"
	v12 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetGodoClientForEnvironment(ctx context.Context, env v12.Environment, kclient client.Client) (*godo.Client, error) {
	// lookup the token
	token, ok := env.Spec.EnvironmentSpecifics["token"]
	var err error
	if !ok || token == "" {
		// maybe token secret?
		tokenSecret, tokSecOk := env.Spec.EnvironmentSpecifics["token-secret"]
		if !tokSecOk || tokenSecret == "" {
			return nil, fmt.Errorf("unable to resolve token, tried 'token' and 'token-secret'")
		}

		token, err = resolveTokenFromSecret(ctx, tokenSecret, kclient)
		if err != nil {
			return nil, err
		}
	}

	// have a token, return godo client
	return godo.NewFromToken(token), nil
}

func GetGodoClient(vmName string, req router.Request) (*godo.Client, error) {
	// lookup the do token
	token := config.ResolveConfigItemName(vmName, req, "token")
	var err error
	if token == "" {
		// check for token secret?
		tokenSecret := config.ResolveConfigItemName(vmName, req, "token-secret")
		if tokenSecret == "" {
			return nil, fmt.Errorf("unable to resolve token/token secret for digitalocean api")
		}

		token, err = resolveTokenFromSecret(req.Ctx, tokenSecret, req.Client)
		if err != nil {
			return nil, err
		}
	}

	return godo.NewFromToken(token), nil
}

func resolveTokenFromSecret(ctx context.Context, secretName string, kclient client.Client) (string, error) {
	secret := &v1.Secret{}
	err := kclient.Get(ctx, client.ObjectKey{
		Namespace: namespace.Resolve(),
		Name:      secretName,
	}, secret)
	if err != nil {
		return "", fmt.Errorf("error retrieving token secret: %s", err.Error())
	}

	if tok, ok := secret.Data["token"]; !ok {
		return "", fmt.Errorf("invalid secret, key 'token' not found")
	} else {
		return string(tok), nil
	}
}
