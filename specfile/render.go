package specfile

import (
	"fmt"

	"github.com/JoeReid/jetload/pkg/template"
	"github.com/brianvoe/gofakeit/v7"
)

// Render takes the bytes of a template file, compiles it and returns the rendered bytes.
//
// Render provides a suite of built-in functions to the template compiler, to enrich the
// templating experience somewhat.
func Render(file []byte) ([]byte, error) {
	comp, err := newCompiler()
	if err != nil {
		return nil, fmt.Errorf("failed to initialise template compiler: %w", err)
	}

	s, err := comp.Render(string(file), nil)
	return []byte(s), err
}

func newCompiler() (*template.Compiler, error) {
	return template.NewCompiler(
		template.WithRootFuncs(basicFuncs),
		template.WithFuncNamespace("time", &timeFuncs{}),
		template.WithFuncNamespace("faker", &fakerFuncs{gofakeit.New(1)}), // Seed of 1 for deterministic results
	)
}
