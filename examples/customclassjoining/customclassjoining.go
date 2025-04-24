package customclassjoining

import (
	"fmt"
	"strings"

	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size string
}

var Button = cva.NewCva(
	"inline-flex items-center justify-center rounded-md",
	cva.WithVariant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-8 rounded-md",
			"medium": "h-10 rounded-md",
			"large":  "h-12 rounded-md",
		},
	),
	cva.WithClassJoiner[Props](func(parts []string) string {
		split := make([]string, 0, len(parts))
		for _, part := range parts {
			for s := range strings.SplitSeq(part, " ") {
				split = append(split, s)
			}
		}

		deduped := make([]string, 0)
		unique := make(map[string]struct{})
		for _, s := range split {
			if _, ok := unique[s]; !ok {
				unique[s] = struct{}{}
				deduped = append(deduped, s)
			}
		}

		return strings.Join(deduped, " ")
	}),
)

func Example() {
	fmt.Println(Button.ClassName(Props{"small"}))
	// inline-flex items-center justify-center rounded-md h-8
	//                  rounded-md was deduped ^
}
