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
	"button",
	cva.WithPredicateVariant(
		func(p Props) bool { return p.Loading },
		"button-loading",
	),
	cva.WithPredicateVariant(
		func(p Props) bool { return p.Disabled },
		"button-disabled",
	),
)

func Example() {
	fmt.Println(Button.ClassName(Props{
		Loading:  true,
		Disabled: false,
	}))
	// button button-loading
}
