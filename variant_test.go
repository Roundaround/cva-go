package cva

import (
	"testing"
)

func TestMatcher(t *testing.T) {
	t.Run("Or", func(t *testing.T) {
		type Props struct {
			Value int
		}

		matcher1 := Matcher[Props]{func(p Props) bool { return p.Value == 1 }}
		matcher2 := Matcher[Props]{func(p Props) bool { return p.Value == 2 }}
		matcher3 := Matcher[Props]{func(p Props) bool { return p.Value == 3 }}

		combined := matcher1.Or(matcher2, matcher3)

		tests := []struct {
			name  string
			props Props
			want  bool
		}{
			{
				name:  "matches-first",
				props: Props{Value: 1},
				want:  true,
			},
			{
				name:  "matches-second",
				props: Props{Value: 2},
				want:  true,
			},
			{
				name:  "matches-third",
				props: Props{Value: 3},
				want:  true,
			},
			{
				name:  "matches-none",
				props: Props{Value: 4},
				want:  false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := combined.fn(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("And", func(t *testing.T) {
		type Props struct {
			Value int
			Flag  bool
		}

		matcher1 := Matcher[Props]{func(p Props) bool { return p.Value > 0 }}
		matcher2 := Matcher[Props]{func(p Props) bool { return p.Flag }}

		combined := matcher1.And(matcher2)

		tests := []struct {
			name  string
			props Props
			want  bool
		}{
			{
				name:  "both-true",
				props: Props{Value: 1, Flag: true},
				want:  true,
			},
			{
				name:  "first-true-second-false",
				props: Props{Value: 1, Flag: false},
				want:  false,
			},
			{
				name:  "first-false-second-true",
				props: Props{Value: 0, Flag: true},
				want:  false,
			},
			{
				name:  "both-false",
				props: Props{Value: 0, Flag: false},
				want:  false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := combined.fn(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("Not", func(t *testing.T) {
		type Props struct {
			Value int
		}

		matcher := Matcher[Props]{func(p Props) bool { return p.Value > 0 }}
		notMatcher := matcher.Not()

		tests := []struct {
			name  string
			props Props
			want  bool
		}{
			{
				name:  "positive-value",
				props: Props{Value: 1},
				want:  false,
			},
			{
				name:  "zero-value",
				props: Props{Value: 0},
				want:  true,
			},
			{
				name:  "negative-value",
				props: Props{Value: -1},
				want:  true,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := notMatcher.fn(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("Then", func(t *testing.T) {
		type Props struct {
			Value int
		}

		matcher := Matcher[Props]{func(p Props) bool { return p.Value > 0 }}
		option := matcher.Then("positive", "number")

		button := New(option)

		tests := []struct {
			name  string
			props Props
			want  string
		}{
			{
				name:  "positive-value",
				props: Props{Value: 1},
				want:  "positive number",
			},
			{
				name:  "zero-value",
				props: Props{Value: 0},
				want:  "",
			},
			{
				name:  "negative-value",
				props: Props{Value: -1},
				want:  "",
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
}

func TestVariant(t *testing.T) {
	t.Run("WithDefault", func(t *testing.T) {
		type Props struct {
			Value int
		}

		variant := NewVariant[Props, int](func(p Props) int { return p.Value })
		variant.WithDefault(42)

		tests := []struct {
			name  string
			props Props
			want  int
		}{
			{
				name:  "non-zero-value",
				props: Props{Value: 1},
				want:  1,
			},
			{
				name:  "zero-value",
				props: Props{Value: 0},
				want:  42,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := variant.get(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("WithValues", func(t *testing.T) {
		type Props struct {
			Value int
		}

		variant := NewVariant[Props, int](func(p Props) int { return p.Value })
		variant.WithValues(1, 2, 3)

		tests := []struct {
			name  string
			props Props
			want  int
		}{
			{
				name:  "valid-value",
				props: Props{Value: 1},
				want:  1,
			},
			{
				name:  "invalid-value",
				props: Props{Value: 4},
				want:  0,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := variant.get(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("Test", func(t *testing.T) {
		type Props struct {
			Value int
		}

		variant := NewVariant[Props, int](func(p Props) int { return p.Value })
		matcher := variant.Test(func(v int) bool { return v > 0 })

		tests := []struct {
			name  string
			props Props
			want  bool
		}{
			{
				name:  "positive-value",
				props: Props{Value: 1},
				want:  true,
			},
			{
				name:  "zero-value",
				props: Props{Value: 0},
				want:  false,
			},
			{
				name:  "negative-value",
				props: Props{Value: -1},
				want:  false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := matcher.fn(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("Is", func(t *testing.T) {
		type Props struct {
			Value int
		}

		variant := NewVariant[Props, int](func(p Props) int { return p.Value })
		matcher := variant.Is(42)

		tests := []struct {
			name  string
			props Props
			want  bool
		}{
			{
				name:  "matching-value",
				props: Props{Value: 42},
				want:  true,
			},
			{
				name:  "non-matching-value",
				props: Props{Value: 43},
				want:  false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := matcher.fn(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("In", func(t *testing.T) {
		type Props struct {
			Value int
		}

		variant := NewVariant[Props, int](func(p Props) int { return p.Value })
		matcher := variant.In(1, 2, 3)

		tests := []struct {
			name  string
			props Props
			want  bool
		}{
			{
				name:  "included-value",
				props: Props{Value: 2},
				want:  true,
			},
			{
				name:  "excluded-value",
				props: Props{Value: 4},
				want:  false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := matcher.fn(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("IsNot", func(t *testing.T) {
		type Props struct {
			Value int
		}

		variant := NewVariant[Props, int](func(p Props) int { return p.Value })
		matcher := variant.IsNot(42)

		tests := []struct {
			name  string
			props Props
			want  bool
		}{
			{
				name:  "matching-value",
				props: Props{Value: 42},
				want:  false,
			},
			{
				name:  "non-matching-value",
				props: Props{Value: 43},
				want:  true,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := matcher.fn(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("NotIn", func(t *testing.T) {
		type Props struct {
			Value int
		}

		variant := NewVariant[Props, int](func(p Props) int { return p.Value })
		matcher := variant.NotIn(1, 2, 3)

		tests := []struct {
			name  string
			props Props
			want  bool
		}{
			{
				name:  "included-value",
				props: Props{Value: 2},
				want:  false,
			},
			{
				name:  "excluded-value",
				props: Props{Value: 4},
				want:  true,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := matcher.fn(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("Map", func(t *testing.T) {
		type Props struct {
			Value int
		}

		variant := NewVariant[Props, int](func(p Props) int { return p.Value })
		option := variant.Map(map[int]string{
			1: "one",
			2: "two",
			3: "three",
		})

		button := New(option)

		tests := []struct {
			name  string
			props Props
			want  string
		}{
			{
				name:  "mapped-value",
				props: Props{Value: 2},
				want:  "two",
			},
			{
				name:  "unmapped-value",
				props: Props{Value: 4},
				want:  "",
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
}

func TestConvenienceFunctions(t *testing.T) {
	t.Run("When", func(t *testing.T) {
		type Props struct {
			Value int
		}

		matcher := Matcher[Props]{func(p Props) bool { return p.Value > 0 }}
		option := When(matcher, "positive", "number")

		button := New(option)

		tests := []struct {
			name  string
			props Props
			want  string
		}{
			{
				name:  "positive-value",
				props: Props{Value: 1},
				want:  "positive number",
			},
			{
				name:  "zero-value",
				props: Props{Value: 0},
				want:  "",
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

	t.Run("Any", func(t *testing.T) {
		type Props struct {
			Value int
		}

		matcher1 := Matcher[Props]{func(p Props) bool { return p.Value == 1 }}
		matcher2 := Matcher[Props]{func(p Props) bool { return p.Value == 2 }}
		matcher3 := Matcher[Props]{func(p Props) bool { return p.Value == 3 }}

		combined := Any(matcher1, matcher2, matcher3)

		tests := []struct {
			name  string
			props Props
			want  bool
		}{
			{
				name:  "matches-first",
				props: Props{Value: 1},
				want:  true,
			},
			{
				name:  "matches-second",
				props: Props{Value: 2},
				want:  true,
			},
			{
				name:  "matches-third",
				props: Props{Value: 3},
				want:  true,
			},
			{
				name:  "matches-none",
				props: Props{Value: 4},
				want:  false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := combined.fn(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})

	t.Run("All", func(t *testing.T) {
		type Props struct {
			Value int
			Flag  bool
		}

		matcher1 := Matcher[Props]{func(p Props) bool { return p.Value > 0 }}
		matcher2 := Matcher[Props]{func(p Props) bool { return p.Flag }}

		combined := All(matcher1, matcher2)

		tests := []struct {
			name  string
			props Props
			want  bool
		}{
			{
				name:  "both-true",
				props: Props{Value: 1, Flag: true},
				want:  true,
			},
			{
				name:  "first-true-second-false",
				props: Props{Value: 1, Flag: false},
				want:  false,
			},
			{
				name:  "first-false-second-true",
				props: Props{Value: 0, Flag: true},
				want:  false,
			},
			{
				name:  "both-false",
				props: Props{Value: 0, Flag: false},
				want:  false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := combined.fn(test.props)
				if got != test.want {
					t.Errorf("got %v, want %v", got, test.want)
				}
			})
		}
	})
}
