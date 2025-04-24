package predicatevariants

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Loading  bool
	Disabled bool
}

var Button = cva.NewCva(
	cva.StaticClasses[Props]("button"),
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
	fmt.Println(Button.ClassName(Props{
		Loading:  true,
		Disabled: false,
	}))
	// Output: button button-loading button-disabled
}
