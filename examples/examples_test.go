package examples

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/Roundaround/cva-go/examples/additionalclasses"
	"github.com/Roundaround/cva-go/examples/compoundvariants"
	"github.com/Roundaround/cva-go/examples/deduping"
	"github.com/Roundaround/cva-go/examples/inheritance"
	"github.com/Roundaround/cva-go/examples/matchers"
	"github.com/Roundaround/cva-go/examples/predicatevariants"
	"github.com/Roundaround/cva-go/examples/simplecase"
	"github.com/Roundaround/cva-go/examples/simplevariant"
	"github.com/Roundaround/cva-go/examples/templintegration"
)

func TestExamples(t *testing.T) {
	t.Run("additionalclasses", func(t *testing.T) {
		tests := []struct {
			name          string
			size          string
			customClasses []string
			want          string
		}{
			{
				name:          "small+single",
				size:          "small",
				customClasses: []string{"bg-red-500"},
				want:          "inline-flex items-center justify-center h-9 px-3 bg-red-500",
			},
			{
				name:          "small+multiple",
				size:          "small",
				customClasses: []string{"bg-red-500", "rounded-md"},
				want:          "inline-flex items-center justify-center h-9 px-3 bg-red-500 rounded-md",
			},
			{
				name:          "medium+single",
				size:          "medium",
				customClasses: []string{"bg-red-500"},
				want:          "inline-flex items-center justify-center h-10 px-4 py-2 bg-red-500",
			},
			{
				name:          "medium+multiple",
				size:          "medium",
				customClasses: []string{"bg-red-500", "rounded-md"},
				want:          "inline-flex items-center justify-center h-10 px-4 py-2 bg-red-500 rounded-md",
			},
			{
				name:          "large+single",
				size:          "large",
				customClasses: []string{"bg-red-500"},
				want:          "inline-flex items-center justify-center h-11 px-8 py-3 bg-red-500",
			},
			{
				name:          "large+multiple",
				size:          "large",
				customClasses: []string{"bg-red-500", "rounded-md"},
				want:          "inline-flex items-center justify-center h-11 px-8 py-3 bg-red-500 rounded-md",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := additionalclasses.Button.Classes(additionalclasses.Props{Size: test.size, Classes: test.customClasses})
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("compoundvariants", func(t *testing.T) {
		base := "inline-flex items-center justify-center"
		small := "h-8"
		medium := "h-10"
		large := "h-12"
		icon := "bg-gray-100 rounded-full aspect-square"
		regular := "bg-gray-100 rounded-md"
		link := "text-blue-500"
		smallIcon := "[&_svg]:size-4"
		mediumIcon := "[&_svg]:size-5"
		largeIcon := "[&_svg]:size-6"

		tests := []struct {
			name  string
			size  string
			style string
			want  []string
		}{
			{
				name:  "small+icon",
				size:  "small",
				style: "icon",
				want:  []string{base, small, icon, smallIcon},
			},
			{
				name:  "medium+icon",
				size:  "medium",
				style: "icon",
				want:  []string{base, medium, icon, mediumIcon},
			},
			{
				name:  "large+icon",
				size:  "large",
				style: "icon",
				want:  []string{base, large, icon, largeIcon},
			},
			{
				name:  "small+regular",
				size:  "small",
				style: "regular",
				want:  []string{base, small, regular},
			},
			{
				name:  "medium+regular",
				size:  "medium",
				style: "regular",
				want:  []string{base, medium, regular},
			},
			{
				name:  "large+regular",
				size:  "large",
				style: "regular",
				want:  []string{base, large, regular},
			},
			{
				name:  "small+link",
				size:  "small",
				style: "link",
				want:  []string{base, small, link},
			},
			{
				name:  "medium+link",
				size:  "medium",
				style: "link",
				want:  []string{base, medium, link},
			},
			{
				name:  "large+link",
				size:  "large",
				style: "link",
				want:  []string{base, large, link},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := compoundvariants.Button.Classes(compoundvariants.Props{Size: test.size, Style: test.style})
				want := strings.Join(test.want, " ")
				if got != want {
					t.Errorf("got %s, want %s", got, want)
				}
			})
		}
	})

	t.Run("deduping", func(t *testing.T) {
		tests := []struct {
			name string
			size string
			want string
		}{
			{
				name: "small",
				size: "small",
				want: "inline-flex items-center justify-center rounded-md h-8",
			},
			{
				name: "medium",
				size: "medium",
				want: "inline-flex items-center justify-center rounded-md h-10",
			},
			{
				name: "large",
				size: "large",
				want: "inline-flex items-center justify-center rounded-md h-12",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := deduping.DedupedClasses(deduping.Props{Size: test.size})
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("inheritance", func(t *testing.T) {
		base := "inline-flex items-center justify-center"
		small := "h-8 px-3"
		medium := "h-10 px-4"
		large := "h-12 px-6"
		primary := "bg-blue-500 text-white"
		secondary := "bg-gray-200 text-gray-800"
		outline := "border border-gray-300 text-gray-800"
		loading := "opacity-50 cursor-not-allowed"
		plusIcon := "rounded-full [&_svg]:size-4"
		settingsIcon := "rounded-full [&_svg]:size-5"
		closeIcon := "rounded-full [&_svg]:size-6"

		t.Run("base-button", func(t *testing.T) {
			tests := []struct {
				name  string
				props inheritance.ButtonProps
				want  []string
			}{
				{
					name:  "small-primary",
					props: inheritance.ButtonProps{Size: "small", Style: "primary"},
					want:  []string{base, small, primary},
				},
				{
					name:  "medium-secondary",
					props: inheritance.ButtonProps{Size: "medium", Style: "secondary"},
					want:  []string{base, medium, secondary},
				},
				{
					name:  "large-outline",
					props: inheritance.ButtonProps{Size: "large", Style: "outline"},
					want:  []string{base, large, outline},
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					got := inheritance.Button.Classes(test.props)
					want := strings.Join(test.want, " ")
					if got != want {
						t.Errorf("got %s, want %s", got, want)
					}
				})
			}
		})

		t.Run("loading-button", func(t *testing.T) {
			tests := []struct {
				name  string
				props inheritance.LoadingButtonProps
				want  []string
			}{
				{
					name: "small-primary",
					props: inheritance.LoadingButtonProps{
						ButtonProps: inheritance.ButtonProps{Size: "small", Style: "primary"},
						Loading:     false,
					},
					want: []string{base, small, primary},
				},
				{
					name: "medium-secondary-loading",
					props: inheritance.LoadingButtonProps{
						ButtonProps: inheritance.ButtonProps{Size: "medium", Style: "secondary"},
						Loading:     true,
					},
					want: []string{base, medium, secondary, loading},
				},
				{
					name: "large-outline-loading",
					props: inheritance.LoadingButtonProps{
						ButtonProps: inheritance.ButtonProps{Size: "large", Style: "outline"},
						Loading:     true,
					},
					want: []string{base, large, outline, loading},
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					got := inheritance.LoadingButton.Classes(test.props)
					want := strings.Join(test.want, " ")
					if got != want {
						t.Errorf("got %s, want %s", got, want)
					}
				})
			}
		})

		t.Run("icon-button", func(t *testing.T) {
			tests := []struct {
				name  string
				props inheritance.IconButtonProps
				want  []string
			}{
				{
					name: "small-primary-plus",
					props: inheritance.IconButtonProps{
						ButtonProps: inheritance.ButtonProps{Size: "small", Style: "primary"},
						Icon:        "plus",
					},
					want: []string{base, small, primary, plusIcon},
				},
				{
					name: "medium-secondary-settings",
					props: inheritance.IconButtonProps{
						ButtonProps: inheritance.ButtonProps{Size: "medium", Style: "secondary"},
						Icon:        "settings",
					},
					want: []string{base, medium, secondary, settingsIcon},
				},
				{
					name: "large-outline-close",
					props: inheritance.IconButtonProps{
						ButtonProps: inheritance.ButtonProps{Size: "large", Style: "outline"},
						Icon:        "close",
					},
					want: []string{base, large, outline, closeIcon},
				},
			}

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					got := inheritance.IconButton.Classes(test.props)
					want := strings.Join(test.want, " ")
					if got != want {
						t.Errorf("got %s, want %s", got, want)
					}
				})
			}
		})
	})

	t.Run("matchers", func(t *testing.T) {
		// Note: Because twmerge.Merge is non-deterministic, we'll create a basic HTML string with the
		// class attribute set directly as the output of twmerge.Merge to generate the expected outputs
		// at runtime. Thankfully due to a caching mechanism in the twmerge package, the output of
		// twmerge.Merge is deterministic for a given input across a single execution of the test suite.

		tests := []struct {
			name  string
			size  matchers.Size
			theme matchers.Theme
			elem  matchers.Element
			want  string
		}{
			{
				name:  "small+danger+button",
				size:  matchers.SizeSmall,
				theme: matchers.ThemeDanger,
				elem:  matchers.ElementButton,
				want:  twmerge.Merge("px-4 py-1 px-2 bg-red-500 text-white rounded-md"),
			},
			{
				name:  "small+danger+link",
				size:  matchers.SizeSmall,
				theme: matchers.ThemeDanger,
				elem:  matchers.ElementLink,
				want:  twmerge.Merge("px-4 py-1 px-2 bg-red-500 text-white rounded-md text-blue-500 hover:text-blue-600 bg-transparent"),
			},
			{
				name:  "small+danger+icon",
				size:  matchers.SizeSmall,
				theme: matchers.ThemeDanger,
				elem:  matchers.ElementIcon,
				want:  twmerge.Merge("px-4 py-1 px-2 bg-red-500 text-white rounded-full"),
			},
			{
				name:  "small+primary+button",
				size:  matchers.SizeSmall,
				theme: matchers.ThemePrimary,
				elem:  matchers.ElementButton,
				want:  twmerge.Merge("px-4 py-1 px-2 bg-blue-500 text-white rounded-md"),
			},
			{
				name:  "small+primary+link",
				size:  matchers.SizeSmall,
				theme: matchers.ThemePrimary,
				elem:  matchers.ElementLink,
				want:  twmerge.Merge("px-4 py-1 px-2 bg-blue-500 text-white rounded-md text-blue-500 hover:text-blue-600 bg-transparent"),
			},
			{
				name:  "small+primary+icon",
				size:  matchers.SizeSmall,
				theme: matchers.ThemePrimary,
				elem:  matchers.ElementIcon,
				want:  twmerge.Merge("px-4 py-1 px-2 bg-blue-500 text-white rounded-full"),
			},
			{
				name:  "large+danger+button",
				size:  matchers.SizeLarge,
				theme: matchers.ThemeDanger,
				elem:  matchers.ElementButton,
				want:  twmerge.Merge("px-4 py-1 px-6 py-2 bg-red-500 text-white rounded-md font-bold"),
			},
			{
				name:  "large+danger+link",
				size:  matchers.SizeLarge,
				theme: matchers.ThemeDanger,
				elem:  matchers.ElementLink,
				want:  twmerge.Merge("px-4 py-1 px-6 py-2 bg-red-500 text-white rounded-md font-bold text-blue-500 hover:text-blue-600 bg-transparent"),
			},
			{
				name:  "large+danger+icon",
				size:  matchers.SizeLarge,
				theme: matchers.ThemeDanger,
				elem:  matchers.ElementIcon,
				want:  twmerge.Merge("px-4 py-1 px-6 py-2 bg-red-500 text-white rounded-full font-bold [&_svg]:size-5"),
			},
			{
				name:  "large+primary+button",
				size:  matchers.SizeLarge,
				theme: matchers.ThemePrimary,
				elem:  matchers.ElementButton,
				want:  twmerge.Merge("px-4 py-1 px-6 py-2 bg-blue-500 text-white rounded-md"),
			},
			{
				name:  "large+primary+link",
				size:  matchers.SizeLarge,
				theme: matchers.ThemePrimary,
				elem:  matchers.ElementLink,
				want:  twmerge.Merge("px-4 py-1 px-6 py-2 bg-blue-500 text-white rounded-md text-blue-500 hover:text-blue-600 bg-transparent"),
			},
			{
				name:  "large+primary+icon",
				size:  matchers.SizeLarge,
				theme: matchers.ThemePrimary,
				elem:  matchers.ElementIcon,
				want:  twmerge.Merge("px-4 py-1 px-6 py-2 bg-blue-500 text-white rounded-full [&_svg]:size-5"),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := twmerge.Merge(matchers.Button.Classes(matchers.Props{Size: test.size, Theme: test.theme, Element: test.elem}))
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("predicatevariants", func(t *testing.T) {
		tests := []struct {
			name     string
			loading  bool
			disabled bool
			want     string
		}{
			{
				name:     "loading+disabled",
				loading:  true,
				disabled: true,
				want:     "button button-loading button-disabled",
			},
			{
				name:     "loading+!disabled",
				loading:  true,
				disabled: false,
				want:     "button button-loading button-disabled",
			},
			{
				name:     "!loading+disabled",
				loading:  false,
				disabled: true,
				want:     "button button-disabled",
			},
			{
				name:     "!loading+!disabled",
				loading:  false,
				disabled: false,
				want:     "button",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := predicatevariants.Button.Classes(predicatevariants.Props{Loading: test.loading, Disabled: test.disabled})
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("simplecase", func(t *testing.T) {
		base := "inline-flex items-center justify-center"
		small := "h-9 px-3"
		medium := "h-10 px-4 py-2 rounded-md"
		large := "h-11 px-8 py-3 rounded-md"

		tests := []struct {
			name string
			size string
			want string
		}{
			{
				name: "small",
				size: "small",
				want: base + " " + small,
			},
			{
				name: "medium",
				size: "medium",
				want: base + " " + medium,
			},
			{
				name: "large",
				size: "large",
				want: base + " " + large,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := simplecase.Button.Classes(simplecase.Props{Size: test.size})
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("simplevariant", func(t *testing.T) {
		base := "inline-flex items-center justify-center"
		small := "h-9 px-3"
		medium := "h-10 px-4 py-2 rounded-md"
		large := "h-11 px-8 py-3 rounded-md"

		tests := []struct {
			name string
			size string
			want string
		}{
			{
				name: "small",
				size: "small",
				want: base + " " + small,
			},
			{
				name: "medium",
				size: "medium",
				want: base + " " + medium,
			},
			{
				name: "large",
				size: "large",
				want: base + " " + large,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := simplevariant.Button.Classes(simplevariant.Props{Size: test.size})
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})

	t.Run("templintegration", func(t *testing.T) {
		// Note: Because twmerge.Merge is non-deterministic, we'll create a basic HTML string with the
		// class attribute set directly as the output of twmerge.Merge to generate the expected outputs
		// at runtime. Thankfully due to a caching mechanism in the twmerge package, the output of
		// twmerge.Merge is deterministic for a given input across a single execution of the test suite.

		base := "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors"
		small := "h-9 px-3"
		medium := "h-10 px-4 py-2"
		large := "h-11 px-8 py-3"
		bg := "bg-red-500"

		goldenButton := func(classes ...string) string {
			return fmt.Sprintf(`<button class="%s"></button>`, twmerge.Merge(classes...))
		}

		tests := []struct {
			name    string
			size    templintegration.Size
			classes []string
			want    string
		}{
			{
				name:    "small",
				size:    templintegration.Small,
				classes: []string{},
				want:    goldenButton(base, small),
			},
			{
				name:    "medium",
				size:    templintegration.Medium,
				classes: []string{},
				want:    goldenButton(base, medium),
			},
			{
				name:    "large",
				size:    templintegration.Large,
				classes: []string{},
				want:    goldenButton(base, large),
			},
			{
				name:    "small+bg",
				size:    templintegration.Small,
				classes: []string{bg},
				want:    goldenButton(base, small, bg),
			},
			{
				name:    "medium+bg",
				size:    templintegration.Medium,
				classes: []string{bg},
				want:    goldenButton(base, medium, bg),
			},
			{
				name:    "large+bg",
				size:    templintegration.Large,
				classes: []string{bg},
				want:    goldenButton(base, large, bg),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				buf := new(bytes.Buffer)
				props := templintegration.ButtonProps{
					Classes: test.classes,
					Size:    test.size,
				}
				templintegration.Button(props).Render(context.Background(), buf)
				got := buf.String()
				if got != test.want {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})
}
