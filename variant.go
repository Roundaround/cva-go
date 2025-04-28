package cva

import (
	"slices"
)

type Matcher[P any] struct {
	fn func(p P) bool
}

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

func (m Matcher[P]) Not() Matcher[P] {
	return Matcher[P]{func(p P) bool {
		return !m.fn(p)
	}}
}

func (m Matcher[P]) Then(classes ...string) Option[P] {
	return PredicateVariant(m.fn, classes...)
}

func NewVariant[P any, V comparable](getter func(p P) V) *Variant[P, V] {
	var defaultVal V
	return &Variant[P, V]{getter, defaultVal, false, nil}
}

type Variant[P any, V comparable] struct {
	getter     func(p P) V
	defaultVal V
	hasDefault bool
	values     []V
}

func (v *Variant[P, V]) WithDefault(val V) *Variant[P, V] {
	v.defaultVal = val
	v.hasDefault = true
	return v
}

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

func (v Variant[P, V]) Test(fn func(V) bool) Matcher[P] {
	return Matcher[P]{func(p P) bool {
		return fn(v.get(p))
	}}
}

func (v Variant[P, V]) Is(val V) Matcher[P] {
	return Matcher[P]{func(p P) bool {
		return v.get(p) == val
	}}
}

func (v Variant[P, V]) In(vals ...V) Matcher[P] {
	return Matcher[P]{func(p P) bool {
		return slices.Contains(vals, v.get(p))
	}}
}

func (v Variant[P, V]) IsNot(val V) Matcher[P] {
	return v.Is(val).Not()
}

func (v Variant[P, V]) NotIn(vals ...V) Matcher[P] {
	return v.In(vals...).Not()
}

func (v Variant[P, V]) Map(m map[V]string) Option[P] {
	return Classes(func(p P) []string {
		if classes, ok := m[v.get(p)]; ok {
			return []string{classes}
		}
		return nil
	})
}

func When[P any](m Matcher[P], classes ...string) Option[P] {
	return m.Then(classes...)
}

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
