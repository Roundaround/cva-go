package cva

import (
	"strconv"
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
				StaticClasses[Props]("button"),
				Variant(
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
				StaticClasses[Props]("button"),
				Variant(
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
			StaticClasses[Props]("button"),
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
			StaticClasses[Props]("button"),
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
			button := NewCva(StaticClasses[Props]("button"))
			got := button.ClassName(Props{})
			want := "button"
			if got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})

		t.Run("multiple_classes", func(t *testing.T) {
			type Props struct{}
			button := NewCva(StaticClasses[Props]("button", "base"))
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
			StaticClasses[Props]("button"),
			PropsClasses(func(p Props) []string {
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
			StaticClasses[Props]("button"),
			Variant(
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
			Context[ButtonProps](ctx),
			StaticClasses[ButtonProps]("button"),
			Variant(
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

	t.Run("inherit", func(t *testing.T) {
		t.Run("prop_mapping", func(t *testing.T) {
			type BaseProps struct {
				Size string
			}

			type ExtendedProps struct {
				Size  int
				Color string
			}

			sizes := []string{"small", "medium", "large"}

			base := NewCva(
				StaticClasses[BaseProps]("button"),
				Variant(
					func(p BaseProps) string { return p.Size },
					map[string]string{
						"small":  "button-small",
						"medium": "button-medium",
						"large":  "button-large",
					},
				),
			)

			extended := NewCva(
				Inherit(
					base,
					func(p ExtendedProps) BaseProps { return BaseProps{Size: sizes[p.Size]} },
				),
				Variant(
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
					got := extended.ClassName(test.props)
					if got != test.want {
						t.Errorf("got %s, want %s", got, test.want)
					}
				})
			}
		})

		t.Run("context_inheritance", func(t *testing.T) {
			type BaseProps struct {
				Size string
			}

			type ExtendedProps struct {
				Size string
			}

			ctx := NewCvaContext().WithClassJoiner(func(parts []string) string {
				return strings.Join(parts, "::")
			})

			base := NewCva(
				Context[BaseProps](ctx),
				StaticClasses[BaseProps]("button"),
				Variant(
					func(p BaseProps) string { return p.Size },
					map[string]string{
						"small":  "button-small",
						"medium": "button-medium",
						"large":  "button-large",
					},
				),
			)

			extended := NewCva(
				Inherit(base, func(p ExtendedProps) BaseProps { return BaseProps(p) }),
				Variant(
					func(p ExtendedProps) string { return p.Size },
					map[string]string{
						"small":  "extended-small",
						"medium": "extended-medium",
						"large":  "extended-large",
					},
				),
			)

			tests := []struct {
				name  string
				props ExtendedProps
				want  string
			}{
				{
					name:  "small",
					props: ExtendedProps{Size: "small"},
					want:  "button::button-small::extended-small",
				},
				{
					name:  "medium",
					props: ExtendedProps{Size: "medium"},
					want:  "button::button-medium::extended-medium",
				},
				{
					name:  "large",
					props: ExtendedProps{Size: "large"},
					want:  "button::button-large::extended-large",
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					got := extended.ClassName(test.props)
					if got != test.want {
						t.Errorf("got %s, want %s", got, test.want)
					}
				})
			}
		})
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

			button := NewCva(
				StaticClasses[Props]("button"),
				Variant(
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
					got := button.ClassName(test.props)
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

			button := NewCva(
				StaticClasses[Props]("button"),
				Variant(
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
					got := button.ClassName(test.props)
					if got != test.want {
						t.Errorf("got %s, want %s", got, test.want)
					}
				})
			}
		})
	})
}

func TestMemoize(t *testing.T) {
	t.Run("basic_memoization", func(t *testing.T) {
		callCount := 0
		expensiveFn := func(x int) int {
			callCount++
			return x * 2
		}

		memoizedFn := Memoize(expensiveFn)

		// First call should compute the result
		if got := memoizedFn(5); got != 10 {
			t.Errorf("got %d, want %d", got, 10)
		}
		if callCount != 1 {
			t.Errorf("callCount = %d, want %d", callCount, 1)
		}

		// Second call with same input should use cached result
		if got := memoizedFn(5); got != 10 {
			t.Errorf("got %d, want %d", got, 10)
		}
		if callCount != 1 {
			t.Errorf("callCount = %d, want %d", callCount, 1)
		}

		// Call with different input should compute new result
		if got := memoizedFn(6); got != 12 {
			t.Errorf("got %d, want %d", got, 12)
		}
		if callCount != 2 {
			t.Errorf("callCount = %d, want %d", callCount, 2)
		}
	})

	t.Run("struct_memoization", func(t *testing.T) {
		type TestStruct struct {
			A int
			B string
		}

		callCount := 0
		expensiveFn := func(s TestStruct) string {
			callCount++
			return s.B + "-" + strconv.Itoa(s.A)
		}

		memoizedFn := Memoize(expensiveFn)

		// First call should compute the result
		input1 := TestStruct{A: 1, B: "test"}
		if got := memoizedFn(input1); got != "test-1" {
			t.Errorf("got %s, want %s", got, "test-1")
		}
		if callCount != 1 {
			t.Errorf("callCount = %d, want %d", callCount, 1)
		}

		// Second call with same input should use cached result
		if got := memoizedFn(input1); got != "test-1" {
			t.Errorf("got %s, want %s", got, "test-1")
		}
		if callCount != 1 {
			t.Errorf("callCount = %d, want %d", callCount, 1)
		}

		// Call with different input should compute new result
		input2 := TestStruct{A: 2, B: "test"}
		if got := memoizedFn(input2); got != "test-2" {
			t.Errorf("got %s, want %s", got, "test-2")
		}
		if callCount != 2 {
			t.Errorf("callCount = %d, want %d", callCount, 2)
		}
	})

	t.Run("map_memoization", func(t *testing.T) {
		type TestKey struct {
			X int
			Y int
		}

		callCount := 0
		expensiveFn := func(k TestKey) int {
			callCount++
			return k.X + k.Y
		}

		memoizedFn := Memoize(expensiveFn)

		// First call should compute the result
		input1 := TestKey{X: 1, Y: 2}
		if got := memoizedFn(input1); got != 3 {
			t.Errorf("got %d, want %d", got, 3)
		}
		if callCount != 1 {
			t.Errorf("callCount = %d, want %d", callCount, 1)
		}

		// Second call with same input should use cached result
		if got := memoizedFn(input1); got != 3 {
			t.Errorf("got %d, want %d", got, 3)
		}
		if callCount != 1 {
			t.Errorf("callCount = %d, want %d", callCount, 1)
		}

		// Call with different input should compute new result
		input2 := TestKey{X: 3, Y: 4}
		if got := memoizedFn(input2); got != 7 {
			t.Errorf("got %d, want %d", got, 7)
		}
		if callCount != 2 {
			t.Errorf("callCount = %d, want %d", callCount, 2)
		}
	})
}

func TestDefaultClassJoiner(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  string
	}{
		{
			name:  "empty",
			parts: []string{},
			want:  "",
		},
		{
			name:  "single_part",
			parts: []string{"button"},
			want:  "button",
		},
		{
			name:  "multiple_parts",
			parts: []string{"button", "primary", "large"},
			want:  "button primary large",
		},
		{
			name:  "parts_with_spaces",
			parts: []string{"button primary", "large"},
			want:  "button primary large",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := defaultClassJoiner(test.parts)
			if got != test.want {
				t.Errorf("defaultClassJoiner(%v) = %q, want %q", test.parts, got, test.want)
			}
		})
	}
}

func TestDedupingClassJoiner(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  string
	}{
		{
			name:  "empty",
			parts: []string{},
			want:  "",
		},
		{
			name:  "single_part",
			parts: []string{"button"},
			want:  "button",
		},
		{
			name:  "multiple_parts_no_duplicates",
			parts: []string{"button", "primary", "large"},
			want:  "button primary large",
		},
		{
			name:  "parts_with_duplicates",
			parts: []string{"button", "primary", "button", "large"},
			want:  "button primary large",
		},
		{
			name:  "parts_with_spaces_and_duplicates",
			parts: []string{"button primary", "large", "button", "primary large"},
			want:  "button primary large",
		},
		{
			name:  "parts_with_multiple_duplicates",
			parts: []string{"button", "primary", "button", "large", "primary", "button"},
			want:  "button primary large",
		},
		{
			name:  "parts_with_empty_strings",
			parts: []string{"button", "", "primary", " ", "large"},
			want:  "button primary large",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := dedupingClassJoiner(test.parts)
			if got != test.want {
				t.Errorf("dedupingClassJoiner(%v) = %q, want %q", test.parts, got, test.want)
			}
		})
	}
}
