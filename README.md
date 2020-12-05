A library for specifying config via hierarchical overrides.

You provide genconf:

1. An arbitrary hierarchy of JSON config:

```
fixtures/
├── default.json
└── env
    ├── prod
    │   ├── color
    │   │   ├── blue.json
    │   │   └── green.json
    │   └── type
    │       └── foo.json
    └── qa.json
```

2. An environment to decide which config to use:

```go
map[string]string{"env": "prod", "color": "green", "type": "foo"}
```

3. An ordering in which to apply the config:

```go
[]string{"env", "color", "type"}
```

Then genconf will deserialize and merge your config into your struct.

## Example

```go
package main

import (
  "fmt"

  "github.com/jpittis/genconf"
)

// Your config struct.
type SomeTarget struct {
  A string
  B string
  C string
  D string
}

func main() {
  // Your config can start with default values.
  target := SomeTarget{
    C: "baz",
  }

  genconf.Load(
    // Path to config tree.
    "./fixtures",
    // Order to apply dimensions of config.
    []string{"env", "type", "color"},
    // Which environment to generate config for.
    map[string]string{"env": "prod", "color": "green", "type": "foo"},
    // Unmarshal into this struct.
    &target,
  )

  fmt.Println("Loaded: ", target)
}
```

```
$ go run main.go
Loaded:  {foo bar baz lol}
```

## Ideas

- Support YAML and protobuf via generic unmarshalling.
- Provide fancier merge operations with recursion for nested config.
- Provide an executable to test arbitrary environments.
- Provide an executable to list every combination of environment and output.
- Assert on invalid inputs like duplicate tuples.
