package conversiontest

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// fakeResource has one attribute of each kind the checks treat differently. Each is named after its kind, except id,
// which the checks compare although it is computed.
type fakeResource struct{}

func (fakeResource) Metadata(context.Context, resource.MetadataRequest, *resource.MetadataResponse) {}

func (fakeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":                             schema.StringAttribute{Computed: true},
		"some_user_required_string":      schema.StringAttribute{Required: true},
		"some_user_optional_int":         schema.Int64Attribute{Optional: true},
		"some_computed_string":           schema.StringAttribute{Computed: true},
		"some_user_optional_string_list": schema.ListAttribute{Optional: true, ElementType: types.StringType},
		"some_user_optional_string_set":  schema.SetAttribute{Optional: true, ElementType: types.StringType},
		"some_user_optional_string_map":  schema.MapAttribute{Optional: true, ElementType: types.StringType},
		"some_user_optional_list_nested": schema.ListNestedAttribute{
			Optional:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: fakeElementAttributes()},
		},
		"some_user_required_object": schema.SingleNestedAttribute{Required: true, Attributes: map[string]schema.Attribute{
			"first_alternative": schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]schema.Attribute{"some_user_required_string": schema.StringAttribute{Required: true}},
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("second_alternative")),
				},
			},
			"second_alternative": schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]schema.Attribute{"some_user_required_string": schema.StringAttribute{Required: true}},
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("first_alternative")),
				},
			},
		}},
	}}
}

func (fakeResource) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {}

