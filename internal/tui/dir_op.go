package tui

import (
	"fmt"
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
	parentDirStat, _ := os.Stat(parentDir)
	parentDirMode := parentDirStat.Mode().String()
	parentDirModTime := parentDirStat.ModTime()
	parentDirSize := parentDirStat.Size()
	entries = append(entries, entry{
		name:     "..",
		isDir:    true,
		fullPath: parentDir,
		mode:     parentDirMode,
		size: parentDirSize,
		modTime: parentDirModTime,

		displayName: "../",
		metadataRaw: fmt.Sprintf("%-11s %7d %s", parentDirMode, parentDirSize, formatLsTime(parentDirModTime)),
	})

	for _, de := range dirEntries {
		deInfo, _ := de.Info()

		mode := deInfo.Mode().String()
		modTime := deInfo.ModTime()
		size := deInfo.Size()
		isDir := de.IsDir()
		name := de.Name()

		var displayName string
		if isDir {
			displayName = fmt.Sprintf("%s/", name)
		} else {
			displayName = fmt.Sprintf("%s", name)
		}

		entries = append(entries, entry{
			name:     name,
			isDir:    isDir,
			fullPath: filepath.Join(targetDir, de.Name()),
			mode:     mode,
			size:     size,
			modTime:  modTime,

			displayName: displayName,
			metadataRaw: fmt.Sprintf("%-11s %7d %s", mode, size, formatLsTime(modTime)),
		})
	}

	return entries, targetDir, err
}
