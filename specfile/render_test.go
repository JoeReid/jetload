package specfile

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	os.Setenv("TEST_RENDER_SENTINEL", "some value")
	t.Cleanup(func() {
		os.Unsetenv("TEST_RENDER_SENTINEL")
	})

	var tests = []struct {
		name     string
		template string
		check    func(t *testing.T, output string)
	}{
		{
			name:     "seq",
			template: `{{ range seq 5 }}{{ . }}{{ end }}`,
			check:    equalString("01234"),
		},
		{
			name:     "time.Now",
			template: `{{ time.Now | time.Format "2006-01-02T15:04:05Z07:00" }}`,
			check:    timeStringWithin(time.Now(), 1*time.Second, time.RFC3339),
		},
		{
			name:     "faker.UUID",
			template: `{{ faker.UUID }}`,
			check:    equalString("b80bacdc-c556-426d-b31c-2dea8f7eed48"), // fixed seed
		},
		{
			name:     "env",
			template: `{{ env "TEST_RENDER_SENTINEL" }}`,
			check:    equalString("some value"),
		},
	}

	compiler, err := newCompiler()
	require.NoError(t, err)

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			_, err := compiler.Compile(tt.template)
			require.NoError(t, err)

			output, err := compiler.Render(tt.template, nil)
			require.NoError(t, err)
			tt.check(t, output)
		})
	}
}

func equalString(expected string) func(t *testing.T, actual string) {
	return func(t *testing.T, actual string) {
		t.Helper()

		assert.Equal(t, expected, actual)
	}
}

func timeStringWithin(expected time.Time, window time.Duration, format string) func(t *testing.T, actual string) {
	return func(t *testing.T, actual string) {
		t.Helper()

		parsed, err := time.Parse(format, actual)
		require.NoError(t, err)
		assert.WithinDuration(t, expected, parsed, window)
	}
}
