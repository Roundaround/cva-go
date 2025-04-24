package twmergejoiner

import (
	"fmt"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/Roundaround/cva-go"
)

type Props struct {
	Size string
}

var Button = cva.NewCva(
	"inline-flex items-center justify-center py-1",
	cva.WithVariant(
		func(p Props) string { return p.Size },
		map[string]string{
			"small":  "h-9 px-3",
			"medium": "h-10 px-4 py-2",
			"large":  "h-11 px-8 py-3",
		},
	),
	cva.WithClassJoiner[Props](func(parts []string) string {
		return twmerge.Merge(parts...)
	}),
)

func Example() {
	fmt.Println(Button.ClassName(Props{"small"}))
	// inline-flex items-center justify-center py-1 h-9 px-3
	//                          py-1 from base ^

	fmt.Println(Button.ClassName(Props{"medium"}))
	// inline-flex items-center justify-center h-10 px-4 py-2
	//                  py-1 replaced with medium's py-2 ^
}
