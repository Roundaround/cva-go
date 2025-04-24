package cva

import (
	"strings"
)

var defaultClassJoiner = func(parts []string) string {
	return strings.Join(parts, " ")
}

// The default class joiner is `strings.Join(parts, " ")`. Replace this value
// with your own function to customize the default class joining behavior.
var FallbackClassJoiner = defaultClassJoiner

type classProducer[P any] interface {
	apply(P) []string
}

type mapVariant[P any, V comparable] struct {
	getter     func(P) V
	classesMap map[V][]string
}

func (v mapVariant[P, V]) apply(p P) []string {
	key := v.getter(p)
	if classes, ok := v.classesMap[key]; ok {
		return classes
	}
	return nil
}

type VariantCompound[V1 comparable, V2 comparable] struct {
	V1      V1
	V2      V2
	Classes []string
}

type compoundKey[V1 comparable, V2 comparable] struct {
	V1 V1
	V2 V2
}

type compoundVariant[P any, V1 comparable, V2 comparable] struct {
	getter     func(P) (V1, V2)
	classesMap map[compoundKey[V1, V2]][]string
}

func (v compoundVariant[P, V1, V2]) apply(p P) []string {
	v1, v2 := v.getter(p)
	if classes, ok := v.classesMap[compoundKey[V1, V2]{v1, v2}]; ok {
		return classes
	}
	return nil
}

type predicateVariant[P any] struct {
	test    func(P) bool
	classes []string
}

func (v predicateVariant[P]) apply(p P) []string {
	if v.test(p) {
		return v.classes
	}
	return nil
}

type staticClassList[P any] struct {
	getter func(P) []string
}

func (v staticClassList[P]) apply(p P) []string {
	return v.getter(p)
}

// Create a new `Cva` value.
//
// The `base` argument is the base class list for the component, which will
// always be applied, regardless of the component's props.
//
// The `opts` argument is a list of options to configure the `Cva`. See
// `WithVariant`, `WithCompoundVariant`, `WithPredicateVariant`, and
// `WithClasses` for examples.
func NewCva[P any, S string | []string](base S, opts ...Option[P]) *Cva[P] {
	var normalized []string
	if baseString, ok := any(base).(string); ok {
		normalized = []string{baseString}
	} else {
		normalized = any(base).([]string)
	}

	c := &Cva[P]{base: normalized}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// A `Cva` value is a class name generator for a component.
//
// The `P` type parameter is the type of the component's props.
type Cva[P any] struct {
	base        []string
	variants    []classProducer[P]
	classJoiner func(parts []string) string
}

// Generate the class list for the component based on the props.
func (c *Cva[P]) ClassName(props P) string {
	parts := c.base
	for _, v := range c.variants {
		parts = append(parts, v.apply(props)...)
	}
	joiner := c.classJoiner
	if joiner == nil {
		joiner = FallbackClassJoiner
	}
	return joiner(parts)
}

// An `Option` is a function that configures a `Cva`, usually by defining
// variants.
type Option[P any] func(*Cva[P])

// Define a variant as a map of values to class lists.
//
// The `classesMap` argument accepts both `map[V]string` and `map[V][]string`,
// where the key is the variant value for which the associated class list in the
// map value will be applied.
func WithVariant[P any, V comparable, M map[V]string | map[V][]string](
	getter func(P) V,
	classesMap M,
) Option[P] {
	return func(c *Cva[P]) {
		var normalized map[V][]string
		if mapOfSlices, ok := any(classesMap).(map[V][]string); ok {
			normalized = mapOfSlices
		} else {
			normalized = make(map[V][]string)
			for k, v := range any(classesMap).(map[V]string) {
				normalized[k] = []string{v}
			}
		}

		tv := mapVariant[P, V]{getter, normalized}
		c.variants = append(c.variants, tv)
	}
}

// Define a variant as a set of variant value pairs and associated class lists.
//
// The `getter` function should return a tuple of values corresponding to the
// compound key.
//
// The `compounds` argument should be a list of `VariantCompound` values, which
// can be created using the `WithCompound` helper.
func WithCompoundVariant[P any, V1 comparable, V2 comparable](
	getter func(P) (V1, V2),
	compounds ...VariantCompound[V1, V2],
) Option[P] {
	return func(c *Cva[P]) {
		classesMap := make(map[compoundKey[V1, V2]][]string)
		for _, compound := range compounds {
			classesMap[compoundKey[V1, V2]{compound.V1, compound.V2}] = compound.Classes
		}
		cv := compoundVariant[P, V1, V2]{getter, classesMap}
		c.variants = append(c.variants, cv)
	}
}

// Create a `VariantCompound` value.
//
// The `v1` and `v2` arguments should be the variant values to be assigned to
// the `VariantCompound` and are used as a unique identifier for the compound.
func WithCompound[V1 comparable, V2 comparable](
	v1 V1,
	v2 V2,
	classes ...string,
) VariantCompound[V1, V2] {
	return VariantCompound[V1, V2]{V1: v1, V2: v2, Classes: classes}
}

// Define a variant that applies a class list based on a predicate function.
func WithPredicateVariant[P any](
	test func(P) bool,
	classes ...string,
) Option[P] {
	return func(c *Cva[P]) {
		pv := predicateVariant[P]{test, classes}
		c.variants = append(c.variants, pv)
	}
}

// Apply all the classes returned from the supplied getter function.
func WithClasses[P any](
	getter func(P) []string,
) Option[P] {
	return func(c *Cva[P]) {
		nv := staticClassList[P]{getter}
		c.variants = append(c.variants, nv)
	}
}

// Define a custom class joiner function. By default, the class joiner is
// `strings.Join(parts, " ")`.
func WithClassJoiner[P any](joiner func(parts []string) string) Option[P] {
	return func(c *Cva[P]) {
		c.classJoiner = joiner
	}
}
