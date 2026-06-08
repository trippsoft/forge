// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"

	"github.com/trippsoft/forge/pkg/secret"
	"github.com/zclconf/go-cty/cty"
)

type sensitiveStringType struct {
	constraints TypeConstraints
}

// WithConstraints implements [Type].
func (s *sensitiveStringType) WithConstraints(constraints ...TypeConstraint) Type {
	if s == nil {
		return nil
	}

	s.constraints = append(s.constraints, constraints...)
	return s
}

// CtyType implements [Type].
func (s *sensitiveStringType) CtyType() cty.Type {
	return cty.String
}

// Convert implements [Type].
func (s *sensitiveStringType) Convert(value cty.Value) (cty.Value, error) {
	if s == nil {
		return cty.NilVal, errors.New("sensitive string type is nil")
	}

	v, err := convertCtyType(value, s.CtyType())
	if err == nil {
		s.AddToFilter(v)
	}

	return v, err
}

// Validate implements [Type].
func (s *sensitiveStringType) Validate(value cty.Value) error {
	if s == nil {
		return errors.New("sensitive string type is nil")
	}

	return s.constraints.Validate(value)
}

// ToProtobuf implements [Type].
func (s *sensitiveStringType) ToProtobuf() (*TypePB, error) {
	if s == nil {
		return nil, errors.New("sensitive string type is nil")
	}

	constraints, err := s.constraints.ToProtobuf()
	if err != nil {
		return nil, err
	}

	return &TypePB{
		Type: &TypePB_Simple{
			Simple: &SimpleTypePB{
				Data:        SimpleTypeDataPB_SENSITIVE_STRING,
				Constraints: constraints,
			},
		},
	}, nil
}

// AddToFilter adds the sensitive string value to the secret filter.
func (s *sensitiveStringType) AddToFilter(value cty.Value) {
	if value.IsNull() {
		return // No need to add null values to the filter.
	}

	v := value.AsString()
	if v != "" {
		secret.SecretFilter.AddSecret(v)
	}
}

// String represents the sensitive string type as a friendly string.
func (s *sensitiveStringType) String() string {
	return "sensitive string"
}

// SensitiveString returns a Type that represents a string that should be treated as sensitive.
func SensitiveString() Type {
	return &sensitiveStringType{}
}
