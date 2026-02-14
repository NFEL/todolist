package cmd

import (
	"context"
	"fmt"
	"graph-interview/internal/api"
	"graph-interview/internal/cfg"

	flag "github.com/spf13/pflag"
)

func ApiCmd(ctx context.Context, flagsStr []string) error {
	apiCmd := flag.NewFlagSet("api", flag.ContinueOnError)

	configPath := apiCmd.String("config", "cfg.toml", "Host to run the API server on")

	if err := apiCmd.Parse(flagsStr); err != nil {
		return fmt.Errorf("failed parsing cmd: %w", err)
	}

	if err := cfg.LoadConfig(*configPath, apiCmd); err != nil {
		return fmt.Errorf("failed loading config: %w", err)
	}

	return api.ServeREST(ctx)
}
