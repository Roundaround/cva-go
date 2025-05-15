package cva

import (
	"strconv"
	"testing"
)

func TestCva(t *testing.T) {
	t.Run("Classes", func(t *testing.T) {
		t.Run("string_value", func(t *testing.T) {
			type Props struct {
				CustomClasses string
			}

			button := New(
				Base[Props]("button"),
				Classes(func(p Props) string {
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
					props: Props{CustomClasses: "custom-1 custom-2"},
					want:  "button custom-1 custom-2",
				},
				{
					name:  "no-custom-classes",
					props: Props{CustomClasses: ""},
					want:  "button",
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					got := button.Classes(test.props)
					if got != test.want {
						t.Errorf("got %s, want %s", got, test.want)
					}
				})
			}
		})

		t.Run("slice_value", func(t *testing.T) {
			type Props struct {
				CustomClasses []string
			}

			button := New(
				Base[Props]("button"),
				Classes(func(p Props) []string {
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
					got := button.Classes(test.props)
					if got != test.want {
						t.Errorf("got %s, want %s", got, test.want)
					}
				})
			}
		})
	})

	t.Run("MapVariant", func(t *testing.T) {
		t.Run("string_value", func(t *testing.T) {
			type Props struct {
				Size string
			}

			button := New(
				Base[Props]("button"),
				MapVariant(
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
					got := button.Classes(test.props)
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

			button := New(
				Base[Props]("button"),
				MapVariant(
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
					got := button.Classes(test.props)
					if got != test.want {
						t.Errorf("got %s, want %s", got, test.want)
					}
				})
			}
		})

		t.Run("unknown_values", func(t *testing.T) {
			type Props struct {
				Size string
			}

			got := New(MapVariant(
				func(p Props) string { return p.Size },
				map[string]string{
					"small":  "small",
					"medium": "medium",
					"large":  "large",
				},
			)).Classes(Props{Size: "unknown"})
			want := ""
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})

		t.Run("zero_values", func(t *testing.T) {
			// Because all the values come from getter functions, we actually don't
			// want any built-in zero value handling. Confirm that a zero value is
			// treated just like any other.

			type Props struct {
				Size string
			}

			t.Run("without_a_mapping_for_zero", func(t *testing.T) {
				got := New(MapVariant(
					func(p Props) string { return p.Size },
					map[string]string{
						"small":  "small",
						"medium": "medium",
						"large":  "large",
					},
				)).Classes(Props{})
				want := ""
				if got != want {
					t.Errorf("got %s, want %s", got, want)
				}
			})

			t.Run("with_a_mapping_for_zero", func(t *testing.T) {
				got := New(MapVariant(
					func(p Props) string { return p.Size },
					map[string]string{
						"small":  "small",
						"medium": "medium",
						"large":  "large",
						"":       "zero",
					},
				)).Classes(Props{})
				want := "zero"
				if got != want {
					t.Errorf("got %s, want %s", got, want)
				}
			})
		})

		t.Run("empty_map", func(t *testing.T) {
			type Props struct {
				Size string
			}

			got := New(MapVariant(
				func(p Props) string { return p.Size },
				map[string]string{},
			)).Classes(Props{})
			want := ""
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})

		t.Run("nil_map", func(t *testing.T) {
			type Props struct {
				Size string
			}

			var m map[string]string
			got := New(MapVariant(
				func(p Props) string { return p.Size },
				m,
			)).Classes(Props{})
			want := ""
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})

		t.Run("map_mutation", func(t *testing.T) {
			t.Run("for_string_maps", func(t *testing.T) {
				type Props struct {
					Size string
				}

				// Create a map and pass it to Variant
				classesMap := map[string]string{
					"small":  "button-small",
					"medium": "button-medium",
				}

				button := New(
					Base[Props]("button"),
					MapVariant(
						func(p Props) string { return p.Size },
						classesMap,
					),
				)

				// Mutate the original map after creating the Cva
				classesMap["small"] = "mutated-small"
				classesMap["large"] = "button-large"

				tests := []struct {
					name  string
					props Props
					want  string
				}{
					{
						name:  "small",
						props: Props{Size: "small"},
						want:  "button button-small", // Should still use original value
					},
					{
						name:  "medium",
						props: Props{Size: "medium"},
						want:  "button button-medium",
					},
					{
						name:  "large",
						props: Props{Size: "large"},
						want:  "button", // Should not have the new large variant
					},
				}

				for _, test := range tests {
					t.Run(test.name, func(t *testing.T) {
						got := button.Classes(test.props)
						if got != test.want {
							t.Errorf("got %s, want %s", got, test.want)
						}
					})
				}
			})

			t.Run("for_slice_maps", func(t *testing.T) {
				type Props struct {
					Size string
				}

				// Create a map and pass it to Variant
				classesMap := map[string][]string{
					"small":  {"button-small", "text-sm"},
					"medium": {"button-medium", "text-base"},
				}

				button := New(
					Base[Props]("button"),
					MapVariant(
						func(p Props) string { return p.Size },
						classesMap,
					),
				)

				// Mutate the original map after creating the Cva
				classesMap["small"] = []string{"mutated-small", "text-lg"}
				classesMap["large"] = []string{"button-large", "text-xl"}

				tests := []struct {
					name  string
					props Props
					want  string
				}{
					{
						name:  "small",
						props: Props{Size: "small"},
						want:  "button button-small text-sm", // Should still use original values
					},
					{
						name:  "medium",
						props: Props{Size: "medium"},
						want:  "button button-medium text-base",
					},
					{
						name:  "large",
						props: Props{Size: "large"},
						want:  "button", // Should not have the new large variant
					},
				}

				for _, test := range tests {
					t.Run(test.name, func(t *testing.T) {
						got := button.Classes(test.props)
						if got != test.want {
							t.Errorf("got %s, want %s", got, test.want)
						}
					})
				}
			})
		})
	})

	t.Run("CompoundVariant", func(t *testing.T) {
		type Props struct {
			Size  string
			Color string
		}

		button := New(
			Base[Props]("button"),
			CompoundVariant(
				func(p Props) (string, string) { return p.Size, p.Color },
				NewCompound("small", "red", "button-small-red"),
				NewCompound("large", "blue", "button-large-blue"),
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
				got := button.Classes(test.props)
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("PredicateVariant", func(t *testing.T) {
		type Props struct {
			IsDisabled bool
			IsLoading  bool
		}

		button := New(
			Base[Props]("button"),
			PredicateVariant(
				func(p Props) bool { return p.IsDisabled },
				"button-disabled",
			),
			PredicateVariant(
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
				got := button.Classes(test.props)
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("Base", func(t *testing.T) {
		t.Run("single_class", func(t *testing.T) {
			type Props struct{}
			button := New(Base[Props]("button"))
			got := button.Classes(Props{})
			want := "button"
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})

		t.Run("multiple_classes", func(t *testing.T) {
			type Props struct{}
			button := New(Base[Props]("button", "base"))
			got := button.Classes(Props{})
			want := "button base"
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})
	})

	t.Run("Static", func(t *testing.T) {
		t.Run("single_class", func(t *testing.T) {
			type Props struct{}
			button := New(Static[Props]("button"))
			got := button.Classes(Props{})
			want := "button"
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})

		t.Run("multiple_classes", func(t *testing.T) {
			type Props struct{}
			button := New(Static[Props]("button", "base"))
			got := button.Classes(Props{})
			want := "button base"
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})

		t.Run("after_variant", func(t *testing.T) {
			type Props struct {
				Size string
			}

			button := New(
				MapVariant(
					func(p Props) string { return p.Size },
					map[string]string{
						"small":  "button-small",
						"medium": "button-medium",
						"large":  "button-large",
					},
				),
				Static[Props]("button-primary"),
			)

			tests := []struct {
				name string
				size string
				want string
			}{
				{
					name: "small",
					size: "small",
					want: "button-small button-primary",
				},
				{
					name: "medium",
					size: "medium",
					want: "button-medium button-primary",
				},
				{
					name: "large",
					size: "large",
					want: "button-large button-primary",
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					got := button.Classes(Props{Size: test.size})
					if got != test.want {
						t.Errorf("got %s, want %s", got, test.want)
					}
				})
			}
		})
	})

	t.Run("Inherit", func(t *testing.T) {
		type BaseProps struct {
			Size string
		}

		type ExtendedProps struct {
			Size  int
			Color string
		}

		sizes := []string{"small", "medium", "large"}

		base := New(
			Base[BaseProps]("button"),
			MapVariant(
				func(p BaseProps) string { return p.Size },
				map[string]string{
					"small":  "button-small",
					"medium": "button-medium",
					"large":  "button-large",
				},
			),
		)

		extended := New(
			Inherit(
				base,
				func(p ExtendedProps) BaseProps { return BaseProps{Size: sizes[p.Size]} },
			),
			MapVariant(
				func(p ExtendedProps) string { return p.Color },
				map[string]string{
					"red":   "button-red",
					"blue":  "button-blue",
					"green": "button-green",
				},
			),
		)

		tests := []struct {
			name  string
			props ExtendedProps
			want  string
		}{
			{
				name:  "small-red",
				props: ExtendedProps{Size: 0, Color: "red"},
				want:  "button button-small button-red",
			},
			{
				name:  "medium-blue",
				props: ExtendedProps{Size: 1, Color: "blue"},
				want:  "button button-medium button-blue",
			},
			{
				name:  "large-green",
				props: ExtendedProps{Size: 2, Color: "green"},
				want:  "button button-large button-green",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := extended.Classes(test.props)
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("composing_multiple", func(t *testing.T) {
		type Props struct {
			Size    string
			Color   string
			Classes []string
		}

		button := New(
			Base[Props]("base"),
			MapVariant(func(p Props) string { return p.Size }, map[string]string{
				"small":  "small",
				"medium": "medium",
				"large":  "large",
			}),
			MapVariant(func(p Props) string { return p.Color }, map[string]string{
				"red":   "red",
				"blue":  "blue",
				"green": "green",
			}),
			Static[Props]("static"),
			CompoundVariant(
				func(p Props) (string, string) { return p.Size, p.Color },
				NewCompound("small", "red", "small-red"),
				NewCompound("medium", "blue", "medium-blue"),
				NewCompound("large", "green", "large-green"),
			),
			PredicateVariant(func(p Props) bool { return p.Size != "small" }, "predicate"),
			Classes(func(p Props) []string { return p.Classes }),
		)

		tests := []struct {
			name  string
			props Props
			want  string
		}{
			{
				name:  "small-red",
				props: Props{Size: "small", Color: "red"},
				want:  "base small red static small-red",
			},
			{
				name:  "small-blue",
				props: Props{Size: "small", Color: "blue"},
				want:  "base small blue static",
			},
			{
				name:  "small-green",
				props: Props{Size: "small", Color: "green"},
				want:  "base small green static",
			},
			{
				name:  "small-red-custom",
				props: Props{Size: "small", Color: "red", Classes: []string{"custom"}},
				want:  "base small red static small-red custom",
			},
			{
				name:  "small-blue-custom",
				props: Props{Size: "small", Color: "blue", Classes: []string{"custom"}},
				want:  "base small blue static custom",
			},
			{
				name:  "small-green-custom",
				props: Props{Size: "small", Color: "green", Classes: []string{"custom"}},
				want:  "base small green static custom",
			},
			{
				name:  "medium-red",
				props: Props{Size: "medium", Color: "red"},
				want:  "base medium red static predicate",
			},
			{
				name:  "medium-blue",
				props: Props{Size: "medium", Color: "blue"},
				want:  "base medium blue static medium-blue predicate",
			},
			{
				name:  "medium-green",
				props: Props{Size: "medium", Color: "green"},
				want:  "base medium green static predicate",
			},
			{
				name:  "medium-red-custom",
				props: Props{Size: "medium", Color: "red", Classes: []string{"custom"}},
				want:  "base medium red static predicate custom",
			},
			{
				name:  "medium-blue-custom",
				props: Props{Size: "medium", Color: "blue", Classes: []string{"custom"}},
				want:  "base medium blue static medium-blue predicate custom",
			},
			{
				name:  "medium-green-custom",
				props: Props{Size: "medium", Color: "green", Classes: []string{"custom"}},
				want:  "base medium green static predicate custom",
			},
			{
				name:  "large-red",
				props: Props{Size: "large", Color: "red"},
				want:  "base large red static predicate",
			},
			{
				name:  "large-blue",
				props: Props{Size: "large", Color: "blue"},
				want:  "base large blue static predicate",
			},
			{
				name:  "large-green",
				props: Props{Size: "large", Color: "green"},
				want:  "base large green static large-green predicate",
			},
			{
				name:  "large-red-custom",
				props: Props{Size: "large", Color: "red", Classes: []string{"custom"}},
				want:  "base large red static predicate custom",
			},
			{
				name:  "large-blue-custom",
				props: Props{Size: "large", Color: "blue", Classes: []string{"custom"}},
				want:  "base large blue static predicate custom",
			},
			{
				name:  "large-green-custom",
				props: Props{Size: "large", Color: "green", Classes: []string{"custom"}},
				want:  "base large green static large-green predicate custom",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := button.Classes(test.props)
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("repeated_classes", func(t *testing.T) {
		// More specifically, we want to test that we simply join the classes
		// together, without any built-in deduplication.

		type Props struct {
			Size string
		}

		button := New(
			Base[Props]("button"),
			MapVariant(
				func(p Props) string { return p.Size },
				map[string]string{"small": "button", "medium": "button"},
			),
		)

		got := button.Classes(Props{Size: "small"})
		want := "button button"
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("clean_whitespace", func(t *testing.T) {
		type Props struct {
			i int
		}

		button := New(
			Base[Props](" start "),
			MapVariant(
				func(p Props) int { return p.i },
				map[int]string{
					0: "",           // Empty produces an extra "joining" space
					1: "tab	 space", // Tabs get converted to spaces & collapsed
				},
			),
			Static[Props]("end "),
		)

		wants := []string{
			"start end",
			"start tab space end",
		}

		for i, want := range wants {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				got := button.Classes(Props{i})
				if got != want {
					t.Errorf("got %s, want %s", got, want)
				}
			})
		}
	})
}
