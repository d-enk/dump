# dump

Deep multi-line dump for any object

Package dump provides a functions for dumping Go values in a human-readable format.

`Dumper` type that can be used to build up a string representation of a value,

`Log` function that logs the output to stderr with a customizable prefix.

Example usage:

```go
import "github.com/d-enk/dump"

func main() {
  dump.Log("Title")(
    "Nested:", "some",
  )(
   []any{
     map[any]any{"key": "val"},
     []int{1, 2, 3},
     struct {
       Field        any
       privateField any
     }{
       Field:        `1`,
       privateField: 1,
     },
  },
  )
}
```

```go

`Title`
    `Nested:` `some`
        [
         {
          `key`: `val`,
         },
         [
          1,
          2,
          3,
         ],
         {
           Field: `1`
           privateField: 1
         },
        ]
```
