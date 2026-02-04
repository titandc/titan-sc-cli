# Titan Small Cloud CLI

A command-line interface for managing resources on [Titan Small Cloud](https://sc.titandc.net).

> **Version 4.x** — This CLI targets the Titan SC API v2 and is the recommended version. It replaces the deprecated v3.x CLI which targeted API v1 and is now limited to a subset of endpoints. We strongly recommend migrating to v4.x for full feature support.

## Features

- **Server Management** — Start, stop, restart, reset, and configure servers
- **Template Management** — List available OS templates and user images
- **Network Management** — Create and manage private networks, attach/detach servers
- **Disaster Recovery (DRP)** — Manage disaster recovery between main and secondary sites
- **Snapshot Management** — Create, restore, rotate, and delete server snapshots
- **IP Management** — Attach, detach, and configure IP addresses and reverse DNS
- **KVM Access** — Remote console access via KVM IP with web URL
- **SSH Key Management** — Add, list, and delete SSH keys
- **API Token Management** — Create, update, and delete API tokens
- **Subscription & Billing** — View subscriptions and billing information
- **Event History** — Track events on servers and companies
- Human-readable output with colored formatting
- JSON output for scripting and automation
- Shell completion for Bash, Zsh, Fish, and PowerShell

## Installation

### Download Pre-compiled Binary

Download the latest [release](https://github.com/titandc/titan-sc-cli/releases) for your operating system and architecture.

### Build from Source

Requires Go 1.24 or later.

```sh
git clone https://github.com/titandc/titan-sc-cli.git
cd titan-sc-cli
make deps   # Download dependencies
make build  # Build the binary
./titan-sc --help
```

To install system-wide (Linux/macOS):

```sh
sudo make install  # Installs to /usr/local/bin
```

### Available Make Targets

```sh
make help       # Show all available targets
make deps       # Download and tidy dependencies
make build      # Build the CLI binary
make install    # Install to /usr/local/bin (requires sudo)
make clean      # Remove build artifacts
```

## Configuration

### Generate an API Token

1. Log in to the [Titan SC dashboard](https://sc.titandc.net)
2. Click the user icon (top right) → **API keys**
3. Create a new token with an optional expiration date
4. Save the generated key

### Setup

Run the setup command to configure the CLI:

```sh
titan-sc setup --token "your-api-token"
```

This validates your token and saves the configuration to:
- **Linux/macOS**: `~/.titan/config`
- **Windows**: `%APPDATA%\titan\config`

For custom API endpoints:

```sh
titan-sc setup --token "your-api-token" --uri "https://custom-api.example.com"
```

### Alternative: Environment Variable

Set the `TITAN_API_TOKEN` environment variable to override any configuration file:

```sh
export TITAN_API_TOKEN="your-api-token"
```

## Usage

### Output Formats

By default, the CLI outputs human-readable formatted tables with colors.

```sh
titan-sc server list
```

For JSON output (useful for scripting):

```sh
titan-sc server list --json
titan-sc server list -j
```

To disable colors in human-readable mode:

```sh
titan-sc server list --no-color
```

### Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `server` | `srv` | Manage servers |
| `template` | | Manage OS templates and user images |
| `network` | `net` | Manage private networks |
| `snapshot` | `snap` | Manage server snapshots |
| `ip` | | Manage IP addresses |
| `kvmip` | `kvm` | Manage KVM remote console access |
| `ssh-key` | | Manage SSH keys |
| `api-token` | `token` | Manage API tokens |
| `subscription` | `sub` | View billing subscriptions |
| `company` | `co` | View company information |
| `history` | `hist` | List events on servers or companies |
| `user` | | View user information |
| `version` | | Show CLI and API version |
| `setup` | | Configure CLI credentials |
| `completion` | | Generate shell completion script |

Use `titan-sc [command] --help` for detailed usage.

### Server Commands

```sh
titan-sc server list                       # List all servers
titan-sc server show --server-oid <oid>    # Show server details
titan-sc server start --server-oid <oid>   # Start a server
titan-sc server stop --server-oid <oid>    # Stop a server
titan-sc server restart --server-oid <oid> # Restart a server
titan-sc server hardstop --server-oid <oid> # Force stop a server
titan-sc server rename --server-oid <oid> --name <name>
titan-sc server addons list --server-oid <oid>  # List available addons
titan-sc server iso mount --server-oid <oid> --uri <url>
titan-sc server iso umount --server-oid <oid>
titan-sc server iso show --server-oid <oid>
titan-sc server reset --server-oid <oid> --template-oid <oid>
titan-sc server drp status --server-oid <oid>
titan-sc server drp failover-soft --server-oid <oid>
titan-sc server drp failover-hard ...      # Force failover (dangerous)
titan-sc server drp resync ...             # Resync after split-brain
```

### Template Commands

```sh
titan-sc template list                     # List all templates
titan-sc template show --oid <oid>         # Show template details
```

### Network Commands

```sh
titan-sc network list                      # List all networks
titan-sc network show --network-oid <oid>  # Show network details
titan-sc network create --name <name>      # Create a new network
titan-sc network delete --network-oid <oid> # Delete a network
titan-sc network rename --network-oid <oid> --name <name>
titan-sc network attach --server-oid <oid> --network-oid <oid>
titan-sc network detach --server-oid <oid> --network-oid <oid>
titan-sc network drp enable --network-oid <oid>
titan-sc network drp disable --network-oid <oid> --yes-i-understand-network-will-be-unavailable
```

### Snapshot Commands

```sh
titan-sc snapshot list --server-oid <oid>     # List snapshots
titan-sc snapshot create --server-oid <oid>   # Create snapshot
titan-sc snapshot restore --snapshot-oid <oid>
titan-sc snapshot delete --snapshot-oid <oid>
titan-sc snapshot rotate --server-oid <oid>   # Rotate snapshots
```

To force snapshot creation when quota is reached:

```sh
titan-sc snapshot create --server-oid <oid> --yes-i-agree-to-erase-oldest-snapshot
```

### IP Commands

```sh
titan-sc ip list                           # List available IPs
titan-sc ip attach --server-oid <oid> --ip <ip>
titan-sc ip detach --server-oid <oid> --ip <ip>
titan-sc ip reverse --ip <ip> --reverse <hostname>
```

### KVM Commands

```sh
titan-sc kvmip show --server-oid <oid>     # Show KVM info with web URL
titan-sc kvmip start --server-oid <oid>    # Start KVM session
titan-sc kvmip stop --server-oid <oid>     # Stop KVM session
```

### API Token Commands

```sh
titan-sc api-token list                    # List all API tokens
titan-sc api-token show --oid <oid>        # Show token details
titan-sc api-token create --name <name>    # Create token (never expires)
titan-sc api-token create --name <name> --expire-days 30
titan-sc api-token update --oid <oid> --name <name>
titan-sc api-token update --oid <oid> --no-expire
titan-sc api-token delete --oid <oid>
```

### Other Commands

```sh
titan-sc ssh-key list                      # List SSH keys
titan-sc ssh-key show --oid <oid>          # Show SSH key details
titan-sc ssh-key add --name <n> --value <k>
titan-sc ssh-key delete --oid <oid>

titan-sc subscription list                 # List subscriptions
titan-sc subscription show --subscription-oid <oid>

titan-sc company list                      # List companies
titan-sc company show                      # Show default company

titan-sc history --server-oid <oid>        # Server event history
titan-sc history                           # Company event history

titan-sc user info                         # Show user information
titan-sc version cli                       # Show CLI version
titan-sc version api                       # Show API version
```

### DRP Commands (Disaster Recovery Plan)

DRP provides disaster recovery between **main** and **secondary** sites. Servers and networks with DRP enabled are replicated in real-time to the secondary site.

#### Server DRP

```sh
# Check DRP status for a server
titan-sc server drp status --server-oid <oid>

# Perform soft failover (server must be stopped)
titan-sc server drp failover-soft --server-oid <oid>

# Force failover (DANGEROUS - causes data loss)
titan-sc server drp failover-hard --server-oid <oid> --target-site main --yes-i-understand-i-will-lose-data

# Resync after split-brain (DANGEROUS - overwrites target data)
titan-sc server drp resync --server-oid <oid> --authoritative-site main --yes-i-understand-i-will-lose-data
```

#### Network DRP

```sh
# Enable DRP for a network
titan-sc network drp enable --network-oid <oid>

# Disable DRP for a network (requires confirmation flag)
titan-sc network drp disable --network-oid <oid> --yes-i-understand-network-will-be-unavailable
```

> **Note**: If you disable DRP for a network and your servers fail over to the secondary site, they will lose private network connectivity until DRP is re-enabled.

#### DRP Status Values

| Status | Description |
|--------|-------------|
| Healthy | DRP is active and synchronized |
| Pending | Operation in progress |
| Disabled/Error | DRP is disabled or encountered an error |
| Split-Brain | Both sites are out of sync, manual intervention required |

> **Note**: The `--target-site` and `--authoritative-site` flags accept `main` or `secondary` as values.

## Shell Completion

Generate completion scripts for your shell:

```sh
titan-sc completion [bash|zsh|fish|powershell]
```

## Examples

### List all your servers

```sh
titan-sc server list
```

### Start all stopped servers

One-liner with `jq` and `xargs`:

```sh
titan-sc srv ls -j | jq -r '.[] | select(.state == "stopped") | .oid' | xargs -L1 titan-sc srv start -s
```

### Force create a snapshot

```sh
titan-sc snapshot create --server-oid ${SERVER_OID} --yes-i-agree-to-erase-oldest-snapshot
```

The `--yes-i-agree-to-erase-oldest-snapshot` flag automatically erases the oldest snapshot when the quota has been reached.

### Restore a snapshot

```sh
titan-sc snapshot restore --snapshot-oid ${SNAPSHOT_OID}
```

> **Warning**: The server must be stopped before restoring. This operation **erases all data** on the server's disk and replaces it with the snapshot content. It is highly recommended to create a fresh snapshot before restoring an old one to allow rollback.

### Rotate snapshots

Create a new snapshot and automatically delete the oldest one if quota is reached:

```sh
titan-sc snapshot rotate --server-oid ${SERVER_OID}
```

Use `--force` to skip confirmation and delete the oldest snapshot automatically:

```sh
titan-sc snapshot rotate --server-oid ${SERVER_OID} --force
```

## Shell Completion (Detailed)

### Bash

```sh
# Current session
source <(titan-sc completion bash)

# Permanent (Linux)
titan-sc completion bash > /etc/bash_completion.d/titan-sc

# Permanent (macOS)
titan-sc completion bash > /usr/local/etc/bash_completion.d/titan-sc
```

### Zsh

```sh
# Current session
source <(titan-sc completion zsh)

# Permanent
titan-sc completion zsh > ~/.titan/_titan-sc
echo 'fpath=(~/.titan $fpath)' >> ~/.zshrc
echo 'autoload -Uz compinit && compinit' >> ~/.zshrc
```

### Fish

```sh
# Current session
titan-sc completion fish | source

# Permanent
titan-sc completion fish > ~/.config/fish/completions/titan-sc.fish
```

### PowerShell

```powershell
# Current session
titan-sc completion powershell | Out-String | Invoke-Expression

# Permanent (add to $PROFILE)
titan-sc completion powershell >> $PROFILE
```


### Legacy Snapshot Support

For backward compatibility with existing v3.x scripts, snapshot commands support legacy `--server-uuid` and `--snapshot-uuid` flags. When using these flags, the CLI automatically targets API v1.

> **Note**: Legacy support is deprecated and will be removed in a future version. New scripts should use the API v2 flags (`--server-oid`, `--snapshot-oid`).

| Command | API v2 (recommended) | API v1 (legacy) |
|---------|---------------------|-----------------|
| `list` | `--server-oid` | `--server-uuid` |
| `create` | `--server-oid` | `--server-uuid` |
| `rotate` | `--server-oid` | `--server-uuid` |
| `delete` | `--snapshot-oid` | `--server-uuid` + `--snapshot-uuid` |
| `restore` | `--snapshot-oid` | `--snapshot-uuid` |

**Examples:**

```sh
# API v2 (recommended)
titan-sc snapshot rotate --server-oid sc-abc123 --force
titan-sc snapshot delete --snapshot-oid snap-xyz789

# API v1 (legacy, for existing scripts)
titan-sc snapshot rotate --server-uuid 12345678-1234-1234-1234-123456789abc --force
titan-sc snapshot delete --server-uuid 12345678-... --snapshot-uuid 87654321-...
```

## License

See [LICENSE](LICENSE) for details.
