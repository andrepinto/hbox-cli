package init

import (
	"errors"
	"github.com/andrepinto/hbox-cli/internal/configs"
	"github.com/andrepinto/hbox-cli/internal/utils"
	"github.com/spf13/afero"
	"os"
)

var (
	errAlreadyInitialized = errors.New("Project already initialized. Remove " + utils.Bold(configs.ConfigPath) + " to reinitialize.")
)

func Run(fsys afero.Fs) error {

	if _, err := fsys.Stat(configs.ConfigPath); err == nil {
		return errAlreadyInitialized
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := configs.WriteConfig(fsys); err != nil {
		return err
	}

	return nil
}
