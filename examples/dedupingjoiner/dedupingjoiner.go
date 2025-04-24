package dedupingjoiner

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size string
}

var dedupingContext = cva.NewCvaContext().WithDedupingClassJoiner()

var Button = cva.NewCva(
	cva.WithContext[Props](dedupingContext),
	cva.WithStaticClasses[Props]("inline-flex items-center justify-center rounded-md"),
	cva.WithVariant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-8 rounded-md",
			"medium": "h-10 rounded-md",
			"large":  "h-12 rounded-md",
		},
	),
)

func Example() {
	fmt.Println(Button.ClassName(Props{"small"}))
	// Output: inline-flex items-center justify-center rounded-md h-8
}
