// Package schemautil checks generated Terraform schemas.
package schemautil

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// ValidateSchema builds the resource's schema and returns every problem with it: errors reported while
// building it, problems the framework finds in it, and any attribute that is none of Required, Optional
// or Computed, which Terraform itself refuses to load.
func ValidateSchema(ctx context.Context, r resource.Resource) diag.Diagnostics {
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	diags := resp.Diagnostics
	diags.Append(resp.Schema.ValidateImplementation(ctx)...)
	requireMode(&diags, "", resp.Schema.Attributes)
	return diags
}

func requireMode(diags *diag.Diagnostics, prefix string, attributes map[string]schema.Attribute) {
	names := make([]string, 0, len(attributes))
	for name := range attributes {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		attribute := attributes[name]
		path := prefix + name
		if !attribute.IsRequired() && !attribute.IsOptional() && !attribute.IsComputed() {
			diags.AddError("Invalid schema attribute", fmt.Sprintf("%s must be Required, Optional or Computed", path))
		}
		switch nested := attribute.(type) {
		case schema.SingleNestedAttribute:
			requireMode(diags, path+".", nested.Attributes)
		case schema.ListNestedAttribute:
			requireMode(diags, path+".", nested.NestedObject.Attributes)
		case schema.SetNestedAttribute:
			requireMode(diags, path+".", nested.NestedObject.Attributes)
		case schema.MapNestedAttribute:
			requireMode(diags, path+".", nested.NestedObject.Attributes)
		}
	}
}
