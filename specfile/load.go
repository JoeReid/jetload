package specfile

import (
	"fmt"
	"io/fs"
)

// LoadPaths loads all files from the given paths. It handles both single files and directories.
func LoadPaths(fsys fs.FS, paths ...string) ([]File, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no paths provided")
	}

	files := make([]File, 0)
	for _, path := range paths {
		f, err := loadPath(fsys, path)
		if err != nil {
			return nil, err
		}

		files = append(files, f...)
	}

	return files, nil
}

func loadPath(fsys fs.FS, path string) ([]File, error) {
	info, err := fs.Stat(fsys, path)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)
	if info.IsDir() {
		dir, err := fs.ReadDir(fsys, path)
		if err != nil {
			return nil, err
		}

		for _, entry := range dir {
			if entry.IsDir() {
				// Let's not decend into directories for now.
				// Maybe we can add this behavior as a flag in the future.
				continue
			}

			paths = append(paths, path+"/"+entry.Name())
		}
	} else {
		paths = append(paths, path)
	}

	files := make([]File, 0)
	for _, p := range paths {
		bytes, err := fs.ReadFile(fsys, p)
		if err != nil {
			return nil, err
		}

		rendered, err := Render(bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to render file %s: %w", p, err)
		}

		file, err := Parse(rendered)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file %s: %w", p, err)
		}

		files = append(files, *file)
	}

	return files, nil
}
