package matchers

import (
	"fmt"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/Roundaround/cva-go"
)

type Size int
type Theme int
type Element int

const (
	SizeSmall Size = iota
	SizeLarge
)

const (
	ThemeDanger Theme = iota
	ThemePrimary
)

const (
	ElementButton Element = iota
	ElementLink
	ElementIcon
)

type Props struct {
	Size    Size
	Theme   Theme
	Element Element
}

var size = cva.NewVariant(func(p Props) Size { return p.Size })
var theme = cva.NewVariant(func(p Props) Theme { return p.Theme })
var elem = cva.NewVariant(func(p Props) Element { return p.Element })

var Button = cva.New(
	// Apply classes no matter what
	cva.Base[Props]("px-4 py-1"),

	// Apply classes based on 1:1 mappings to variant values
	size.Map(map[Size]string{
		SizeSmall: "px-2",
		SizeLarge: "px-6 py-2",
	}),
	theme.Map(map[Theme]string{
		ThemeDanger:  "bg-red-500 text-white",
		ThemePrimary: "bg-blue-500 text-white",
	}),

	// Construct matchers using Variant.* (Is, In, IsNot, NotIn, or Test)
	elem.Is(ElementIcon).Then("rounded-full"),
	elem.IsNot(ElementIcon).Then("rounded-md"),
	// elem.In(ElementButton, ElementLink).Then("rounded-md"),

	// Compose complex matchers using Matcher.* (Or, And, Not)
	size.Is(SizeLarge).And(theme.Is(ThemeDanger)).Then("font-bold"),
	elem.Is(ElementIcon).And(size.IsNot(SizeSmall)).Then("[&_svg]:size-5"),

	// Use Any, All, and When if method chaining is not your cup of tea
	cva.When(elem.Is(ElementLink), "text-blue-500 hover:text-blue-600 bg-transparent"),
)

func Example() {
	fmt.Println(twmerge.Merge(Button.Classes(Props{SizeLarge, ThemeDanger, ElementLink})))
	// Output: px-6 py-2 bg-red-500 text-white rounded-full text-blue-500 hover:text-blue-600 bg-transparent
}
