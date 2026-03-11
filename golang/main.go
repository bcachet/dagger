// A Dagger module for Go projects.
//
// Provides build, lint (golangci-lint), and vulnerability scanning (govulncheck)
// functions with integrated caching for fast, repeatable CI workflows.
package main

import (
	"context"
	"dagger/golang/internal/dagger"
)

type Golang struct {
	Container *dagger.Container
}

const (
	mountPath = "/app"
	outDir    = "/out/"
)

func New(
	ctx context.Context,
	// Source directory of the Go project
	// +defaultPath="./"
	source *dagger.Directory,
	// Go toolchain container image
	// +defaultAddress="docker.io/library/golang:1.25"
	container *dagger.Container,
	// golangci-lint version to use
	// +default="2.10.1"
	golangci_lint string,
	// govulncheck version to use
	// +default="1.1.4"
	govulncheck string,
) (*Golang, error) {
	ctr := container.
		WithWorkdir(mountPath).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod_cache")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("go-build_cache")).
		WithEnvVariable("GOCACHE", "/root/.cache/go-build").
		WithFile(mountPath+"/go.mod", source.File("./go.mod"))

	exists, err := source.Exists(ctx, "./go.sum")
	if err != nil {
		return nil, err
	}
	if exists {
		ctr = ctr.WithFile(mountPath+"/go.sum", source.File("./go.sum"))
	}

	g := &Golang{
		Container: ctr.
			WithExec([]string{"go", "mod", "download"}).
			WithMountedDirectory(mountPath, source),
	}
	return g.WithGolangciLint(golangci_lint).WithGovulncheck(govulncheck), nil
}

// Build compiles the Go project and places the output binary in /out/.
func (g *Golang) Build(
	// Additional arguments passed to `go build` (e.g. `-tags netgo ./cmd/myapp`)
	// +optional
	args []string,
) *dagger.Directory {
	command := append([]string{"go", "build", "-o", outDir}, args...)
	return g.Container.
		WithExec(command).
		Directory(outDir)
}

// WithGolangciLint installs the specified version of golangci-lint in the container.
func (g *Golang) WithGolangciLint(version string) *Golang {
	return &Golang{
		Container: g.Container.
			WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v" + version}).
			WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume("golangci-lint_cache")),
	}
}

// Lint runs golangci-Lint on the source code
// +check
func (g *Golang) Lint(ctx context.Context) error {
	_, err := g.Container.
		WithExec([]string{"golangci-lint", "run", "--timeout", "5m"}).
		Sync(ctx)
	return err
}

// WithGovulncheck installs the specified version of govulncheck in the container.
func (g *Golang) WithGovulncheck(version string) *Golang {
	return &Golang{
		Container: g.Container.
			WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@v" + version}),
	}
}

// VulnCheck runs govulncheck on the source code
// +check
func (g *Golang) VulnCheck(ctx context.Context) error {
	_, err := g.Container.
		WithExec([]string{"govulncheck", "./..."}).
		Sync(ctx)
	return err

}

// WithEnv sets an environment variable in the container.
func (g *Golang) WithEnv(name, value string) *Golang {
	return &Golang{
		Container: g.Container.WithEnvVariable(name, value),
	}
}

// WithCGOEnabled enables CGO in the container.
func (g *Golang) WithCgoEnabled() *Golang {
	return g.WithEnv("CGO_ENABLED", "1")
}

// WithCGODisabled disables CGO in the container.
func (g *Golang) WithCgoDisabled() *Golang {
	return g.WithEnv("CGO_ENABLED", "0")
}
