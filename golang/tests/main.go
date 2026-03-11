// Tests for the Golang Dagger module.
//
// Tests cover Build, Lint, and VulnCheck functionality
// using the golang/example project as the test subject.
package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/tests/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type Tests struct {
	Source *dagger.Directory
}

func New(
	// Source directory of the Go project to run tests against
	// +defaultPath="../example"
	source *dagger.Directory,
) *Tests {
	return &Tests{
		Source: source,
	}
}

// All runs all tests in sequence.
func (t *Tests) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(t.Build)
	p.Go(t.Lint)
	p.Go(t.VulnCheck)
	p.Go(t.WithCgoDisabled)
	p.Go(t.WithCgoEnabled)

	return p.Wait()
}

// Build tests that the Golang module can build the example Go project.
func (t *Tests) Build(ctx context.Context) error {
	entries, err := dag.Golang(dagger.GolangOpts{Source: t.Source}).
		Build(dagger.GolangBuildOpts{}).
		Entries(ctx)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("expected build output, got empty directory")
	}
	return nil
}

// Lint tests that the Golang module can lint the example Go project.
func (t *Tests) Lint(ctx context.Context) error {
	return dag.Golang(dagger.GolangOpts{Source: t.Source}).
		Lint(ctx)
}

// VulnCheck tests that the Golang module can run vulnerability checks on the example Go project.
func (t *Tests) VulnCheck(ctx context.Context) error {
	return dag.Golang(dagger.GolangOpts{Source: t.Source}).
		VulnCheck(ctx)
}

// lddOutput returns the output of ldd on the given binary file using the golang container.
func lddOutput(ctx context.Context, binary *dagger.File) (string, error) {
	return dag.Container().
		From("golang:1.25").
		WithFile("/app", binary).
		WithExec([]string{"sh", "-c", "ldd /app 2>&1 || true"}).
		Stdout(ctx)
}

// WithCgoDisabled tests that a CGO-disabled binary is statically linked and does not require libc.
func (t *Tests) WithCgoDisabled(ctx context.Context) error {
	binary := dag.Golang(dagger.GolangOpts{Source: t.Source}).
		WithCgoDisabled().
		Build(dagger.GolangBuildOpts{}).
		File("/example")

	out, err := lddOutput(ctx, binary)
	if err != nil {
		return err
	}
	if !strings.Contains(out, "not a dynamic executable") {
		return fmt.Errorf("binary is dynamically linked, expected no libc dependency: %s", out)
	}
	return nil
}

// WithCgoEnabled tests that a CGO-enabled binary is dynamically linked against libc.
func (t *Tests) WithCgoEnabled(ctx context.Context) error {
	binary := dag.Golang(dagger.GolangOpts{Source: t.Source}).
		WithCgoEnabled().
		Build(dagger.GolangBuildOpts{}).
		File("/example")

	out, err := lddOutput(ctx, binary)
	if err != nil {
		return err
	}
	if !strings.Contains(out, "libc") {
		return fmt.Errorf("binary is statically linked, expected libc dependency: %s", out)
	}
	return nil
}
