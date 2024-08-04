# sqlc-gen-perl

## Usage

```yaml
version: '2'
plugins:
- name: perl
  wasm:
    url: TODO
    sha256: TODO
sql:
- schema: schema.sql
  queries: query.sql
  engine: postgresql
  codegen:
  - plugin: perl
    out: db
```

## Building from source

Assuming you have the Go toolchain set up, from the project root you can simply `make all`.

```sh
make all
```

This will produce a standalone binary and a WASM blob in the `bin` directory.
They don't depend on each other, they're just two different plugin styles. You can
use either with sqlc, but we recommend WASM and all of the configuration examples
here assume you're using a WASM plugin.

To use a local WASM build with sqlc, just update your configuration with a `file://`
URL pointing at the WASM blob in your `bin` directory:

```yaml
plugins:
- name: perl
  wasm:
    url: file:///path/to/bin/sqlc-gen-perl.wasm
    sha256: ""
```

As-of sqlc v1.24.0 the `sha256` is optional, but without it sqlc won't cache your
module internally which will impact performance.

