package gcp

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cenkalti/backoff"
	log "github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v1"
)

// standard instance name
const instanceName = "kf-github-action"

func getBackoff(maxTimeout time.Duration) *backoff.ExponentialBackOff {
	backOff := backoff.NewExponentialBackOff()
	backOff.MaxElapsedTime = maxTimeout
	return backOff
}

// ListInstances takes in the GCP project and zone; returns the List of all instances there
func ListInstances(ctx context.Context, project string, zone string, computeService *compute.Service) ([]string, error) {
	if project == "" {
		project = os.Getenv("PROJECT")
	}
	if zone == "" {
		zone = os.Getenv("ZONE")
	}
	var instanceNames []string
	// List instances
	resp, err := computeService.Instances.List(project, zone).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Error authenticating to GCE Instances Service")
	}
	if resp != nil {
		for _, instance := range resp.Items {
			fmt.Printf("Instance name: %s\n", instance.Name)
			instanceNames = append(instanceNames, instance.Name)
		}
	}
	return instanceNames, nil
}

// CreateInstance takes in the project and zone to create a standard VM if it doesn't exist and returns an error
func CreateInstance(ctx context.Context, instanceName string, project string, zone string, computeService *compute.Service) error {
	if project == "" {
		project = os.Getenv("PROJECT")
	}
	if zone == "" {
		zone = os.Getenv("ZONE")
	}
	machineType := "zones/" + zone + "/machineTypes/n2-standard-8"
	diskParams := &compute.AttachedDiskInitializeParams{
		DiskSizeGb:  20,
		SourceImage: "projects/cpsg-ai-kubeflow/global/images/github-action-image-from-snapshot",
	}
	persistentDisk := &compute.AttachedDisk{
		DeviceName:       "persistent-" + instanceName,
		Boot:             true,
		InitializeParams: diskParams,
		AutoDelete:       true,
		Type:             "PERSISTENT",
	}
	networkAccessConfig := &compute.AccessConfig{
		Name:        "External NAT",
		Type:        "ONE_TO_ONE_NAT",
		NetworkTier: "PREMIUM",
	}
	networkInterface := &compute.NetworkInterface{
		Network:       "global/networks/default",
		AccessConfigs: []*compute.AccessConfig{networkAccessConfig},
	}
	serviceAccount := &compute.ServiceAccount{
		Email: "default",
		Scopes: []string{
			"https://www.googleapis.com/auth/devstorage.read_write",
			"https://www.googleapis.com/auth/logging.write",
		},
	}
	var instanceLabel = "github-action"
	metadataLabel := &compute.MetadataItems{
		Key:   "ci-instance",
		Value: &instanceLabel,
	}
	startupScript := ""
	metadataStartupScript := &compute.MetadataItems{
		Key:   "startup-script",
		Value: &startupScript,
	}
	// Insert Op to create a new VM if one doesn't exist
	_, err := computeService.Instances.Get(project, zone, instanceName).Context(ctx).Do()
	if err != nil {
		log.Infof("Creating VM: %v", instanceName)
		instance := &compute.Instance{
			Name:              instanceName,
			MachineType:       machineType,
			Disks:             []*compute.AttachedDisk{persistentDisk},
			NetworkInterfaces: []*compute.NetworkInterface{networkInterface},
			ServiceAccounts:   []*compute.ServiceAccount{serviceAccount},
			Metadata: &compute.Metadata{
				Items: []*compute.MetadataItems{metadataLabel, metadataStartupScript}},
		}
		_, err := computeService.Instances.Insert(project, zone, instance).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("Error creating new instance: %v", err)
		}
		waitForCreateOp := func() error {
			getOp, err := computeService.Instances.Get(project, zone, instanceName).Context(ctx).Do()
			if getOp.Status == "RUNNING" || err != nil {
				log.Infof("VM Instance creation status: %v", getOp.Status)
				return nil
			}
			log.Infof("VM Instance creation pending, VM status: %v", getOp.Status)
			return fmt.Errorf("VM Instance creation pending")
		}
		createBackoff := getBackoff(5 * time.Minute)
		err = backoff.Retry(waitForCreateOp, createBackoff)
		if err != nil {
			return err
		}
		log.Infof("VM Instance Creation Succeeded")
	}
	return nil
}

// DeleteInstance - Used to delete an Instance
func DeleteInstance(ctx context.Context, instanceName string, project string, zone string, computeService *compute.Service) error {
	if project == "" {
		project = os.Getenv("PROJECT")
	}
	if zone == "" {
		zone = os.Getenv("ZONE")
	}
	_, err := computeService.Instances.Delete(project, zone, instanceName).Context(ctx).Do()
	if err != nil {
		return err
	}
	waitForDeleteOp := func() error {
		getOp, err := computeService.Instances.Get(project, zone, instanceName).Context(ctx).Do()
		if getOp != nil {
			if getOp.Status == "TERMINATED" || err != nil {
				log.Infof("VM Instance deletion Status: %v", getOp.Status)
				return nil
			}
		} else {
			return nil
		}
		log.Infof("VM Instance deletion pending, VM status: %v", getOp.Status)
		return fmt.Errorf("VM Instance deletion pending")
	}
	deleteBackoff := getBackoff(5 * time.Minute)
	err = backoff.Retry(waitForDeleteOp, deleteBackoff)
	if err != nil {
		return err
	}
	log.Infof("VM Instance Deletion Succeeded")
	return nil
}
