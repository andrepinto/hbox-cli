package serve

import (
	"context"
	"fmt"
	"github.com/andrepinto/hbox-cli/internal/configs"
	"github.com/andrepinto/hbox-cli/internal/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

func Run(ctx context.Context, slug string, envFilePath string, verifyJWT bool, fsys afero.Fs) error {

	{
		if err := configs.CheckConfigFile(fsys); err != nil {
			return err
		}

		if err := configs.LoadConfigFS(fsys); err != nil {
			return err
		}

		if err := utils.AssertDockerIsRunning(); err != nil {
			return err
		}

		if err := utils.ValidateFunctionSlug(slug); err != nil {
			return err
		}

	}

	{

		_ = utils.Docker.ContainerRemove(ctx, configs.RelayContainerName, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		_, err = utils.Docker.NetworkInspect(ctx, configs.NetId, types.NetworkInspectOptions{})
		if err != nil && !client.IsErrNotFound(err) {
			return err
		} else if client.IsErrNotFound(err) {
			_, err = utils.Docker.NetworkCreate(
				ctx,
				configs.NetId,
				types.NetworkCreate{
					CheckDuplicate: true,
					Labels: map[string]string{
						configs.DockerLabelProjectID: configs.Config.ProjectId,
						configs.DockerLabelID:        configs.Config.ProjectId,
					},
				},
			)
			if err != nil {
				return err
			}
		}

		err = utils.DockerPullImageIfNotCached(context.Background(), utils.GetRegistryImageUrl(configs.RelayImage))
		if err != nil {
			return err
		}

		if verifyJWT {
			configs.Env = append(configs.Env, "VERIFY_JWT=true")
		} else {
			configs.Env = append(configs.Env, "VERIFY_JWT=false")
		}

		if _, err := utils.DockerRun(
			ctx,
			configs.RelayContainerName,
			&container.Config{
				Image: utils.GetRegistryImageUrl(configs.RelayImage),
				Env:   configs.Env,
				Labels: map[string]string{
					configs.DockerLabelProjectID: configs.Config.ProjectId,
					configs.DockerLabelID:        configs.Config.ProjectId,
				},
			},
			&container.HostConfig{
				Binds:        []string{filepath.Join(cwd, configs.FunctionsDir) + ":" + configs.RelayFuncDir + ":ro,z"},
				NetworkMode:  container.NetworkMode(configs.NetId),
				PortBindings: configs.Ports,
			},
		); err != nil {
			return err
		}

		go func() {
			<-ctx.Done()
			if ctx.Err() != nil {
				if err := utils.Docker.ContainerRemove(context.Background(), configs.RelayContainerName, types.ContainerRemoveOptions{
					RemoveVolumes: true,
					Force:         true,
				}); err != nil {
					fmt.Fprintln(os.Stderr, "Failed to remove container:", configs.RelayContainerName, err)
				}
			}
		}()
	}

	localFuncDir := filepath.Join(configs.FunctionsDir, slug)
	dockerFuncPath := configs.RelayFuncDir + "/" + slug + "/index.ts"

	_, err := utils.DockerExecOnce(ctx, configs.RelayContainerName, nil, []string{
		"deno", "cache", dockerFuncPath,
	})

	if err != nil {
		return err
	}

	{
		fmt.Println("Serving " + utils.Bold(localFuncDir))
		exec, err := utils.Docker.ContainerExecCreate(
			ctx,
			configs.RelayContainerName,
			types.ExecConfig{
				Env: []string{},
				Cmd: []string{
					"deno", "run", "--no-check=remote", "--allow-all", "--watch", "--no-clear-screen", dockerFuncPath,
				},
				AttachStderr: true,
				AttachStdout: true,
			},
		)
		if err != nil {
			return err
		}

		resp, err := utils.Docker.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
		if err != nil {
			return err
		}

		if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, resp.Reader); err != nil {
			return err
		}
	}

	fmt.Println("Stopped serving " + utils.Bold(localFuncDir))
	return nil

}
