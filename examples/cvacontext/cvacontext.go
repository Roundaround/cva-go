package cvacontext

import (
	"fmt"
	"strings"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size string
}

// Create a context with a custom class joiner
var ctx = cva.NewCvaContext().WithClassJoiner(func(parts []string) string {
	return strings.Join(strings.Split(strings.Join(parts, " "), " "), "::")
})

// Create a button with the context
var Button = cva.NewCva(
	cva.Context[Props](ctx),
	cva.StaticClasses[Props]("inline-flex items-center justify-center"),
	cva.Variant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-9 px-3",
			"medium": "h-10 px-4 py-2",
			"large":  "h-11 px-8 py-3",
		},
	),
)

// Create another button with the same context
var Button2 = cva.NewCva(
	cva.Context[Props](ctx),
	cva.StaticClasses[Props]("button-base"),
	cva.Variant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "button-small",
			"medium": "button-medium",
			"large":  "button-large",
		},
	),
)

// Create a button with a different context
var altCtx = cva.NewCvaContext().WithClassJoiner(func(parts []string) string {
	return strings.Join(parts, " ")
})

var Button3 = cva.NewCva(
	cva.Context[Props](altCtx),
	cva.StaticClasses[Props]("button-base"),
	cva.Variant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "button-small",
			"medium": "button-medium",
			"large":  "button-large",
		},
	),
)

func Example() {
	fmt.Println(Button.ClassName(Props{"small"}))
	// Output: inline-flex::items-center::justify-center::h-9::px-3

	fmt.Println(Button2.ClassName(Props{"medium"}))
	// Output: button-base::button-medium

	fmt.Println(Button3.ClassName(Props{"medium"}))
	// Output: button-base button-medium
}