func (fakeResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {}

func (fakeResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}

func (fakeResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {}

func newFakeResource() resource.Resource {
	return fakeResource{}
}

func fakeElementAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{"some_user_required_string": schema.StringAttribute{Required: true}}
}

type fakeModel struct {
	Id                         types.String     `tfsdk:"id"`
	SomeUserRequiredString     types.String     `tfsdk:"some_user_required_string"`
	SomeUserOptionalInt        types.Int64      `tfsdk:"some_user_optional_int"`
	SomeComputedString         types.String     `tfsdk:"some_computed_string"`
	SomeUserOptionalStringList types.List       `tfsdk:"some_user_optional_string_list"`
	SomeUserOptionalStringSet  types.Set        `tfsdk:"some_user_optional_string_set"`
	SomeUserOptionalStringMap  types.Map        `tfsdk:"some_user_optional_string_map"`
	SomeUserOptionalListNested types.List       `tfsdk:"some_user_optional_list_nested"`
	SomeUserRequiredObject     *fakeObjectModel `tfsdk:"some_user_required_object"`
}

type fakeElementModel struct {
	SomeUserRequiredString types.String `tfsdk:"some_user_required_string"`
}

var fakeElementType = schema.NestedAttributeObject{Attributes: fakeElementAttributes()}.Type()

type fakeObjectModel struct {
	FirstAlternative  *fakeAlternativeModel `tfsdk:"first_alternative"`
	SecondAlternative *fakeAlternativeModel `tfsdk:"second_alternative"`
}

type fakeAlternativeModel struct {
	SomeUserRequiredString types.String `tfsdk:"some_user_required_string"`
}

type fakeSDK struct {
	Id                     string
	SomeUserRequiredString string
	SomeUserOptionalInt    int64
	SomeComputedString     string
	StringList             []string
	StringSet              []string
	StringMap              map[string]string
	ListNested             []string
	FirstAlternative       string
	SecondAlternative      string
}

// fakeConversions are complete unless a test leaves out the attribute at the path it names.
type fakeConversions struct {
	toSDKSkips         string
	updateFromSDKSkips string
}

func (conversions fakeConversions) toSDK(model *fakeModel, ctx context.Context, diags *diag.Diagnostics) *fakeSDK {
	sdk := &fakeSDK{Id: model.Id.ValueString(), SomeUserOptionalInt: model.SomeUserOptionalInt.ValueInt64()}
	if conversions.toSDKSkips != "some_user_required_string" {
		sdk.SomeUserRequiredString = model.SomeUserRequiredString.ValueString()
	}
	if conversions.toSDKSkips != "some_user_optional_string_list" {
		diags.Append(model.SomeUserOptionalStringList.ElementsAs(ctx, &sdk.StringList, false)...)
	}
	if conversions.toSDKSkips != "some_user_optional_string_set" {
		diags.Append(model.SomeUserOptionalStringSet.ElementsAs(ctx, &sdk.StringSet, false)...)
	}
	if conversions.toSDKSkips != "some_user_optional_string_map" {
		diags.Append(model.SomeUserOptionalStringMap.ElementsAs(ctx, &sdk.StringMap, false)...)
	}
	var elements []fakeElementModel
	diags.Append(model.SomeUserOptionalListNested.ElementsAs(ctx, &elements, false)...)
	for _, element := range elements {
		if conversions.toSDKSkips == "some_user_optional_list_nested[*].some_user_required_string" {
			sdk.ListNested = append(sdk.ListNested, "")
		} else {
			sdk.ListNested = append(sdk.ListNested, element.SomeUserRequiredString.ValueString())
		}
	}
	object := model.SomeUserRequiredObject
	if object.FirstAlternative != nil {
		sdk.FirstAlternative = object.FirstAlternative.SomeUserRequiredString.ValueString()
	}
	if object.SecondAlternative != nil &&
		conversions.toSDKSkips != "some_user_required_object.second_alternative.some_user_required_string" {
		sdk.SecondAlternative = object.SecondAlternative.SomeUserRequiredString.ValueString()
	}
	return sdk
}

func (conversions fakeConversions) updateFromSDK(
	model *fakeModel, ctx context.Context, sdk *fakeSDK, diags *diag.Diagnostics,
) {
	model.Id = types.StringValue(sdk.Id)
	model.SomeUserRequiredString = types.StringValue(sdk.SomeUserRequiredString)
	if conversions.updateFromSDKSkips != "some_user_optional_int" {
		model.SomeUserOptionalInt = types.Int64Value(sdk.SomeUserOptionalInt)
	}
	model.SomeComputedString = types.StringValue(sdk.SomeComputedString)
	var conversionDiags diag.Diagnostics
	model.SomeUserOptionalStringList, conversionDiags = types.ListValueFrom(ctx, types.StringType, sdk.StringList)
	diags.Append(conversionDiags...)
	model.SomeUserOptionalStringSet, conversionDiags = types.SetValueFrom(ctx, types.StringType, sdk.StringSet)
	diags.Append(conversionDiags...)
	model.SomeUserOptionalStringMap, conversionDiags = types.MapValueFrom(ctx, types.StringType, sdk.StringMap)
	diags.Append(conversionDiags...)
	if conversions.updateFromSDKSkips != "some_user_optional_list_nested" {
		var elements []fakeElementModel
		for _, value := range sdk.ListNested {
			element := fakeElementModel{SomeUserRequiredString: types.StringValue(value)}
			if conversions.updateFromSDKSkips == "some_user_optional_list_nested[*].some_user_required_string" {
				element.SomeUserRequiredString = types.StringNull()
			}
			elements = append(elements, element)
		}
		model.SomeUserOptionalListNested, conversionDiags = types.ListValueFrom(ctx, fakeElementType, elements)
		diags.Append(conversionDiags...)
	}
	model.SomeUserRequiredObject = &fakeObjectModel{}
	if sdk.FirstAlternative != "" {
		model.SomeUserRequiredObject.FirstAlternative = &fakeAlternativeModel{
			SomeUserRequiredString: types.StringValue(sdk.FirstAlternative),
		}
	}
	if sdk.SecondAlternative != "" {
		model.SomeUserRequiredObject.SecondAlternative = &fakeAlternativeModel{}
		if conversions.updateFromSDKSkips != "some_user_required_object.second_alternative.some_user_required_string" {
			model.SomeUserRequiredObject.SecondAlternative.SomeUserRequiredString = types.StringValue(sdk.SecondAlternative)
		}
	}
}

// failureRecorder collects the failures a check reports.
type failureRecorder struct {
	testing.TB
	failures []string
}

func (recorder *failureRecorder) Errorf(format string, args ...any) {
	recorder.failures = append(recorder.failures, fmt.Sprintf(format, args...))
}

func (recorder *failureRecorder) Fatalf(format string, args ...any) {
	recorder.Errorf(format, args...)
	panic(recorder)
}

// failures runs the check and returns what it reported.
func failures(t *testing.T, check func(testing.TB)) (failures []string) {
	recorder := &failureRecorder{TB: t}
	defer func() {
		if recovered := recover(); recovered != nil && recovered != recorder {
			panic(recovered)
		}
		failures = recorder.failures
	}()
	check(recorder)
	return
}

func TestChecks(t *testing.T) {
	for name, test := range map[string]struct {
		conversions fakeConversions
		check       func(testing.TB, fakeConversions)
		want        []string
	}{
		"toSDK complete": {
			check: func(t testing.TB, conversions fakeConversions) { CheckToSDK(t, newFakeResource, conversions.toSDK) },
		},
		"toSDK ignores an attribute": {
			conversions: fakeConversions{toSDKSkips: "some_user_required_string"},
			check:       func(t testing.TB, conversions fakeConversions) { CheckToSDK(t, newFakeResource, conversions.toSDK) },
			want:        []string{"toSDK ignores some_user_required_string"},
		},
		"toSDK ignores an attribute of an alternative": {
			conversions: fakeConversions{
				toSDKSkips: "some_user_required_object.second_alternative.some_user_required_string",
			},
			check: func(t testing.TB, conversions fakeConversions) { CheckToSDK(t, newFakeResource, conversions.toSDK) },
			want:  []string{"toSDK ignores some_user_required_object.second_alternative.some_user_required_string"},
		},
		"toSDK ignores a list": {
			conversions: fakeConversions{toSDKSkips: "some_user_optional_string_list"},
			check:       func(t testing.TB, conversions fakeConversions) { CheckToSDK(t, newFakeResource, conversions.toSDK) },
			want:        []string{"toSDK ignores some_user_optional_string_list"},
		},
		"toSDK ignores a set": {
			conversions: fakeConversions{toSDKSkips: "some_user_optional_string_set"},
			check:       func(t testing.TB, conversions fakeConversions) { CheckToSDK(t, newFakeResource, conversions.toSDK) },
			want:        []string{"toSDK ignores some_user_optional_string_set"},
		},
		"toSDK ignores a map": {
			conversions: fakeConversions{toSDKSkips: "some_user_optional_string_map"},
			check:       func(t testing.TB, conversions fakeConversions) { CheckToSDK(t, newFakeResource, conversions.toSDK) },
			want:        []string{"toSDK ignores some_user_optional_string_map"},
		},
		"toSDK ignores an attribute of a list element": {
			conversions: fakeConversions{toSDKSkips: "some_user_optional_list_nested[*].some_user_required_string"},
			check:       func(t testing.TB, conversions fakeConversions) { CheckToSDK(t, newFakeResource, conversions.toSDK) },
			want:        []string{"toSDK ignores some_user_optional_list_nested[*].some_user_required_string"},
		},
		"toSDK ignores an excluded attribute": {
			conversions: fakeConversions{toSDKSkips: "some_user_required_string"},
			check: func(t testing.TB, conversions fakeConversions) {
				CheckToSDK(t, newFakeResource, conversions.toSDK, Exclude("some_user_required_string"))
			},
		},
		"excluded attribute does not exist": {
			check: func(t testing.TB, conversions fakeConversions) {
				CheckToSDK(t, newFakeResource, conversions.toSDK, Exclude("missing"))
			},
			want: []string{`Exclude("missing"): no such attribute`},
		},
		"updateFromSDK complete": {
			check: func(t testing.TB, conversions fakeConversions) {
				CheckUpdateFromSDK(t, newFakeResource, conversions.updateFromSDK)
			},
		},
		"updateFromSDK does not set an attribute": {
			conversions: fakeConversions{updateFromSDKSkips: "some_user_optional_int"},
			check: func(t testing.TB, conversions fakeConversions) {
				CheckUpdateFromSDK(t, newFakeResource, conversions.updateFromSDK)
			},
			want: []string{"updateFromSDK does not set some_user_optional_int"},
		},
		"updateFromSDK does not set a list nested attribute": {
			conversions: fakeConversions{updateFromSDKSkips: "some_user_optional_list_nested"},
			check: func(t testing.TB, conversions fakeConversions) {
				CheckUpdateFromSDK(t, newFakeResource, conversions.updateFromSDK)
			},
			want: []string{"updateFromSDK does not set some_user_optional_list_nested"},
		},
		"round trip complete": {
			check: func(t testing.TB, conversions fakeConversions) {
				CheckRoundTrip(t, newFakeResource, conversions.toSDK, conversions.updateFromSDK)
			},
		},
		"round trip does not set an attribute": {
			conversions: fakeConversions{updateFromSDKSkips: "some_user_optional_int"},
			check: func(t testing.TB, conversions fakeConversions) {
				CheckRoundTrip(t, newFakeResource, conversions.toSDK, conversions.updateFromSDK)
			},
			want: []string{"updateFromSDK does not set some_user_optional_int"},
		},
		"round trip loses an attribute of a new nested object": {
			conversions: fakeConversions{
				updateFromSDKSkips: "some_user_required_object.second_alternative.some_user_required_string",
			},
			check: func(t testing.TB, conversions fakeConversions) {
				CheckRoundTrip(t, newFakeResource, conversions.toSDK, conversions.updateFromSDK)
			},
			want: []string{"some_user_required_object.second_alternative.some_user_required_string does not round trip"},
		},
		"round trip loses an attribute of a list element": {
			conversions: fakeConversions{
				updateFromSDKSkips: "some_user_optional_list_nested[*].some_user_required_string",
			},
			check: func(t testing.TB, conversions fakeConversions) {
				CheckRoundTrip(t, newFakeResource, conversions.toSDK, conversions.updateFromSDK)
			},
			want: []string{"some_user_optional_list_nested[*].some_user_required_string does not round trip"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := failures(t, func(recorder testing.TB) { test.check(recorder, test.conversions) })
			if !slices.Equal(got, test.want) {
				t.Errorf("got failures %q, want %q", got, test.want)
			}
		})
	}
}

// schemaResource is a fakeResource with the given attributes.
type schemaResource struct {
	fakeResource
	attributes map[string]schema.Attribute
}

func (fakeSchema schemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: fakeSchema.attributes}
}

