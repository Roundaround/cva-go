package cva

import (
	"strings"
)

func Classes(s ...string) []string {
	return s
}

var defaultClassJoiner = func(parts []string) string {
	return strings.Join(parts, " ")
}

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

type identityApplicator[P any] struct {
	getter func(P) []string
}

func (v identityApplicator[P]) apply(p P) []string {
	return v.getter(p)
}

func NewCva[P any, S string | []string](base S, opts ...Option[P]) *Cva[P] {
	var baseSlice []string
	if baseString, ok := any(base).(string); ok {
		baseSlice = []string{baseString}
	} else {
		baseSlice = any(base).([]string)
	}
	c := &Cva[P]{base: baseSlice}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type Cva[P any] struct {
	base        []string
	variants    []classProducer[P]
	classJoiner func(parts []string) string
}

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

type Option[P any] func(*Cva[P])

func WithVariant[P any, V comparable](
	getter func(P) V,
	classesMap map[V][]string,
) Option[P] {
	return func(c *Cva[P]) {
		tv := mapVariant[P, V]{getter, classesMap}
		c.variants = append(c.variants, tv)
	}
}

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

func WithCompound[V1 comparable, V2 comparable](
	v1 V1,
	v2 V2,
	classes []string,
) VariantCompound[V1, V2] {
	return VariantCompound[V1, V2]{V1: v1, V2: v2, Classes: classes}
}

func WithPredicateVariant[P any](
	test func(P) bool,
	classes ...string,
) Option[P] {
	return func(c *Cva[P]) {
		pv := predicateVariant[P]{test, classes}
		c.variants = append(c.variants, pv)
	}
}

func WithClasses[P any](
	getter func(P) []string,
) Option[P] {
	return func(c *Cva[P]) {
		nv := identityApplicator[P]{getter}
		c.variants = append(c.variants, nv)
	}
}

func WithClassJoiner[P any](joiner func(parts []string) string) Option[P] {
	return func(c *Cva[P]) {
		c.classJoiner = joiner
	}
}
