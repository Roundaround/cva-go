package cva

import (
	"strings"
)

// Context holds configuration for CVA instances.
type Context struct {
	classJoiner func(parts []string) string
}

// NewCvaContext creates a new Context with default configuration.
func NewCvaContext() *Context {
	return &Context{
		classJoiner: defaultClassJoiner,
	}
}

// WithClassJoiner returns a new Context with the specified class joiner.
func (ctx *Context) WithClassJoiner(joiner func(parts []string) string) *Context {
	return &Context{
		classJoiner: joiner,
	}
}

// WithDefaultClassJoiner returns a new Context with the default class joiner.
//
// The default class joiner concatenates parts with a space using strings.Join.
func (ctx *Context) WithDefaultClassJoiner() *Context {
	return ctx.WithClassJoiner(defaultClassJoiner)
}

// WithDedupingClassJoiner returns a new Context with the deduping class joiner.
//
// The deduping class joiner splits parts by spaces, deduplicates them, and
// concatenates them with a space using strings.Join, preserving the order of
// the first occurrence of each class.
func (ctx *Context) WithDedupingClassJoiner() *Context {
	return ctx.WithClassJoiner(dedupingClassJoiner)
}

// NewCva creates a new Cva instance.
//
// The opts argument is a list of options to configure the Cva. See
// WithContext, WithStaticClasses, and WithVariant for some examples.
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
	ctx       *Context
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

// WithContext sets the context for the Cva instance.
func WithContext[P any](ctx *Context) Option[P] {
	return func(c *Cva[P]) {
		c.ctx = ctx
	}
}

// WithStaticClasses defines a static class list for the component to be applied
// regardless of the component's props.
func WithStaticClasses[P any](classes ...string) Option[P] {
	return func(c *Cva[P]) {
		c.producers = append(c.producers, staticClassList[P]{func(P) []string { return classes }})
	}
}

// WithVariant defines a variant as a map of values to class lists.
//
// The classesMap argument accepts both map[V]string and map[V][]string,
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
		c.producers = append(c.producers, tv)
	}
}

// WithCompoundVariant defines a variant as a set of variant value pairs and associated class lists.
//
// The getter function should return a tuple of values corresponding to the
// compound key. Whenever this exact tuple of values is encountered, the
// associated class list will be applied.
//
// The compounds argument should be a list of VariantCompound values, which
// can be created using the WithCompound helper.
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
		c.producers = append(c.producers, cv)
	}
}

// WithCompound creates a VariantCompound value. For use in WithCompoundVariant.
//
// The v1 and v2 arguments should be the variant values to match against returned
// by the getter function in WithCompoundVariant. Each pair of values should be
// unique, and adding the same pair of values more than once will result in the
// last occurrence being used.
func WithCompound[V1 comparable, V2 comparable](
	v1 V1,
	v2 V2,
	classes ...string,
) VariantCompound[V1, V2] {
	return VariantCompound[V1, V2]{V1: v1, V2: v2, Classes: classes}
}

// VariantCompound is a set of variant value pairs and associated class lists,
// used in conjunction with WithCompoundVariant.
type VariantCompound[V1 comparable, V2 comparable] struct {
	V1      V1
	V2      V2
	Classes []string
}

// WithPredicateVariant defines a variant that applies a class list based on a predicate function.
func WithPredicateVariant[P any](
	test func(P) bool,
	classes ...string,
) Option[P] {
	return func(c *Cva[P]) {
		pv := predicateVariant[P]{test, classes}
		c.producers = append(c.producers, pv)
	}
}

// WithPropsClasses applies all the classes returned from the supplied getter function.
func WithPropsClasses[P any](
	getter func(P) []string,
) Option[P] {
	return func(c *Cva[P]) {
		nv := staticClassList[P]{getter}
		c.producers = append(c.producers, nv)
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
