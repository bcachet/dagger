// Tests for the Protobuf Dagger module.
//
// Tests cover Lint and Generate functionality
// using the protobuf/example/v1 project as the test subject.
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
	// Source directory of the Protobuf project to run tests against
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

	p.Go(t.Lint)
	p.Go(t.Generate)

	return p.Wait()
}

// Lint tests that the Protobuf module can lint the example proto files.
func (t *Tests) Lint(ctx context.Context) error {
	_, err := dag.Protobuf(dagger.ProtobufOpts{Source: t.Source}).
		Lint(dagger.ProtobufLintOpts{}).
		Sync(ctx)
	return err
}

// Generate tests that the Protobuf module can generate code from the example proto files.
func (t *Tests) Generate(ctx context.Context) error {
	entries, err := dag.Protobuf(dagger.ProtobufOpts{Source: t.Source}).
		Generate(dagger.ProtobufGenerateOpts{}).
		Directory("/out/").
		Entries(ctx)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("expected generate output, got empty directory")
	}
	return nil
}
