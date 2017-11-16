# mmtop

This is an experimental project to rewrite [mmtop](https://github.com/osheroff/mmtop) in Go.

## Setup

To build this tool you need [govendor](https://github.com/kardianos/govendor).

``` bash
govendor sync
govendor build ./cmd/sqlgo
```

You need to supply a configuration file:

``` ini
[localhost]
username = root
password = secret

[localhost2]
hostname = localhost
username = root
password = more-secret
port = 3308
```
