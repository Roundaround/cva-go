package cva

import (
	"strings"
)

// CvaContext holds configuration for Cva instances.
type CvaContext struct {
	classJoiner func(parts []string) string
}

// NewCvaContext creates a new CvaContext with default configuration.
func NewCvaContext() *CvaContext {
	return &CvaContext{
		classJoiner: defaultClassJoiner,
	}
}

// WithClassJoiner returns a new CvaContext with the specified class joiner.
func (ctx *CvaContext) WithClassJoiner(joiner func(parts []string) string) *CvaContext {
	return &CvaContext{
		classJoiner: joiner,
	}
}

// WithDefaultClassJoiner returns a new CvaContext with the default class joiner.
//
// The default class joiner concatenates parts with a space using strings.Join.
func (ctx *CvaContext) WithDefaultClassJoiner() *CvaContext {
	return ctx.WithClassJoiner(defaultClassJoiner)
}

// WithDedupingClassJoiner returns a new CvaContext with the deduping class joiner.
//
// The deduping class joiner splits parts by spaces, deduplicates them, and
// concatenates them with a space using strings.Join, preserving the order of
// the first occurrence of each class.
func (ctx *CvaContext) WithDedupingClassJoiner() *CvaContext {
	return ctx.WithClassJoiner(dedupingClassJoiner)
}

