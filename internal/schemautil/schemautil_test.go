package schemautil

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestValidateSchemaReportsAnAttributeWithNoMode(t *testing.T) {
	r := schemaResource{attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{Required: true},
		"source": schema.SingleNestedAttribute{Optional: true, Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{},
		}},
	}}

	diags := ValidateSchema(context.Background(), r)

	expectError(t, diags, "source.region must be Required, Optional or Computed")
}

func TestValidateSchemaReportsErrorsFromBuildingTheSchema(t *testing.T) {
	r := schemaResource{attributes: map[string]schema.Attribute{}, buildError: "cannot build"}

	diags := ValidateSchema(context.Background(), r)

	expectError(t, diags, "cannot build")
}

func TestValidateSchemaAcceptsAValidSchema(t *testing.T) {
	r := schemaResource{attributes: map[string]schema.Attribute{
		"id":   schema.StringAttribute{Computed: true},
		"name": schema.StringAttribute{Required: true},
	}}

	if diags := ValidateSchema(context.Background(), r); diags.HasError() {
		t.Error(diags)
	}
}

func expectError(t *testing.T, diags diag.Diagnostics, detail string) {
	t.Helper()
	for _, d := range diags.Errors() {
		if d.Detail() == detail {
			return
		}
	}
	t.Errorf("no error %q in %v", detail, diags)
}

// schemaResource is a resource whose Schema returns the given attributes, after reporting an error if asked.
type schemaResource struct {
	attributes map[string]schema.Attribute
	buildError string
}

func (r schemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	if r.buildError != "" {
		resp.Diagnostics.AddError("Invalid schema", r.buildError)
	}
	resp.Schema = schema.Schema{Attributes: r.attributes}
}
func (schemaResource) Metadata(context.Context, resource.MetadataRequest, *resource.MetadataResponse) {
}
func (schemaResource) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {}
func (schemaResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse)       {}
func (schemaResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}
func (schemaResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {}
