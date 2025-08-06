package hclspec

import (
	"github.com/hashicorp/hcl/v2/hcldec"
)

// SensitiveSpec wraps an hcldec.Spec to indicate that it has a sensitive value.
type SensitiveSpec struct {
	hcldec.Spec
}

// IsSensitive checks if the given spec is a SensitiveSpec.
func IsSensitive(spec hcldec.Spec) bool {
	_, isSensitive := spec.(*SensitiveSpec)
	return isSensitive
}
