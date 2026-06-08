// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type listType struct {
	elementType Type
	constraints TypeConstraints
}

// WithConstraints implements [Type].
func (l *listType) WithConstraints(constraints ...TypeConstraint) Type {
	if l == nil {
		return nil
	}

	l.constraints = append(l.constraints, constraints...)
	return l
}

// CtyType implements [Type].
func (l *listType) CtyType() cty.Type {
	if l == nil {
		return cty.NilType
	}

	return cty.List(l.elementType.CtyType())
}

// Convert implements [Type].
func (l *listType) Convert(value cty.Value) (cty.Value, error) {
	if l == nil {
		return cty.NilVal, errors.New("list type is nil")
	}

	converted, err := convertCtyType(value, l.CtyType())
	if err != nil {
		var e error
		converted, e = l.elementType.Convert(value)
		if e != nil {
			return cty.NilVal, err
		}

		return cty.ListVal([]cty.Value{converted}), nil
	}

	if converted.IsNull() {
		return converted, nil
	}

	it := converted.ElementIterator()
	values := make([]cty.Value, 0, converted.LengthInt())
	sensitiveString, isSensitiveString := l.elementType.(*sensitiveStringType)
	for it.Next() {
		_, elem := it.Element()
		if elem.IsNull() {
			continue // Skip null elements.
		}

		if isSensitiveString {
			sensitiveString.AddToFilter(elem)
		}

		values = append(values, elem)
	}

	if len(values) == 0 {
		return cty.ListValEmpty(l.elementType.CtyType()), nil
	}

	return cty.ListVal(values), nil
}

// Validate implements [Type].
func (l *listType) Validate(value cty.Value) error {
	if l == nil {
		return errors.New("list type is nil")
	}

	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	err := l.constraints.Validate(value)

	it := value.ElementIterator()
	for it.Next() {
		index, elem := it.Element()
		e := l.elementType.Validate(elem)
		if e != nil {
			i, _ := index.AsBigFloat().Int64()
			err = errors.Join(err, fmt.Errorf("element at index %d: %w", i, e))
		}
	}

	return err
}

// ToProtobuf implements [Type].
func (l *listType) ToProtobuf() (*TypePB, error) {
	if l == nil {
		return nil, errors.New("list type is nil")
	}

	elementType, err := l.elementType.ToProtobuf()
	if err != nil {
		return nil, err
	}

	return &TypePB{
		Type: &TypePB_List{
			List: &ListTypePB{
				ElementType: elementType,
			},
		},
	}, nil
}

// String represents the list type as a friendly string.
func (l *listType) String() string {
	return fmt.Sprintf("list of %s", l.elementType)
}

// List represents a Type for a list of elements of the given type.
func List(elementType Type) Type {
	return &listType{elementType: elementType}
}

// ToType converts a protobuf ListTypePB to a listType instance.
func (l *ListTypePB) ToType() (Type, error) {
	if l == nil {
		return nil, errors.New("ListTypePB is nil")
	}

	if l.ElementType == nil {
		return nil, errors.New("ElementType in ListTypePB is nil")
	}

	elementType, err := l.ElementType.ToType()
	if err != nil {
		return nil, err
	}

	constraints := make(TypeConstraints, 0, len(l.Constraints))
	for _, c := range l.Constraints {
		if c == nil {
			continue
		}

		constraint, err := c.ToTypeConstraint()
		if err != nil {
			return nil, fmt.Errorf("failed to convert constraint: %w", err)
		}

		constraints = append(constraints, constraint)
	}

	return List(elementType).WithConstraints(constraints...), nil
}
