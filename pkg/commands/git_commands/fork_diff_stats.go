package git_commands

import (
	"bytes"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/spf13/afero"
)

// Fork: totals for the Changes panel, covering staged, unstaged and untracked
// changes against HEAD.
type DiffStats struct {
	LinesAdded   int
	LinesDeleted int
	FilesAdded   int
	FilesChanged int
	FilesRemoved int
}

func (self *FileLoader) GetDiffStats(files []*models.File) DiffStats {
	stats := DiffStats{}

	fileDiffs, err := self.getFileDiffs()
	if err != nil {
		self.Log.Error(err)
	}
	for _, diff := range fileDiffs {
		stats.LinesAdded += diff.LinesAdded
		stats.LinesDeleted += diff.LinesDeleted
	}

	for _, file := range files {
		// git diff --numstat ignores untracked files, so we count their lines ourselves
		if file.ShortStatus == "??" && !strings.HasSuffix(file.Path, "/") {
			stats.LinesAdded += self.countUntrackedFileLines(file.Path)
		}

		if file.Added {
			stats.FilesAdded++
		} else if file.Deleted {
			stats.FilesRemoved++
		} else {
			stats.FilesChanged++
		}
	}

	return stats
}

// Binary files (detected like git does, by a NUL byte) count as zero lines.
func (self *FileLoader) countUntrackedFileLines(path string) int {
	content, err := afero.ReadFile(self.Fs, path)
	if err != nil {
		self.Log.Error(err)
		return 0
	}

	if bytes.IndexByte(content, 0) != -1 {
		return 0
	}

	lineCount := bytes.Count(content, []byte{'\n'})
	if len(content) > 0 && content[len(content)-1] != '\n' {
		lineCount++
	}

	return lineCount
}
