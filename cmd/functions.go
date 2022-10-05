package cmd

import (
	"github.com/andrepinto/hbox-cli/internal/functions/new"
	"github.com/andrepinto/hbox-cli/internal/functions/serve"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var (
	functionsCmd = &cobra.Command{
		Use:   "functions",
		Short: "Manage hbox Edge functions",
	}

	functionsServeCmd = &cobra.Command{
		Use:   "serve <Function name>",
		Short: "Serve a Function locally",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, _ := signal.NotifyContext(cmd.Context(), os.Interrupt)
			return serve.Run(ctx, args[0], "envFilePath", false, afero.NewOsFs())
		},
	}

	functionsNewCmd = &cobra.Command{
		Use:   "new <Function name>",
		Short: "Create a new Function locally",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, _ := signal.NotifyContext(cmd.Context(), os.Interrupt)
			return new.Run(ctx, args[0], afero.NewOsFs())
		},
	}
)

func init() {
	functionsCmd.AddCommand(functionsServeCmd)
	functionsCmd.AddCommand(functionsNewCmd)
	rootCmd.AddCommand(functionsCmd)
}
