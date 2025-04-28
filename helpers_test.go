package cva

import (
	"strconv"
	"testing"
)

func TestMemoize(t *testing.T) {
	t.Run("scalar", func(t *testing.T) {
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

	t.Run("struct", func(t *testing.T) {
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

	t.Run("map", func(t *testing.T) {
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

func TestDedupeClasses(t *testing.T) {
	tests := []struct {
		name    string
		classes []string
		want    string
	}{
		{
			name:    "empty",
			classes: []string{},
			want:    "",
		},
		{
			name:    "single_part",
			classes: []string{"button"},
			want:    "button",
		},
		{
			name:    "multiple_parts_no_duplicates",
			classes: []string{"button", "primary", "large"},
			want:    "button primary large",
		},
		{
			name:    "parts_with_duplicates",
			classes: []string{"button", "primary", "button", "large"},
			want:    "button primary large",
		},
		{
			name:    "parts_with_spaces_and_duplicates",
			classes: []string{"button primary", "large", "button", "primary large"},
			want:    "button primary large",
		},
		{
			name:    "parts_with_multiple_duplicates",
			classes: []string{"button", "primary", "button", "large", "primary", "button"},
			want:    "button primary large",
		},
		{
			name:    "parts_with_empty_strings",
			classes: []string{"button", "", "primary", " ", "large"},
			want:    "button primary large",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := DedupeClasses(test.classes...)
			if got != test.want {
				t.Errorf("DedupeClasses(%v) = %q, want %q", test.classes, got, test.want)
			}
		})
	}
}

func TestJoinClasses(t *testing.T) {
	tests := []struct {
		classes []string
		want    string
	}{
		{
			classes: []string{},
			want:    "",
		},
		{
			classes: []string{"    simple       "},
			want:    "simple",
		},
		{
			classes: []string{" start tab	 space  end\n"},
			want:    "start tab space end",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := JoinClasses(test.classes...)
			if got != test.want {
				t.Errorf("joinClasses(%v) = %q, want %q", test.classes, got, test.want)
			}
		})
	}
}
