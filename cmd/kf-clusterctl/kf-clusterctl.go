package main

import (
	"context"
	"os"

	"github.com/CiscoAI/kf-cluster-api/pkg/gcp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const defaultLevel = log.InfoLevel

// Flags for the kind command
type Flags struct {
	LogLevel string
	File     string
}

// NewCommand creates the root cobra command
func NewCommand() *cobra.Command {
	flags := &Flags{}
	cmd := &cobra.Command{
		Use:   "kf-clusterctl",
		Short: "kf-clusterctl create KF Clusters",
		Long: `kf-clusterctl - a CLI tool to create KF Clusters

Usage:
	'kf-clusterctl create -f "kfcluster-gcp.yaml"' - creates a KF Cluster from spec.
	'kf-clusterctl delete' - deletes a KF Cluster from spec.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringVar(&flags.LogLevel, "loglevel", "info", "Default Log Level")
	cmd.Flags().StringVar(&flags.File, "file", "", "KF Cluster spec file")
	return cmd
}

func runE(flags *Flags, cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		log.Fatalf("kf-clusterctl needs an argument: `kf-clusterctl create` or `kf-clusterctl delete`")
	}
	// handle logLevel logic
	level := defaultLevel
	parsed, err := log.ParseLevel(flags.LogLevel)
	if err != nil {
		log.Warnf("Invalid log level '%s', defaulting to '%s'", flags.LogLevel, level)
	} else {
		level = parsed
	}
	log.SetLevel(level)

	ctx := context.Background()
	if args[0] == "create" {
		// Get compute engine client
		computeService, err := gcp.GetClient(ctx)
		if err != nil {
			return err
		}
		// Create new VM instance
		err = gcp.CreateInstance(ctx, "kf-github-action", "", "", computeService)
		if err != nil {
			return err
		}
	}
	if args[0] == "delete" {
		// Get compute engine client
		computeService, err := gcp.GetClient(ctx)
		if err != nil {
			return err
		}
		// Delete the VM instance
		err = gcp.DeleteInstance(ctx, "kf-github-action", "", "", computeService)
		if err != nil {
			return err
		}
	}
	return nil
}

// Run runs the `kf-clusterctl` root command
func Run() error {
	return NewCommand().Execute()
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
		ForceColors:     true,
	})
	if err := Run(); err != nil {
		os.Exit(1)
	}
}
