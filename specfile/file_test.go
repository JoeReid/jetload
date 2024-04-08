package specfile

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		file      *File
		expectErr bool
	}{
		{
			name: "basic",
			file: &File{
				Stream: "USERS",
				Messages: []FileMessage{
					{
						Subject: "profile.4c31c006-b120-407d-bcca-0e83bd465ce7",
						JSON:    `{"user_id": "4c31c006-b120-407d-bcca-0e83bd465ce7", "username": "Alice"}`,
					},
				},
			},
			expectErr: false,
		},
		{
			name: "template",
			file: &File{
				Stream: "USERS",
				Messages: []FileMessage{
					{
						Subject: "profile.b80bacdc-c556-426d-b31c-2dea8f7eed48",
						JSON:    `{"user_id": "b80bacdc-c556-426d-b31c-2dea8f7eed48", "username": "CookerThrower"}`,
					},
					{
						Subject: "profile.bb5cf92e-8e96-4a5b-991e-60b234f453fe",
						JSON:    `{"user_id": "bb5cf92e-8e96-4a5b-991e-60b234f453fe", "username": "EagerPencil"}`,
					},
					{
						Subject: "profile.4b7e13fb-215f-4390-a2a6-50b47600bb96",
						JSON:    `{"user_id": "4b7e13fb-215f-4390-a2a6-50b47600bb96", "username": "LycheeCalm"}`,
					},
				},
			},
			expectErr: false,
		},
		{
			name:      "invalid",
			file:      nil,
			expectErr: true,
		},
		{
			name:      "unparsable",
			file:      nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			in, err := os.ReadFile("testdata/parse/" + tt.name)
			require.NoError(t, err)

			rendered, err := Render(in)
			require.NoError(t, err, "failed to render file before parsing")

			gotFile, err := Parse(rendered)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.file, gotFile)
		})
	}
}

func TestLoadAll(t *testing.T) {
	tests := []struct {
		name      string
		paths     []string
		expectErr bool
		numFiles  int
	}{
		{
			name:      "single yaml",
			paths:     []string{"loadall/basic"},
			expectErr: false,
			numFiles:  1,
		},
		{
			name:      "single template",
			paths:     []string{"loadall/template"},
			expectErr: false,
			numFiles:  1,
		},
		{
			name:      "multiple files",
			paths:     []string{"loadall/basic", "loadall/template"},
			expectErr: false,
			numFiles:  2,
		},
		{
			name:      "dir",
			paths:     []string{"loadall"},
			expectErr: false,
			numFiles:  2,
		},
		{
			name:      "no files",
			paths:     []string{},
			expectErr: true,
			numFiles:  0,
		},
		{
			name:      "missing file",
			paths:     []string{"loadall/missing"},
			expectErr: true,
			numFiles:  0,
		},
	}

	fsys := os.DirFS("./testdata")

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			files, err := LoadPaths(fsys, tt.paths...)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Len(t, files, tt.numFiles)
		})
	}
}
