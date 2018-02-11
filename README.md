[![Build Status](https://travis-ci.org/cwarner818/tt.svg?branch=master)](https://travis-ci.org/cwarner818/tt)

IOTΛ Tangle Tools
=================

This program contains tools to get information about the IOTΛ tangle that may
not be easily accessible otherwise. 

Configuration
-------------

All options can be set either on the command line with long (`--address`) and
short (`-a`) flags, as well as having saved values in a configuration file
located in `$HOME/.tt.yaml`. Note that other configuration file formats are
available but YAML is probably the easiest to use. If you're feeling adventures
rty making a `.tt.json` or `.tt.toml` file. HCL and Java properties config file
formats are also supported. 

An example `.tt.yaml` would be:
```yaml
node: http://nodes.iota.fm:80
timeout: 30s
```

Usage
-----
`tt <command> <flags>` is the standard usage. You can use `tt help <command>`
for command specific help.

Commands
--------

`confirms` -- returns confirmation percentage for a random sampling of
transactions sent to the specified address.