// checkShape runs CheckUpdateFromSDK with model M and the attributes, so newChecker checks their shape.
func checkShape[M any](t testing.TB, attributes map[string]schema.Attribute) {
	newResource := func() resource.Resource { return schemaResource{attributes: attributes} }
	CheckUpdateFromSDK(t, newResource, func(*M, context.Context, *fakeSDK, *diag.Diagnostics) {})
}

func exactlyOneOf(sibling string) []validator.Object {
	return []validator.Object{objectvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName(sibling))}
}

func object(attributes map[string]schema.Attribute, validators ...validator.Object) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{Optional: true, Attributes: attributes, Validators: validators}
}

func TestUnsupportedShapes(t *testing.T) {
	optionalString := schema.StringAttribute{Optional: true}
	for name, test := range map[string]struct {
		check func(testing.TB)
		want  string
	}{
		"set nested attribute": {
			check: func(t testing.TB) {
				type model struct {
					SomeSetNested types.Set `tfsdk:"some_set_nested"`
				}
				checkShape[model](t, map[string]schema.Attribute{"some_set_nested": schema.SetNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{"some_user_optional_string": optionalString},
					},
				}})
			},
			want: "some_set_nested: set nested attributes are not supported",
		},
		"map nested attribute": {
			check: func(t testing.TB) {
				type model struct {
					SomeMapNested types.Map `tfsdk:"some_map_nested"`
				}
				checkShape[model](t, map[string]schema.Attribute{"some_map_nested": schema.MapNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{"some_user_optional_string": optionalString},
					},
				}})
			},
			want: "some_map_nested: map nested attributes are not supported",
		},
		"single nested attribute held as an object value": {
			check: func(t testing.TB) {
				type model struct {
					SomeSingleNested types.Object `tfsdk:"some_single_nested"`
				}
				checkShape[model](t, map[string]schema.Attribute{
					"some_single_nested": object(map[string]schema.Attribute{"some_user_optional_string": optionalString}),
				})
			},
			want: "some_single_nested: model field SomeSingleNested is basetypes.ObjectValue, not a pointer to a struct",
		},
		"alternatives nested in an alternative other than the first": {
			check: func(t testing.TB) {
				type firstAlternative struct {
					SomeUserOptionalString types.String `tfsdk:"some_user_optional_string"`
				}
				type secondAlternative struct {
					FirstNestedAlternative  types.String `tfsdk:"first_nested_alternative"`
					SecondNestedAlternative types.String `tfsdk:"second_nested_alternative"`
				}
				type model struct {
					FirstAlternative  *firstAlternative  `tfsdk:"first_alternative"`
					SecondAlternative *secondAlternative `tfsdk:"second_alternative"`
				}
				nestedAlternative := func(sibling string) schema.StringAttribute {
					return schema.StringAttribute{Optional: true, Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName(sibling)),
					}}
				}
				checkShape[model](t, map[string]schema.Attribute{
					"first_alternative": object(
						map[string]schema.Attribute{"some_user_optional_string": optionalString},
						exactlyOneOf("second_alternative")...,
					),
					"second_alternative": object(
						map[string]schema.Attribute{
							"first_nested_alternative":  nestedAlternative("second_nested_alternative"),
							"second_nested_alternative": nestedAlternative("first_nested_alternative"),
						},
						exactlyOneOf("first_alternative")...,
					),
				})
			},
			want: "second_alternative.first_nested_alternative: alternatives nested in second_alternative are not " +
				"supported, as it is not the first alternative",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := failures(t, test.check)
			if !slices.Equal(got, []string{test.want}) {
				t.Errorf("got failures %q, want %q", got, test.want)
			}
		})
	}
}
