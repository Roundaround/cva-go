# cva-go

<strong>C</strong>lass <strong>V</strong>ariance <strong>A</strong>uthority

_(but for Go!)_

cva-go has no direct affiliation with the original cva project, but is largely
inspired by it's design. Head on over to
[cva's documentation](https://cva.style/) if you need an overview of the project
and what problems it is trying to solve. cva-go is a rough re-implementation of
the high-level concept for the Go programming language, and pairs beautifully
with [templ](https://templ.guide/), [TailwindCSS](https://tailwindcss.com/), and
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

button := cva.NewCva(
  cva.WithStaticClasses[Props]("inline-flex items-center justify-center"),
  cva.WithVariant(
    func(p Props) string { return p.Size },
    map[string][]string{
      "small":  {"h-9 px-3"},
      "medium": {"h-10 px-4 py-2"},
      "large":  {"h-11 px-8 py-3"},
    },
  ),
)

fmt.Println(button.ClassName(Props{"small"}))
// Output: inline-flex items-center justify-center h-9 px-3
```

### Compound variants

The `WithCompoundVariant` helper allows you to apply classes based on a pair of
values. For more than two values, check out [predicate variants](#predicate-variants)

```go
type Props struct {
	Size  string
	Style string
}

button := cva.NewCva(
  cva.WithStaticClasses[Props]("inline-flex items-center justify-center"),
	cva.WithVariant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-8",
			"medium": "h-10",
			"large":  "h-12",
		},
	),
	cva.WithVariant(
		func(p Props) string { return p.Style },
		map[string]string{
			"icon":    "bg-gray-100 rounded-full aspect-square",
			"regular": "bg-gray-100 rounded-md",
			"link":    "text-blue-500",
		},
	),
	cva.WithCompoundVariant(
		func(p Props) (string, string) { return p.Size, p.Style },
		cva.WithCompound("small", "icon", "[&_svg]:size-4"),
		cva.WithCompound("medium", "icon", "[&_svg]:size-5"),
		cva.WithCompound("large", "icon", "[&_svg]:size-6"),
	),
)

fmt.Println(button.ClassName(Props{"small", "icon"}))
// Output: inline-flex items-center justify-center h-8 bg-gray-100 rounded-full
//   aspect-square [&_svg]:size-4
```

### Predicate variants

The `WithPredicateVariant` helper lets you specify a predicate function for each
class list you want to apply. When working with non-string-map values (i.e.
bools) or more than two variant properties, predicate variants can give you more
advanced control over when classes are applied.

```go
type Props struct {
	Loading  bool
	Disabled bool
}

button := cva.NewCva(
	cva.WithStaticClasses[Props]("button"),
	cva.WithPredicateVariant(
		func(p Props) bool { return p.Loading },
		"button-loading",
	),
	cva.WithPredicateVariant(
		func(p Props) bool { return p.Disabled || p.Loading },
		"button-disabled",
	),
)

fmt.Println(button.ClassName(Props{
  Loading:  true,
  Disabled: false,
}))
// Output: button button-loading button-disabled
```

### Customizing the class joining behavior

```go
dedupingContext := cva.NewCvaContext().WithDedupingClassJoiner()

button := cva.NewCva(
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

fmt.Println(button.ClassName(Props{"small"}))
// Output: inline-flex items-center justify-center rounded-md h-8
```

### Using [tailwind-merge-go](https://github.com/Oudwins/tailwind-merge-go) to merge classes

```go
import twmerge "github.com/Oudwins/tailwind-merge-go"

ctx := cva.NewCvaContext().WithClassJoiner(func(parts []string) string {
  return twmerge.Merge(parts...)
})

button := cva.NewCva(
  cva.WithContext[Props](ctx),
  cva.WithStaticClasses[Props]("inline-flex items-center justify-center px-2 py-1"),
  cva.WithVariant(
    func(p Props) string { return p.size },
    map[string][]string{
      "small":  {"h-9 px-3"},
      "medium": {"h-10 px-4 py-2"},
      "large":  {"h-11 px-8 py-3"},
    },
  ),
)

fmt.Println(button.ClassName(Props{"medium"}))
// Output: inline-flex items-center justify-center h-10 px-4 py-2
```

### Memoizing expensive property computations

If for some reason your getter functions are actually computing values (and said
computations are expensive), cva-go also exposes a simple memoization helper as
cva.Memoize.

__Most of the time you will not need this.__ Unless you run into performance
issues in computation time or memory allocation, you can skip right over this.
It is included as a convenience for more advanced use cases.

```go
type Props struct {
	Size int
}

type ExtendedProps struct {
	Size string
}

base := NewCva[Props](
	StaticClasses[Props]("button"),
	Variant(
		Memoize(func(p Props) string {
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

extended := NewCva[ExtendedProps](
	Inherit(
		base,
		Memoize(func(p ExtendedProps) Props {
			return Props{
				// Call to some expensive transformation function
				Size: transformSize(p.Size),
			}
		}),
	),
	Variant(
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

See the [examples directory](https://github.com/Roundaround/cva-go/tree/main/examples) for more usage examples

### TailwindCSS VSCode extension integration

Because the TailwindCSS VSCode extension uses regular expressions to target
which content should be evaluated as Tailwind classes, there's no single
solution to integrating cva-go with the extension.

That being said, with a simple local function, you can create something easy to
target with regular expressions:

```go
// Wrap your class list strings in a tw() call to get Tailwind intellisense
func tw(parts ...string) []string {
  return parts
}

var Button = cva.NewCva(
	cva.WithStaticClasses[Props](tw("inline-flex items-center justify-center")),
	cva.WithVariant(
		func(p Props) Size { return p.size },
		map[Size][]string{
			Small:  tw("h-9 px-3"),
			Medium: tw("h-10 px-4 py-2"),
			Large:  tw("h-11 px-8"),
		},
	),
)
```

Then in your VSCode's settings.json, add or extend the `tailwindCSS.experimental.classRegex` property with the following regular expression, which simply matches any string inside a tw() function call:

```json
{
  "tailwindCSS.experimental.classRegex": [
    ["tw\\((?:\"([^\"]*)\")|(?:`([^`]*)`)\\)"]
  ]
}
```

## Attributions, license, and copyright

Unless otherwise stated, cva-go is licensed under the MIT license. It is largely
inspired by the popular npm package
[class-variance-authority](https://cva.style/), but has no direct affiliation.

The Go programming language and the Go logo are trademarks of Google. The cva-go
project has no affiliation with the Go programming language or Google.

Copyright (c) 2025 Evan Steinkerchner
