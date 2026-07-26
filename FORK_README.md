# LazyGit Fork

This is a personal fork of [LazyGit](https://github.com/jesseduffield/lazygit) with custom modifications.

## Changes from Upstream

### Auto-Navigation to Next File on Stage

When pressing the space key (stage/unstage keybinding) on a **file** in the files panel, the cursor automatically moves down to the next file in the tree. This allows for faster sequential staging/unstaging of files without manually pressing 'j' (down) after each space press.

### Edit Key Works on Folders

Pressing `e` (edit) on a **folder** in the files panel now opens it using the `os.openDirInEditor` config instead of showing an error. This allows configuring a terminal to open at the folder location.

### Tmp Commit Warning in Status Panel

When the HEAD commit's subject is exactly `tmp`, the status panel content is replaced with `⚠️ TMP COMMIT ⚠️` instead of the usual `repo → branch` line.

### Select First File on Startup

On startup, LazyGit now automatically selects the first actual changed file in the files panel instead of the root directory. This means the diff view immediately shows a specific file's changes rather than all changes combined.
