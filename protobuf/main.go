// A generated module for Protobuf functions
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
	"dagger/protobuf/internal/dagger"
)

const (
	MOUNT_PATH = "/app"
	OUT_DIR    = "/out/"
)

type Protobuf struct {
	Container *dagger.Container
}

func New(
	// +defaultPath="./"
	source *dagger.Directory,
	// +defaultAddress="docker.io/bufbuild/buf:1.66"
	container *dagger.Container,
) *Protobuf {
	return &Protobuf{
		Container: container.
			WithWorkdir(MOUNT_PATH).
			WithExec([]string{"buf", "--version"}).
			WithMountedDirectory(MOUNT_PATH, source),
	}
}

func (p *Protobuf) Generate(
	// +optional
	args []string,
) *dagger.Container {
	return p.Container.
		WithExec(append([]string{"buf", "generate", "--output", OUT_DIR}, args...))
}

func (p *Protobuf) Lint(
	// +optional
	args []string,
) *dagger.Container {
	return p.Container.
		WithExec(append([]string{"buf", "lint"}, args...))
}

// Format proto files using buf format and return the formatted directory
func (p *Protobuf) Format(
	// +optional
	args []string,
) *dagger.Directory {
	return p.Container.
		WithExec(append([]string{"buf", "format", "--output", OUT_DIR}, args...)).
		Directory(OUT_DIR)
}
