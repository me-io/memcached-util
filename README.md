## Memcached Utility
> Tiny memcached utility that allows you to backup and restore memcached cache

Useful for stopping/starting/restarting memcache server without sacrificing any cached data.

## Installation

Download and install from the [releases page](https://github.com/me-io/memcached-util/releases)

## Usage

```sh
# Generate the backup in output.json
memcached-util --backup --name "output.json"
# Restores the backup from output.json
memcached-util --restore --name "output.json"
```

**Note**: `--name` is optional and if not given it will generate the backup file in the current directory. Also it can be a full path.

## Contributing

Anyone is welcome to [contribute](CONTRIBUTING.md), however, if you decide to get involved, please take a moment to review the guidelines:

* [Only one feature or change per pull request](CONTRIBUTING.md#only-one-feature-or-change-per-pull-request)
* [Write meaningful commit messages](CONTRIBUTING.md#write-meaningful-commit-messages)
* [Follow the existing coding standards](CONTRIBUTING.md#follow-the-existing-coding-standards)


## License

The code is available under the [MIT license](LICENSE.md).
