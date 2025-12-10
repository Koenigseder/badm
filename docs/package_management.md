# Package management setup

With **BADM** you are able to manage all your packages. This comes in handy if you decide to set up a new installation
with all necessary packages but do not want to manually install each and every package manually.

Also, you can define different **installation scopes* where you are able to manage as granular as you want which
packages
to install.

## Relevant commands

| Command    | Action            |
|------------|-------------------|
| `packages` | Install packages. |

## Setup

1. Create a config file `.badm.yaml` in your dotfiles repository (it's located under `/home/<user>/.dotfiles`)
2. Edit it as you will and define various `packageManagers`, `installScopes` and `scripts`
3. Use `badm packages <installScope>` to install all packages and execute all scripts defined by this installation scope

### Config file `.badm.yaml`

For a better DevEx writing the config file you can find the according **JSON schema
mapping** [badm.schema.json](../badm.schema.json) in this repository which can be used in e.g. IntelliJ for e.g. type
hints.

## Structure

### `packages`

- **`packageManagers`**: A map of package managers and their installation commands. The key is used as reference.
- **`installScopes`**: A map which defines installation variations (e.g., `basics`, `slim`, `full`). The key is used as
  reference.
    - **`packages`**: A map with packages to install using the defined package managers (package managers are referenced
      by the key defined at `packageManagers`). The order is respected.
    - **`dependsOn`**: (Optional) Name of another installation scope that this scope depends on and which gets executed
      prior.
    - **`scripts`**: (Optional) List of scripts to execute after package installation. The order is respected.

### `scripts`

- Defines scripts that can be referenced in installation scopes.
    - **`exec`**: The command or script to execute.
    - **`shell`**: The shell in which the script will be executed (e.g., `bash`).

## Example

An example can be found [here](example.badm.yaml).
