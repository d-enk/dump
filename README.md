# dump

Deep multi-line dump for any object

Package dump provides a functions for dumping Go values in a human-readable format.

`Dumper` type that can be used to write a string representation of a value to `io.Writer`,

`Log` function that logs the output to stderr with a customizable prefix.

[Example usage:](dump_test.go)

```go
package main

import "github.com/d-enk/dump"

func main() {
 const multiline = `multi
line`

 dump.Log(
  multiline,
  []any{
   0, "str", false, nil,
   []any{multiline},
   struct {
    Field        any
    privateField any
   }{
    Field:        `1`,
    privateField: 1,
   },
   map[any]any{"-": multiline},
  },
 )
 // Output:
 // `multi
 // line` [
 //  0,
 //  `str`,
 //  false,
 //  nil,
 //  [
 //    `multi
 //     line`,
 //  ],
 //  {
 //    Field: `1`
 //    privateField: 1
 //  },
 //  {
 //    `-`: |
 //     `multi
 //      line`,
 //  },
 // ]
}
```
