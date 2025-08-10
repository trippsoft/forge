package hclspec

import (
	"fmt"
	"strings"

	"github.com/zclconf/go-cty/cty"
)

// objectConstraint represents a constraint on an object type.
type objectConstraint interface {
	Validate(values map[string]cty.Value) error // Validate checks if the given value satisfies the constraint.
	ValidateSpec(t *objectType) error           // ValidateSpec checks if the constraint is valid for the object type.
}

type ObjectConstraints []objectConstraint

// Validate checks if all constraints in the objectConstraints slice are satisfied by the given value.
func (c ObjectConstraints) Validate(values map[string]cty.Value) error {
	for _, constraint := range c {
		if err := constraint.Validate(values); err != nil {
			return err
		}
	}
	return nil
}

// mutuallyExclusiveGroup represents a group of fields that are mutually exclusive in an object.
type mutuallyExclusiveGroup struct {
	fields []string // List of field names that are mutually exclusive.
}

// MutuallyExclusive creates a new mutuallyExclusiveGroup with the given field names.
func MutuallyExclusive(fields ...string) objectConstraint {
	return &mutuallyExclusiveGroup{fields: fields}
}

// Validate ensures that the mutually exclusive fields are not all present in the given value.
// This should be called after the object aliases are handled and any conversions are completed.
func (m *mutuallyExclusiveGroup) Validate(values map[string]cty.Value) error {

	var foundFields []string
	for _, field := range m.fields {
		if value, ok := values[field]; ok {
			if !value.IsWhollyKnown() || value.IsNull() {
				continue
			}
			foundFields = append(foundFields, field)
		}
	}

	if len(foundFields) > 1 {
		return fmt.Errorf("mutually exclusive fields %q are all present", foundFields)
	}

	return nil
}

// ValidateSpec checks if the constraint is valid for the object type.
func (m *mutuallyExclusiveGroup) ValidateSpec(t *objectType) error {
	for _, field := range m.fields {
		if _, ok := t.fields[field]; !ok {
			return fmt.Errorf("field %q is not defined in the object type", field)
		}
	}
	return nil
}

// requiredTogetherGroup represents a group of fields that are required together in an object.
type requiredTogetherGroup struct {
	fields []string // List of field names that are required together.
}

// RequiredTogether creates a new requiredTogetherGroup with the given field names.
func RequiredTogether(fields ...string) objectConstraint {
	return &requiredTogetherGroup{fields: fields}
}

// Validate checks if all or none of the required fields are present in the given value.
func (r *requiredTogetherGroup) Validate(values map[string]cty.Value) error {

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
		return fmt.Errorf("fields %q are required together, but only %q is/are present", strings.Join(r.fields, ", "), strings.Join(foundFields, ", "))
	}

	return nil
}

// ValidateSpec checks if the constraint is valid for the object type.
func (r *requiredTogetherGroup) ValidateSpec(t *objectType) error {
	for _, field := range r.fields {
		if _, ok := t.fields[field]; !ok {
			return fmt.Errorf("field %q is not defined in the object type", field)
		}
	}
	return nil
}

// requiredOneOfGroup represents a group of fields of which a minimum of one is required.
type requiredOneOfGroup struct {
	fields []string // List of field names of which at least one is required.
}

// RequiredOneOf creates a new requiredOneOfGroup with the given field names.
func RequiredOneOf(fields ...string) objectConstraint {
	return &requiredOneOfGroup{fields: fields}
}

