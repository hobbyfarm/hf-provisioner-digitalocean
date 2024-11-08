package cleanup

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/apis/provisioning.hobbyfarm.io/v1alpha1"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/digitalocean"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/log"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/providerregistration"
	"github.com/ebauman/hf-provisioner-digitalocean/pkg/tags"
	v12 "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	"github.com/hobbyfarm/hf-provisioner-shared/instanceid"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func RunCleanup(kclient client.Client, cleanupPeriodSeconds int, stopCh chan struct{}, errCh chan error) {
	var cleanupPeriod = time.Duration(cleanupPeriodSeconds * 1000000000)

	timer := time.NewTimer(cleanupPeriod)
	ctx, cancel := context.WithCancel(context.Background())

	for {
		select {
		case <-stopCh:
			cancel()
			return
		case <-timer.C:
			if err := executeCleanup(ctx, kclient); err != nil {
				errCh <- err
			}
			timer = time.NewTimer(cleanupPeriod)
		}
	}
}

func executeCleanup(ctx context.Context, kclient client.Client) error {
	log.Infof("starting cleanup at %s", time.Now().Format(time.RFC3339))
	// list all environments
	envs := &v12.EnvironmentList{}
	err := kclient.List(ctx, envs)
	if err != nil {
		return err
	}

	instanceId, err := instanceid.GetOrCreateInstanceId(ctx, kclient)
	if err != nil {
		// unable to resolve instance id
		return err
	}
	log.Debugf("got instance id: %s", instanceId)

	// which ones are DO envs?
	for _, e := range envs.Items {
		log.Debugf("cleaning up environment %s", e.Name)
		if e.Spec.Provider != providerregistration.ProviderName() {
			continue // don't work on things that aren't ours!
		}

		// get a godo client
		log.Debugf("setting up godo client")
		godoClient, err := digitalocean.GetGodoClientForEnvironment(ctx, e, kclient)
		if err != nil {
			return err
		}

		// fetch all droplets from digitalocean
		log.Debugf("fetching droplets from digitalocean")
		droplets, _, err := godoClient.Droplets.List(ctx, &godo.ListOptions{})
		if err != nil {
			return err
		}

		log.Debugf("got %d droplets", len(droplets))

		// fetch all droplets from k8s
		log.Debugf("fetching droplets from hobbyfarm")
		k8sDroplets := &v1alpha1.DropletList{}
		if err := kclient.List(ctx, k8sDroplets); err != nil {
			return err
		}

		log.Debugf("got %d HF droplets", len(k8sDroplets.Items))

		// loop the droplets
		for _, d := range droplets {
			log.Debugf("checking droplet %s", d.Name)
			// does this droplet belong to us?
			if !isThisOurDroplet(instanceId, d) {
				log.Debugf("skipping droplet %s", d.Name)
				continue
			}

			// if it is our droplet, do we have a matching droplet resource in k8s?
			if !shouldThisDropletExist(d, k8sDroplets) {
				log.Debugf("droplet %s should not exist, deleting", d.Name)
				// need to remove droplet
				if _, err := godoClient.Droplets.Delete(ctx, d.ID); err != nil {
					return err
				}
			} else {
				log.Debugf("droplet %s matched with HF droplet", d.Name)
			}
		}
	}

	log.Infof("completed cleanup at %s", time.Now().Format(time.RFC3339))
	return nil
}

func isThisOurDroplet(instanceId string, droplet godo.Droplet) bool {
	if tags.GetTagValue(droplet.Tags, tags.DOInstanceIdPrefix) == instanceId {
		return true
	}

	return false
}

func shouldThisDropletExist(droplet godo.Droplet, k8sDroplets *v1alpha1.DropletList) bool {
	dropletNameFromTag := tags.GetTagValue(droplet.Tags, tags.DODropletNamePrefix)

	for _, kd := range k8sDroplets.Items {
		if kd.Name == dropletNameFromTag {
			return true
		}
	}

	return false
}
