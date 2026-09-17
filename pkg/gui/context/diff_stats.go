package context

import "github.com/jesseduffield/lazygit/pkg/gui/types"

// Fork: display context for the Changes panel shown in place of the stash window.

const DIFF_STATS_CONTEXT_KEY types.ContextKey = "diffStats"

func NewDiffStatsContext(c *ContextCommon) types.Context {
	return NewDisplayContext(DIFF_STATS_CONTEXT_KEY, c.Views().DiffStats, "diffStats")
}
