package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
)

func (self *RefreshHelper) HeadCommit() *models.Commit {
	self.c.Mutexes().LocalCommitsMutex.Lock()
	defer self.c.Mutexes().LocalCommitsMutex.Unlock()

	headCommit, found := lo.Find(self.c.Model().Commits, func(commit *models.Commit) bool {
		return !commit.IsTODO()
	})
	if !found {
		return nil
	}
	return headCommit
}
