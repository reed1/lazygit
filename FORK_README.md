# LazyGit Fork

This is a personal fork of [LazyGit](https://github.com/jesseduffield/lazygit) with custom modifications.

## Running Tests

Every change below is switched off when `LAZYGIT_UPSTREAM_BEHAVIOR=1` is set, so upstream's tests run unmodified against upstream behavior. The `make` test targets (`make test`, `make unit-test`, `make integration-test-*`) set it; when running `go test` directly, set it yourself:

```
LAZYGIT_UPSTREAM_BEHAVIOR=1 go test ./...
```

Day-to-day use leaves it unset.

## Changes from Upstream

### Auto-Navigation to Next File on Stage

When pressing the space key (stage/unstage keybinding) on a **file** in the files panel, the cursor automatically moves down to the next file in the tree. This allows for faster sequential staging/unstaging of files without manually pressing 'j' (down) after each space press.

### Edit Key Works on Folders

Pressing `e` (edit) on a **folder** in the files panel now opens it using the `os.openDirInEditor` config instead of showing an error. This allows configuring a terminal to open at the folder location.

### Head Commit Age in Status Panel

The status panel shows the HEAD commit's age (e.g. `5h`) in place of the `✓` that marks a branch in sync with its upstream.

### Changes Panel Instead of Stash and Command Log

The stash panel at the bottom of the side column is replaced by a non-focusable **Changes** panel summarising all uncommitted changes against `HEAD` (staged, unstaged and untracked):

```
+120 -34
3 added  5 changed  1 removed
```

The first line is the total lines added and deleted; untracked files count every line as added. The second line counts files that are added (including untracked), removed, and otherwise changed.

While a commit list is focused (commits, reflog, or a branch's commits) the panel shows the selected commit's stats instead, and its title carries the short hash (`Changes - 56f99235`). Opening the commit's files keeps those stats; focusing any other side panel switches back to uncommitted changes. Merge commits are compared against their first parent.

The command log panel below the main view is hidden, so the main view extends to the bottom of the screen. Its keybindings are kept but do nothing (`@` opens no menu), and stash is no longer reachable via side-window navigation (`5`, tab).

Code changes are pooled in files of their own (`git diff --diff-filter=A --name-only upstream/master...HEAD` lists them); every patch to an upstream file is marked with a `// fork:` comment (`rg "// fork:"` lists them).

### Select First File on Startup

On startup, LazyGit now automatically selects the first actual changed file in the files panel instead of the root directory. This means the diff view immediately shows a specific file's changes rather than all changes combined.
