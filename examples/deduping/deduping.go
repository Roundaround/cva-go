package deduping

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size string
}

var Button = cva.New(
	cva.Base[Props]("inline-flex items-center justify-center rounded-md"),
	cva.Variant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-8 rounded-md",
			"medium": "h-10 rounded-md",
			"large":  "h-12 rounded-md",
		},
	),
)

func DedupedClasses(p Props) string {
	return cva.DedupeClasses(Button.Classes(p))
}

func Example() {
	fmt.Println(Button.Classes(Props{"small"}))
	// Output: inline-flex items-center justify-center rounded-md h-8 rounded-md

	fmt.Println(DedupedClasses(Props{"small"}))
	// Output: inline-flex items-center justify-center rounded-md h-8
}