// NewCva creates a new Cva instance.
//
// The opts argument is a list of options to configure the Cva. See
// Context, StaticClasses, and Variant for some examples.
func NewCva[P any](opts ...Option[P]) *Cva[P] {
	c := &Cva[P]{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Cva is a class name generator for a component.
//
// The P type parameter is the type of the component's props.
type Cva[P any] struct {
	producers []classProducer[P]
	ctx       *CvaContext
}

// ClassName generates the class list for the component based on the props.
func (c *Cva[P]) ClassName(props P) string {
	parts := make([]string, 0)
	for _, producer := range c.producers {
		parts = append(parts, producer.apply(props)...)
	}
	joiner := defaultClassJoiner
	if c.ctx != nil && c.ctx.classJoiner != nil {
		joiner = c.ctx.classJoiner
	}
	return joiner(parts)
}

// Option is a function that configures a Cva instance.
type Option[P any] func(*Cva[P])

// Context sets the context for the Cva instance.
func Context[P any](ctx *CvaContext) Option[P] {
	return func(c *Cva[P]) {
		c.ctx = ctx
	}
}

// StaticClasses defines a static class list for the component to be applied
// regardless of the component's props.
func StaticClasses[P any](classes ...string) Option[P] {
	return func(c *Cva[P]) {
		c.producers = append(c.producers, staticClassList[P]{func(P) []string { return classes }})
	}
}

// Variant defines a variant as a map of values to class lists.
//
// The classesMap argument accepts both map[V]string and map[V][]string,
// where the key is the variant value for which the associated class list in the
// map value will be applied.
func Variant[P any, V comparable, M map[V]string | map[V][]string](
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
		c.producers = append(c.producers, tv)
	}
}

// CompoundVariant defines a variant as a set of variant value pairs and associated class lists.
//
// The getter function should return a tuple of values corresponding to the
// compound key. Whenever this exact tuple of values is encountered, the
// associated class list will be applied.
//
// The compounds argument should be a list of Compound values, which
// can be created using the NewCompound helper.
func CompoundVariant[P any, V1 comparable, V2 comparable](
	getter func(P) (V1, V2),
	compounds ...Compound[V1, V2],
) Option[P] {
	return func(c *Cva[P]) {
		classesMap := make(map[compoundKey[V1, V2]][]string)
		for _, compound := range compounds {
			classesMap[compoundKey[V1, V2]{compound.V1, compound.V2}] = compound.Classes
		}
		cv := compoundVariant[P, V1, V2]{getter, classesMap}
		c.producers = append(c.producers, cv)
	}
}

// NewCompound creates a Compound value. For use in CompoundVariant.
//
// The v1 and v2 arguments should be the variant values to match against returned
// by the getter function in CompoundVariant. Each pair of values should be
// unique, and adding the same pair of values more than once will result in the
// last occurrence being used.
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

// PredicateVariant defines a variant that applies a class list based on a predicate
// function.
func PredicateVariant[P any](
	test func(P) bool,
	classes ...string,
) Option[P] {
	return func(c *Cva[P]) {
		pv := predicateVariant[P]{test, classes}
		c.producers = append(c.producers, pv)
	}
}

// PropsClasses applies all the classes returned from the supplied getter function.
func PropsClasses[P any](
	getter func(P) []string,
) Option[P] {
	return func(c *Cva[P]) {
		nv := staticClassList[P]{getter}
		c.producers = append(c.producers, nv)
	}
}

// Inherit creates a new Cva that inherits all classes and variants from another Cva.
//
// The base argument is the Cva instance to inherit from. The props argument is
// a function that maps the new props type to the base props type, so that it
// can be passed to all the base Cva's producers.
//
// Example:
//
//	type BaseProps struct {
//		Size string
//	}
//
//	type ExtendedProps struct {
//		Size  string
//		Color string
//	}
//
//	base := NewCva[BaseProps](
//		StaticClasses[BaseProps]("button"),
//		Variant(
//			func(p BaseProps) string { return p.Size },
//			map[string]string{
//				"small":  "button-small",
//				"medium": "button-medium",
//				"large":  "button-large",
//			},
//		),
//	)
//
//	extended := NewCva[ExtendedProps](
//		Inherit[ExtendedProps, BaseProps](
//			base,
//			func(p ExtendedProps) BaseProps { return BaseProps{Size: p.Size} },
//		),
//		Variant(
//			func(p ExtendedProps) string { return p.Color },
//			map[string]string{
//				"red":   "button-red",
//				"blue":  "button-blue",
//				"green": "button-green",
//			},
//		),
//	)
func Inherit[P any, B any](base *Cva[B], props func(P) B) Option[P] {
	return func(c *Cva[P]) {
		for _, producer := range base.producers {
			c.producers = append(c.producers, mappedProducer[P, B]{
				base:    producer,
				propsFn: props,
			})
		}

		if base.ctx != nil {
			c.ctx = base.ctx
		}
	}
}

// Memoize returns a memoized version of the given function.
//
// The memoized function will cache the result of the first call for each unique
// input value. This is useful for expensive computations or when the same input
// is likely to be used multiple times. Might be useful for prop getters or
// transformers that for some reason are expensive to compute or involve copying
// data.
//
// Most of the time you will not need this, especially if all your getter
// functions simply return values.
//
// Example:
//
//	type Props struct {
//		Size int
//	}
//
//	type ExtendedProps struct {
//		Size string
//	}
//
//	var base = NewCva[Props](
//		StaticClasses[Props]("button"),
//		Variant(
//			Memoize(func(p Props) string {
//				// Call to some expensive transformation function
//				return sizeString(p.Size)
//			}),
//			map[string]string{
//				"small":  "button-small",
//				"medium": "button-medium",
//				"large":  "button-large",
//			},
//		),
//	)
//
//	var extended = NewCva[ExtendedProps](
//		Inherit(
//			base,
//			Memoize(func(p ExtendedProps) Props {
//				return Props{
//					// Call to some kind of expensive transformation function
//					Size: transformSize(p.Size),
//				}
//			}),
//		),
//		Variant(
//			func(p ExtendedProps) string { return p.Size },
//			map[string]string{
//				"xs":  "button-small",
//				"md": "button-medium",
//				"lg":  "button-large",
//			},
//		),
//	)
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

func defaultClassJoiner(parts []string) string {
	return strings.Join(parts, " ")
}

func dedupingClassJoiner(parts []string) string {
	split := make([]string, 0, len(parts))
	for _, part := range parts {
		for s := range strings.SplitSeq(part, " ") {
			split = append(split, s)
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

	return strings.Join(deduped, " ")
}

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

type mappedProducer[P any, B any] struct {
	base    classProducer[B]
	propsFn func(P) B
}

func (p mappedProducer[P, B]) apply(props P) []string {
	return p.base.apply(p.propsFn(props))
}
