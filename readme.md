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
memcached-util [--op <backup|restore>]
               [--filename <path-of-backup-file>]
               [--addr <address-with-port>]
```

Detail of options are listed below

| **Option** | **Default** | **Description** |
|--------|------|-------|
| `--addr` | `0.0.0.0:11211` | Memcached address to connect to |
| `--filename` | `mem_backup.json` | Path to the file to generate the backup in or to restore from |
| `--op` | | `backup` or `restore |

### Examples

Generate the backup and store at the given path
```sh
memcached-util --addr "192.168.99.100:11211" --op "backup" --filename "mem_backup.json"
memcached-util --op "backup" --filename "/some/path/mem_backup.json"
```
Restore the backup from the given path
```sh
memcached-util --addr "192.168.99.100:11211" --op "restore" --filename "mem_backup.json"
memcached-util --op "restore" --filename "/some/path/mem_backup.json"
```

## Contributing

Anyone is welcome to [contribute](CONTRIBUTING.md), however, if you decide to get involved, please take a moment to review the guidelines:

* [Only one feature or change per pull request](CONTRIBUTING.md#only-one-feature-or-change-per-pull-request)
* [Write meaningful commit messages](CONTRIBUTING.md#write-meaningful-commit-messages)
* [Follow the existing coding standards](CONTRIBUTING.md#follow-the-existing-coding-standards)


## License

The code is available under the [MIT license](LICENSE.md).
