package tui

import (
	"os"
	"path/filepath"
)

func loadDirEntries(path string) ([]entry, string, error) {
	var entries []entry

	targetDir, err := filepath.Abs(path)
	if err != nil {
		return entries, path, err
	}

	dirEntries, err := os.ReadDir(targetDir)
	if err != nil {
		return entries, targetDir, err
	}

	parentDir := filepath.Dir(targetDir)
	entries = append(entries, entry{
		name:     "..",
		isDir:    true,
		fullPath: parentDir,
	})

	for _, de := range dirEntries {
		entries = append(entries, entry{
			name:     de.Name(),
			isDir:    de.IsDir(),
			fullPath: filepath.Join(targetDir, de.Name()),
		})
	}

	return entries, targetDir, err
}
