package additionalclasses

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size    string
	Classes []string // Additional classes i.e. from parent components
}

var Button = cva.New(
	cva.Base[Props]("inline-flex items-center justify-center"),
	cva.MapVariant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-9 px-3",
			"medium": "h-10 px-4 py-2",
			"large":  "h-11 px-8 py-3",
		},
	),
	cva.Classes(func(p Props) []string { return p.Classes }),
)

func Example() {
	fmt.Println(Button.Classes(Props{"small", []string{"bg-red-500"}}))
	// Output: inline-flex items-center justify-center h-9 px-3 bg-red-500
}
