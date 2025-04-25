# cva-go

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
  cva.Variant(
    func(p Props) string { return p.Size },
    map[string][]string{
      "small":  {"h-9 px-3"},
      "medium": {"h-10 px-4 py-2"},
      "large":  {"h-11 px-8 py-3"},
    },
  ),
)

fmt.Println(button.Classes(Props{"small"}))
// Output: inline-flex items-center justify-center h-9 px-3
```

### Compound variants

The `CompoundVariant` helper allows you to apply classes based on a pair of values. For more than
two values, check out [predicate variants](#predicate-variants)

```go
type Props struct {
	Size  string
	Style string
}

button := cva.New(
  cva.Base[Props]("inline-flex items-center justify-center"),
	cva.Variant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-8",
			"medium": "h-10",
			"large":  "h-12",
		},
	),
	cva.Variant(
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

The `PredicateVariant` helper lets you specify a predicate function for each class list you want to
apply. When working with non-string-map values (i.e. bools) or more than two variant properties,
predicate variants can give you more advanced control over when classes are applied.

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
	cva.Variant(
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
  cva.Variant(
    func(p Props) string { return p.size },
    map[string][]string{
      "small":  {"h-9 px-3"},
      "medium": {"h-10 px-4 py-2"},
      "large":  {"h-11 px-8 py-3"},
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
  cva.Variant(
    func(p Props) string { return p.size },
    map[string][]string{
      "small":  {"h-9 px-3"},
      "medium": {"h-10 px-4 py-2"},
      "large":  {"h-11 px-8 py-3"},
    },
  ),
).Classes // Storing a reference to the Classes func

fmt.Println(button(Props{"medium"}))
// Output: inline-flex items-center justify-center h-10 px-4 py-2
```

### Composition through inheritance

You can compose components through inheritance by using the `Inherit` option. This will copy all the
class lists and variants from the base Cva, allowing you to build upon it.

```go
type ButtonProps struct {
	Size  string
	Style string
}

button := cva.New(
	cva.Base[ButtonProps]("inline-flex items-center justify-center"),
	cva.Variant(
		func(p ButtonProps) string { return p.Size },
		map[string]string{
			"small":  "h-8 px-3",
			"medium": "h-10 px-4",
			"large":  "h-12 px-6",
		},
	),
	cva.Variant(
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
	cva.Variant(
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
	cva.Variant(
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
