package cva

import (
	"testing"
)

func TestExample(t *testing.T) {
	// TODO: Obviously not a real test. Move this example implementation to a new example folder/impl
	// and the README.

	type Variants struct {
		Size  string
		Style string
	}

	type Props struct {
		Variants
		Foo string
		Bar string
		Baz string
	}

	compound := NewVariant(func(p Props) Variants { return p.Variants })

	button := New(
		Base[Props]("button"),
		compound.Is(Variants{"small", "icon"}).Then("small-icon"),
	)

	want := "button small-icon"
	got := button.Classes(Props{Variants: Variants{"small", "icon"}})
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
