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
  "inline-flex items-center justify-center",
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
// "inline-flex items-center justify-center h-9 px-3"
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
  "inline-flex items-center justify-center",
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
// inline-flex items-center justify-center h-8 bg-gray-100 rounded-full
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

fmt.Println(button.ClassName(Props{
  Loading:  true,
  Disabled: false,
}))
// button button-loading
```

### Using [tailwind-merge-go](https://github.com/Oudwins/tailwind-merge-go) to merge classes

```go
import (
  twmerge "github.com/Oudwins/tailwind-merge-go"
  "github.com/Roundaround/cva-go"
)

type Props struct {
  size string
}

button := cva.NewCva(
  "inline-flex items-center justify-center py-1",
  cva.WithVariant(
    func(p Props) string { return p.size },
    map[string][]string{
      "small":  {"h-9 px-3"},
      "medium": {"h-10 px-4 py-2"},
      "large":  {"h-11 px-8 py-3"},
    },
  ),
  cva.WithClassJoiner[Props](func(parts []string) string {
    return twmerge.Merge(parts...)
  }),
)

button.ClassName(Props{"small"})
// "inline-flex items-center justify-center py-1 h-9 px-3"
//                           py-1 from base ^

button.ClassName(Props{"medium"})
// "inline-flex items-center justify-center h-9 px-3 py-2"
//                  py-1 replaced with medium's py-2 ^
```

### Change class joiner for all Cva instances

```go
cva.FallbackClassJoiner = func(parts []string) string {
  return twmerge.Merge(parts...)
}

// . . .

button := cva.NewCva(
  "inline-flex items-center justify-center py-1",
  cva.WithVariant(
    func(p Props) string { return p.size },
    map[string][]string{
      "small":  {"h-9 px-3"},
      "medium": {"h-10 px-4 py-2"},
      "large":  {"h-11 px-8 py-3"},
    },
  ),
  // cva.WithClassJoiner(...) no longer needed; inherited from fallback
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
	tw("inline-flex items-center justify-center"),
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
