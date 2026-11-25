// Package convert bridges the SDK's optional.Optional[T] type and Go
// primitives to the terraform-plugin-framework attr.Value types
// (types.String, types.Int64, types.Bool, ...).
//
// The two type systems line up well: optional.Optional[T]'s zero value is
// "absent", exactly like a null attr.Value. Every resource/data source
// should go through these helpers rather than converting ad hoc, so that
// null/unknown handling stays consistent across the provider.
package convert

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// OptionalStringToTF converts an optional.Optional[string] to types.String,
// mapping an absent value to null.
func OptionalStringToTF(o optional.Optional[string]) types.String {
	if v, ok := o.Get(); ok {
		return types.StringValue(v)
	}
	return types.StringNull()
}

// TFToOptionalString converts a types.String to optional.Optional[string],
// mapping null and unknown to absent.
func TFToOptionalString(v types.String) optional.Optional[string] {
	if v.IsNull() || v.IsUnknown() {
		return optional.Optional[string]{}
	}
	return optional.Of(v.ValueString())
}

// OptionalEnumToTF converts an optional SDK enum -- any defined string type, e.g.
// model.CloudProvider -- to types.String, mapping an absent value to null.
func OptionalEnumToTF[E ~string](o optional.Optional[E]) types.String {
	if v, ok := o.Get(); ok {
		return types.StringValue(string(v))
	}
	return types.StringNull()
}

// TFToOptionalEnum converts a types.String to an optional SDK enum, mapping
// null and unknown to absent.
func TFToOptionalEnum[E ~string](v types.String) optional.Optional[E] {
	if v.IsNull() || v.IsUnknown() {
		return optional.Optional[E]{}
	}
	return optional.Of(E(v.ValueString()))
}

// OptionalInt64ToTF converts an optional.Optional[int] to types.Int64.
func OptionalInt64ToTF(o optional.Optional[int]) types.Int64 {
	if v, ok := o.Get(); ok {
		return types.Int64Value(int64(v))
	}
	return types.Int64Null()
}

// TFToOptionalInt64 converts a types.Int64 to optional.Optional[int].
func TFToOptionalInt64(v types.Int64) optional.Optional[int] {
	if v.IsNull() || v.IsUnknown() {
		return optional.Optional[int]{}
	}
	return optional.Of(int(v.ValueInt64()))
}

// OptionalInt32ToTF converts an optional.Optional[int] to types.Int32.
func OptionalInt32ToTF(o optional.Optional[int]) types.Int32 {
	if v, ok := o.Get(); ok {
		return types.Int32Value(int32(v))
	}
	return types.Int32Null()
}

// TFToOptionalInt32 converts a types.Int32 to optional.Optional[int].
func TFToOptionalInt32(v types.Int32) optional.Optional[int] {
	if v.IsNull() || v.IsUnknown() {
		return optional.Optional[int]{}
	}
	return optional.Of(int(v.ValueInt32()))
}

// OptionalFloat64ToTF converts an optional.Optional[float64] to types.Float64.
func OptionalFloat64ToTF(o optional.Optional[float64]) types.Float64 {
	if v, ok := o.Get(); ok {
		return types.Float64Value(v)
	}
	return types.Float64Null()
}

// TFToOptionalFloat64 converts a types.Float64 to optional.Optional[float64].
func TFToOptionalFloat64(v types.Float64) optional.Optional[float64] {
	if v.IsNull() || v.IsUnknown() {
		return optional.Optional[float64]{}
	}
	return optional.Of(v.ValueFloat64())
}

// OptionalFloat32ToTF converts an optional.Optional[float32] to types.Float32.
func OptionalFloat32ToTF(o optional.Optional[float32]) types.Float32 {
	if v, ok := o.Get(); ok {
		return types.Float32Value(v)
	}
	return types.Float32Null()
}

// TFToOptionalFloat32 converts a types.Float32 to optional.Optional[float32].
func TFToOptionalFloat32(v types.Float32) optional.Optional[float32] {
	if v.IsNull() || v.IsUnknown() {
		return optional.Optional[float32]{}
	}
	return optional.Of(v.ValueFloat32())
}

// OptionalBoolToTF converts an optional.Optional[bool] to types.Bool.
func OptionalBoolToTF(o optional.Optional[bool]) types.Bool {
	if v, ok := o.Get(); ok {
		return types.BoolValue(v)
	}
	return types.BoolNull()
}

