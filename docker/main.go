// A generated module for Container functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/docker/internal/dagger"
)

type Docker struct {
	// +private
	Source *dagger.Directory
	// +private
	BuildArgs []dagger.BuildArg
	// +private
	Secrets []*dagger.Secret
	// +private
	SSH *dagger.Socket
}

func New(
	// +defaultPath="."
	source *dagger.Directory,
) *Docker {
	return &Docker{
		Source: source,
	}
}

func (d *Docker) WithBuildArg(name string, value string) *Docker {
	d.BuildArgs = append(d.BuildArgs, dagger.BuildArg{
		Name:  name,
		Value: value,
	})
	return d
}

// func (d *Docker) WithSecret(secret *dagger.Secret) *Docker {
// 	d.Secrets = append(d.Secrets, secret)
// 	return d
// }

func (d *Docker) WithSSH(socket *dagger.Socket) *Docker {
	d.SSH = socket
	return d
}

func (d *Docker) Build(
	ctx context.Context,
	// +default="Dockerfile"
	file string,
	// +default=""
	target string,
	// +default="linux/amd64"
	platform dagger.Platform,

) *dagger.Container {
	return d.Source.DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: file,
		Target:     target,
		Platform:   platform,
		BuildArgs:  d.BuildArgs,
		Secrets:    d.Secrets,
		SSH:        d.SSH,
	})
}
