package compoundvariants

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size  string
	Style string
}

var Button = cva.New(
	cva.Base[Props]("inline-flex items-center justify-center"),
	cva.Variant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-8",
			"medium": "h-10",
			"large":  "h-12",
		},
	),
	cva.Variant(
		func(p Props) string { return p.Style },
		map[string]string{
			"icon":    "bg-gray-100 rounded-full aspect-square",
			"regular": "bg-gray-100 rounded-md",
			"link":    "text-blue-500",
		},
	),
	cva.CompoundVariant(
		func(p Props) (string, string) { return p.Size, p.Style },
		cva.NewCompound("small", "icon", "[&_svg]:size-4"),
		cva.NewCompound("medium", "icon", "[&_svg]:size-5"),
		cva.NewCompound("large", "icon", "[&_svg]:size-6"),
	),
)

func Example() {
	fmt.Println(Button.Classes(Props{"small", "icon"}))
	// Output: inline-flex items-center justify-center h-8 bg-gray-100 rounded-full
	//   aspect-square [&_svg]:size-4
}
