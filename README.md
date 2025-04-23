## Example usage

__Basic button with a "size" variant property__

```go
type Props struct {
  size string
}

button := cva.NewCva(
  "inline-flex items-center justify-center",
  cva.WithVariant(
    func(p Props) string { return p.size },
    map[string][]string{
      "small":  {"h-9 px-3"},
      "medium": {"h-10 px-4 py-2"},
      "large":  {"h-11 px-8 py-3"},
    },
  ),
)

button.ClassName(Props{"small"})
// ^ "inline-flex items-center justify-center ... h-9 px-3"
```

__Using [github.com/Oudwins/tailwind-merge-go](https://github.com/Oudwins/tailwind-merge-go) as the class joiner__

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
// ^ "inline-flex items-center justify-center py-1 h-9 px-3"

button.ClassName(Props{"medium"})
// ^ "inline-flex items-center justify-center h-9 px-3 py-2"
```
