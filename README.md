# **Titan Small Cloud CLI** (v3.x - Legacy)

> ⚠️ **Deprecation Notice**
>
> This CLI version (v3.x) is based on the **Titan SmallCloud API v1**, which is now **deprecated**.
> It will remain functional with limited features for a transitional period, but **no new features will be added**.
>
> **Users are encouraged to migrate to the new CLI v4.x**, which fully supports the **API v2** and offers
> a comprehensive set of features for managing your Titan SmallCloud resources.
>
> 👉 [Download the new CLI v4.x](https://github.com/titandc/titan-sc-cli/releases)

---

A minimal command line interface written in Go for basic operations on [Titan Small Cloud (SC)](https://sc.titandc.net).

**Supported features (v1 API):**
- List and view server details
- Manage server snapshots (list, create, delete, rotate)

## Installation

You can build the CLI manually using go or you can just download the pre-compiled binary for your operating system and
architecture.

### 1) __Manual build__

#### Dependencies

Golang v1.23 or above is required, follow the [official documentation](https://golang.org/doc/install) to install it on
your system.

#### Instructions

```shell script
git clone https://github.com/titandc/titan-sc-cli.git
cd titan-sc-cli
go build .
./titan-sc -h
```

### 2) __Download pre-compiled static binary__

Go to the latest [release](https://github.com/titandc/titan-sc-cli/releases) and download the tarball for your operating
system and architecture.

## Configuration

### Requirements

You must first generate an API token from the [Titan SC dashboard](https://sc.titandc.net), to do so:

- login to your dashboard, click on the top right user icon and select `API keys`
- create a new token by giving it a name and an optional expiration date
- save the generated key locally

### Automated setup (preferred way)

Run the following command to automatically setup your environment (replace `your token` by the API key previously
generated):

```shell script
./titan-sc setup --token "your token"
```

This will automatically create the configuration file filled with your API key. If using Linux/Mac this will also copy
the binary on your system (see details below).

_Notes_:

- On Windows, the configuration file is located on the local folder (where the binary resides), if you intend to move
  the binary to another place then you should either move the `config` file accordingly or use the **environment
  variable instead** (see details below).
- On Linux/Mac the file resides in `${HOME}/.titan/config` and the `titan-sc` binary is copied in `/usr/local/bin/` for
  global usage.

### Manual setup

#### Using configuration file

The CLI will look for the configuration file (namely `config`) at two different locations ordered by ascending priority:

1) the binary's root folder
2) the hidden folder `.titan` on your home directory

Here is the content of the default configuration file provided (`config.sample`):

```
[default]
token = "your token"
```

You can update it manually by replacing `your token` with the content of your API key and rename it as `config`, here is
an example for Linux/Mac:

```shell script
TOKEN="..."
mkdir -p ~/.titan
sed "s/your token/${TOKEN}/g" ./config.sample > ~/.titan/config
cp ./titan-sc /usr/local/bin/
```

_Note_: The configuration file can be overrided by the environment variable `TITAN_API_TOKEN`, see details below.

#### Using environment variable

The CLI checks if the environment variable `TITAN_API_TOKEN` is defined and use it in priority to grab the API token (
higher priority than configuration files). You can therefore override the confguration file by exporting the environment
variable.

1) __On Linux/Mac__

Export the environment variable (replace `your token` by the content of your API key):

```shell script
export TITAN_API_TOKEN="your token"
```

_Note_: This can also be added to your shell configuration file (eg. for bash: `$HOME/.bashrc`) to automatically export
the variable in your shell environment.

2) __On Windows__

You can follow [this guide](https://www.computerhope.com/issues/ch000549.htm) to create an environment variable on
Windows.

## Command-line completion

Command-line completion for commands and options are supported with the following shells:

- Bash
- Zsh
- Fish
- Powershell

### Bash

First make sure to have the package `bash-completion` installed on your system.

```shell script
source <(titan-sc completion bash)
```

To load completions for each session, execute once:

##### Linux:

```shell script
titan-sc completion bash > /etc/bash_completion.d/titan-sc
```

##### MacOS:

```shell script
titan-sc completion bash > /usr/local/etc/bash_completion.d/titan-sc
```

### Zsh

If shell completion is not already enabled in your environment you will need to enable it. You can execute the following
once:

```shell script
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

To load completions for each session, execute once:

```shell script
titan-sc completion zsh > "${fpath[1]}/_titan-sc"
```

You will need to start a new shell for this setup to take effect.

### Fish

```shell script
titan-sc completion fish | source
```

To load completions for each session, execute once:

```shell script
titan-sc completion fish > ~/.config/fish/completions/titan-sc.fish
```

### Powershell

You need PowerShell version 5.0 or above, which comes with Windows 10 and can be downloaded separately for Windows 7 or
8.1. You can then source the completion file from your PowerShell profile, which is referenced by the `$Profile`
environment variable. Execute `Get-Help about_Profiles` for more info about PowerShell profiles.

Generate the completion file:

```shell script
titan-sc completion powershell > completion_file
```


## Usage

### Basics

Global help:

```
Titan Small Cloud - Command Line Interface

Usage:
  titan-sc [command]

Available Commands:
  completion       Generate completion script
  help             Help about any command
  server           Manage servers.
  setup            Automated config/install.
  snapshot         Manage servers' snapshots.

Flags:
  -C, --color   Enable colorized output.
  -h, --help    help for titan-sc
  -H, --human   Format output for human.

Use "titan-sc [command] --help" for more information about a command.

```

Get (sub)commands help:

```
titan-sc [command] --help
```

The CLI default output is in JSON but you can print a more human readable output by using the flag `--human` or `-H`:

```
titan-sc [command] --human
```

### Examples

* List all your servers:

```
titan-sc server list
```

* Show details of a specific server:

```
titan-sc server show --server-uuid ${SERVER_UUID}
```

*where `${SERVER_UUID}` is the UUID of the targeted server.*

* List snapshots for a server:

```
titan-sc snapshot list --server-uuid ${SERVER_UUID}
```

* Force create a new snapshot for your server:

```
titan-sc snapshot create --server-uuid ${SERVER_UUID} --yes-i-agree-to-erase-oldest-snapshot
```

*where `${SERVER_UUID}` is the UUID of the targeted server. The last option may be used to automatically erase oldest server's snapshot when quota has been reached.*

* Delete a snapshot:

```
titan-sc snapshot delete --server-uuid ${SERVER_UUID} --snapshot-uuid ${SNAPSHOT_UUID}
```

*where `${SERVER_UUID}` is the UUID of the targeted server and `${SNAPSHOT_UUID}` is the UUID of the snapshot to delete.*

* Rotate snapshots (create new and delete oldest if needed):

```
titan-sc snapshot rotate --server-uuid ${SERVER_UUID} --force
```

*where `${SERVER_UUID}` is the UUID of the targeted server. The `--force` flag skips confirmation prompt.*
