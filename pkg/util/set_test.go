// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"reflect"
	"sort"
	"testing"
)

func TestSet_NewSet(t *testing.T) {
	s := NewSet[int]()
	if s.Size() != 0 {
		t.Errorf("Expected empty set size 0, got %d", s.Size())
	}

	s2 := NewSet(1, 2, 3)
	if s2.Size() != 3 {
		t.Errorf("Expected set size 3, got %d", s2.Size())
	}

	s3 := NewSet(1, 2, 2, 3)
	if s3.Size() != 3 {
		t.Errorf("Expected set size 3 (duplicates removed), got %d", s3.Size())
	}
}

func TestSet_Add(t *testing.T) {
	s := NewSet[int]()

	s.Add(1)
	if !s.Contains(1) {
		t.Error("Expected set to contain 1")
	}
	if s.Size() != 1 {
		t.Errorf("Expected size 1, got %d", s.Size())
	}

	s.Add(1)
	if s.Size() != 1 {
		t.Errorf("Expected size to remain 1 after adding duplicate, got %d", s.Size())
	}

	s.Add(2)
	if s.Size() != 2 {
		t.Errorf("Expected size 2, got %d", s.Size())
	}
}

func TestSet_Remove(t *testing.T) {
	s := NewSet(1, 2, 3)

	s.Remove(2)
	if s.Contains(2) {
		t.Error("Expected item 2 to be removed")
	}

	if s.Size() != 2 {
		t.Errorf("Expected size 2, got %d", s.Size())
	}

	s.Remove(5)
	if s.Size() != 2 {
		t.Errorf("Expected size to remain 2 after removing non-existing item, got %d", s.Size())
	}
}

func TestSet_Contains(t *testing.T) {
	s := NewSet("a", "b", "c")

	if !s.Contains("a") {
		t.Error("Expected set to contain 'a'")
	}

	if !s.Contains("b") {
		t.Error("Expected set to contain 'b'")
	}

	if s.Contains("d") {
		t.Error("Expected set not to contain 'd'")
	}
}

func TestSet_Size(t *testing.T) {
	tests := []struct {
		name     string
		items    []int
		expected int
	}{
		{
			name:     "empty set",
			items:    []int{},
			expected: 0,
		},
		{
			name:     "single item",
			items:    []int{1},
			expected: 1,
		},
		{
			name:     "multiple items",
			items:    []int{1, 2, 3},
			expected: 3,
		},
		{
			name:     "duplicates",
			items:    []int{1, 1, 2, 2},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet(tt.items...)
			if s.Size() != tt.expected {
				t.Errorf("Expected size %d, got %d", tt.expected, s.Size())
			}
		})
	}
}

func TestSet_Items(t *testing.T) {
	s := NewSet[int]()
	items := s.Items()
	if len(items) != 0 {
		t.Errorf("Expected empty slice, got %v", items)
	}

	s2 := NewSet(3, 1, 2)
	items2 := s2.Items()
	sort.Ints(items2)

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(items2, expected) {
		t.Errorf("Expected %v, got %v", expected, items2)
	}
}

func TestSet_Clear(t *testing.T) {
	s := NewSet(1, 2, 3)
	s.Clear()

	if !s.IsEmpty() {
		t.Error("Expected set to be empty after clear")
	}

	if s.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", s.Size())
	}
}

func TestSet_IsEmpty(t *testing.T) {
	s := NewSet[int]()
	if !s.IsEmpty() {
		t.Error("Expected new set to be empty")
	}

	s.Add(1)
	if s.IsEmpty() {
		t.Error("Expected set with item to not be empty")
	}

	s.Remove(1)
	if !s.IsEmpty() {
		t.Error("Expected set to be empty after removing all items")
	}
}

func TestSet_Clone(t *testing.T) {
	s := NewSet(1, 2, 3)
	clone := s.Clone()

	if clone.Size() != s.Size() {
		t.Errorf("Expected clone size %d, got %d", s.Size(), clone.Size())
	}

	for _, item := range s.Items() {
		if !clone.Contains(item) {
			t.Errorf("Expected clone to contain item %v", item)
		}
	}

	clone.Add(4)
	if s.Contains(4) {
		t.Error("Expected original set to not contain item added to clone")
	}
}

func TestSet_Union(t *testing.T) {
	result := Union[int]()
	if !result.IsEmpty() {
		t.Error("Expected union of no sets to be empty")
	}

	s1 := NewSet(1, 2, 3)
	result = Union(s1)
	if result.Size() != 3 {
		t.Errorf("Expected union size 3, got %d", result.Size())
	}

	s2 := NewSet(3, 4, 5)
	s3 := NewSet(5, 6, 7)
	result = Union(s1, s2, s3)

	expected := []int{1, 2, 3, 4, 5, 6, 7}
	items := result.Items()
	sort.Ints(items)
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Expected %v, got %v", expected, items)
	}
}

func TestSet_Intersection(t *testing.T) {
	result := Intersection[int]()
	if !result.IsEmpty() {
		t.Error("Expected intersection of no sets to be empty")
	}

	s1 := NewSet(1, 2, 3)
	result = Intersection(s1)
	if result.Size() != 3 {
		t.Errorf("Expected intersection size 3, got %d", result.Size())
	}

	s2 := NewSet(2, 3, 4)
	s3 := NewSet(3, 4, 5)
	result = Intersection(s1, s2, s3)
	if !result.Contains(3) {
		t.Error("Expected intersection to contain 3")
	}

	if result.Size() != 1 {
		t.Errorf("Expected intersection size 1, got %d", result.Size())
	}

	s4 := NewSet(1, 2)
	s5 := NewSet(3, 4)
	result = Intersection(s4, s5)
	if !result.IsEmpty() {
		t.Error("Expected intersection of disjoint sets to be empty")
	}
}

func TestSet_Difference(t *testing.T) {
	s1 := NewSet(1, 2, 3, 4)
	s2 := NewSet(3, 4, 5, 6)

	result := Difference(s1, s2)
	expected := []int{1, 2}
	items := result.Items()
	sort.Ints(items)
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Expected %v, got %v", expected, items)
	}

	empty := NewSet[int]()
	result = Difference(s1, empty)
	if result.Size() != s1.Size() {
		t.Errorf("Expected difference with empty set to equal original, got size %d", result.Size())
	}

	result = Difference(empty, s1)
	if !result.IsEmpty() {
		t.Error("Expected difference of empty set to be empty")
	}
}

func TestSet_SetWithStrings(t *testing.T) {
	s := NewSet("hello", "world")
	s.Add("test")

	if !s.Contains("hello") {
		t.Error("Expected set to contain 'hello'")
	}

	if s.Size() != 3 {
		t.Errorf("Expected size 3, got %d", s.Size())
	}

	s.Remove("world")
	if s.Contains("world") {
		t.Error("Expected 'world' to be removed")
	}
}
