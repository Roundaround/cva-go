package cva

import (
	"strings"
	"testing"
)

func TestCva(t *testing.T) {
	t.Run("mapped_variant", func(t *testing.T) {
		t.Run("string_value", func(t *testing.T) {
			type Props struct {
				Size string
			}

			button := NewCva(
				WithStaticClasses[Props]("button"),
				WithVariant(
					func(p Props) string { return p.Size },
					map[string]string{
						"small":  "button-small",
						"medium": "button-medium",
						"large":  "button-large",
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

		t.Run("slice_value", func(t *testing.T) {
			type Props struct {
				Size string
			}

			button := NewCva(
				WithStaticClasses[Props]("button"),
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
	})

	t.Run("compound_variant", func(t *testing.T) {
		type Props struct {
			Size  string
			Color string
		}

		button := NewCva(
			WithStaticClasses[Props]("button"),
			WithCompoundVariant(
				func(p Props) (string, string) { return p.Size, p.Color },
				WithCompound("small", "red", "button-small-red"),
				WithCompound("large", "blue", "button-large-blue"),
			),
		)

		tests := []struct {
			name  string
			props Props
			want  string
		}{
			{
				name:  "small-red",
				props: Props{Size: "small", Color: "red"},
				want:  "button button-small-red",
			},
			{
				name:  "large-blue",
				props: Props{Size: "large", Color: "blue"},
				want:  "button button-large-blue",
			},
			{
				name:  "unknown-combination",
				props: Props{Size: "medium", Color: "green"},
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

	t.Run("predicate_variant", func(t *testing.T) {
		type Props struct {
			IsDisabled bool
			IsLoading  bool
		}

		button := NewCva(
			WithStaticClasses[Props]("button"),
			WithPredicateVariant(
				func(p Props) bool { return p.IsDisabled },
				"button-disabled",
			),
			WithPredicateVariant(
				func(p Props) bool { return p.IsLoading },
				"button-loading",
			),
		)

		tests := []struct {
			name  string
			props Props
			want  string
		}{
			{
				name:  "disabled",
				props: Props{IsDisabled: true},
				want:  "button button-disabled",
			},
			{
				name:  "loading",
				props: Props{IsLoading: true},
				want:  "button button-loading",
			},
			{
				name:  "disabled-and-loading",
				props: Props{IsDisabled: true, IsLoading: true},
				want:  "button button-disabled button-loading",
			},
			{
				name:  "enabled",
				props: Props{IsDisabled: false, IsLoading: false},
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

	t.Run("static_classes", func(t *testing.T) {
		t.Run("single_class", func(t *testing.T) {
			type Props struct{}
			button := NewCva(WithStaticClasses[Props]("button"))
			got := button.ClassName(Props{})
			want := "button"
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})

		t.Run("multiple_classes", func(t *testing.T) {
			type Props struct{}
			button := NewCva(WithStaticClasses[Props]("button", "base"))
			got := button.ClassName(Props{})
			want := "button base"
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})
	})

	t.Run("with_props_classes", func(t *testing.T) {
		type Props struct {
			CustomClasses []string
		}

		button := NewCva(
			WithStaticClasses[Props]("button"),
			WithPropsClasses(func(p Props) []string {
				return p.CustomClasses
			}),
		)

		tests := []struct {
			name  string
			props Props
			want  string
		}{
			{
				name:  "with-custom-classes",
				props: Props{CustomClasses: []string{"custom-1", "custom-2"}},
				want:  "button custom-1 custom-2",
			},
			{
				name:  "no-custom-classes",
				props: Props{CustomClasses: nil},
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

	t.Run("repeated_classes", func(t *testing.T) {
		type Props struct {
			Size string
		}

		button := NewCva(
			WithStaticClasses[Props]("button"),
			WithVariant(
				func(p Props) string { return p.Size },
				map[string]string{"small": "button", "medium": "button"},
			),
		)

		got := button.ClassName(Props{Size: "small"})
		want := "button button"
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("custom_class_joiner", func(t *testing.T) {
		ctx := NewCvaContext().WithClassJoiner(func(parts []string) string {
			return strings.Join(parts, "::")
		})

		type ButtonProps struct {
			size string
		}

		button := NewCva(
			WithContext[ButtonProps](ctx),
			WithStaticClasses[ButtonProps]("button"),
			WithVariant(
				func(p ButtonProps) string { return p.size },
				map[string]string{
					"small":  "small",
					"medium": "medium",
					"large":  "large",
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
				want:  "button::small",
			},
			{
				name:  "medium",
				props: ButtonProps{"medium"},
				want:  "button::medium",
			},
			{
				name:  "large",
				props: ButtonProps{"large"},
				want:  "button::large",
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
