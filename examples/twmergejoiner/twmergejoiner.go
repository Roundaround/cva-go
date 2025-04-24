package twmergejoiner

import (
	"fmt"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size string
}

var ctx = cva.NewCvaContext().WithClassJoiner(func(parts []string) string {
	return twmerge.Merge(parts...)
})

var Button = cva.NewCva(
	cva.WithContext[Props](ctx),
	cva.WithStaticClasses[Props]("inline-flex items-center justify-center py-1"),
	cva.WithVariant(
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
	// Output: inline-flex items-center justify-center py-1 h-9 px-3

	fmt.Println(Button.ClassName(Props{"medium"}))
	// Output: inline-flex items-center justify-center h-10 px-4 py-2
}
