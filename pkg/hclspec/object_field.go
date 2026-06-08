// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

type objectField struct {
	name         string
	t            Type
	aliases      []string
	required     bool
	defaultValue cty.Value
}

// WithAliases sets the alias list to the specified ones.
//
// This function will overwrite any previous aliases.
func (f *objectField) WithAliases(aliases ...string) *objectField {
	if f == nil {
		return nil
	}

	f.aliases = aliases
	return f
}

// WithDefaultValue sets the default value to the specified one.
//
// This function will overwrite any previous default value.
func (f *objectField) WithDefaultValue(value cty.Value) *objectField {
	f.defaultValue = value
	return f
}

// Validate checks if the provided value satisfies the field's type and constraints.
func (f *objectField) Validate(value cty.Value) error {
	if f == nil {
		return fmt.Errorf("field is nil")
	}

	if f.t == nil {
		return fmt.Errorf("field %q has no type defined", f.name)
	}

	if value.IsNull() {
		if f.required {
			return fmt.Errorf("missing required field %q", f.name)
		}

		return nil
	}

	var err error
	if e := f.t.Validate(value); e != nil {
		err = errors.Join(err, fmt.Errorf("field %q validation failed: %w", f.name, e))
	}

	return err
}

// ToProtobuf converts the object field to its protobuf representation.
func (f *objectField) ToProtobuf() (string, *ObjectFieldPB, error) {
	if f == nil {
		return "", nil, fmt.Errorf("object field is nil")
	}

	typ, err := f.t.ToProtobuf()
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert field %q type to protobuf: %w", f.name, err)
	}

	defaultValue, err := json.Marshal(f.defaultValue, cty.DynamicPseudoType)
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert field %q default value to protobuf: %w", f.name, err)
	}

	return f.name, &ObjectFieldPB{
		Type:         typ,
		Aliases:      f.aliases,
		Required:     f.required,
		DefaultValue: defaultValue,
	}, nil
}

// RequiredField creates a new required object field with the given name and type.
//
// The default value will be set to a null value of the specified type.
// It is initialized with no aliases.
func RequiredField(name string, t Type) *objectField {
	return &objectField{
		name:         name,
		t:            t,
		required:     true,
		defaultValue: cty.NullVal(t.CtyType()),
	}
}

// OptionalField creates a new optional object field with the given name and type.
//
// The default value will be set to a null value of the specified type.
// It is initialized with no aliases.
func OptionalField(name string, t Type) *objectField {
	return &objectField{
		name:         name,
		t:            t,
		required:     false,
		defaultValue: cty.NullVal(t.CtyType()),
	}
}

// ToObjectField converts a protobuf ObjectFieldPB to an objectField instance.
func (f *ObjectFieldPB) ToObjectField(name string) (*objectField, error) {
	if f == nil {
		return nil, fmt.Errorf("ObjectFieldPB is nil")
	}

	typ, err := f.Type.ToType()
	if err != nil {
		return nil, fmt.Errorf("failed to convert field %q type from protobuf: %w", name, err)
	}

	defaultValue, err := json.Unmarshal(f.DefaultValue, cty.DynamicPseudoType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert field %q default value from protobuf: %w", name, err)
	}

	return &objectField{
		name:         name,
		t:            typ,
		aliases:      f.Aliases,
		required:     f.Required,
		defaultValue: defaultValue,
	}, nil
}
