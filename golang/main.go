package main

import (
	"context"
	"dagger/golang/internal/dagger"
)

type Golang struct {
	Container *dagger.Container
}

const (
	MOUNT_PATH = "/app"
	OUT_DIR    = "/out/"
)

func New(
	// +defaultPath="./"
	source *dagger.Directory,
	// Container
	// +defaultAddress="docker.io/library/golang:1.25"
	container *dagger.Container,
	// golangci-lint version to use
	// +default="2.10.1"
	golangci_lint string,
	// govulncheck version to use
	// +default="1.1.4"
	govulncheck string,

) *Golang {
	ctr := container.
		WithWorkdir(MOUNT_PATH).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod_cache")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("go-build_cache")).
		WithEnvVariable("GOCACHE", "/root/.cache/go-build").
		WithFile(MOUNT_PATH+"/go.mod", source.File("./go.mod"))

	exists, err := source.Exists(context.Background(), "./go.sum")
	if err == nil && exists {
		ctr = ctr.WithFile(MOUNT_PATH+"/go.sum", source.File("./go.sum"))
	}

	g := &Golang{
		Container: ctr.
			WithExec([]string{"go", "mod", "download"}).
			WithMountedDirectory(MOUNT_PATH, source),
	}
	return g.WithGolangciLint(golangci_lint).WithGovulncheck(govulncheck)
}

// Build go binary/library
func (g *Golang) Build(ctx context.Context,
	// +optional
	args []string) *dagger.Container {
	command := append([]string{"go", "build", "-o", OUT_DIR}, args...)
	return g.Container.
		WithExec(command)
}

func (g *Golang) WithGolangciLint(
	version string,
) *Golang {
	g.Container = g.Container.
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v" + version}).
		WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume("golangci-lint_cache"))
	return g
}

// Lint runs golangci-Lint on the source code
// +check
func (g *Golang) Lint(ctx context.Context,
) *dagger.Container {
	return g.Container.
		WithExec([]string{"golangci-lint", "run", "--timeout", "5m"})
}

func (g *Golang) WithGovulncheck(
	version string,
) *Golang {
	g.Container = g.Container.
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@v" + version})
	return g
}

// VulnCheck runs govulncheck on the source code
// +check
func (g *Golang) VulnCheck(ctx context.Context) *dagger.Container {

	return g.Container.
		WithExec([]string{"govulncheck", "./..."})
}
