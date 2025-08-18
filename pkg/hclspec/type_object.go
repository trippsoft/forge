// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"
	"slices"

	"github.com/trippsoft/forge/pkg/errorwrap"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

type objectType struct {
	rawFields   []*objectField
	fields      map[string]*objectField
	constraints ObjectConstraints
}

// Object creates a new object type with the specified fields.
func Object(fields ...*objectField) *objectType {
	fieldMap := make(map[string]*objectField, len(fields))
	for _, field := range fields {
		if field == nil {
			continue
		}

		fieldMap[field.name] = field
	}

	return &objectType{
		rawFields: fields,
		fields:    fieldMap,
	}
}

// WithConstraints sets the constraints on the object to the specified ones.
func (o *objectType) WithConstraints(constraints ...ObjectConstraint) *objectType {
	o.constraints = constraints
	return o
}

// CtyType implements Type.
func (o *objectType) CtyType() cty.Type {
	if o == nil {
		return cty.NilType
	}

	if len(o.fields) == 0 {
		return cty.EmptyObject
	}

	fieldTypes := make(map[string]cty.Type, len(o.fields))
	for name, field := range o.fields {
		fieldTypes[name] = field.t.CtyType()
	}

	return cty.Object(fieldTypes)
}

// Convert implements Type.
func (o *objectType) Convert(value cty.Value) (cty.Value, error) {
	if o == nil {
		return cty.NilVal, errors.New("object type is nil")
	}

	if !value.IsWhollyKnown() {
		return cty.NilVal, errors.New("cannot convert unknown value")
	}

	if value.IsNull() {
		return cty.NullVal(o.CtyType()), nil // Make sure the null value is of the correct type
	}

	if !value.Type().IsObjectType() && !value.Type().IsMapType() {
		return cty.NilVal, fmt.Errorf("cannot convert %q to %q", value.Type().FriendlyName(), o.CtyType().FriendlyName())
	}

	valueMap := value.AsValueMap()

	values, err := o.convertMap(valueMap)
	if err != nil {
		return cty.NilVal, err
	}

	return cty.ObjectVal(values), nil
}

func (o *objectType) convertMap(values map[string]cty.Value) (map[string]cty.Value, error) {
	if o == nil {
		return nil, errors.New("object type is nil")
	}

	if values == nil {
		values = make(map[string]cty.Value)
	}

	resultFields := make(map[string]cty.Value, len(o.fields))
	validKeys := []string{}

	for name, field := range o.fields {
		fieldValue := field.defaultValue
		foundAs := []string{}
		validKeys = append(validKeys, name)

		if value, ok := values[name]; ok {
			value, err := field.t.Convert(value)
			if err != nil {
				return nil, fmt.Errorf("cannot convert field %q: %w", name, err)
			}
			foundAs = append(foundAs, name)
			fieldValue = value
		}

		for _, alias := range field.aliases {
			validKeys = append(validKeys, alias)
			if value, ok := values[alias]; ok {
				value, err := field.t.Convert(value)
				if err != nil {
					return nil, fmt.Errorf("cannot convert field %q (alias %q): %w", name, alias, err)
				}
				foundAs = append(foundAs, alias)
				fieldValue = value
			}
		}

		if len(foundAs) > 1 {
			return nil, fmt.Errorf("field %q is defined multiple times as %v", name, foundAs)
		}

		resultFields[name] = fieldValue
	}

	invalidIndexes := []string{}
	for key := range values {
		if !slices.Contains(validKeys, key) {
			invalidIndexes = append(invalidIndexes, key)
		}
	}

	if len(invalidIndexes) > 0 {
		return nil, fmt.Errorf("invalid indexes found: %v", invalidIndexes)
	}

	return resultFields, nil
}

// Validate implements Type.
func (o *objectType) Validate(value cty.Value) error {
	if o == nil {
		return fmt.Errorf("object type is nil")
	}

	if !value.IsWhollyKnown() {
		return fmt.Errorf("cannot convert unknown value")
	}

	if value.IsNull() {
		return nil // null values are assumed to be valid
	}

	if !value.Type().IsObjectType() && !value.Type().IsMapType() {
		return fmt.Errorf("cannot convert %q to %q", value.Type().FriendlyName(), o.CtyType().FriendlyName())
	}

	valueMap := value.AsValueMap()

	return o.validateMap(valueMap)
}

