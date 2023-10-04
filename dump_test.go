package dump

import (
	"log"
	"os"
)

func Example() {
	info = log.New(os.Stdout, "", log.Lmsgprefix).Println

	Log("Title")(
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
	// Output:
	// `Title`
	//     `Nested:` `some`
	//         [
	//          {
	//           `key`: `val`,
	//          },
	//          [
	//           1,
	//           2,
	//           3,
	//          ],
	//          {
	//            Field: `1`
	//            privateField: 1
	//          },
	//         ]
}
