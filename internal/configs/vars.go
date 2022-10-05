package configs

import "github.com/docker/go-connections/nat"

var (
	RelayContainerName = "relay"
	Env                = []string{
		"DENO_ORIGIN=http://localhost:8000",
	}

	Ports = nat.PortMap{
		"8081/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "8081",
			},
		}}
	NetId = "hbox_network_"
)

const (
	ConfigPath   = "hbox/config.toml"
	RelayFuncDir = "/home/deno/functions"
	FunctionsDir = "hbox/functions"
	RelayImage   = "andrepinto/hbox-relay:0.1.0"
)

const (
	DockerLabelProjectID = "hbox.project"
	DockerLabelID        = "com.docker.project"
)
