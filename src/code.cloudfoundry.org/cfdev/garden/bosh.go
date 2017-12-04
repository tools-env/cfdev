package garden

import (
	"fmt"

	"code.cloudfoundry.org/garden"
)

func DeployBosh(client garden.Client) error {
	containerSpec := garden.ContainerSpec{
		Handle:     "deploy-bosh",
		Privileged: true,
		Network:    "10.246.0.0/16",
		Image: garden.ImageRef{
			URI: "/var/vcap/director/cache/deploy-bosh.tar",
		},
		BindMounts: []garden.BindMount{
			{
				SrcPath: "/var/vcap",
				DstPath: "/var/vcap",
				Mode:    garden.BindMountModeRW,
			},
			{
				SrcPath: "/var/vcap/director/cache",
				DstPath: "/var/vcap/director/cache",
				Mode:    garden.BindMountModeRO,
			},
		},
	}

	container, err := client.Create(containerSpec)
	if err != nil {
		return err
	}

	process, err := container.Run(garden.ProcessSpec{
		ID:   "deploy-bosh",
		Path: "/usr/bin/deploy-bosh",
		User: "root",
	}, garden.ProcessIO{})

	if err != nil {
		return err
	}

	exitCode, err := process.Wait()
	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("process exited with status %v", exitCode)
	}

	client.Destroy("deploy-bosh")

	return nil
}
