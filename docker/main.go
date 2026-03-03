// A Dagger module for building Docker images.
//
// Provides a builder that wraps Dagger's DockerBuild with support for
// build arguments, build secrets, and SSH forwarding.
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

// New creates a Docker builder from the given source directory.
func New(
	// Source directory containing the Dockerfile
	// +defaultPath="."
	source *dagger.Directory,
) *Docker {
	return &Docker{
		Source: source,
	}
}

// WithBuildArg adds a build argument passed to docker build.
func (d *Docker) WithBuildArg(name string, value string) *Docker {
	d.BuildArgs = append(d.BuildArgs, dagger.BuildArg{
		Name:  name,
		Value: value,
	})
	return d
}

// WithSecret mounts a secret for use during the build (e.g. RUN --mount=type=secret).
func (d *Docker) WithSecret(id string, secret *dagger.Secret) *Docker {
	d.Secrets = append(d.Secrets, Secret{
		ID:     id,
		Secret: secret,
	})
	return d
}

// namedSecrets resolves each secret's plaintext and re-registers it under its ID.
func (d *Docker) namedSecrets(ctx context.Context) ([]*dagger.Secret, error) {
	secrets := make([]*dagger.Secret, 0, len(d.Secrets))
	for _, s := range d.Secrets {
		plaintext, err := s.Secret.Plaintext(ctx)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve %s secret: %w", s.ID, err)
		}
		secret := dag.SetSecret(s.ID, plaintext)
		secrets = append(secrets, secret)
	}
	return secrets, nil
}

// WithSSH forwards the given SSH socket into the build for use with RUN --mount=type=ssh.
func (d *Docker) WithSSH(socket *dagger.Socket) *Docker {
	d.SSH = socket
	return d
}

// Build runs docker build and returns the resulting container image.
func (d *Docker) Build(
	ctx context.Context,
	// Dockerfile path relative to the source directory
	// +default="Dockerfile"
	file string,
	// Build stage target to stop at
	// +default=""
	target string,
	// Target platform for the image
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
