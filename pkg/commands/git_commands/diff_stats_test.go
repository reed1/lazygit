package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFileLoaderGetDiffStats(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"diff", "--numstat", "-z", "HEAD"},
			"4\t1\tmodified.txt\x002\t0\tstaged.txt\x000\t3\tdeleted.txt\x00-\t-\timage.png",
			nil,
		)

	fs := afero.NewMemMapFs()
	assert.NoError(t, afero.WriteFile(fs, "untracked.txt", []byte("one\ntwo\nthree"), 0o644))
	assert.NoError(t, afero.WriteFile(fs, "untracked.bin", []byte("a\x00b\n"), 0o644))

	newFile := func(path string, shortStatus string) *models.File {
		file := &models.File{Path: path}
		models.SetStatusFields(file, shortStatus)
		return file
	}
	files := []*models.File{
		newFile("modified.txt", "MM"),
		newFile("staged.txt", "A "),
		newFile("deleted.txt", " D"),
		newFile("image.png", "M "),
		newFile("untracked.txt", "??"),
		newFile("untracked.bin", "??"),
		newFile("untracked-dir/", "??"),
	}

	loader := &FileLoader{
		GitCommon: buildGitCommon(commonDeps{appState: &config.AppState{}, userConfig: &config.UserConfig{}, fs: fs}),
		cmd:       oscommands.NewDummyCmdObjBuilder(runner),
	}

	assert.Equal(t, DiffStats{
		LinesAdded:   9,
		LinesDeleted: 4,
		FilesAdded:   4,
		FilesChanged: 2,
		FilesRemoved: 1,
	}, loader.GetDiffStats(files))
}

func TestCommitCommandsGetCommitDiffStats(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"show", "--format=", "--no-renames", "--diff-merges=first-parent", "--raw", "--numstat", "-z", "abc123"},
			":000000 100644 000000000 89fb3e525 A\x00new.go\x00"+
				":100644 100644 653944314 6f16e8a8c M\x00:odd name.go\x00"+
				":100644 000000 bb36ea03d 000000000 D\x00old.go\x00"+
				":100644 100644 aaaaaaaaa bbbbbbbbb M\x00image.png\x00"+
				"69\t0\tnew.go\x0015\t3\t:odd name.go\x000\t12\told.go\x00-\t-\timage.png\x00",
			nil,
		)

	instance := buildCommitCommands(commonDeps{runner: runner})

	stats, err := instance.GetCommitDiffStats("abc123")
	assert.NoError(t, err)
	assert.Equal(t, DiffStats{
		LinesAdded:   84,
		LinesDeleted: 15,
		FilesAdded:   1,
		FilesChanged: 2,
		FilesRemoved: 1,
	}, stats)
}
