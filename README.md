# Statusctl

**created and maintained by: [Nick Gerace](https://nickgerace.dev)**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fnickgerace%2Fstatusctl%2Fbadge&style=flat)](https://actions-badge.atrox.dev/nickgerace/statusctl/goto)
[![Go Report Card](https://goreportcard.com/badge/github.com/nickgerace/statusctl)](https://goreportcard.com/report/github.com/nickgerace/statusctl)

CLI tool to help keep track of your Git repositories.

## Disclaimer

This tool is being tested and is considered use at your own risk until it reaches version 1.0.

## Requirements

1. Go 1.13+
2. Make

## Installation

There will be further instructions to come.
For now, you can get started by executing the following commands.

The "install" target builds the binary and moves it to your */usr/local/bin* directory.
Moving the binary may require root (sudo) access.

After building and moving the binary, you can safely remove the repository.

```bash
git clone --depth=1 https://github.com/nickgerace/statusctl.git
make -f statusctl/Makefile install
rm -ri statusctl
```

## Usage

Execute without any arguments or flags to get started!

```bash
statusctl
```

## Uninstallation

If you still have this repository cloned, you can use the uninstall target.

```bash
make uninstall
```

Otherwise, remove the generated, configuration file and directory..

```bash
rm -ri ~/.config/statusctl/
```

Afterwards, remove the binary from your executables directory.
This may require root (sudo) access.

```bash
sudo rm -i /usr/local/bin/statusctl
```

## License

This repository is under the MIT License.
