// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

// Spec represents a specification for user HCL inputs.
type Spec struct {
	object *objectType
}

// Convert converts the input values to match the spec.
//
// This includes type conversion, default value assignment, and handling of aliases.
// This function should be called before passing the value into the Validate function.
func (s *Spec) Convert(value cty.Value) (cty.Value, error) {
	if s == nil {
		return cty.NilVal, fmt.Errorf("spec is nil")
	}

	if s.object == nil {
		return cty.NilVal, fmt.Errorf("object type is nil")
	}

	converted, err := s.object.Convert(value)
	if err != nil {
		return cty.NilVal, err
	}

	return converted, nil
}

// Validate validates input against the spec.
//
// This function should be called after Convert to ensure the values are in the correct format.
// The validation checks that required fields are present and that all constraints are satisfied.
func (s *Spec) Validate(value cty.Value) error {
	if s == nil {
		return fmt.Errorf("spec is nil")
	}

	if s.object == nil {
		return fmt.Errorf("object type is nil")
	}

	return s.object.Validate(value)
}

// ToProtobuf converts the Spec to its protobuf representation.
func (s *Spec) ToProtobuf() (*SpecPB, error) {
	if s == nil {
		return nil, fmt.Errorf("spec is nil")
	}

	if s.object == nil {
		return nil, fmt.Errorf("object type is nil")
	}

	typ, err := s.object.ToProtobuf()
	if err != nil {
		return nil, err
	}

	obj, ok := typ.Type.(*TypePB_Object)
	if !ok {
		return nil, fmt.Errorf("expected object type protobuf, got %T", typ.Type)
	}

	return &SpecPB{
		Object: obj.Object,
	}, nil
}

// NewSpec creates a new Spec instance.
func NewSpec(object *objectType) *Spec {
	return &Spec{
		object: object,
	}
}

// ToSpec converts a protobuf SpecPB to a Spec instance.
func (s *SpecPB) ToSpec() (*Spec, error) {
	if s == nil {
		return nil, fmt.Errorf("SpecPB is nil")
	}

	if s.Object == nil {
		return nil, fmt.Errorf("Object in SpecPB is nil")
	}

	obj, err := s.Object.ToType()
	if err != nil {
		return nil, err
	}

	return NewSpec(obj), nil
}