// TFToOptionalBool converts a types.Bool to optional.Optional[bool].
func TFToOptionalBool(v types.Bool) optional.Optional[bool] {
	if v.IsNull() || v.IsUnknown() {
		return optional.Optional[bool]{}
	}
	return optional.Of(v.ValueBool())
}

// StringSliceToTFList converts a []string to a types.List.
func StringSliceToTFList(values []string) types.List {
	if len(values) == 0 {
		return types.ListNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(values))
	for _, s := range values {
		elements = append(elements, types.StringValue(s))
	}
	return types.ListValueMust(types.StringType, elements)
}

// TFListToStringSlice converts a []types.String to a []string.
func TFListToStringSlice(values types.List) []string {
	if values.IsNull() || values.IsUnknown() {
		return nil
	}
	elements := make([]string, 0, len(values.Elements()))
	for _, value := range values.Elements() {
		elements = append(elements, value.(types.String).ValueString())
	}
	return elements
}

// StringSliceToTFSet converts a []string to a types.Set.
func StringSliceToTFSet(values []string) types.Set {
	if len(values) == 0 {
		return types.SetNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(values))
	for _, s := range values {
		elements = append(elements, types.StringValue(s))
	}
	return types.SetValueMust(types.StringType, elements)
}

// TFSetToStringSlice converts a []types.String to a []string.
func TFSetToStringSlice(values types.Set) []string {
	if values.IsNull() || values.IsUnknown() {
		return nil
	}
	elements := make([]string, 0, len(values.Elements()))
	for _, value := range values.Elements() {
		elements = append(elements, value.(types.String).ValueString())
	}
	return elements
}

// StringMapToTF converts a map[string]string to a types.Map.
func StringMapToTF(m map[string]string) types.Map {
	if len(m) == 0 {
		return types.MapNull(types.StringType)
	}
	out := make(map[string]attr.Value, len(m))
	for k, v := range m {
		out[k] = types.StringValue(v)
	}
	return types.MapValueMust(types.StringType, out)
}

// TFToStringMap converts a types.Map to a map[string]string.
func TFToStringMap(m types.Map) map[string]string {
	out := make(map[string]string, len(m.Elements()))
	for k, v := range m.Elements() {
		out[k] = v.(types.String).ValueString()
	}
	return out
}

// TFToOptionalObjectByMapping converts a nested object to an optional.Optional, if the terraform object is present an
// optional containing the mapped value is returned, otherwise nil is returned.
func TFToOptionalObjectByMapping[M any, S any](m *M, mapper func(*M) S) optional.Optional[S] {
	if m == nil {
		return optional.Optional[S]{}
	}
	return optional.Of(mapper(m))
}

// OptionalObjectToTFByMapping converts an optional SDK object to a new nested object. If the SDK object is present it
// is mapped by the mapper function, otherwise the value is nil.
func OptionalObjectToTFByMapping[S any, M any](o optional.Optional[S], mapper func(*M, *S)) *M {
	v, ok := o.Get()
	if !ok {
		return nil
	}
	m := new(M)
	mapper(m, &v)
	return m
}

// TFListToSlice converts a list of nested objects with toSDK.
// M is the model type of the elements in the TF types.List
// S is the Slice element type
func TFListToSlice[M any, S any](ctx context.Context, list types.List, toSDK func(*M) S, diags *diag.Diagnostics) []S {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var sources []M
	diags.Append(list.ElementsAs(ctx, &sources, false)...)

	var result []S
	for _, element := range sources {
		result = append(result, toSDK(&element))
	}
	return result
}

// SliceToTFList converts a slice of SDK objects to nested objects.
// S is the Slice element type
// M is the model type of the elements in the returned types.List
func SliceToTFList[S any, M any](ctx context.Context, slice []S, fromSDK func(*S, *diag.Diagnostics) M, elementType attr.Type, diags *diag.Diagnostics) types.List {
	if len(slice) == 0 {
		return types.ListNull(elementType)
	}

	sources := make([]M, 0, len(slice))
	for _, element := range slice {
		sdk := fromSDK(&element, diags)
		sources = append(sources, sdk)
	}

	list, listDiags := types.ListValueFrom(ctx, elementType, sources)
	diags.Append(listDiags...)
	return list
}

// EnumToTF converts an SDK enum -- any defined string type, e.g. model.CloudProvider -- to types.String
func EnumToTF[E ~string](o E) types.String {
	return types.StringValue(string(o))
}

// TFToEnum converts a types.String to an SDK enum
func TFToEnum[E ~string](v types.String) E {
	return E(v.ValueString())
}
