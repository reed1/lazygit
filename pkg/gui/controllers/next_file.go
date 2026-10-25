package controllers

import "github.com/jesseduffield/lazygit/pkg/gui/types"

func (self *FilesController) nextFileBinding(opts types.KeybindingsOpts) *types.Binding {
	return &types.Binding{
		Key:         opts.GetKey(opts.Config.Files.NextFile),
		Handler:     self.selectNextFile,
		Description: "Next file",
		Tooltip:     "Select the next file, skipping directories, without staging or unstaging the current one.",
	}
}

func (self *FilesController) selectNextFile() error {
	self.moveToNextFile()
	return nil
}
