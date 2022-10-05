package new

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/andrepinto/hbox-cli/internal/configs"
	"github.com/andrepinto/hbox-cli/internal/utils"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

var (
	//go:embed templates/index.ts
	index string
)

func Run(ctx context.Context, slug string, fsys afero.Fs) error {
	funcDir := filepath.Join(configs.FunctionsDir, slug)
	{
		if err := utils.ValidateFunctionSlug(slug); err != nil {
			return err
		}
		if _, err := fsys.Stat(funcDir); !errors.Is(err, os.ErrNotExist) {
			return errors.New("Function " + utils.Aqua(slug) + " already exists locally.")
		}
	}

	if err := utils.MkdirIfNotExistFS(fsys, funcDir); err != nil {
		return err
	}
	if err := afero.WriteFile(fsys, filepath.Join(funcDir, "index.ts"), []byte(index), 0644); err != nil {
		return err
	}

	fmt.Println("Created new Function at " + utils.Bold(funcDir))
	return nil
}