func (o *objectType) validateMap(values map[string]cty.Value) error {
	if o == nil {
		return fmt.Errorf("object type is nil")
	}

	if values == nil {
		values = make(map[string]cty.Value)
	}

	var err error
	for name, field := range o.fields {
		value, ok := values[name]
		if !ok {
			err = errors.Join(err, fmt.Errorf("missing field %q", name))
		}

		if e := field.validate(value); e != nil {
			err = errors.Join(err, e)
		}
	}

	e := o.constraints.Validate(values)
	return errors.Join(err, e)
}

// ValidateSpec implements Type.
func (o *objectType) ValidateSpec() error {
	if o == nil {
		return errors.New("object type is nil")
	}

	var err error
	if len(o.fields) != len(o.rawFields) {
		fieldNamesDefined := make(map[string]int, len(o.fields))
		for _, field := range o.rawFields {
			if field == nil {
				err = errors.Join(err, errors.New("nil field definition found"))
				continue
			}

			if _, ok := fieldNamesDefined[field.name]; !ok {
				fieldNamesDefined[field.name] = 0
			}

			fieldNamesDefined[field.name]++
		}

		for name, count := range fieldNamesDefined {
			if count > 1 {
				err = errors.Join(err, fmt.Errorf("field %q is defined multiple times", name))
			}
		}

		if err != nil {
			return err
		}
	}

	definedNames := map[string][]string{}
	for name, field := range o.fields {
		e := field.validateSpec()
		err = errors.Join(err, e)

		definedNames[name] = append(definedNames[name], name)

		for _, alias := range field.aliases {
			definedNames[alias] = append(definedNames[alias], name)
		}
	}

	for name, definitions := range definedNames {
		if len(definitions) > 1 {
			err = errors.Join(err, fmt.Errorf("field %q is defined multiple times (aliases: %v)", name, definitions))
		}
	}

	if o.constraints == nil {
		return err // No constraints to validate
	}

	e := o.constraints.ValidateSpec(o)
	err = errors.Join(err, e)

	return err
}

// ObjectConstraint represents a constraint on an object type.
type ObjectConstraint interface {
	// Validate checks if the given value satisfies the constraint.
	//
	// Implementations are expected to return all errors produced by the validation, if possible.
	Validate(values map[string]cty.Value) error

	// ValidateSpec checks if the constraint is valid for the object type.
	ValidateSpec(t *objectType) error
}

type ObjectConstraints []ObjectConstraint

// Validate checks if all constraints are satisfied by the given value.
//
// This function may return more than one error joined together.
func (c ObjectConstraints) Validate(values map[string]cty.Value) error {
	var err error
	for _, constraint := range c {
		if constraint == nil {
			continue // Skip nil constraints
		}

		if e := constraint.Validate(values); e != nil {
			for _, unwrapped := range errorwrap.UnwrapErrors(e) {
				err = errors.Join(err, fmt.Errorf("validation failed: %w", unwrapped))
			}
		}
	}

	return err
}

// ValidateSpec checks if all constraints are valid for the object type.
//
// This function may return more than one error joined together.
func (c ObjectConstraints) ValidateSpec(t *objectType) error {
	var err error
	for _, constraint := range c {
		if constraint == nil {
			err = errors.Join(err, errors.New("object type has a nil constraint"))
			continue
		}

		if e := constraint.ValidateSpec(t); e != nil {
			for _, unwrapped := range errorwrap.UnwrapErrors(e) {
				err = errors.Join(err, fmt.Errorf("constraint validation failed: %w", unwrapped))
			}
		}
	}

	return err
}

// mutuallyExclusiveGroup represents a group of fields that are mutually exclusive in an object.
type mutuallyExclusiveGroup struct {
	fields []string // List of field names that are mutually exclusive.
}

// MutuallyExclusive creates a new mutuallyExclusiveGroup with the given field names.
func MutuallyExclusive(fields ...string) ObjectConstraint {
	return &mutuallyExclusiveGroup{fields: fields}
}

