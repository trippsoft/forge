// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type setType struct {
	elementType Type
	constraints TypeConstraints
}

// WithConstraints implements [Type].
func (s *setType) WithConstraints(constraints ...TypeConstraint) Type {
	if s == nil {
		return nil
	}

	s.constraints = append(s.constraints, constraints...)
	return s
}

// CtyType implements [Type].
func (s *setType) CtyType() cty.Type {
	if s == nil {
		return cty.NilType
	}

	return cty.Set(s.elementType.CtyType())
}

// Convert implements [Type].
func (s *setType) Convert(value cty.Value) (cty.Value, error) {
	if s == nil {
		return cty.NilVal, errors.New("set type is nil")
	}

	converted, err := convertCtyType(value, s.CtyType())
	if err != nil {
		var e error
		converted, e = s.elementType.Convert(value)
		if e != nil {
			return cty.NilVal, err
		}

		return cty.SetVal([]cty.Value{converted}), nil
	}

	if converted.IsNull() {
		return converted, nil
	}

	it := converted.ElementIterator()
	values := make([]cty.Value, 0, converted.LengthInt())
	sensitiveString, isSensitiveString := s.elementType.(*sensitiveStringType)
	for it.Next() {
		_, elem := it.Element()
		if elem.IsNull() {
			continue
		}

		if isSensitiveString {
			sensitiveString.AddToFilter(elem)
		}

		values = append(values, elem)
	}

	if len(values) == 0 {
		return cty.SetValEmpty(s.elementType.CtyType()), nil
	}

	return cty.SetVal(values), nil
}

// Validate implements [Type].
func (s *setType) Validate(value cty.Value) error {
	if s == nil {
		return errors.New("set type is nil")
	}

	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	err := s.constraints.Validate(value)
	if err != nil {
		return err
	}

	it := value.ElementIterator()
	for it.Next() {
		_, elem := it.Element()
		e := s.elementType.Validate(elem)
		if e != nil {
			err = errors.Join(err, fmt.Errorf("invalid set element: %w", e))
		}
	}

	return err
}

// ToProtobuf implements [Type].
func (s *setType) ToProtobuf() (*TypePB, error) {
	if s == nil {
		return nil, errors.New("set type is nil")
	}

	elementType, err := s.elementType.ToProtobuf()
	if err != nil {
		return nil, err
	}

	return &TypePB{
		Type: &TypePB_Set{
			Set: &SetTypePB{
				ElementType: elementType,
			},
		},
	}, nil
}

// String represents the set type as a friendly string.
func (s *setType) String() string {
	if s == nil {
		return "nil"
	}

	return fmt.Sprintf("set of %s", s.elementType)
}

// Set returns a Type representing a set of the given element type.
func Set(elementType Type) Type {
	return &setType{elementType: elementType}
}

// ToType converts a protobuf SetTypePB to a setType instance.
func (s *SetTypePB) ToType() (Type, error) {
	if s == nil {
		return nil, errors.New("SetTypePB is nil")
	}

	if s.ElementType == nil {
		return nil, errors.New("ElementType in SetTypePB is nil")
	}

	elemType, err := s.ElementType.ToType()
	if err != nil {
		return nil, err
	}

	constraints := make(TypeConstraints, 0, len(s.Constraints))
	for _, c := range s.Constraints {
		if c == nil {
			continue
		}

		constraint, err := c.ToTypeConstraint()
		if err != nil {
			return nil, fmt.Errorf("failed to convert constraint: %w", err)
		}

		constraints = append(constraints, constraint)
	}

	return Set(elemType).WithConstraints(constraints...), nil
}
