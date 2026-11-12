package presentation

import (
	"os"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// RS_GIT_TMP_MIRROR=1 replaces the status panel line with a warning.
func ShowTmpMirrorWarning() bool {
	return !utils.UpstreamBehavior() && os.Getenv("RS_GIT_TMP_MIRROR") == "1"
}

func FormatTmpMirrorWarning() string {
	return style.FgRed.SetBold().Sprint("⚠️ TMP_MIRROR ⚠️")
}
