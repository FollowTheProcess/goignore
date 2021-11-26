package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
)

// writeToIgnoreFile writes 'data' to cwd/.gitignore
func (a *App) writeToIgnoreFile(cwd string, data []byte) error {
	ignorePath := filepath.Join(cwd, ".gitignore")

	if _, err := a.fs.Stat(ignorePath); errors.Is(err, fs.ErrNotExist) {
		// No .gitignore file currently, good to go
		err := a.fs.WriteFile(ignorePath, data, 0o755)
		if err != nil {
			return fmt.Errorf("could not write to %q: %w", ignorePath, err)
		}
		return nil
	}

	return fmt.Errorf("%q already exists", ignorePath)
}
