# dump

Deep multi-line dump for any object

Package dump provides a functions for dumping Go values in a human-readable format.

[`Dumper`](dump.go#L22) type that can be used to write a string representation of a value to `io.Writer`

[`Log`](dump.go#L30) function that logs the result of dump working to stderr
and return `Log` function with nested level prefix

[`Dump`](dump.go#L42) function that return string result of dump working

[Example usage:](dump_test.go)

```go
package main

import "github.com/d-enk/dump"

func main() {
 const multiline = `multi
line`

 dump.Log("Title")(
  "Nested",
 )(
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
 // Title
 //   Nested
 //     `multi
 //     line` [
 //      0,
 //      `str`,
 //      false,
 //      nil,
 //      [
 //        `multi
 //         line`,
 //      ],
 //      {
 //        Field: `1`
 //        privateField: 1
 //      },
 //      {
 //        `-`: |
 //         `multi
 //          line`,
 //      },
 //     ]
}
```
