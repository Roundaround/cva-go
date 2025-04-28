package cva

import (
	"slices"
)

// Cva is a class name generator for a component.
//
// The P type parameter is the type of the component's props.
type Cva[P any] struct {
	producers []func(P) []string
}

// Classes generates the class list for the component based on the props.
func (c *Cva[P]) Classes(props P) string {
	parts := make([]string, 0)
	for _, producer := range c.producers {
		parts = append(parts, producer(props)...)
	}
	return JoinClasses(parts...)
}

// New creates a new Cva instance.
func New[P any](opts ...Option[P]) *Cva[P] {
	c := &Cva[P]{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Option is a function that configures a Cva instance.
type Option[P any] func(*Cva[P])

// Classes applies all the classes returned from the supplied getter function.
func Classes[P any, S string | []string](fn func(P) S) Option[P] {
	var nFn func(P) []string
	if sliceFn, ok := any(fn).(func(P) []string); ok {
		nFn = sliceFn
	} else {
		nFn = func(p P) []string {
			return []string{any(fn).(func(P) string)(p)}
		}
	}

	return func(c *Cva[P]) {
		c.producers = append(c.producers, nFn)
	}
}

// Static defines a static class list for the component to be applied regardless of the component's
// props.
func Static[P any](classes ...string) Option[P] {
	return Classes(func(P) []string { return classes })
}

// Base defines a static class list for the component to be applied regardless of the component's
// props. Alias for Static, and included for consistency with the original cva API.
func Base[P any](classes ...string) Option[P] {
	return Classes(func(P) []string { return classes })
}

// MapVariant defines an inline variant as a map of values to class lists.
//
// The classesMap argument accepts both map[V]string and map[V][]string, where the key is the
// variant value for which the associated class list in the map value will be applied.
func MapVariant[P any, V comparable, S string | []string](
	getter func(P) V,
	classesMap map[V]S,
) Option[P] {
	nMap := make(map[V][]string)
	if sliceMap, ok := any(classesMap).(map[V][]string); ok {
		for k, v := range sliceMap {
			nMap[k] = slices.Clone(v)
		}
	} else {
		for k, v := range any(classesMap).(map[V]string) {
			nMap[k] = []string{v}
		}
	}

	return Classes(func(p P) []string {
		key := getter(p)
		if classes, ok := nMap[key]; ok {
			return classes
		}
		return nil
	})
}

// NewCompound creates a Compound value for use in CompoundVariant.
//
// The v1 and v2 arguments should be the variant values to match against returned by the getter
// function in CompoundVariant. Each pair of values should be unique, and adding the same pair of
// values more than once will result in the last occurrence being used.
func NewCompound[V1 comparable, V2 comparable](
	v1 V1,
	v2 V2,
	classes ...string,
) Compound[V1, V2] {
	return Compound[V1, V2]{V1: v1, V2: v2, Classes: classes}
}

// Compound is a set of variant value pairs and associated class lists,
// used in conjunction with CompoundVariant.
type Compound[V1 comparable, V2 comparable] struct {
	V1      V1
	V2      V2
	Classes []string
}

type pair[V1 comparable, V2 comparable] struct {
	V1 V1
	V2 V2
}

// CompoundVariant defines an inline variant as a set of variant value pairs and associated class
// lists.
//
// The getter function should return a tuple of values corresponding to the compound key. Whenever
// this exact tuple of values is encountered, the associated class list will be applied.
//
// The compounds argument should be a list of Compound values, which can be created using the
// NewCompound helper.
func CompoundVariant[P any, V1 comparable, V2 comparable](
	getter func(P) (V1, V2),
	compounds ...Compound[V1, V2],
) Option[P] {
	classesMap := make(map[pair[V1, V2]][]string)
	for _, compound := range compounds {
		classesMap[pair[V1, V2]{compound.V1, compound.V2}] = compound.Classes
	}

	return Classes(func(p P) []string {
		v1, v2 := getter(p)
		if classes, ok := classesMap[pair[V1, V2]{v1, v2}]; ok {
			return classes
		}
		return nil
	})
}

// PredicateVariant defines an inline variant that applies a class list based on a predicate
// function.
func PredicateVariant[P any](
	test func(P) bool,
	classes ...string,
) Option[P] {
	return Classes(func(p P) []string {
		if test(p) {
			return classes
		}
		return nil
	})
}

// Inherit creates a new Cva that inherits all classes and variants from another Cva.
//
// The base argument is the Cva instance to inherit from. The props argument is a function that
// maps the new props type to the base props type, so that it can be passed to all the base Cva's
// producers.
func Inherit[P any, B any](base *Cva[B], baseMapper func(P) B) Option[P] {
	return func(c *Cva[P]) {
		mapped := make([]func(P) []string, len(base.producers))
		for i, producer := range base.producers {
			mapped[i] = func(p P) []string {
				return producer(baseMapper(p))
			}
		}
		c.producers = append(c.producers, mapped...)
	}
}
