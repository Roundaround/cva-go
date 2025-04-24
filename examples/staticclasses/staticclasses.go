package staticclasses

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size          string
	CustomClasses []string // Additional classes i.e. from parent components
}

var Button = cva.NewCva(
	cva.WithStaticClasses[Props]("inline-flex items-center justify-center"),
	cva.WithVariant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-9 px-3",
			"medium": "h-10 px-4 py-2",
			"large":  "h-11 px-8 py-3",
		},
	),
	cva.WithPropsClasses(func(p Props) []string {
		return p.CustomClasses
	}),
)

func Example() {
	fmt.Println(Button.ClassName(Props{"small", []string{"bg-red-500"}}))
	// Output: inline-flex items-center justify-center h-9 px-3 bg-red-500
}
