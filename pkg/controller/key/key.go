package key

import (
	"fmt"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/digitalocean/godo"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/digitalocean"
	v1 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EnsureDeletedFinalizer is the name of the finalizer for ensuring DO keys are deleted
// Used in routing (pkg/controller/routes) as well as EnsureDeleted
const EnsureDeletedFinalizer = "provisioners.hobbyfarm.io/digitalocean-key-deleted"

// EnsureDeleted is a finalizer func that ensures that the corresponding
// digitalocean ssh key has been deleted.
func EnsureDeleted(req router.Request, _ router.Response) error {
	key := req.Object.(*v1alpha1.Key)

	var doKey = godo.Key{}
	if err := json.Unmarshal(key.Status.Key.Raw, &doKey); err != nil || (doKey.ID == 0 && doKey.Fingerprint == "") {
		logrus.Warnf("Could not retrieve DO key from object JSON. This may result "+
			"in an orphan SSH key in DigitalOcean. Key name was %s", key.Name)
		return nil
	}

	if !v1alpha1.RetryDeleteKey.CanRetry(key) {
		logrus.Errorf("Exceeded max attempts to retry %s, dropping object", v1alpha1.RetryDeleteKey.Action)
		return nil
	}

	dClient, err := digitalocean.GetGodoClient(key.Spec.Machine, req)
	if err != nil {
		v1alpha1.RetryDeleteKey.Failuref(key, "error building digitalocean client: %s", err.Error())
		_ = req.Client.Status().Update(req.Ctx, key)
		return fmt.Errorf("building digitalocean client: %s", err.Error())
	}

	byIdResp, err := dClient.Keys.DeleteByID(req.Ctx, doKey.ID)
	if byIdResp.StatusCode != http.StatusNotFound && err != nil {
		v1alpha1.RetryDeleteKey.Failuref(key, "error deleting key in digitalocean (by id): %s", err.Error())
		_ = req.Client.Status().Update(req.Ctx, key)
		return fmt.Errorf("deleting key in digitalocean (by id): %s", err.Error())
	}

	if byIdResp.StatusCode == http.StatusNotFound {
		// tried by id, no dice. try again by fingerprint
		byFpResp, err := dClient.Keys.DeleteByFingerprint(req.Ctx, doKey.Fingerprint)
		if byFpResp.StatusCode != http.StatusNotFound && err != nil {
			v1alpha1.RetryDeleteKey.Failuref(key, "error deleting key in digitalocean (by fingerprint): %s", err.Error())
			_ = req.Client.Status().Update(req.Ctx, key)
			return fmt.Errorf("deleting key in digitalocean (by fingerprint): %s", err.Error())
		}
	}

	v1alpha1.RetryDeleteKey.Success(key, "key deleted")
	_ = req.Client.Status().Update(req.Ctx, key)

	return nil
}

// EnsureStatus ensures that the conditions slice on the Status struct
// has been initialized. Named EnsureStatus in case other required values
// need creation in the future (e.g. not just conditions)
func EnsureStatus(req router.Request, resp router.Response) error {
	key := req.Object.(*v1alpha1.Key)

	if len(key.Status.Conditions) == 0 {
		v1alpha1.ConditionKeyCreated.SetStatus(key, "unknown")
	}

	resp.Objects(key)

	return nil
}

// NotYetCreated is a filtering middleware to only pass to CreateKey
// those keys that have not yet been created in DO
// criterion is ConditionKeyCreated.Status == "unknown"
func NotYetCreated(next router.Handler) router.Handler {
	return router.HandlerFunc(func(req router.Request, resp router.Response) error {
		key := req.Object.(*v1alpha1.Key)

		if v1alpha1.ConditionKeyCreated.GetStatus(key) == "unknown" {
			return next.Handle(req, resp)
		}

		return nil
	})
}

func Created(next router.Handler) router.Handler {
	return router.HandlerFunc(func(req router.Request, resp router.Response) error {
		key := req.Object.(*v1alpha1.Key)

		if v1alpha1.ConditionKeyCreated.GetStatus(key) == "true" {
			return next.Handle(req, resp)
		}

		return nil
	})
}

// CreateKey creates a corresponding ssh-key in DigitalOcean for the given
// key. Updates the ConditionKeyCreated condition depending on results
// from DO call.
func CreateKey(req router.Request, _ router.Response) error {
	key := req.Object.(*v1alpha1.Key)

	dClient, err := digitalocean.GetGodoClient(key.Spec.Machine, req)
	if err != nil {
		v1alpha1.ConditionKeyExists.SetStatus(key, "false")
		v1alpha1.ConditionKeyExists.SetError(key, "error creating digitalocean client", err)
		return fmt.Errorf("getting key from digitalocean: %s", err.Error())
	}

	var kcr = godo.KeyCreateRequest{}
	if err := json.Unmarshal(key.Spec.Key.Raw, &kcr); err != nil {
		return fmt.Errorf("unmarshalling keycreaterequest: %s", err.Error())
	}
	createdKey, _, err := dClient.Keys.Create(req.Ctx, &kcr)
	if err != nil {
		v1alpha1.ConditionKeyCreated.SetStatus(key, "false")
		v1alpha1.ConditionKeyCreated.SetError(key, "digitalocean request failed", err)
	} else {
		jKey, err := json.Marshal(createdKey)
		if err != nil {
			return fmt.Errorf("marshalling created key: %s", err.Error())
		}
		key.Status.Key.Raw = jKey
		v1alpha1.ConditionKeyCreated.SetStatus(key, "true")
	}

	err = req.Client.Status().Update(req.Ctx, key)
	if err != nil {
		return fmt.Errorf("updating key status: %s", err.Error())
	}

	return nil
}

func WriteVM(req router.Request, resp router.Response) error {
	key := req.Object.(*v1alpha1.Key)

	var vm = v1.VirtualMachine{}
	err := req.Client.Get(req.Ctx, client.ObjectKey{
		Name:      key.Spec.Machine,
		Namespace: key.Namespace,
	}, &vm)
	if err != nil {
		return fmt.Errorf("error retrieving vm %s: %s", key.Spec.Machine, err.Error())
	}

	vm.Spec.SecretName = key.Spec.Secret

	return req.Client.Update(req.Ctx, &vm)
}
