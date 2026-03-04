// Tests for the Docker Dagger module.
//
// Tests cover Build, BuildWithBuildArg, and BuildWithSecret functionality
// using the docker/example project as the test subject.
package main

import (
	"context"
	"fmt"

	"dagger/tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type Tests struct {
	Source *dagger.Directory
}

func New(
	// Source directory containing the example Dockerfile
	// +defaultPath="../example"
	source *dagger.Directory,
) *Tests {
	return &Tests{
		Source: source,
	}
}

// All runs all tests in parallel.
func (t *Tests) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(t.Build)
	p.Go(t.BuildWithBuildArg)
	p.Go(t.BuildWithSecret)

	return p.Wait()
}

// Build tests that the Docker module can build the example Dockerfile.
func (t *Tests) Build(ctx context.Context) error {
	_, err := dag.Docker(dagger.DockerOpts{Source: t.Source}).
		Build(dagger.DockerBuildOpts{}).
		Sync(ctx)
	return err
}

// BuildWithBuildArg tests that build arguments are passed correctly to the Dockerfile.
func (t *Tests) BuildWithBuildArg(ctx context.Context) error {
	content, err := dag.Docker(dagger.DockerOpts{Source: t.Source}).
		WithBuildArg("foo", "custom-value").
		Build(dagger.DockerBuildOpts{Target: "builder"}).
		File("/tmp/foo").
		Contents(ctx)
	if err != nil {
		return err
	}
	if content != "custom-value\n" {
		return fmt.Errorf("expected /tmp/foo to contain %q, got %q", "custom-value\n", content)
	}
	return nil
}

// BuildWithSecret tests that secrets are mounted and accessible during the build.
func (t *Tests) BuildWithSecret(ctx context.Context) error {
	secret := dag.SetSecret("hidden", "test-secret")
	content, err := dag.Docker(dagger.DockerOpts{Source: t.Source}).
		WithSecret("hidden", secret).
		Build(dagger.DockerBuildOpts{Target: "builder"}).
		File("/tmp/revealed").
		Contents(ctx)
	if err != nil {
		return err
	}
	if content != "test-secret" {
		return fmt.Errorf("expected /tmp/revealed to contain test-secret")
	}
	return nil
}
