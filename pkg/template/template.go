package template

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"
)

type Compiler struct {
	funcs template.FuncMap
}

func (c *Compiler) Compile(source string) (*template.Template, error) {
	return template.New("").Funcs(c.funcs).Parse(source)
}

func (c *Compiler) Render(source string, data any) (string, error) {
	tmpl, err := c.Compile(source)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func NewCompiler(opts ...CompilerOption) (*Compiler, error) {
	c := &Compiler{
		funcs: make(template.FuncMap),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

type CompilerOption func(*Compiler) error

// WithRootFuncs registers a map of functions accessible in the template at the top level
//
// This should be used sparingly. Instead, consider using WithFuncNamespace to group functions together.
func WithRootFuncs(funcs template.FuncMap) CompilerOption {
	return func(c *Compiler) error {
		for k, v := range funcs {
			if _, dup := c.funcs[k]; dup {
				return fmt.Errorf("name already registered: %s", k)
			}
			c.funcs[k] = v
		}

		return nil
	}
}

// WithFuncNamespace registers a namespace of functions accessible in the template as namespace.FuncName
//
// This is acheived by providing a struct pointer with the functions as methods.
// The argument is provided as an any type, but reflection is used to check that it is a pointer to a
// struct with no public fields and only methods.
//
// This feels a bit hacky, but it's the best way I could find to achieve this.
func WithFuncNamespace(name string, structPtr any) CompilerOption {
	return func(c *Compiler) error {
		ptr := reflect.TypeOf(structPtr)
		if ptr.Kind() != reflect.Ptr {
			return fmt.Errorf("expected a pointer, got %s", ptr.Kind())
		}

		str := ptr.Elem()

		if str.Kind() != reflect.Struct {
			return fmt.Errorf("expected a struct, got %s", str.Kind())
		}

		for i := 0; i < str.NumField(); i++ {
			if str.Field(i).PkgPath == "" {
				return fmt.Errorf("struct should not have any public fields")
			}
		}

		if ptr.NumMethod() == 0 {
			return fmt.Errorf("struct should have at least one method")
		}

		if _, dup := c.funcs[name]; dup {
			return fmt.Errorf("name already registered: %s", name)
		}

		c.funcs[name] = func() any { return structPtr }
		return nil
	}
}
