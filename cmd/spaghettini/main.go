package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/kofuk/spaghettini/server"
	"github.com/spf13/cobra"
)

type Options struct {
	logLevel string
	address  string
	source   string
}

func newCommand() *cobra.Command {
	options := &Options{}

	command := &cobra.Command{
		Use: "spaghettini",
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.source == "" {
				return fmt.Errorf("eval option is required")
			}
			if !slices.Contains([]string{"debug", "info", "warn", "error"}, options.logLevel) {
				return fmt.Errorf("invalid log level: %s", options.logLevel)
			}

			serverOptions := server.ServerOptions{
				LogLevel: server.LogLevel(options.logLevel),
				Addr:     options.address,
				Source:   options.source,
			}

			server, err := server.NewServer(serverOptions)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}

			return server.Start(cmd.Context())
		},
	}

	flags := command.Flags()

	flags.StringVarP(&options.logLevel, "log-level", "v", "info", "set the log level (debug, info, warn, error)")
	flags.StringVarP(&options.address, "address", "a", "0.0.0.0:8080", "the address to listen on")
	flags.StringVarP(&options.source, "eval", "e", "", "evaluate the given code")

	return command
}

func main() {
	command := newCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
