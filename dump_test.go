package dump_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/d-enk/dump"
)

func ExampleDump() {
	const multiline = `multi
line`

	fmt.Println(dump.Dump(
		multiline,
		[]any{
			0, "str", false, nil,
			[]any{multiline},
			map[any]any{"-": multiline},
		},
	))
	// Output:
	// `multi
	//  line` [
	//   0,
	//   `str`,
	//   false,
	//   nil,
	//   [
	//     `multi
	//      line`,
	//   ],
	//   {
	//     `-`: |
	//      `multi
	//       line`,
	//   },
	// ]
}

func ExampleLog() {
	dump.LogWriter = os.Stdout // only for test

	l := dump.Log("Title")(
		"Nested:", "some",
	)
	l("Next nested")("...")
	l(
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
	// Output:
	// `Title`
	//   `Nested:` `some`
	//     `Next nested`
	//       `...`
	//     [
	//       {
	//         `key`: `val`,
	//       },
	//       [
	//         1,
	//         2,
	//         3,
	//       ],
	//       {
	//         Field: `1`
	//         privateField: 1
	//       },
	//     ]
}

func ExampleDumper() {
	builder := &strings.Builder{}
	dumper := dump.New(builder)
	dumper.WithPrefix("---").Dumpln("With prefix")
	dumper.Dumpln(map[any]any{
		"1": []any{
			false,
			nil,
			"",
		},
	})
	dumper.Dump(
		1, 2, 3,
		struct{ A4, A5, A6 int }{4, 5, 6},
	)
	fmt.Println(builder.String())
	// Output:
	// ---`With prefix`
	// {
	//   `1`: [
	//     false,
	//     nil,
	//     ``,
	//   ],
	// }
	// 1 2 3 {
	//   A4: 4
	//   A5: 5
	//   A6: 6
	// }
}
