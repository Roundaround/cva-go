package variantapi

import (
	"fmt"

	"github.com/Roundaround/cva-go"
)

// ButtonProps demonstrates different types of variants
type ButtonProps struct {
	// String variant
	Size string
	// Boolean variant
	IsLoading bool
	// Integer variant
	Count int
	// Custom enum variant
	Status Status
}

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
	StatusWarning Status = "warning"
)

var Button = cva.New(
	// Base classes
	cva.Base[ButtonProps]("inline-flex items-center justify-center rounded-md font-medium"),

	// Simple string variant with Map
	cva.MapVariant(
		func(p ButtonProps) string { return p.Size },
		map[string]string{
			"small":  "h-9 px-3 text-sm",
			"medium": "h-10 px-4 py-2 text-base",
			"large":  "h-11 px-8 py-3 text-lg",
		},
	),

	// Boolean variant with Is
	cva.NewVariant(func(p ButtonProps) bool { return p.IsLoading }).
		Is(true).
		Then("opacity-50 cursor-not-allowed"),

	// Integer variant with Test
	cva.NewVariant(func(p ButtonProps) int { return p.Count }).
		Test(func(count int) bool { return count > 0 }).
		Then("relative"),

	// Integer variant with In
	cva.NewVariant(func(p ButtonProps) int { return p.Count }).
		In(1, 2, 3).
		Then("font-bold"),

	// Custom enum variant with default color
	cva.NewVariant(func(p ButtonProps) Status { return p.Status }).
		WithValues(StatusSuccess, StatusError, StatusWarning).
		Map(map[Status]string{
			StatusSuccess: "bg-green-500 text-white",
			StatusError:   "bg-red-500 text-white",
			StatusWarning: "bg-yellow-500 text-black",
		}),

	// Default color when no status is set
	cva.NewVariant(func(p ButtonProps) Status { return p.Status }).
		Is("").
		Then("bg-green-500 text-white"),
)

func Example() {
	// Example 1: Basic usage with size
	fmt.Println(Button.Classes(ButtonProps{Size: "medium"}))
	// Output: inline-flex items-center justify-center rounded-md font-medium h-10 px-4 py-2 text-base bg-green-500 text-white

	// Example 2: Loading state
	fmt.Println(Button.Classes(ButtonProps{Size: "medium", IsLoading: true}))
	// Output: inline-flex items-center justify-center rounded-md font-medium h-10 px-4 py-2 text-base opacity-50 cursor-not-allowed bg-green-500 text-white

	// Example 3: Count-based styling
	fmt.Println(Button.Classes(ButtonProps{Size: "medium", Count: 2}))
	// Output: inline-flex items-center justify-center rounded-md font-medium h-10 px-4 py-2 text-base relative font-bold bg-green-500 text-white

	// Example 4: Status with error
	fmt.Println(Button.Classes(ButtonProps{Size: "medium", Status: StatusError}))
	// Output: inline-flex items-center justify-center rounded-md font-medium h-10 px-4 py-2 text-base bg-red-500 text-white
}
