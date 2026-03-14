package main

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/jonathan5p/lazys3/internal/s3"
	"github.com/jonathan5p/lazys3/internal/ui"
	"github.com/spf13/cobra"
)

type config struct {
	region   string
	endpoint string
	profile  string
}

func newRootCmd() *cobra.Command {
	cfg := &config{}

	cmd := &cobra.Command{
		Use:   "lazys3",
		Short: "A terminal UI for browsing AWS S3",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cfg)
		},
	}

	cmd.PersistentFlags().StringVar(&cfg.region, "region", "", "AWS region")
	cmd.PersistentFlags().StringVar(&cfg.endpoint, "endpoint", "", "Custom S3 endpoint URL")
	cmd.PersistentFlags().StringVar(&cfg.profile, "profile", "", "AWS profile name")

	return cmd
}

func run(cfg *config) error {
	ctx := context.Background()
	client, err := s3.NewAWSClient(ctx, cfg.region, cfg.endpoint, cfg.profile)
	if err != nil {
		return fmt.Errorf("init s3 client: %w", err)
	}
	app := ui.NewApp(client, 0, 0)
	_, err = tea.NewProgram(app).Run()
	return err
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
