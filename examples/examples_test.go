package main

import (
	"testing"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/Roundaround/cva-go"
)

func TestExamples(t *testing.T) {
	t.Run("example_1", func(t *testing.T) {
		type ButtonProps struct {
			size string
		}

		base := "inline-flex items-center justify-center rounded-md text-sm font-medium"
		small := "h-9 px-3"
		medium := "h-10 px-4 py-2"
		large := "h-11 px-8 py-3"

		button := cva.NewCva(
			base,
			cva.WithVariant(
				func(p ButtonProps) string { return p.size },
				map[string][]string{
					"small":  {small},
					"medium": {medium},
					"large":  {large},
				},
			),
		)

		tests := []struct {
			name  string
			props ButtonProps
			want  string
		}{
			{
				name:  "small",
				props: ButtonProps{"small"},
				want:  base + " " + small,
			},
			{
				name:  "medium",
				props: ButtonProps{"medium"},
				want:  base + " " + medium,
			},
			{
				name:  "large",
				props: ButtonProps{"large"},
				want:  base + " " + large,
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

	t.Run("example_2", func(t *testing.T) {
		type ButtonProps struct {
			size string
		}

		base := "inline-flex items-center justify-center py-1"
		small := "h-9 px-3"
		medium := "h-10 px-4 py-2"
		large := "h-11 px-8 py-3"

		button := cva.NewCva(
			base,
			cva.WithVariant(
				func(p ButtonProps) string { return p.size },
				map[string][]string{
					"small":  {small},
					"medium": {medium},
					"large":  {large},
				},
			),
			cva.WithClassJoiner[ButtonProps](func(parts []string) string {
				return twmerge.Merge(parts...)
			}),
		)

		tests := []struct {
			name  string
			props ButtonProps
			want  string
		}{
			{
				name:  "small",
				props: ButtonProps{"small"},
				want:  twmerge.Merge(base, small),
			},
			{
				name:  "medium",
				props: ButtonProps{"medium"},
				want:  twmerge.Merge(base, medium),
			},
			{
				name:  "large",
				props: ButtonProps{"large"},
				want:  twmerge.Merge(base, large),
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
