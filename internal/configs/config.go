package configs

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/andrepinto/hbox-cli/internal/utils"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"text/template"
)

var Config config

type (
	config struct {
		ProjectId string `toml:"project_id"`
		Api       api
	}

	api struct {
		Port uint
	}
)

func CheckConfigFile(fsys afero.Fs) error {
	if _, err := fsys.Stat(ConfigPath); errors.Is(err, os.ErrNotExist) {
		return errors.New("Cannot find " + utils.Bold(ConfigPath) + " in the current directory. Have you set up the project with " + utils.Aqua("hbox init") + "?")
	} else if err != nil {
		return err
	}

	return nil
}

func LoadConfigFS(fsys afero.Fs) error {
	if _, err := toml.DecodeFS(afero.NewIOFS(fsys), ConfigPath, &Config); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Missing config: %w", err)
	} else if err != nil {
		return fmt.Errorf("Failed to read config: %w", err)
	}
	if Config.ProjectId == "" {
		return errors.New("Project_id is empty")
	} else {
		NetId = NetId + Config.ProjectId
	}

	return nil
}

func WriteConfig(fsys afero.Fs) error {
	// Using current directory name as project id
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir := filepath.Base(cwd)

	var initConfigBuf bytes.Buffer
	var tmpl *template.Template

	tmpl = initConfigTemplate

	if err := tmpl.Execute(
		&initConfigBuf,
		struct{ ProjectId string }{ProjectId: dir},
	); err != nil {
		return err
	}

	if err := utils.MkdirIfNotExistFS(fsys, filepath.Dir(ConfigPath)); err != nil {
		return err
	}

	if err := afero.WriteFile(fsys, ConfigPath, initConfigBuf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
