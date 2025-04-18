# BADM - Born Again Dotfile Manager

BADM is a new (yes, again a new one) Dotfile manager. It uses Git as its backbone.

:warning: This Dotfiles manager is still a **work in progress**, some features are not implemented right now. Have a
look at [Planned features / Known issues](#planned-features--known-issues) to see what is planned to be added.

## Setup

If it's your first time you are using this tool, you might want to follow these setup steps:

### Create a fresh Dotfiles repository

In case you do not have a repository to manage your Dotfiles (or you want to start over) it is pretty straight forward
to get going:

1. Create new Dotfiles repository
    1. Create a remote repository (e.g. on GitHub, GitLab, ...). Make sure you are able to push and pull to that
       repository.
    2. `badm new <REMOTE_REPO_CLONE_URL>`
2. Start managing your Dotfiles
    1. `badm add <PATH_TO_FILE> <PATH_TO_FILE2>` &rarr; Your Dotfiles get automatically pushed to your remote repository
    2. `badm rm <PATH_TO:FILE> <PATH_TO_FILE2>` &rarr; Your files get removed from the Dotfiles repository and restored
       to its original location. **No files get lost!**
    3. `badm save` &rarr; Write any changes to the local Dotfiles repository to the remote one. Needed after each change
       to local Dotfiles.

## Commands

| Command      | Action                                                                                                                                                                          |
|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `add`        | Add a Dotfile (or multiple) to your remote Dotfile repository. It gets replaced with a symlink. Choose the `--override` flag to override already existing files on your system. |
| `rm`         | Remove a Dotfile (or multiple) from your remote Dotfile repository. The symlink gets replaced with the original file.                                                           |
| `save`       | Write any changes to the local Dotfiles repository to the remote one. Needed after each change to local Dotfiles.                                                               |
| `fetch`      | Manually fetch the current remote repository stage and persist remote changes onto the local system.                                                                            |
| `new`        | Initialize a fresh local BADM repository, where new Dotfiles can be added afterwards.                                                                                           |
| `get`        | Pull an already existing BADM remote repository and persist everything on the system.                                                                                           |
| `help`       | Print helpful information.                                                                                                                                                      |
| `version`    | Check your current version.                                                                                                                                                     |
| `completion` | Generate the autocompletion script for the specified shell.                                                                                                                     |

## Planned features / Known issues

- Remove dead symlinks (file got remotely deleted)
- Undo all BADM changes to the system
- Execute custom scripts on wish (e.g. pacman, yay, ...)
- Actual configuration features