// Validate checks if at least one of the required fields is present in the given value.
func (r *requiredOneOfGroup) Validate(values map[string]cty.Value) error {
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

// ValidateSpec checks if the constraint is valid for the object type.
func (r *requiredOneOfGroup) ValidateSpec(t *objectType) error {
	for _, field := range r.fields {
		if _, ok := t.fields[field]; !ok {
			return fmt.Errorf("field %q is not defined in the object type", field)
		}
	}
	return nil
}

// ObjectCondition represents a condition that can be applied to an object.
type ObjectCondition interface {
	IsMet(values map[string]cty.Value) bool // IsMet checks if the condition is satisfied by the given values.
	Description() string                    // Description provides a human-readable description of the condition for error messages.
	ValidateSpec(t *objectType) error       // ValidateSpec checks if the condition is valid for the object type.
}

// fieldIsNotNullCondition represents a condition that checks if a specific field is not null.
type fieldIsNotNullCondition struct {
	field string // The name of the field to check.
}

func FieldIsNotNull(field string) ObjectCondition {
	return &fieldIsNotNullCondition{field: field}
}

// IsMet checks if the condition is satisfied by the given values.
func (c *fieldIsNotNullCondition) IsMet(values map[string]cty.Value) bool {
	if value, ok := values[c.field]; ok {
		return value.IsWhollyKnown() && !value.IsNull()
	}

	return false
}

// Description provides a human-readable description of the condition for error messages.
func (c *fieldIsNotNullCondition) Description() string {
	return fmt.Sprintf("field %q is not null", c.field)
}

// ValidateSpec checks if the condition is valid for the object type.
func (c *fieldIsNotNullCondition) ValidateSpec(t *objectType) error {
	if _, ok := t.fields[c.field]; !ok {
		return fmt.Errorf("field %q is not defined in the object type", c.field)
	}
	return nil
}

// fieldIsNullCondition represents a condition that checks if a specific field is null.
type fieldIsNullCondition struct {
	field string // The name of the field to check.
}

// FieldIsNull creates a new fieldIsNullCondition for the given field.
func FieldIsNull(field string) ObjectCondition {
	return &fieldIsNullCondition{field: field}
}

// IsMet checks if the condition is satisfied by the given values.
func (c *fieldIsNullCondition) IsMet(values map[string]cty.Value) bool {
	if value, ok := values[c.field]; ok {
		return value.IsWhollyKnown() && value.IsNull()
	}

	return false
}

// Description provides a human-readable description of the condition for error messages.
func (c *fieldIsNullCondition) Description() string {
	return fmt.Sprintf("field %q is null", c.field)
}

// ValidateSpec checks if the condition is valid for the object type.
func (c *fieldIsNullCondition) ValidateSpec(t *objectType) error {
	if _, ok := t.fields[c.field]; !ok {
		return fmt.Errorf("field %q is not defined in the object type", c.field)
	}
	return nil
}

// fieldEqualsCondition represents a condition that checks if a specific field is equal to a given value.
type fieldEqualsCondition struct {
	field string    // The name of the field to check.
	value cty.Value // The value to compare against.
}

// FieldEquals creates a new fieldEqualsCondition for the given field and value.
func FieldEquals(field string, value cty.Value) ObjectCondition {
	return &fieldEqualsCondition{field: field, value: value}
}

// IsMet checks if the condition is satisfied by the given values.
func (c *fieldEqualsCondition) IsMet(values map[string]cty.Value) bool {
	if value, ok := values[c.field]; ok {
		return value.IsWhollyKnown() && value.Equals(c.value).True()
	}

	return false
}

// Description provides a human-readable description of the condition for error messages.
func (c *fieldEqualsCondition) Description() string {
	return fmt.Sprintf("field %q is equal to \"%v\"", c.field, c.value)
}

// ValidateSpec checks if the condition is valid for the object type.
func (c *fieldEqualsCondition) ValidateSpec(t *objectType) error {
	if _, ok := t.fields[c.field]; !ok {
		return fmt.Errorf("field %q is not defined in the object type", c.field)
	}
	return nil
}

// conditionalConstraint represents a constraint that is applied conditionally based on the evaluation of an ObjectCondition.
type conditionalConstraint struct {
	condition  ObjectCondition  // The condition that must be met for the constraint to apply.
	constraint objectConstraint // The specific constraint to apply if the condition is met.
}

// ConditionalConstraint creates a new conditionalConstraint for the given condition and constraint.
func ConditionalConstraint(condition ObjectCondition, constraint objectConstraint) objectConstraint {
	return &conditionalConstraint{
		condition:  condition,
		constraint: constraint,
	}
}

// Validate checks if the condition is satisfied by the given values and applies the constraint if it is.
func (c *conditionalConstraint) Validate(values map[string]cty.Value) error {
	if c.condition.IsMet(values) {
		err := c.constraint.Validate(values)
		if err != nil {
			return fmt.Errorf("conditional constraint failed: when %s, %w", c.condition.Description(), err)
		}
	}
	return nil
}

// ValidateSpec checks if the constraint is valid for the object type.
func (c *conditionalConstraint) ValidateSpec(t *objectType) error {
	if err := c.condition.ValidateSpec(t); err != nil {
		return err
	}
	return c.constraint.ValidateSpec(t)
}
