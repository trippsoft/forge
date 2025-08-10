// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

// Spec represents a specification for inputs.
type Spec struct {
	object *objectType
}

// NewSpec creates a new Spec instance.
func NewSpec(object *objectType) *Spec {
	return &Spec{
		object: object,
	}
}

// Convert converts the input values to the spec's object type.
func (s *Spec) Convert(values map[string]cty.Value) (map[string]cty.Value, error) {

	if s.object == nil {
		return nil, fmt.Errorf("object type is nil")
	}

	convertedValues, err := s.object.convert(values)
	if err != nil {
		return nil, err
	}

	return convertedValues, nil
}

// Validate validates input against the spec.
func (s *Spec) Validate(values map[string]cty.Value) error {

	if s.object == nil {
		return fmt.Errorf("object type is nil")
	}

	return s.object.validate(values)
}

// ValidateSpec validates the spec is valid.
func (s *Spec) ValidateSpec() []error {

	if s.object == nil {
		return []error{fmt.Errorf("object type is nil")}
	}

	return s.object.ValidateSpec()
}
