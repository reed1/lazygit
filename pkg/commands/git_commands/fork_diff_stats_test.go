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
