package cva

import (
	"regexp"
	"strings"
)

var whitespaceRe = regexp.MustCompile(`\s+`)

// Memoize returns a memoized version of the given function.
//
// The memoized function will cache the result of the first call for each unique
// input value. This is useful for expensive computations or when the same input
// is likely to be used multiple times. Might be useful for prop getters or
// transformers that for some reason are expensive to compute or involve copying
// data.
func Memoize[P comparable, R any](fn func(P) R) func(P) R {
	var lastProps P
	var lastResult R
	var hasLast bool

	return func(p P) R {
		if !hasLast || lastProps != p {
			lastProps = p
			lastResult = fn(p)
			hasLast = true
		}
		return lastResult
	}
}

// DedupeClasses deduplicates classes from the given list and joins them all
// back together with spaces.
//
// Provided as a convience for those who are not combining cva-go with
// TailwindCSS & github.com/Oudwins/tailwind-merge-go.
func DedupeClasses(classes ...string) string {
	split := make([]string, 0, len(classes))
	for _, part := range classes {
		for s := range strings.SplitSeq(part, " ") {
			s = strings.TrimSpace(s)
			if s != "" {
				split = append(split, s)
			}
		}
	}

	deduped := make([]string, 0)
	unique := make(map[string]struct{})
	for _, s := range split {
		if _, ok := unique[s]; !ok {
			unique[s] = struct{}{}
			deduped = append(deduped, s)
		}
	}

	return JoinClasses(deduped...)
}

// JoinClasses joins the given classes together with spaces, trims whitespace, and converts all
// whitespace to single spaces.
func JoinClasses(classes ...string) string {
	return whitespaceRe.ReplaceAllString(strings.TrimSpace(strings.Join(classes, " ")), " ")
}
