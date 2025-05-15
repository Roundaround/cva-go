package cva

import (
	"slices"
)

// Matcher is a chainable predicate function that can be used to match against a property.
type Matcher[P any] struct {
	fn func(p P) bool
}

// Or returns a new Matcher that matches if any of the given matchers match.
func (m Matcher[P]) Or(others ...Matcher[P]) Matcher[P] {
	return Matcher[P]{func(p P) bool {
		if m.fn(p) {
			return true
		}
		for _, other := range others {
			if other.fn(p) {
				return true
			}
		}
		return false
	}}
}

// And returns a new Matcher that matches if all of the given matchers match.
func (m Matcher[P]) And(others ...Matcher[P]) Matcher[P] {
	return Matcher[P]{func(p P) bool {
		if !m.fn(p) {
			return false
		}
		for _, other := range others {
			if !other.fn(p) {
				return false
			}
		}
		return true
	}}
}

// Not returns a new Matcher that matches if the original matcher does not match.
func (m Matcher[P]) Not() Matcher[P] {
	return Matcher[P]{func(p P) bool {
		return !m.fn(p)
	}}
}

// Then returns a new Option that applies the given classes if the matcher matches.
func (m Matcher[P]) Then(classes ...string) Option[P] {
	return PredicateVariant(m.fn, classes...)
}

// NewVariant creates a new Variant that can be used to create Cva Options.
func NewVariant[P any, V comparable](getter func(p P) V) *Variant[P, V] {
	var defaultVal V
	return &Variant[P, V]{getter, defaultVal, false, nil}
}

// Variant is a helper struct that can be used to create Cva Options with its Matcher-producing
// methods like Test, Is, In, IsNot, and NotIn.
type Variant[P any, V comparable] struct {
	getter     func(p P) V
	defaultVal V
	hasDefault bool
	values     []V
}

// WithDefault sets the default value for the variant.
func (v *Variant[P, V]) WithDefault(val V) *Variant[P, V] {
	v.defaultVal = val
	v.hasDefault = true
	return v
}

// WithValues sets the values for the variant.
func (v *Variant[P, V]) WithValues(vals ...V) *Variant[P, V] {
	v.values = vals
	return v
}

func (v Variant[P, V]) get(p P) V {
	var zero V
	val := v.getter(p)

	if v.values != nil && !slices.Contains(v.values, val) {
		val = zero
	}

	if !v.hasDefault || val != zero {
		return val
	}
	return v.defaultVal
}

// Test returns a new Matcher that matches if the variant value matches the given predicate function.
func (v Variant[P, V]) Test(fn func(V) bool) Matcher[P] {
	return Matcher[P]{func(p P) bool {
		return fn(v.get(p))
	}}
}

// Is returns a new Matcher that matches if the variant value is equal to the given value.
func (v Variant[P, V]) Is(val V) Matcher[P] {
	return Matcher[P]{func(p P) bool {
		return v.get(p) == val
	}}
}

// In returns a new Matcher that matches if the variant value is in the given list of values.
func (v Variant[P, V]) In(vals ...V) Matcher[P] {
	return Matcher[P]{func(p P) bool {
		return slices.Contains(vals, v.get(p))
	}}
}

// IsNot returns a new Matcher that matches if the variant value is not equal to the given value.
func (v Variant[P, V]) IsNot(val V) Matcher[P] {
	return v.Is(val).Not()
}

// NotIn returns a new Matcher that matches if the variant value is not in the given list of values.
func (v Variant[P, V]) NotIn(vals ...V) Matcher[P] {
	return v.In(vals...).Not()
}

// Map returns a new Option that applies the given classes if the variant value is in the given map.
func (v Variant[P, V]) Map(m map[V]string) Option[P] {
	return Classes(func(p P) []string {
		if classes, ok := m[v.get(p)]; ok {
			return []string{classes}
		}
		return nil
	})
}

// When returns a new Option that applies the given classes if the given matcher matches.
//
// This is a convience method that is equivalent to calling Then on the matcher.
func When[P any](m Matcher[P], classes ...string) Option[P] {
	return m.Then(classes...)
}

// Any returns a new Matcher that matches if any of the given matchers match.
//
// This is a convience method that is equivalent to chaining matchers with Matcher.Or.
func Any[P any](matchers ...Matcher[P]) Matcher[P] {
	return Matcher[P]{func(p P) bool {
		for _, m := range matchers {
			if m.fn(p) {
				return true
			}
		}
		return false
	}}
}

// All returns a new Matcher that matches if all of the given matchers match.
//
// This is a convience method that is equivalent to chaining matchers with Matcher.And.
func All[P any](matchers ...Matcher[P]) Matcher[P] {
	return Matcher[P]{func(p P) bool {
		for _, m := range matchers {
			if !m.fn(p) {
				return false
			}
		}
		return true
	}}
}
