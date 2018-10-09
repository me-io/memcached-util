## Memcached Utility
> Tiny memcached utility that allows you to backup and restore memcached cache

[![Go Report Card](https://goreportcard.com/badge/github.com/me-io/memcached-util)](https://goreportcard.com/report/github.com/me-io/memcached-util)
[![GoDoc](https://godoc.org/github.com/me-io/memcached-util?status.svg)](https://godoc.org/github.com/me-io/memcached-util)
[![Blog URL](https://img.shields.io/badge/Author-blog-green.svg?style=flat-square)](https://meabed.com)
[![Build Status](https://travis-ci.org/me-io/memcached-util.svg?branch=master)](https://travis-ci.org/me-io/memcached-util)

Useful for stopping/starting/restarting memcached server without sacrificing any cached data.

## Installation

Download and install from the [releases page](https://github.com/me-io/memcached-util/releases)

## Usage

Use the below signature to generate or restore the backup

```sh
memcached-util [--backup|--restore]
               [--name <path-of-backup-file>]
               [--host <memcached-host>]
               [--port <memcached-port>]
```

Detail of options are listed below

| **Option** | **Default** | **Description** |
|--------|------|-------|
| `--backup` |  | Generate the backup |
| `--restore` | | Restore the backup |
| `--name` | mem_backup.json | Path to the file to generate the backup in or to restore from |
| `--host` | `0.0.0.0` | Memcached host to connect to |
| `--port` | `11211` | Memcached port to connect to |

### Examples

Generate the backup and store at the given path
```sh
memcached-util --host "192.168.99.100" --port "11211" --backup --name "mem_backup.json"
# Store at the given path
memcached-util --backup --name "/some/path/mem_backup.json"
```
Restore the backup from the given path
```sh
memcached-util --host "192.168.99.100" --port "11211" --restore --name "mem_backup.json"
memcached-util --restore --name "/some/path/mem_backup.json"
```

## Contributing

Anyone is welcome to [contribute](CONTRIBUTING.md), however, if you decide to get involved, please take a moment to review the guidelines:

* [Only one feature or change per pull request](CONTRIBUTING.md#only-one-feature-or-change-per-pull-request)
* [Write meaningful commit messages](CONTRIBUTING.md#write-meaningful-commit-messages)
* [Follow the existing coding standards](CONTRIBUTING.md#follow-the-existing-coding-standards)


## License

The code is available under the [MIT license](LICENSE.md).
