// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import "github.com/zclconf/go-cty/cty"

// StepIteratorType defines the type of step iterator.
type StepIteratorType int

const (
	StepIteratorTypeSingle StepIteratorType = iota // Single indicates a single execution step iterator.
	StepIteratorTypeList                           // List indicates a list-based step iterator.
	StepIteratorTypeMap                            // Map indicates a map-based step iterator.
)

// StepIteration represents a single iteration in a step iterator.
type StepIteration struct {
	label string
	index cty.Value
	value cty.Value
}

// StepIterationResult holds the result of a single step iteration.
type StepIterationResult struct {
	iteration *StepIteration
	result    cty.Value
}

var (
	StepIterationSingle = &StepIteration{
		label: "",
		index: cty.NilVal,
		value: cty.NilVal,
	}
	StepIteratorSingle StepIterator = &SingleStepIterator{}
)

// StepIterator defines an interface for iterating over a step when the loop block is configured on a step.
type StepIterator interface {
	// Type returns the type of the step iterator.
	Type() StepIteratorType

	// Next returns the current iteration and advances the iterator to the next item.
	//
	// It returns nil if the current item is beyond the end of the collection.
	// It returns true if there are more items to iterate over, false otherwise.
	Next() (*StepIteration, bool)
}

// SingleStepIterator implements a step iterator for steps with no loop block.
type SingleStepIterator struct{}

// Type implements StepIterator.
func (s *SingleStepIterator) Type() StepIteratorType {
	return StepIteratorTypeSingle
}

// Next implements StepIterator.
func (s *SingleStepIterator) Next() (*StepIteration, bool) {
	return StepIterationSingle, false
}

// MultipleStepIterator implements a step iterator for steps with loop blocks.
type MultipleStepIterator struct {
	iteratorType StepIteratorType
	currentIndex int
	items        []*StepIteration
}

// Type implements StepIterator.
func (m *MultipleStepIterator) Type() StepIteratorType {
	return m.iteratorType
}

// Next implements StepIterator.
func (m *MultipleStepIterator) Next() (*StepIteration, bool) {
	if m.currentIndex >= len(m.items) {
		return nil, false
	}

	iteration := m.items[m.currentIndex]
	m.currentIndex++

	return iteration, true
}
