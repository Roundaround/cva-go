package simplecase

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size string
}

var Button = cva.New(
	cva.Base[Props]("inline-flex items-center justify-center"),
	cva.MapVariant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-9 px-3",
			"medium": "h-10 px-4 py-2 rounded-md",
			"large":  "h-11 px-8 py-3 rounded-md",
		},
	),
)

func Example() {
	fmt.Println(Button.Classes(Props{"medium"}))
	// Output: inline-flex items-center justify-center h-10 px-4 py-2 rounded-md
}
