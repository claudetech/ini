# ini file for Go

This module is a ini file parser for Golang.
It tries to have an API as close as possible to the
standard library.

Full documentation is available at

https://godoc.org/github.com/claudetech/ini

## Usage

Here is a sample usage, using the wrapper to read files.

```go
package main

import (
  "fmt"
  "github.com/claudetech/ini"
)

func main() {
  var conf ini.Config
  if err := ini.DecodeFile("/path/to/ini", &conf); err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(conf)
}
```

The type `ini.Config` is only an alias for

```
map [string]map[string]string
```

Note that you can pass any interface to receive the result
as long as it is supported by the [mapstructure](https://github.com/mitchellh/mapstructure) package.

The more general version uses the `ini.Decoder` structure.
A `ini.Decoder` can be created with `ini.NewDecoder` and takes
anything that responds to the `io.Reader` interface.

```go
package main

import (
  "fmt"
  "github.com/claudetech/ini"
  "os"
)

func exitError(err error) {
  fmt.Println(err)
  os.Exit(1)
}

func main() {
  var conf ini.Config
  file, err := os.Open("/etc/php/php.ini")
  if err != nil {
    exitError(err)
  }
  d := ini.NewDecoder(file)
  d.IdRegexp("^[a-z][a-z0-9_\\. -]+$")
  if err = d.Decode(&conf); err != nil {
    exitError(err)
  }
  fmt.Println(conf)
}
```

## Configuration

There are several configuration options for
the decoder.

```go
type Options struct {
  IdRegexp     string // default: "^[a-z][a-z0-9_]+$"
  SepChars     []byte // default: []byte{'='}
  CommentChars []byte // default: []byte{';'}
  LowCaseIds   bool   // default: true
}
```

You can pass an `ini.Options` to `ini.NewDecoderWithOptions`,
or you can set it directly on the `ini.Decoder` object through
the setters of the same name.

For example, if you want `;` and `#` to stand for comments,
`:` and `=` to be separators, and sections keys to contain
anything, and the keys not to be lower-cased,
you could write the following.

```go
options := ini.Options{".*", []byte{'=', ':'}, []byte{';', '#'}, false}
d := ini.NewDecoderWithOptions(file, options)
```
