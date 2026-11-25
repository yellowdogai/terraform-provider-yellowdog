// Package validators holds schema validators that the generated schemas call.
package validators

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ydidPattern is the platform's YdidValidator.YDID_PATTERN: a type code, up to two optional 6 hex digit segments, a
// UUID, then any number of :<n> suffixes.
var ydidPattern = regexp.MustCompile(`^ydid:([a-z]+):(?:[0-9A-F]{6}:)?(?:[0-9A-F]{6}:)?[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}(?::[1-9][0-9]*)*$`)

// Ydid returns a validator that the value is a YellowDog ID (ydid). With no type codes, a ydid of any type passes;
// otherwise its type code, e.g. "role" in ydid:role:..., must be one of them.
func Ydid(typeCodes ...string) validator.String {
	return ydidValidator{typeCodes: typeCodes}
}

type ydidValidator struct {
	typeCodes []string
}

var _ validator.String = ydidValidator{}

func (v ydidValidator) Description(_ context.Context) string {
	return "must be " + v.expected()
}

func (v ydidValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ydidValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Often unknown at plan time, as it comes from another resource's id
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	// The platform lets a blank value through, leaving it to any other constraint on the attribute
	if strings.TrimSpace(value) == "" {
		return
	}
	match := ydidPattern.FindStringSubmatch(value)
	if match == nil || (len(v.typeCodes) > 0 && !slices.Contains(v.typeCodes, match[1])) {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid YellowDog ID", fmt.Sprintf("must be %s, got %q", v.expected(), value))
	}
}

// expected describes the ids that pass, e.g. "a role ydid (ydid:role:...)"
func (v ydidValidator) expected() string {
	switch len(v.typeCodes) {
	case 0:
		return "a YellowDog ID (ydid:...)"
	case 1:
		return fmt.Sprintf("a %s ydid (ydid:%s:...)", v.typeCodes[0], v.typeCodes[0])
	default:
		examples := make([]string, len(v.typeCodes))
		for i, code := range v.typeCodes {
			examples[i] = "ydid:" + code + ":..."
		}
		return fmt.Sprintf("a %s ydid (%s)", strings.Join(v.typeCodes, " or "), strings.Join(examples, ", "))
	}
}
