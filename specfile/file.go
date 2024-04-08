package specfile

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type File struct {
	Stream   string        `yaml:"stream" validate:"required"`
	Messages []FileMessage `yaml:"messages" validate:"required,min=1"`
}

type FileMessage struct {
	Subject string `yaml:"subject" validate:"required"`
	JSON    string `yaml:"json" validate:"required,json"`
}

func Parse(data []byte) (*File, error) {
	var f File
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	if err := validate.Struct(f); err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}

	return &f, nil
}

var validate = validator.New(validator.WithRequiredStructEnabled())
