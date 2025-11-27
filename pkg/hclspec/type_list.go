// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type listType struct {
	elementType Type // The type of elements in the list.
}

// CtyType implements Type.
func (l *listType) CtyType() cty.Type {
	if l == nil {
		return cty.NilType
	}

	return cty.List(l.elementType.CtyType())
}

// Convert implements Type.
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

// Validate implements Type.
func (l *listType) Validate(value cty.Value) error {
	if l == nil {
		return errors.New("list type is nil")
	}

	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	it := value.ElementIterator()
	var err error
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

// ValidateSpec implements Type.
func (l *listType) ValidateSpec() error {
	if l == nil {
		return errors.New("list type is nil")
	}

	return l.elementType.ValidateSpec()
}

// String represents the list type as a friendly string.
func (l *listType) String() string {
	return l.CtyType().FriendlyName()
}

// List creates a new listType representing a list of the given element type.
func List(elementType Type) Type {
	return &listType{elementType: elementType}
}
