package main

import (
	"context"
	"dagger/golang/internal/dagger"
	"fmt"
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
	// Container registry
	// +default="docker.io"
	registry string,
	// Golang image tag to use
	// +default="1.25"
	go_version string,
	// golangci-lint version to use
	// +default="2.10.1"
	lint_version string,
	// govulncheck version to use
	// +default="1.1.4"
	govulncheck_version string,

) *Golang {
	ctr := dag.Container().
		From(fmt.Sprintf("%s/library/golang:%s", registry, go_version)).
		WithWorkdir(MOUNT_PATH).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod_cache")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("go-build_cache")).
		WithEnvVariable("GOCACHE", "/root/.cache/go-build").
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@v" + govulncheck_version}).
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v" + lint_version}).
		WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume("golangci-lint_cache")).
		WithFile(MOUNT_PATH+"/go.mod", source.File("./go.mod"))

	exists, err := source.Exists(context.Background(), "./go.sum")
	if err == nil && exists {
		ctr = ctr.WithFile(MOUNT_PATH+"/go.sum", source.File("./go.sum"))
	}

	return &Golang{
		Container: ctr.
			WithExec([]string{"go", "mod", "download"}).
			WithMountedDirectory(MOUNT_PATH, source),
	}
}

func (g *Golang) Build(ctx context.Context,
	// +optional
	args []string) *dagger.Container {
	command := append([]string{"go", "build", "-o", OUT_DIR, MOUNT_PATH}, args...)
	return g.Container.
		WithExec(command)
}

// Lint runs golangci-Lint on the source code
// +check
func (g *Golang) Lint(ctx context.Context,
) *dagger.Container {
	return g.Container.
		WithExec([]string{"golangci-lint", "run", "--timeout", "5m"})
}

// VulnCheck runs govulncheck on the source code
// +check
func (g *Golang) VulnCheck(ctx context.Context) *dagger.Container {

	return g.Container.
		WithExec([]string{"govulncheck", "./..."})
}
