// A Dagger module for Protobuf projects.
//
// Provides generate, lint, and format functions backed by buf,
// with an configurable container image for repeatable CI workflows.
package main

import (
	"context"
	"dagger/protobuf/internal/dagger"
)

const (
	mountPath = "/app"
	outDir    = "/out/"
)

type Protobuf struct {
	Container *dagger.Container
}

// New creates a Protobuf builder from the given source directory.
func New(
	// Source directory containing the buf.yaml and proto files
	// +defaultPath="./"
	source *dagger.Directory,
	// buf container image to use
	// +defaultAddress="docker.io/bufbuild/buf:1.66"
	container *dagger.Container,
) *Protobuf {
	return &Protobuf{
		Container: container.
			WithWorkdir(mountPath).
			WithExec([]string{"buf", "--version"}).
			WithMountedDirectory(mountPath, source),
	}
}

// Generate runs buf generate and returns the container with generated output under /out/.
func (p *Protobuf) Generate(
	// Additional arguments passed to `buf generate`
	// +optional
	args []string,
) *dagger.Container {
	return p.Container.
		WithExec(append([]string{"buf", "generate", "--output", outDir}, args...))
}

// Lint runs buf lint on the proto files.
// +check
func (p *Protobuf) Lint(
	ctx context.Context,
	// Additional arguments passed to `buf lint`
	// +optional
	args []string,
) error {
	_, err := p.Container.
		WithExec(append([]string{"buf", "lint"}, args...)).
		Sync(ctx)
	return err
}

// Format runs buf format and returns the formatted proto files.
func (p *Protobuf) Format(
	// Additional arguments passed to `buf format`
	// +optional
	args []string,
) *dagger.Directory {
	return p.Container.
		WithExec(append([]string{"buf", "format", "--output", outDir}, args...)).
		Directory(outDir)
}
