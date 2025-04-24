package context

import (
	"fmt"
	"strings"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size string
}

func Example() {
	// Create a context with a custom class joiner
	ctx := cva.NewCvaContext().WithClassJoiner(func(parts []string) string {
		return strings.Join(parts, "::")
	})

	// Create a button with the context
	button := cva.NewCva(
		cva.WithContext[Props](ctx),
		cva.WithStaticClasses[Props]("inline-flex items-center justify-center"),
		cva.WithVariant(
			func(p Props) string { return p.Size },
			map[string]string{
				"small":  "h-9 px-3",
				"medium": "h-10 px-4 py-2",
				"large":  "h-11 px-8 py-3",
			},
		),
	)

	// Create another button with the same context
	button2 := cva.NewCva(
		cva.WithContext[Props](ctx),
		cva.WithStaticClasses[Props]("button-base"),
		cva.WithVariant(
			func(p Props) string { return p.Size },
			map[string]string{
				"small":  "button-small",
				"medium": "button-medium",
				"large":  "button-large",
			},
		),
	)

	fmt.Println(button.ClassName(Props{"small"}))
	// Output: inline-flex::items-center::justify-center::h-9::px-3

	fmt.Println(button2.ClassName(Props{"medium"}))
	// Output: button-base::button-medium

	// Create a button with a different context
	altCtx := cva.NewCvaContext().WithClassJoiner(func(parts []string) string {
		return strings.Join(parts, " ")
	})

	button3 := cva.NewCva(
		cva.WithContext[Props](altCtx),
		cva.WithStaticClasses[Props]("button-base"),
		cva.WithVariant(
			func(p Props) string { return p.Size },
			map[string]string{
				"small":  "button-small",
				"medium": "button-medium",
				"large":  "button-large",
			},
		),
	)

	fmt.Println(button3.ClassName(Props{"medium"}))
	// Output: button-base button-medium
}