// Validate implements ObjectConstraint.
func (m *mutuallyExclusiveGroup) Validate(values map[string]cty.Value) error {
	if m == nil {
		return errors.New("mutually exclusive group is nil")
	}

	var foundFields []string
	for _, field := range m.fields {
		if value, ok := values[field]; ok {
			if !value.IsWhollyKnown() || value.IsNull() {
				continue
			}
			foundFields = append(foundFields, field)
		}
	}

	fieldNames := ""
	for i, field := range m.fields {
		if i > 0 {
			fieldNames += ", "
		}
		fieldNames = fmt.Sprintf("%s%q", fieldNames, field)
	}

	if len(foundFields) > 1 {
		return fmt.Errorf("mutually exclusive fields %s are all present", fieldNames)
	}

	return nil
}

// ValidateSpec implements ObjectConstraint.
func (m *mutuallyExclusiveGroup) ValidateSpec(t *objectType) error {
	if m == nil {
		return errors.New("mutually exclusive group is nil")
	}

	var err error
	for _, field := range m.fields {
		if _, ok := t.fields[field]; !ok {
			err = errors.Join(err, fmt.Errorf("field %q is not defined in the object type", field))
		}
	}

	return err
}

// requiredTogetherGroup represents a group of fields that are required together in an object.
type requiredTogetherGroup struct {
	fields []string // List of field names that are required together.
}

// RequiredTogether creates a new requiredTogetherGroup with the given field names.
func RequiredTogether(fields ...string) ObjectConstraint {
	return &requiredTogetherGroup{fields: fields}
}

// Validate implements ObjectConstraint.
func (r *requiredTogetherGroup) Validate(values map[string]cty.Value) error {
	if r == nil {
		return errors.New("required together group is nil")
	}

	var foundFields []string
	for _, field := range r.fields {
		if value, ok := values[field]; ok {
			if !value.IsWhollyKnown() || value.IsNull() {
				continue
			}
			foundFields = append(foundFields, field)
		}
	}

	if len(foundFields) > 0 && len(foundFields) != len(r.fields) {
		return fmt.Errorf("fields %q are required together, but only %q is/are present", r.fields, foundFields)
	}

	return nil
}

// ValidateSpec implements ObjectConstraint.
func (r *requiredTogetherGroup) ValidateSpec(t *objectType) error {
	if r == nil {
		return errors.New("required together group is nil")
	}

	var err error
	for _, field := range r.fields {
		if _, ok := t.fields[field]; !ok {
			err = errors.Join(err, fmt.Errorf("field %q is not defined in the object type", field))
		}
	}

	return err
}

type requiredOneOfGroup struct {
	fields []string // List of field names of which at least one is required.
}

// RequiredOneOf creates a constraint requiring one of the specified fields to be present.
func RequiredOneOf(fields ...string) ObjectConstraint {
	return &requiredOneOfGroup{fields: fields}
}

// Validate implements ObjectConstraint.
func (r *requiredOneOfGroup) Validate(values map[string]cty.Value) error {
	if r == nil {
		return errors.New("required one of group is nil")
	}

	var foundFields []string
	for _, field := range r.fields {
		if value, ok := values[field]; ok {
			if !value.IsWhollyKnown() || value.IsNull() {
				continue
			}
			foundFields = append(foundFields, field)
		}
	}

	if len(foundFields) == 0 {
		return fmt.Errorf("at least one of the fields %q is required", r.fields)
	}

	return nil
}

// ValidateSpec implements ObjectConstraint.
func (r *requiredOneOfGroup) ValidateSpec(t *objectType) error {
	if r == nil {
		return errors.New("required one of group is nil")
	}

	var err error
	for _, field := range r.fields {
		if _, ok := t.fields[field]; !ok {
			err = errors.Join(err, fmt.Errorf("field %q is not defined in the object type", field))
		}
	}

	return err
}

// ObjectCondition is used by constraints that require a specific condition to be met to apply.
type ObjectCondition interface {
	IsMet(values map[string]cty.Value) bool // IsMet checks if the condition is satisfied by the given values.
	Description() string                    // Description provides a human-readable description of the condition for error messages.
	ValidateSpec(t *objectType) error       // ValidateSpec checks if the condition is valid for the object type.
}

