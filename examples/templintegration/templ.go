package templintegration

import (
	"context"
	"os"
)

func Example() {
	props := ButtonProps{
		Classes: []string{"bg-red-500"},
		Size:    Small,
	}

	Button(props).Render(context.Background(), os.Stdout)
	// Output: <button class="inline-flex items-center ..."></button> (x2)
}
