package simplecase

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size string
}

var Button = cva.NewCva(
	cva.StaticClasses[Props]("inline-flex items-center justify-center"),
	cva.Variant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-9 px-3",
			"medium": "h-10 px-4 py-2",
			"large":  "h-11 px-8 py-3",
		},
	),
)

func Example() {
	fmt.Println(Button.ClassName(Props{"small"}))
	// Output: inline-flex items-center justify-center h-9 px-3
}
