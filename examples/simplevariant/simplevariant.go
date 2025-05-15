package simplevariant

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size string
}

var size = cva.NewVariant(func(p Props) string { return p.Size })

var Button = cva.New(
	cva.Base[Props]("inline-flex items-center justify-center"),
	size.Map(map[string]string{
		"small":  "h-9 px-3",
		"medium": "h-10 px-4 py-2",
		"large":  "h-11 px-8 py-3",
	}),
	size.IsNot("small").Then("rounded-md"),
)

func Example() {
	fmt.Println(Button.Classes(Props{"medium"}))
	// Output: inline-flex items-center justify-center h-10 px-4 py-2 rounded-md
}
