package helpers

import (
	"fmt"

	"github.com/jesseduffield/lazycore/pkg/boxlayout"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

// Fork: the stash window's slot is taken over by the Changes panel, which
// summarises all uncommitted changes against HEAD.

const forkDiffStatsTitle = "Changes"

func forkReplaceStashWithDiffStats(windows map[string]boxlayout.Dimensions) map[string]boxlayout.Dimensions {
	if dimensions, ok := windows["stash"]; ok {
		windows["diffStats"] = dimensions
		delete(windows, "stash")
	}

	return windows
}

func (self *RefreshHelper) forkRefreshDiffStats() {
	stats := self.c.Git().Loaders.FileLoader.GetDiffStats(self.c.Model().Files)

	self.c.OnUIThread(func() error {
		view := self.c.Views().DiffStats
		view.Title = forkDiffStatsTitle
		self.c.SetViewContent(view, formatDiffStats(stats))
		return nil
	})
}

func formatDiffStats(stats git_commands.DiffStats) string {
	lineStats := fmt.Sprintf("%s %s",
		style.FgGreen.SetBold().Sprintf("+%d", stats.LinesAdded),
		style.FgRed.SetBold().Sprintf("-%d", stats.LinesDeleted),
	)

	fileStats := fmt.Sprintf("%s added  %s changed  %s removed",
		style.FgGreen.Sprint(stats.FilesAdded),
		style.FgYellow.Sprint(stats.FilesChanged),
		style.FgRed.Sprint(stats.FilesRemoved),
	)

	return lineStats + "\n" + fileStats
}
