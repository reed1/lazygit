package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// Fork: the commit list controllers point the Changes panel at their selected
// commit; focusing any other side panel switches it back to uncommitted changes.
func (gui *Gui) attachDiffStats() {
	refreshHelper := gui.helpers.Refresh

	for _, ctx := range gui.State.Contexts.Flatten() {
		if ctx.GetKind() != types.SIDE_CONTEXT || lo.Contains(helpers.CommitDiffStatsContextKeys, ctx.GetKey()) {
			continue
		}
		ctx.AddOnFocusFn(func(types.OnFocusOpts) {
			gui.c.OnWorker(func(gocui.Task) error {
				refreshHelper.RenderUncommittedDiffStats()
				return nil
			})
		})
	}
}