type fieldPresentCondition struct {
	field string // The name of the field to check.
}

// FieldPresent creates a condition that checks if a specified field is present.
func FieldPresent(field string) ObjectCondition {
	return &fieldPresentCondition{field: field}
}

// IsMet implements ObjectCondition.
func (c *fieldPresentCondition) IsMet(values map[string]cty.Value) bool {
	if c == nil {
		return false
	}

	if value, ok := values[c.field]; ok {
		return value.IsWhollyKnown() && !value.IsNull()
	}

	return false
}

// Description implements ObjectCondition.
func (c *fieldPresentCondition) Description() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("field %q is present", c.field)
}

// ValidateSpec implements ObjectCondition.
func (c *fieldPresentCondition) ValidateSpec(t *objectType) error {
	if c == nil {
		return errors.New("field present condition is nil")
	}

	if _, ok := t.fields[c.field]; !ok {
		return fmt.Errorf("field %q is not defined in the object type", c.field)
	}

	return nil
}

type fieldNotPresentCondition struct {
	field string // The name of the field to check.
}

// FieldNotPresent creates a condition that checks if a specified field is not present.
func FieldNotPresent(field string) ObjectCondition {
	return &fieldNotPresentCondition{field: field}
}

// IsMet implements ObjectCondition.
func (c *fieldNotPresentCondition) IsMet(values map[string]cty.Value) bool {
	if c == nil {
		return false
	}

	if value, ok := values[c.field]; ok {
		return value.IsWhollyKnown() && value.IsNull()
	}

	return false
}

// Description implements ObjectCondition.
func (c *fieldNotPresentCondition) Description() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("field %q is not present", c.field)
}

// ValidateSpec implements ObjectCondition.
func (c *fieldNotPresentCondition) ValidateSpec(t *objectType) error {
	if c == nil {
		return errors.New("field not present condition is nil")
	}

	if _, ok := t.fields[c.field]; !ok {
		return fmt.Errorf("field %q is not defined in the object type", c.field)
	}

	return nil
}

type fieldEqualsCondition struct {
	field string    // The name of the field to check.
	value cty.Value // The value to compare against.
}

// FieldEquals creates a condition that checks if a specified field has a specified value.
func FieldEquals(field string, value cty.Value) ObjectCondition {
	return &fieldEqualsCondition{field: field, value: value}
}

// IsMet implements ObjectCondition.
func (c *fieldEqualsCondition) IsMet(values map[string]cty.Value) bool {
	if c == nil {
		return false
	}

	if value, ok := values[c.field]; ok {
		return value.IsWhollyKnown() && value.Equals(c.value).True()
	}

	return false
}

// Description implements ObjectCondition.
func (c *fieldEqualsCondition) Description() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("field %q is equal to %s", c.field, util.FormatCtyValueToString(c.value, 0, 0))
}

// ValidateSpec implements ObjectCondition.
func (c *fieldEqualsCondition) ValidateSpec(t *objectType) error {
	if c == nil {
		return errors.New("field equals condition is nil")
	}

	if _, ok := t.fields[c.field]; !ok {
		return fmt.Errorf("field %q is not defined in the object type", c.field)
	}
	return nil
}

type conditionalConstraint struct {
	condition  ObjectCondition
	constraint ObjectConstraint
}

// ConditionalConstraint creates a constraint that only applies if the specified condition is met.
func ConditionalConstraint(condition ObjectCondition, constraint ObjectConstraint) ObjectConstraint {
	return &conditionalConstraint{
		condition:  condition,
		constraint: constraint,
	}
}

// Validate implements ObjectConstraint.
func (c *conditionalConstraint) Validate(values map[string]cty.Value) error {
	if c.condition.IsMet(values) {
		err := c.constraint.Validate(values)
		if err != nil {
			return fmt.Errorf("conditional constraint failed: when %s, %w", c.condition.Description(), err)
		}
	}

	return nil
}

// ValidateSpec implements ObjectConstraint.
func (c *conditionalConstraint) ValidateSpec(t *objectType) error {
	var err error
	if e := c.condition.ValidateSpec(t); e != nil {
		err = e
	}

	return errors.Join(err, c.constraint.ValidateSpec(t))
}
