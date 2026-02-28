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
	"fmt"
)

type Docker struct {
	// +private
	Source *dagger.Directory
	// +private
	BuildArgs []dagger.BuildArg
	// +private
	Secrets []Secret
	// +private
	SSH *dagger.Socket
}

type Secret struct {
	ID     string
	Secret *dagger.Secret
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

func (d *Docker) WithSecret(id string, secret *dagger.Secret) *Docker {
	d.Secrets = append(d.Secrets, Secret{
		ID:     id,
		Secret: secret,
	})
	return d
}

// +private
func (d *Docker) namedSecrets(ctx context.Context) ([]*dagger.Secret, error) {
	secrets := make([]*dagger.Secret, 0, len(d.Secrets))
	for _, s := range d.Secrets {
		plaintext, err := s.Secret.Plaintext(ctx)
		if err != nil {
			return nil, fmt.Errorf("Cannot retrieve %s secret: %w", s.ID, err)
		}
		secret := dag.SetSecret(s.ID, plaintext)
		secrets = append(secrets, secret)
	}
	return secrets, nil
}

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

) (*dagger.Container, error) {
	secrets, err := d.namedSecrets(ctx)
	if err != nil {
		return nil, err
	}
	return d.Source.DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: file,
		Target:     target,
		Platform:   platform,
		BuildArgs:  d.BuildArgs,
		Secrets:    secrets,
		SSH:        d.SSH,
	}), nil
}
