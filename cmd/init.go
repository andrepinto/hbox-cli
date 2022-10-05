package cmd

import (
	"fmt"
	_init "github.com/andrepinto/hbox-cli/internal/init"
	"github.com/andrepinto/hbox-cli/internal/utils"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a local project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fsys := afero.NewOsFs()
		if err := _init.Run(fsys); err != nil {
			return err
		}

		fmt.Println("Finished " + utils.Aqua("hbox init") + ".")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
