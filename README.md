# cva-go

[![Go Reference](https://pkg.go.dev/badge/github.com/Roundaround/cva-go.svg)](https://pkg.go.dev/github.com/Roundaround/cva-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/Roundaround/cva-go)](https://goreportcard.com/report/github.com/Roundaround/cva-go)
[![Coverage Status](https://coveralls.io/repos/github/Roundaround/cva-go/badge.svg?branch=main)](https://coveralls.io/github/Roundaround/cva-go?branch=main)

<strong>C</strong>lass <strong>V</strong>ariance <strong>A</strong>uthority

_(but for Go!)_

cva-go has no direct affiliation with the original cva project, but is largely inspired by it's
design. Head on over to [cva's documentation](https://cva.style/) if you need an overview of the
project and what problems it is trying to solve. cva-go is a rough re-implementation of the
high-level concept for the Go programming language, and pairs beautifully with
[templ](https://templ.guide/), [TailwindCSS](https://tailwindcss.com/), and
[tailwind-merge-go](https://github.com/Oudwins/tailwind-merge-go).

## Getting started

```sh
go get github.com/Roundaround/cva-go
```

## Example usage

### Basic button with a "size" variant property

```go
type Props struct {
	Size  string
}

button := cva.New(
	cva.Base[Props]("inline-flex items-center justify-center"),
	cva.MapVariant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-9 px-3",
			"medium": "h-10 px-4 py-2 rounded-md",
			"large":  "h-11 px-8 py-3 rounded-md",
		},
	),
)

fmt.Println(button.Classes(Props{"medium"}))
// Output: inline-flex items-center justify-center h-10 px-4 py-2 rounded-md
```

### Explicitly defined variants

If you prefer, you can define your variants up front so they can be used in multiple options:

```go
type Props struct {
	Size  string
}

size := cva.NewVariant(func(p Props) string { return p.Size })

button := cva.New(
	cva.Base[Props]("inline-flex items-center justify-center"),
	size.Map(map[string]string{
		"small":  "h-9 px-3",
		"medium": "h-10 px-4 py-2",
		"large":  "h-11 px-8 py-3",
	}),
	size.IsNot("small").Then("rounded-md"),
)

fmt.Println(button.Classes(Props{"medium"}))
// Output: inline-flex items-center justify-center h-10 px-4 py-2 rounded-md
```

If your data is not fully sanitized, you can also take advantage of the `WithValues` and
`WithDefault` methods. Using `WithValues` simply maps any value to the zero value of the variant's
type when it doesn't explicitly match one of the supplied values, and `WithDefault` maps the zero
value to the specified value.

```go
type Props struct {
	Size  string
}

size := cva.NewVariant(func(p Props) string { return p.Size }).
	WithValues("small", "medium", "large").
	WithDefault("medium")

button := cva.New(
	cva.Base[Props]("inline-flex items-center justify-center"),
	size.Map(map[string]string{
		"small":  "h-9 px-3",
		"medium": "h-10 px-4 py-2",
		"large":  "h-11 px-8 py-3",
	}),
	size.IsNot("small").Then("rounded-md"),
)

fmt.Println(button.Classes(Props{"custom"}))
// Invalid value, so this is converted to "", then to the default ("medium")
// Output: inline-flex items-center justify-center h-10 px-4 py-2 rounded-md
```

### Compound variants

The `CompoundVariant` helper allows you to apply classes based on a pair of values. When defining
your compound values, the `NewCompound` helper will provide strong typing for its params.
For more than two values, check out [predicate variants](#predicate-variants).

```go
type Props struct {
	Size  string
	Style string
}

button := cva.New(
	cva.Base[Props]("inline-flex items-center justify-center"),
	cva.MapVariant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-8",
			"medium": "h-10",
			"large":  "h-12",
		},
	),
	cva.MapVariant(
		func(p Props) string { return p.Style },
		map[string]string{
			"icon":    "bg-gray-100 rounded-full aspect-square",
			"regular": "bg-gray-100 rounded-md",
			"link":    "text-blue-500",
		},
	),
	cva.CompoundVariant(
		func(p Props) (string, string) { return p.Size, p.Style },
		cva.NewCompound("small", "icon", "[&_svg]:size-4"),
		cva.NewCompound("medium", "icon", "[&_svg]:size-5"),
		cva.NewCompound("large", "icon", "[&_svg]:size-6"),
	),
)

fmt.Println(button.Classes(Props{"small", "icon"}))
// Output: inline-flex items-center justify-center h-8 bg-gray-100 rounded-full
//   aspect-square [&_svg]:size-4
```

### Predicate variants

The `PredicateVariant` helper lets you specify a predicate function for each class list you want
to apply. When working with non-string-map values (i.e. bools) or more than two variant
properties, predicate variants can give you more advanced control over when classes are applied.

```go
type Props struct {
	Loading  bool
	Disabled bool
}

button := cva.NewCva(
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

fmt.Println(button.Classes(Props{
	Loading:  true,
	Disabled: false,
}))
// Output: button button-loading button-disabled
```

### Deduplicating classes

Out of the box, cva-go will not deduplicate classes for you directly as part of your component
definitions. It does however expose a helper function for deduping classes as a post-processing
step, should you need it.

```go
button := cva.New(
	cva.Base[Props]("inline-flex items-center justify-center rounded-md"),
	cva.MapVariant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-8 rounded-md", // rounded-md is repeated from the base
			"medium": "h-10 rounded-md",
			"large":  "h-12 rounded-md",
		},
	),
)

fmt.Println(button.Classes(Props{"small"}))
// Output: inline-flex items-center justify-center rounded-md h-8 rounded-md

fmt.Println(cva.DedupeClasses(button.Classes(Props{"small"})))
// Output: inline-flex items-center justify-center rounded-md h-8
```

### Using [tailwind-merge-go](https://github.com/Oudwins/tailwind-merge-go) to merge classes

```go
import twmerge "github.com/Oudwins/tailwind-merge-go"

button := cva.New(
	cva.Base[Props]("inline-flex items-center justify-center px-2 py-1"),
	cva.MapVariant(
		func(p Props) string { return p.size },
		map[string]string{
			"small":  "h-9 px-3",
			"medium": "h-10 px-4 py-2",
			"large":  "h-11 px-8 py-3",
		},
	),
)

fmt.Println(twmerge.Merge(button.Classes(Props{"medium"})))
// Output: inline-flex items-center justify-center h-10 px-4 py-2
```

### Simplifying your component usage

If you find that obtaining the list of classes at the call site is too noisy or verbose (which can
easily be the case when integrating with templating systems like templ), you may wish to store a
reference to the `Classes` function rather than the component struct itself.

```go
button := cva.New(
	cva.Base[Props]("inline-flex items-center justify-center"),
	cva.MapVariant(
		func(p Props) string { return p.size },
		map[string]string{
			"small":  "h-9 px-3",
			"medium": "h-10 px-4 py-2",
			"large":  "h-11 px-8 py-3",
		},
	),
).Classes // Storing a reference to the Classes func

fmt.Println(button(Props{"medium"}))
// Output: inline-flex items-center justify-center h-10 px-4 py-2
```

Or going one step further, if you're using the `DedupeClasses`, `twmerge.Merge`, or some other
post-processing function, you can wrap it all together all at once:

```go
buttonCva := cva.New(
	cva.Base[Props]("inline-flex items-center justify-center rounded-md"),
	cva.MapVariant(
		func(p Props) string { return p.size },
		map[string]string{
			"small":  "h-9 px-3 rounded-md",
			"medium": "h-10 px-4 py-2 rounded-md",
			"large":  "h-11 px-8 py-3 rounded-md",
		},
	),
)

button := func(p Props) {
	return cva.DedupeClasses(buttonCva.Classes(p))
}

fmt.Println(button(Props{"medium"}))
// Output: inline-flex items-center justify-center rounded-md h-10 px-4 py-2
```

### Complex variant definitions

Sometimes simple maps are not quite enough, and you might find yourself needing more complex
matching logic for your classes. To help facilitate this, cva-go ships with a `Matcher` struct
that can be used to build out details boolean expressions.

```go
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

size := cva.NewVariant(func(p Props) Size { return p.Size })
theme := cva.NewVariant(func(p Props) Theme { return p.Theme })
elem := cva.NewVariant(func(p Props) Element { return p.Element })

button := cva.New(
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

fmt.Println(button.Classes(Props{SizeLarge, ThemeDanger, ElementLink}))

// Of course, you might also want to use tailwind-merge-go for these more complex cases:
fmt.Println(twmerge.Merge(button.Classes(Props{SizeLarge, ThemeDanger, ElementLink})))
```

### Composition through inheritance

You can compose components through inheritance by using the `Inherit` option. This will copy all
the class lists and variants from the base Cva, allowing you to build upon it.

```go
type ButtonProps struct {
	Size  string
	Style string
}

button := cva.New(
	cva.Base[ButtonProps]("inline-flex items-center justify-center"),
	cva.MapVariant(
		func(p ButtonProps) string { return p.Size },
		map[string]string{
			"small":  "h-8 px-3",
			"medium": "h-10 px-4",
			"large":  "h-12 px-6",
		},
	),
	cva.MapVariant(
		func(p ButtonProps) string { return p.Style },
		map[string]string{
			"primary":   "bg-blue-500 text-white",
			"secondary": "bg-gray-200 text-gray-800",
			"outline":   "border border-gray-300 text-gray-800",
		},
	),
)

type LoadingButtonProps struct {
	ButtonProps
	Loading bool
}

loadingButton := cva.New(
	cva.Inherit(
		button,
		func(p LoadingButtonProps) ButtonProps { return p.ButtonProps },
	),
	cva.PredicateVariant(
		func(p LoadingButtonProps) bool { return p.Loading },
		"opacity-50 cursor-not-allowed",
	),
)

fmt.Println(loadingButton.Classes({
	ButtonProps{Size: "medium", Style: "primary"},
	Loading: true,
}))
// Output: inline-flex items-center justify-center h-10 px-4 bg-blue-500
//   text-white opacity-50 cursor-not-allowed
```

### Memoizing expensive property computations

If for some reason your getter functions are actually computing values (and said computations are
expensive), cva-go also exposes a simple memoization helper as cva.Memoize.

__Most of the time you will not need this.__ Unless you run into performance issues in computation
time or memory allocation, you can skip right over this. It is included as a convenience for more
advanced use cases.

```go
type Props struct {
	Size int
}

type ExtendedProps struct {
	Size string
}

base := cva.New[Props](
	cva.Base[Props]("button"),
	cva.MapVariant(
		cva.Memoize(func(p Props) string {
			// Call to some expensive transformation function
			return sizeString(p.Size)
		}),
		map[string]string{
			"small":  "button-small",
			"medium": "button-medium",
			"large":  "button-large",
		},
	),
)

extended := cva.New[ExtendedProps](
	cva.Inherit(
		base,
		cva.Memoize(func(p ExtendedProps) Props {
			return Props{
				// Call to some expensive transformation function
				Size: transformSize(p.Size),
			}
		}),
	),
	cva.MapVariant(
		func(p ExtendedProps) string { return p.Size },
		map[string]string{
			"xs":  "button-small",
			"md": "button-medium",
			"lg":  "button-large",
		},
	),
)
```

### Additional examples

See the [examples directory](https://github.com/Roundaround/cva-go/tree/main/examples) for more
usage examples

### TailwindCSS VSCode extension integration

If you're working with VSCode and TailwindCSS, you can get incredible intellisense/language server
support for your Tailwind classes, including autocomplete, tooltips, and color indicators using the
[Tailwind CSS IntelliSense](https://marketplace.visualstudio.com/items?itemName=bradlc.vscode-tailwindcss)
extension. To integrate it with your Go project, add/extend the following entries to your VSCode
settings.json:

```json
{
	"tailwindCSS.includeLanguages": {
		"templ": "html",
		"go": "javascript"
	},
	"tailwindCSS.classFunctions": [
		"cva.New"
	],
}
```

## Attributions, license, and copyright

Unless otherwise stated, cva-go is licensed under the MIT license. It is largely inspired by the
popular npm package [class-variance-authority](https://cva.style/), but has no direct affiliation.

The Go programming language and the Go logo are trademarks of Google. The cva-go project has no
affiliation with the Go programming language or Google.

Copyright (c) 2025 Evan Steinkerchner
