// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"
)

// ToType converts a SimpleTypePB to its Type representation.
func (p *SimpleTypePB) ToType() (Type, error) {
	if p == nil {
		return nil, errors.New("simple type protobuf is nil")
	}

	constraints := make(TypeConstraints, 0, len(p.Constraints))
	for _, c := range p.Constraints {
		if c == nil {
			continue
		}

		constraint, err := c.ToTypeConstraint()
		if err != nil {
			return nil, fmt.Errorf("failed to convert constraint: %w", err)
		}

		constraints = append(constraints, constraint)
	}

	switch p.Data {
	case SimpleTypeDataPB_BOOL:
		return Bool().WithConstraints(constraints...), nil
	case SimpleTypeDataPB_STRING:
		return String().WithConstraints(constraints...), nil
	case SimpleTypeDataPB_SENSITIVE_STRING:
		return SensitiveString().WithConstraints(constraints...), nil
	default:
		return nil, fmt.Errorf("unsupported simple type: %v", p.Data)
	}
}
