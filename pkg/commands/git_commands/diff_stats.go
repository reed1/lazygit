package git_commands

import (
	"bytes"
	"fmt"
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

	numstat, err := self.cmd.New(
		NewGitCmd("diff").
			Arg("--numstat", "-z", self.findRenamesArg(), "HEAD").
			ToArgv(),
	).DontLog().RunWithOutput()
	if err != nil {
		self.Log.Error(err)
	} else if lineStats, err := parseDiffStats(numstat); err != nil {
		self.Log.Error(err)
	} else {
		stats.LinesAdded = lineStats.LinesAdded
		stats.LinesDeleted = lineStats.LinesDeleted
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

func (self *GitCommon) findRenamesArg() string {
	return fmt.Sprintf("--find-renames=%d%%", self.UserConfig().Git.RenameSimilarityThreshold)
}

// Merge commits are compared against their first parent, so the stats cover
// what the merge brought in.
func (self *CommitCommands) GetCommitDiffStats(hash string) (DiffStats, error) {
	cmdArgs := NewGitCmd("show").
		Arg("--format=", self.findRenamesArg(), "--diff-merges=first-parent", "--raw", "--numstat", "-z").
		Arg(hash).
		ToArgv()

	output, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return DiffStats{}, err
	}

	return parseDiffStats(output)
}

// Parses NUL-separated --raw and --numstat entries. A --raw entry is
// ":<modes> <hashes> <status>" followed by one path, or by the source and
// destination paths for a rename or copy. A --numstat entry is
// "<added>\t<deleted>\t<path>", or "<added>\t<deleted>\t" followed by the
// source and destination paths for a rename or copy.
func parseDiffStats(output string) (DiffStats, error) {
	stats := DiffStats{}

	tokens := strings.Split(strings.TrimSuffix(output, "\x00"), "\x00")
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if token == "" {
			continue
		}

		if strings.HasPrefix(token, ":") {
			fields := strings.Fields(token)
			status := fields[len(fields)-1][:1]
			switch status {
			case "A":
				stats.FilesAdded++
				i++
			case "C":
				stats.FilesAdded++
				i += 2
			case "D":
				stats.FilesRemoved++
				i++
			case "M", "T":
				stats.FilesChanged++
				i++
			case "R":
				stats.FilesChanged++
				i += 2
			default:
				return DiffStats{}, fmt.Errorf("unexpected diff status %q in %q", status, token)
			}
			continue
		}

		numstat := strings.SplitN(token, "\t", 3)
		if len(numstat) != 3 {
			return DiffStats{}, fmt.Errorf("unexpected numstat entry %q", token)
		}
		// binary files report "-" and count as zero lines
		if added, err := strconv.Atoi(numstat[0]); err == nil {
			stats.LinesAdded += added
		}
		if deleted, err := strconv.Atoi(numstat[1]); err == nil {
			stats.LinesDeleted += deleted
		}
		if numstat[2] == "" {
			i += 2
		}
	}

	return stats, nil
}
