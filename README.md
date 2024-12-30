# gohole

A simple DNS server / hosts file manager for blocking access to specified domains.

## Installation

Prebuilt binaries are available in [Releases](https://github.com/ducng99/gohole/releases/latest)

Store the file in a new directory.

## Usage

### Add a new source

A source is an URL to a file with hosts entries. To add to gohole, run

```sh
gohole add <url>
```

### Update source(s)

Sources can be updated by running

```sh
gohole update
```

or to update a specific source, run

```sh
gohole update <id>
```

where `id` is from the command below.

### Listing current sources

To get a list of all added sources, run

```sh
gohole ls
```

### Hosts file

> [!WARNING]
> Windows users should not use this method, large hosts file can cause system hang on start up.
> See [DNS server](#dns-server)

To update the local hosts file with domain entries from sources, run

```sh
gohole hosts
```

### DNS server

To start gohole as a DNS server, run

```sh
gohole dns start
```

You can then use `127.0.0.1` as your DNS resolver.

For Windows, see [this](https://www.windowscentral.com/how-change-your-pcs-dns-settings-windows-10#section-how-to-change-dns-settings-using-control-panel-on-windows-11)
For Linux, see [this](https://www.makeuseof.com/find-and-change-dns-server-on-linux/#how-to-change-dns-server-on-linux) - you should know how already mate
For MacOS, see [this](https://support.apple.com/guide/mac-help/change-dns-settings-on-mac-mh14127/mac)

#### Autostart DNS server

On Windows, to automagically start gohole DNS server on system start up, run

```sh
gohole dns autostart
```

This will create a new task in Task Scheduler and run it.

### Help

More commands can be found by running

```sh
gohole help
```

## License

See [LICENSE](./LICENSE)
