package cva

import (
	"testing"
)

func TestCva(t *testing.T) {
	t.Run("single_variant", func(t *testing.T) {
		type Props struct {
			Size string
		}

		button := NewCva(
			"button",
			WithVariant(
				func(p Props) string { return p.Size },
				map[string][]string{
					"small":  {"button-small"},
					"medium": {"button-medium"},
					"large":  {"button-large"},
				},
			),
		)

		tests := []struct {
			name  string
			props Props
			want  string
		}{
			{
				name:  "small",
				props: Props{Size: "small"},
				want:  "button button-small",
			},
			{
				name:  "medium",
				props: Props{Size: "medium"},
				want:  "button button-medium",
			},
			{
				name:  "large",
				props: Props{Size: "large"},
				want:  "button button-large",
			},
			{
				name:  "unknown",
				props: Props{Size: "not-a-size"},
				want:  "button",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := button.ClassName(test.props)
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})
}
