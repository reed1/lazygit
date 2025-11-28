# LazyGit Fork

This is a personal fork of [LazyGit](https://github.com/jesseduffield/lazygit) with custom modifications.

## Changes from Upstream

### Auto-Navigation to Next File on Stage

When pressing the space key (stage/unstage keybinding) on a **file** in the files panel, the cursor automatically moves down to the next file in the tree. This allows for faster sequential staging/unstaging of files without manually pressing 'j' (down) after each space press.

### Edit Key Works on Folders

Pressing `e` (edit) on a **folder** in the files panel now opens it using the `os.openDirInEditor` config instead of showing an error. This allows configuring a terminal to open at the folder location:
