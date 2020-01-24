package gcp

import (
	"context"
	"fmt"

	compute "google.golang.org/api/compute/v1"
)

// GetClient authenticates to GCP and fetches the Instance List
func GetClient(ctx context.Context) (*compute.Service, error) {
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error autenticating to the GCP service account")
	}
	return computeService, nil
}
