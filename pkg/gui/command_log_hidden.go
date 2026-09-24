package gui

import "github.com/jesseduffield/lazygit/pkg/utils"

// Fork: the command log panel is never shown, so its menu and focus handlers
// are no-ops. Keybindings are kept to stay close to upstream.
var commandLogHidden = !utils.UpstreamBehavior()
