
# ChangeLog

## 3.0.1

- Handle snapshots dates as miliseconds timestamps.
- Format dates as ISO o human output

## 3.0.0

- Implement 'run' middleware package to fully separate 'api' and 'cmd'
- Add new interactive 'snapshot rotate' command
- Add new flag to 'snapshot create' to force rotation

## 2.1.5

- Add servers notifications output

## 2.1.4

- Add script example to rotate snapshots.
- Ensure JSON as default output for snapshot related commands.

## 2.1.3

- Fix create servers on premium accounts
- Update commands usage

## 2.1.2

- Fix snapshot listing with multi arrays.

## 2.1.1

- Add new fields on server's hypervisor description.
- Add missing option to setup documentation.

## 2.1.0

- Add command-line completion
- Use a flag for all options
- Add/del/list PNAT rules
- Updated documentation (completion section)

## 2.0.1

- Fix servers textual rendering.
- Fix networks list.
- Add managed infos to companies.
- Add firewall data to networks list.

## 2.0.0

- List server addons with pricing.
- List available OS templates.
- List user's SSH keys.
- Create new servers (regular & managed servers).
- Delete/reset existing servers.
- Add/remove/list firewall rules on managed networks.

## 1.3.0

- Allow to set API URI from config file
- Improve error management (pretty print)
- Add new 'set-gw' & 'unset-gw' commands
- Add managed services activation
- Enable managed networks creation (using CIDR)

## 1.2.0

- Read optional URI from environment var & configuration file
- Add server name & reverse update
- Add global option to colorize servers list output
- New attach/detach IP commands on server
- New command to list IPs available on a company
- New commands to show api & cli versions

## 1.1.0

- Add KVM IP infos on server details
- Add plan & reverse on server details
- Add notes field on server details
- Add pending actions on server details
- New load/unload ISO commands on server

## 1.0.0

- First version released, based on Titan SC API v1.

