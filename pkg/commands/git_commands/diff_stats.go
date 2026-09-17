package git_commands

import (
	"bytes"
	"strconv"
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

// Merge commits are compared against their first parent, so the stats cover
// what the merge brought in.
func (self *CommitCommands) GetCommitDiffStats(hash string) (DiffStats, error) {
	cmdArgs := NewGitCmd("show").
		Arg("--format=", "--no-renames", "--diff-merges=first-parent", "--raw", "--numstat", "-z").
		Arg(hash).
		ToArgv()

	output, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return DiffStats{}, err
	}

	return parseCommitDiffStats(output), nil
}

// Parses --raw entries (":<modes> <hashes> <status>" followed by a path token)
// and --numstat entries ("<added>\t<deleted>\t<path>"), all NUL-separated.
func parseCommitDiffStats(output string) DiffStats {
	stats := DiffStats{}

	tokens := strings.Split(output, "\x00")
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if strings.HasPrefix(token, ":") {
			fields := strings.Fields(token)
			switch fields[len(fields)-1] {
			case "A":
				stats.FilesAdded++
			case "D":
				stats.FilesRemoved++
			default:
				stats.FilesChanged++
			}
			i++
			continue
		}

		numstat := strings.SplitN(token, "\t", 3)
		if len(numstat) != 3 {
			continue
		}
		// binary files report "-" and count as zero lines
		if added, err := strconv.Atoi(numstat[0]); err == nil {
			stats.LinesAdded += added
		}
		if deleted, err := strconv.Atoi(numstat[1]); err == nil {
			stats.LinesDeleted += deleted
		}
	}

	return stats
}
