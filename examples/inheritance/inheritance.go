package inheritance

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

// Base button with size and style variants
type BaseButtonProps struct {
	Size  string
	Style string
}

var BaseButton = cva.NewCva(
	cva.StaticClasses[BaseButtonProps]("inline-flex items-center justify-center"),
	cva.Variant(
		func(p BaseButtonProps) string { return p.Size },
		map[string]string{
			"small":  "h-8 px-3",
			"medium": "h-10 px-4",
			"large":  "h-12 px-6",
		},
	),
	cva.Variant(
		func(p BaseButtonProps) string { return p.Style },
		map[string]string{
			"primary":   "bg-blue-500 text-white",
			"secondary": "bg-gray-200 text-gray-800",
			"outline":   "border border-gray-300 text-gray-800",
		},
	),
)

// Loading button that inherits from base button and adds loading state
type LoadingButtonProps struct {
	Size    string
	Style   string
	Loading bool
}

var LoadingButton = cva.NewCva(
	cva.Inherit(
		BaseButton,
		func(p LoadingButtonProps) BaseButtonProps {
			return BaseButtonProps{
				Size:  p.Size,
				Style: p.Style,
			}
		},
	),
	cva.PredicateVariant(
		func(p LoadingButtonProps) bool { return p.Loading },
		"opacity-50 cursor-not-allowed",
	),
)

// Icon button that inherits from base button and adds icon-specific styles
type IconButtonProps struct {
	Size  string
	Style string
	Icon  string
}

var IconButton = cva.NewCva(
	cva.Inherit(
		BaseButton,
		func(p IconButtonProps) BaseButtonProps {
			return BaseButtonProps{
				Size:  p.Size,
				Style: p.Style,
			}
		},
	),
	cva.Variant(
		func(p IconButtonProps) string { return p.Icon },
		map[string]string{
			"plus":     "rounded-full [&_svg]:size-4",
			"settings": "rounded-full [&_svg]:size-5",
			"close":    "rounded-full [&_svg]:size-6",
		},
	),
)

func Example() {
	// Base button examples
	fmt.Println("Base Button Examples:")
	fmt.Println(BaseButton.ClassName(BaseButtonProps{Size: "medium", Style: "primary"}))
	fmt.Println(BaseButton.ClassName(BaseButtonProps{Size: "small", Style: "secondary"}))
	fmt.Println(BaseButton.ClassName(BaseButtonProps{Size: "large", Style: "outline"}))
	fmt.Println()

	// Loading button examples
	fmt.Println("Loading Button Examples:")
	fmt.Println(LoadingButton.ClassName(LoadingButtonProps{Size: "medium", Style: "primary", Loading: true}))
	fmt.Println(LoadingButton.ClassName(LoadingButtonProps{Size: "small", Style: "secondary", Loading: false}))
	fmt.Println()

	// Icon button examples
	fmt.Println("Icon Button Examples:")
	fmt.Println(IconButton.ClassName(IconButtonProps{Size: "medium", Style: "primary", Icon: "plus"}))
	fmt.Println(IconButton.ClassName(IconButtonProps{Size: "small", Style: "secondary", Icon: "settings"}))
	fmt.Println(IconButton.ClassName(IconButtonProps{Size: "large", Style: "outline", Icon: "close"}))
}
