package helpers

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazycore/pkg/boxlayout"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// Fork: the stash window's slot is taken over by the Changes panel, which
// summarises all uncommitted changes against HEAD, or the selected commit's
// changes while browsing commits.

const diffStatsTitle = "Changes"

func replaceStashWithDiffStats(windows map[string]boxlayout.Dimensions) map[string]boxlayout.Dimensions {
	if dimensions, ok := windows["stash"]; ok {
		windows["diffStats"] = dimensions
		delete(windows, "stash")
	}

	return windows
}

func (self *RefreshHelper) refreshDiffStats() {
	if self.diffStatsFollowsCommit() {
		return
	}

	self.RenderUncommittedDiffStats()
}

// While a commit list (or the files of a commit opened from it) is the side
// context, the Changes panel shows the selected commit's stats instead.
var CommitDiffStatsContextKeys = []types.ContextKey{
	context.LOCAL_COMMITS_CONTEXT_KEY,
	context.REFLOG_COMMITS_CONTEXT_KEY,
	context.SUB_COMMITS_CONTEXT_KEY,
	context.COMMIT_FILES_CONTEXT_KEY,
}

func (self *RefreshHelper) diffStatsFollowsCommit() bool {
	return lo.Contains(CommitDiffStatsContextKeys, self.c.Context().CurrentSide().GetKey())
}

func (self *RefreshHelper) RenderUncommittedDiffStats() {
	stats := self.c.Git().Loaders.FileLoader.GetDiffStats(self.c.Model().Files)

	self.c.OnUIThread(func() error {
		self.setDiffStatsContent(diffStatsTitle, formatDiffStats(stats))
		return nil
	})
}

func (self *RefreshHelper) RenderCommitDiffStats(commitContext interface{ GetSelected() *models.Commit }) {
	commit := commitContext.GetSelected()
	// update-ref todos in an interactive rebase have no hash
	if commit == nil || commit.Hash() == "" {
		self.setDiffStatsContent(diffStatsTitle, "")
		return
	}

	hash := commit.Hash()
	self.c.OnWorker(func(gocui.Task) error {
		stats, err := self.c.Git().Commit.GetCommitDiffStats(hash)

		self.c.OnUIThread(func() error {
			selected := commitContext.GetSelected()
			if selected == nil || selected.Hash() != hash || !self.diffStatsFollowsCommit() {
				return nil
			}

			title := fmt.Sprintf("%s - %s", diffStatsTitle, utils.ShortHash(hash))
			if err != nil {
				self.c.Log.Error(err)
				self.setDiffStatsContent(title, "")
				return nil
			}

			self.setDiffStatsContent(title, formatDiffStats(stats))
			return nil
		})
		return nil
	})
}

func (self *RefreshHelper) setDiffStatsContent(title string, content string) {
	view := self.c.Views().DiffStats
	view.Title = title
	self.c.SetViewContent(view, content)
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
