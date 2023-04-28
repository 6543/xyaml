# eXtend YAML

[![Tests](https://ci.codeberg.org/api/badges/6543/xyaml/status.svg)](https://ci.codeberg.org/6543/xyaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/codeberg.org/6543/xyaml?status.svg)](https://godoc.org/codeberg.org/6543/xyaml)
[![Go Report Card](https://goreportcard.com/badge/codeberg.org/6543/xyaml)](https://goreportcard.com/report/codeberg.org/6543/xyaml)

<a href="https://codeberg.org/6543/xyaml">
    <img alt="Get it on Codeberg" src="https://codeberg.org/Codeberg/GetItOnCodeberg/media/branch/main/get-it-on-neon-blue.png" height="60">
</a>

is a library to extend [`gopkg.in/yaml.v3`](https://github.com/go-yaml/yaml/tree/v3)
to allow merging [sequences](https://github.com/yaml/yaml/issues/48) and [arrays](https://github.com/yaml/yaml/issues/35)

## Features

- [x] merge sequences
  - [x] single alias
  - [x] array of alias
- [ ] merge maps
  - [ ] overwrite maps

## How to use

`go get codeberg.org/6543/xyaml`

and just replace your

```go
err := yaml.Unmarshal(in, out)
```

with

```go
err := xyaml.Unmarshal(in, out)
```

## Examples

### merge sequences

```yml
array1: &my_array_alias
- foo
- bar

array2:
- <<: *my_array_alias
- baz
```

will be interpreted as:

```yml
array1: &my_array_alias
- foo
- bar

array2:
- foo
- bar
- baz
```
