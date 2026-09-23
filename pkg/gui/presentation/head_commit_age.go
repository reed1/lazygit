package presentation

import (
	"os"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Like BranchStatus, but with LAZYGIT_REED_FORK=1 the in-sync checkmark is
// replaced by the head commit's age.
func StatusPanelBranchStatus(
	branch *models.Branch,
	headCommit *models.Commit,
	itemOperation types.ItemOperation,
	tr *i18n.TranslationSet,
	now time.Time,
	userConfig *config.UserConfig,
) string {
	status := BranchStatus(branch, itemOperation, tr, now, userConfig)
	if os.Getenv("LAZYGIT_REED_FORK") == "1" && headCommit != nil && status == style.FgGreen.Sprint("✓") {
		return style.FgCyan.Sprint(utils.UnixToTimeAgo(headCommit.UnixTimestamp))
	}
	return status
}
