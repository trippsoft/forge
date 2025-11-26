// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package core

// ReadOnlySet is a generic interface representing a read-only set of unique items of type T.
type ReadOnlySet[T comparable] interface {
	// Contains checks if the set contains the specified item.
	Contains(item T) bool
	// Size returns the number of unique items in the set.
	Size() int
	// Items returns a slice of all items in the set.
	// The order of items in the slice is not guaranteed to be consistent.
	Items() []T
	// IsEmpty checks if the set is empty.
	IsEmpty() bool
}

// Set is a generic set data structure that holds unique items of type T.
//
// It provides methods to add, remove, check for existence, and perform set operations.
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet creates a new Set with the provided items.
//
// If no items are provided, it initializes an empty set.
// Duplicate items are ignored, ensuring all items in the set are unique.
// The type T must be comparable, meaning it can be used as a key in a map.
func NewSet[T comparable](items ...T) *Set[T] {
	itemsMap := make(map[T]struct{}, len(items))
	for _, item := range items {
		itemsMap[item] = struct{}{}
	}

	return &Set[T]{
		items: itemsMap,
	}
}

// Add adds an item to the set. If the item already exists, it does nothing.
//
// The item must be of type T, which is comparable.
func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

// Remove removes an item from the set. If the item does not exist, it does nothing.
//
// The item must be of type T, which is comparable.
func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

// Contains checks if the set contains the specified item.
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// Size returns the number of unique items in the set.
func (s *Set[T]) Size() int {
	return len(s.items)
}

// Items returns a slice of all items in the set.
//
// The order of items in the slice is not guaranteed to be consistent.
// If the set is empty, it returns an empty slice.
func (s *Set[T]) Items() []T {
	items := make([]T, 0, len(s.items))
	for item := range s.items {
		items = append(items, item)
	}

	return items
}

// Clear removes all items from the set, leaving it empty.
func (s *Set[T]) Clear() {
	s.items = make(map[T]struct{})
}

// IsEmpty checks if the set is empty.
func (s *Set[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Clone creates a shallow copy of the set and returns it as a new Set instance.
func (s *Set[T]) Clone() *Set[T] {
	clone := NewSet[T]()
	for item := range s.items {
		clone.Add(item)
	}

	return clone
}

// Union returns a new Set that is the union of all provided sets.
//
// It contains all unique items from each set.
// If no sets are provided, it returns an empty set.
func Union[T comparable](sets ...*Set[T]) *Set[T] {
	unionSet := NewSet[T]()
	for _, set := range sets {
		for item := range set.items {
			unionSet.Add(item)
		}
	}

	return unionSet
}

// Intersection returns a new Set that is the intersection of all provided sets.
//
// It contains only the items that are present in all sets.
// If no sets are provided, it returns an empty set.
func Intersection[T comparable](sets ...*Set[T]) *Set[T] {
	if len(sets) == 0 {
		return NewSet[T]()
	}

	intersectionSet := NewSet(sets[0].Items()...)
	for _, set := range sets[1:] {
		for item := range intersectionSet.items {
			if !set.Contains(item) {
				intersectionSet.Remove(item)
			}
		}
	}

	return intersectionSet
}

// Difference returns a new Set that contains items from the first set that are not in the second set.
//
// It effectively computes the set difference: set1 - set2.
func Difference[T comparable](set1, set2 *Set[T]) *Set[T] {
	differenceSet := NewSet[T]()
	for item := range set1.items {
		if !set2.Contains(item) {
			differenceSet.Add(item)
		}
	}

	return differenceSet
}
