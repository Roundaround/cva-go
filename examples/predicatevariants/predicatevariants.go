package predicatevariants

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Loading  bool
	Disabled bool
}

var Button = cva.New(
	cva.Base[Props]("button"),
	cva.PredicateVariant(
		func(p Props) bool { return p.Loading },
		"button-loading",
	),
	cva.PredicateVariant(
		func(p Props) bool { return p.Disabled || p.Loading },
		"button-disabled",
	),
)

func Example() {
	fmt.Println(Button.Classes(Props{
		Loading:  true,
		Disabled: false,
	}))
	// Output: button button-loading button-disabled
}
