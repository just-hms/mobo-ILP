package set

import (
	"fmt"

	"golang.org/x/exp/maps"
)

// Set represents a collection of unique string values.
type Set[K comparable] map[K]struct{}

// New creates a new Set with the given values.
func New[K comparable](values ...K) Set[K] {
	var s = make(Set[K], len(values))
	s.Add(values...)
	return s
}

// Add adds one or more values to the set.
func (s Set[K]) Add(values ...K) {
	for _, v := range values {
		s[v] = struct{}{}
	}
}

// Delete removes one or more values from the set.
func (s Set[K]) Delete(values ...K) {
	for _, v := range values {
		delete(s, v)
	}
}

// Has returns true if the set contains the given value.
func (s Set[K]) Has(value K) bool {
	_, exists := s[value]
	return exists
}

// Clone returns a new set with the cloned values
func (s Set[K]) Clone() Set[K] {
	var cloned = make(Set[K], len(s))
	cloned.Add(s.List()...)
	return cloned
}

// Intersect returns a new set containing the values that exist in both sets.
func Intersect[K comparable](s1, s2 Set[K]) Set[K] {
	res := Set[K]{}
	for v := range s1 {
		if s2.Has(v) {
			res.Add(v)
		}
	}
	return res
}

// Union add the second set's element to the first one
func Union[K comparable](s1, s2 Set[K]) Set[K] {
	for v := range s2 {
		s1.Add(v)
	}
	return s1
}

// IsSubset returns true if s1 is a subset of s2.
func IsSubset[K comparable](s1, s2 Set[K]) bool {
	if len(s1) > len(s2) {
		return false
	}

	for v := range s1 {
		if !s2.Has(v) {
			return false
		}
	}
	return true
}

// IsSuperset returns true if s1 is a superset of s2.
func IsSuperset[K comparable](s1, s2 Set[K]) bool {
	return IsSubset(s2, s1)
}

func (s Set[K]) List() []K {
	list := make([]K, 0, len(s))
	for v := range s {
		list = append(list, v)
	}
	return list
}

func Difference[K comparable](s1, s2 Set[K]) Set[K] {
	res := Set[K]{}
	for v := range s1 {
		if !s2.Has(v) {
			res.Add(v)
		}
	}
	return res
}

func Equal[K comparable](s1, s2 Set[K]) bool {
	return IsSubset(s1, s2) && IsSuperset(s1, s2)
}

func (s Set[K]) String() string {
	keys := maps.Keys(s)
	res := ""
	prefix := "{"
	for _, k := range keys {
		res += prefix
		res += fmt.Sprint(k)
		prefix = ","
	}
	res += "}"
	return res
}
