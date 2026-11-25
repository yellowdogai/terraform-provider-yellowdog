// Package conversiontest checks a resource model's hand-written conversions to and from its SDK type against the
// resource's schema, so an attribute a conversion misses fails a unit test, naming the attribute path.
//
// Attribute paths join names with "." and write an element of a list, set or map as "[*]", e.g.
// "constraints[*].numeric_attribute_constraint.min".
package conversiontest

import (
	"context"
	"hash/fnv"
	"math/big"
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Option changes what the checks cover.
type Option func(*schemaChecker)

// Exclude skips the attributes at these paths, and everything nested in them, in every check.
func Exclude(paths ...string) Option {
	return func(checker *schemaChecker) {
		checker.excluded = append(checker.excluded, paths...)
	}
}

// CheckToSDK fails for each configurable attribute, or id, whose value toSDK does not carry into the SDK value.
func CheckToSDK[M any, S any](
	t testing.TB,
	newResource func() resource.Resource,
	toSDK func(*M, context.Context, *diag.Diagnostics) S,
	options ...Option,
) {
	t.Helper()
	checker := newChecker[M](t, newResource, options)
	checked := map[string]bool{}
	for _, base := range checker.bases() {
		want := callToSDK(checker, toSDK, load[M](checker, checker.value(base, "")))
		for _, leaf := range checker.leaves {
			if checked[leaf] || checker.isExcluded(leaf) || !checker.isPresent(leaf, base) {
				continue
			}
			checked[leaf] = true
			got := callToSDK(checker, toSDK, load[M](checker, checker.value(base, leaf)))
			if reflect.DeepEqual(want, got) {
				t.Errorf("toSDK ignores %s", leaf)
			}
		}
	}
}

// CheckUpdateFromSDK fails for each attribute updateFromSDK leaves unset when given the zero SDK value.
func CheckUpdateFromSDK[M any, S any](
	t testing.TB,
	newResource func() resource.Resource,
	updateFromSDK func(*M, context.Context, S, *diag.Diagnostics),
	options ...Option,
) {
	t.Helper()
	checker := newChecker[M](t, newResource, options)
	var sdk S
	if sdkType := reflect.TypeFor[S](); sdkType.Kind() == reflect.Pointer {
		sdk = reflect.New(sdkType.Elem()).Interface().(S)
	}
	model := poisoned[M](checker)
	callUpdateFromSDK(checker, updateFromSDK, model, sdk)
	for _, unsetPath := range checker.unset(model) {
		t.Errorf("updateFromSDK does not set %s", unsetPath)
	}
}

// CheckRoundTrip fails for each attribute updateFromSDK leaves unset when given toSDK's result, and for each
// configurable attribute, or id, whose value does not survive the round trip.
func CheckRoundTrip[M any, S any](
	t testing.TB,
	newResource func() resource.Resource,
	toSDK func(*M, context.Context, *diag.Diagnostics) S,
	updateFromSDK func(*M, context.Context, S, *diag.Diagnostics),
	options ...Option,
) {
	t.Helper()
	checker := newChecker[M](t, newResource, options)
	reported := map[string]bool{}
	for _, base := range checker.bases() {
		want := checker.value(base, "")
		model := poisoned[M](checker)
		callUpdateFromSDK(checker, updateFromSDK, model, callToSDK(checker, toSDK, load[M](checker, want)))

		var failures []string
		unset := checker.unset(model)
		for _, unsetPath := range unset {
			failures = append(failures, "updateFromSDK does not set "+unsetPath)
		}
		for _, differentPath := range checker.notRoundTripped(want, model) {
			isUnset := func(unsetPath string) bool { return isWithin(differentPath, unsetPath) }
			if !slices.ContainsFunc(unset, isUnset) {
				failures = append(failures, differentPath+" does not round trip")
			}
		}
		for _, failure := range failures {
			if !reported[failure] {
				reported[failure] = true
				t.Errorf("%s", failure)
			}
		}
	}
}

// schemaChecker holds what the checks derive from a resource's schema.
type schemaChecker struct {
	t      testing.TB
	ctx    context.Context
	schema schema.Schema

	// paths of every attribute in the schema
	paths []string

	// paths of the primitive or collection attributes the checks compare, i.e. those configurable, or named id
	leaves []string

	// paths of the single nested attributes
	singleNested []string

	// alternatives, each the paths of the attributes an ExactlyOneOf validator allows exactly one of
	groups [][]string

	excluded []string
}

// newChecker fails the test if the schema or model M has a shape the checks would not fully cover.
func newChecker[M any](t testing.TB, newResource func() resource.Resource, options []Option) *schemaChecker {
	t.Helper()
	ctx := context.Background()
	var response resource.SchemaResponse
	newResource().Schema(ctx, resource.SchemaRequest{}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("schema: %v", response.Diagnostics)
	}
	checker := &schemaChecker{t: t, ctx: ctx, schema: response.Schema}
	for _, option := range options {
		option(checker)
	}
	checker.walkSchema(nil, response.Schema.Attributes)
	for _, excluded := range checker.excluded {
		if !slices.Contains(checker.paths, excluded) {
			t.Fatalf("Exclude(%q): no such attribute", excluded)
		}
	}
	checker.checkGroups()
	checker.checkModel("", reflect.TypeFor[M]())
	return checker
}

func (checker *schemaChecker) walkSchema(parent *path.Expression, attributes map[string]schema.Attribute) {
	for _, name := range sortedKeys(attributes) {
		attribute := attributes[name]
		expression := path.MatchRoot(name)
		if parent != nil {
			expression = parent.AtName(name)
		}
		attributePath := expressionString(expression)
		checker.paths = append(checker.paths, attributePath)
		checker.addExactlyOneOfGroups(expression, attribute)

		switch nested := attribute.(type) {
		case schema.SingleNestedAttribute:
			checker.singleNested = append(checker.singleNested, attributePath)
			checker.walkSchema(&expression, nested.Attributes)
		case schema.ListNestedAttribute:
			element := expression.AtAnyListIndex()
			checker.walkSchema(&element, nested.NestedObject.Attributes)
		case schema.SetNestedAttribute:
			// a change inside an element changes the element's identity, so a round trip diff could not name the leaf
			checker.t.Fatalf("%s: set nested attributes are not supported", attributePath)
		case schema.MapNestedAttribute:
			// no resource uses one, so none is tested
			checker.t.Fatalf("%s: map nested attributes are not supported", attributePath)
		default:
			if isConfigurable(attribute) || name == "id" {
				checker.leaves = append(checker.leaves, attributePath)
			}
		}
	}
}

// addExactlyOneOfGroups finds ExactlyOneOf validators by reflection, as their type is in an internal package.
func (checker *schemaChecker) addExactlyOneOfGroups(expression path.Expression, attribute schema.Attribute) {
	validators := reflect.ValueOf(attribute).FieldByName("Validators")
	if !validators.IsValid() {
		return
	}
	for index := range validators.Len() {
		validator := validators.Index(index).Elem()
		if validator.Kind() != reflect.Struct || validator.Type().Name() != "ExactlyOneOfValidator" {
			continue
		}
		members := []string{expressionString(expression)}
		for _, other := range validator.FieldByName("PathExpressions").Interface().(path.Expressions) {
			members = append(members, expressionString(expression.Merge(other).Resolve()))
		}
		sort.Strings(members)
		members = slices.Compact(members)
		isSame := func(group []string) bool { return slices.Equal(group, members) }
		if !slices.ContainsFunc(checker.groups, isSame) {
			checker.groups = append(checker.groups, members)
		}
	}
}

// checkGroups fails if a group is nested in an alternative other than the first of another group, as bases never
// chooses both.
func (checker *schemaChecker) checkGroups() {
	for _, outer := range checker.groups {
		for _, alternative := range outer[1:] {
			for _, inner := range checker.groups {
				for _, member := range inner {
					if member != alternative && isWithin(member, alternative) {
						checker.t.Fatalf("%s: alternatives nested in %s are not supported, as it is not the first alternative",
							member, alternative)
					}
				}
			}
		}
	}
}

// checkModel fails if a single nested attribute's model field is not a pointer to a struct, as unset only looks
// inside those.
func (checker *schemaChecker) checkModel(prefix string, structType reflect.Type) {
	for index := range structType.NumField() {
		field := structType.Field(index)
		attributePath := joinPath(prefix, field.Tag.Get("tfsdk"))
		if !slices.Contains(checker.singleNested, attributePath) {
			continue
		}
		if field.Type.Kind() != reflect.Pointer || field.Type.Elem().Kind() != reflect.Struct {
			checker.t.Fatalf("%s: model field %s is %s, not a pointer to a struct", attributePath, field.Name, field.Type)
		}
		checker.checkModel(attributePath, field.Type.Elem())
	}
}

// bases are the sets of chosen alternatives to build values with: one per alternative, with the first alternative of
// every other group.
func (checker *schemaChecker) bases() []map[string]bool {
	if len(checker.groups) == 0 {
		return []map[string]bool{{}}
	}
	var bases []map[string]bool
	seen := map[string]bool{}
	for groupIndex, group := range checker.groups {
		for _, member := range group {
			base := map[string]bool{member: true}
			for otherIndex, other := range checker.groups {
				if otherIndex != groupIndex {
					base[other[0]] = true
				}
			}
			key := strings.Join(sortedKeys(base), ",")
			if !seen[key] {
				seen[key] = true
				bases = append(bases, base)
			}
		}
	}
	return bases
}

// isPresent is whether the leaf is set in a value built with the chosen alternatives.
func (checker *schemaChecker) isPresent(leaf string, chosen map[string]bool) bool {
	for _, group := range checker.groups {
		for _, member := range group {
			if !chosen[member] && isWithin(leaf, member) {
				return false
			}
		}
	}
	return true
}

func (checker *schemaChecker) isAlternative(attributePath string) bool {
	isMember := func(group []string) bool { return slices.Contains(group, attributePath) }
	return slices.ContainsFunc(checker.groups, isMember)
}

func (checker *schemaChecker) isExcluded(attributePath string) bool {
	isWithinExcluded := func(excludedPath string) bool { return isWithin(attributePath, excludedPath) }
	return slices.ContainsFunc(checker.excluded, isWithinExcluded)
}

// value builds the resource's value with the chosen alternatives set and the others null. Each leaf gets a value
// derived from its path, and the leaf at changed path gets a different one.
func (checker *schemaChecker) value(chosen map[string]bool, changed string) tftypes.Value {
	builder := valueBuilder{ctx: checker.ctx, chosen: chosen, changed: changed, isAlternative: checker.isAlternative}
	return builder.object(nil, checker.schema.Attributes, checker.schema.Type().TerraformType(checker.ctx))
}

// poisonedValue builds the resource's value with every nested object set, including every alternative, and
// everything else unknown.
func (checker *schemaChecker) poisonedValue() tftypes.Value {
	builder := valueBuilder{ctx: checker.ctx, poison: true, isAlternative: checker.isAlternative}
	return builder.object(nil, checker.schema.Attributes, checker.schema.Type().TerraformType(checker.ctx))
}

type valueBuilder struct {
	ctx           context.Context
	chosen        map[string]bool
	changed       string
	poison        bool
	isAlternative func(string) bool
}

func (builder *valueBuilder) object(
	parent *path.Expression, attributes map[string]schema.Attribute, objectType tftypes.Type,
) tftypes.Value {
	values := map[string]tftypes.Value{}
	for name, attribute := range attributes {
		expression := path.MatchRoot(name)
		if parent != nil {
			expression = parent.AtName(name)
		}
		values[name] = builder.attribute(expression, name, attribute)
	}
	return tftypes.NewValue(objectType, values)
}

func (builder *valueBuilder) attribute(
	expression path.Expression, name string, attribute schema.Attribute,
) tftypes.Value {
	attributePath := expressionString(expression)
	attributeType := attribute.GetType().TerraformType(builder.ctx)
	if !builder.poison && builder.isAlternative(attributePath) && !builder.chosen[attributePath] {
		return tftypes.NewValue(attributeType, nil)
	}
	if nested, ok := attribute.(schema.SingleNestedAttribute); ok {
		return builder.object(&expression, nested.Attributes, attributeType)
	}
	if builder.poison {
		return tftypes.NewValue(attributeType, tftypes.UnknownValue)
	}
	if nested, ok := attribute.(schema.ListNestedAttribute); ok {
		element := expression.AtAnyListIndex()
		elementType := attributeType.(tftypes.List).ElementType
		return tftypes.NewValue(attributeType, []tftypes.Value{
			builder.object(&element, nested.NestedObject.Attributes, elementType),
		})
	}
	if !isConfigurable(attribute) && name != "id" {
		return tftypes.NewValue(attributeType, nil)
	}
	return leafValue(attributePath, attributeType, attributePath == builder.changed)
}

func leafValue(attributePath string, valueType tftypes.Type, changed bool) tftypes.Value {
	switch collection := valueType.(type) {
	case tftypes.List:
		return tftypes.NewValue(collection, []tftypes.Value{leafValue(attributePath, collection.ElementType, changed)})
	case tftypes.Set:
		return tftypes.NewValue(collection, []tftypes.Value{leafValue(attributePath, collection.ElementType, changed)})
	case tftypes.Map:
		return tftypes.NewValue(collection, map[string]tftypes.Value{
			"key": leafValue(attributePath, collection.ElementType, changed),
		})
	}
	switch {
	case valueType.Is(tftypes.String):
		if changed {
			return tftypes.NewValue(valueType, attributePath+"-changed")
		}
		return tftypes.NewValue(valueType, attributePath)
	case valueType.Is(tftypes.Bool):
		return tftypes.NewValue(valueType, !changed)
	case valueType.Is(tftypes.Number):
		// whole numbers, so they suit an attribute of any number type
		hash := fnv.New32a()
		_, _ = hash.Write([]byte(attributePath))
		number := int64(hash.Sum32()%1000) + 1
		if changed {
			number++
		}
		return tftypes.NewValue(valueType, new(big.Float).SetInt64(number))
	default:
		panic("conversiontest: unsupported attribute type " + valueType.String() + " at " + attributePath)
	}
}

func load[M any](checker *schemaChecker, value tftypes.Value) *M {
	checker.t.Helper()
	model := new(M)
	if diags := (tfsdk.Plan{Schema: checker.schema, Raw: value}).Get(checker.ctx, model); diags.HasError() {
		checker.t.Fatalf("load model: %v", diags)
	}
	return model
}

func poisoned[M any](checker *schemaChecker) *M {
	return load[M](checker, checker.poisonedValue())
}

func callToSDK[M any, S any](
	checker *schemaChecker, toSDK func(*M, context.Context, *diag.Diagnostics) S, model *M,
) S {
	checker.t.Helper()
	var diags diag.Diagnostics
	sdk := toSDK(model, checker.ctx, &diags)
	if diags.HasError() {
		checker.t.Fatalf("toSDK: %v", diags)
	}
	return sdk
}

func callUpdateFromSDK[M any, S any](
	checker *schemaChecker, updateFromSDK func(*M, context.Context, S, *diag.Diagnostics), model *M, sdk S,
) {
	checker.t.Helper()
	var diags diag.Diagnostics
	updateFromSDK(model, checker.ctx, sdk, &diags)
	if diags.HasError() {
		checker.t.Fatalf("updateFromSDK: %v", diags)
	}
}

// unset are the paths of the model's attributes still unknown from poisoning. A nested object left entirely
// unknown is reported by its own path.
func (checker *schemaChecker) unset(model any) []string {
	unset, _ := checker.unsetFields("", reflect.ValueOf(model).Elem())
	return unset
}

// unsetFields returns the unset attribute paths in the struct, and how many attributes it holds.
func (checker *schemaChecker) unsetFields(prefix string, structValue reflect.Value) ([]string, int) {
	var unset []string
	total := 0
	for index := range structValue.NumField() {
		attributePath := joinPath(prefix, structValue.Type().Field(index).Tag.Get("tfsdk"))
		if checker.isExcluded(attributePath) {
			continue
		}
		field := structValue.Field(index)
		if value, ok := field.Interface().(attr.Value); ok {
			total++
			if value.IsUnknown() {
				unset = append(unset, attributePath)
			}
			continue
		}
		if field.Kind() == reflect.Pointer && !field.IsNil() {
			nestedUnset, nestedTotal := checker.unsetFields(attributePath, field.Elem())
			total += nestedTotal
			if nestedTotal > 0 && len(nestedUnset) == nestedTotal {
				unset = append(unset, attributePath)
			} else {
				unset = append(unset, nestedUnset...)
			}
		}
	}
	return unset, total
}

// notRoundTripped compares the model with the value it came from, returning the paths of the compared leaves that
// differ, and of the nested objects that are null on only one side.
func (checker *schemaChecker) notRoundTripped(want tftypes.Value, model any) []string {
	checker.t.Helper()
	state := tfsdk.State{Schema: checker.schema}
	if diags := state.Set(checker.ctx, model); diags.HasError() {
		checker.t.Fatalf("write model to state: %v", diags)
	}
	diffs, err := want.Diff(state.Raw)
	if err != nil {
		checker.t.Fatalf("diff: %v", err)
	}

	var objects, leaves []string
	for _, diff := range diffs {
		diffPath := attributePathString(diff.Path)
		holdsDiff := func(leaf string) bool { return isWithin(diffPath, leaf) }
		if leafIndex := slices.IndexFunc(checker.leaves, holdsDiff); leafIndex >= 0 {
			if !checker.isExcluded(checker.leaves[leafIndex]) {
				leaves = append(leaves, checker.leaves[leafIndex])
			}
			continue
		}
		isComparedLeafWithin := func(leaf string) bool { return isWithin(leaf, diffPath) && !checker.isExcluded(leaf) }
		if isNull(diff.Value1) != isNull(diff.Value2) && slices.ContainsFunc(checker.leaves, isComparedLeafWithin) {
			objects = append(objects, diffPath)
		}
	}

	result := slices.Clone(objects)
	for _, leaf := range leaves {
		isWithinObject := func(object string) bool { return isWithin(leaf, object) }
		if !slices.ContainsFunc(objects, isWithinObject) {
			result = append(result, leaf)
		}
	}
	sort.Strings(result)
	return slices.Compact(result)
}

func isNull(value *tftypes.Value) bool {
	return value == nil || value.IsNull()
}

func isConfigurable(attribute schema.Attribute) bool {
	return attribute.IsRequired() || attribute.IsOptional()
}

// isWithin is whether attributePath is the path ancestor, or nested in it.
func isWithin(attributePath, ancestor string) bool {
	return attributePath == ancestor ||
		strings.HasPrefix(attributePath, ancestor+".") ||
		strings.HasPrefix(attributePath, ancestor+"[")
}

func joinPath(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

func expressionString(expression path.Expression) string {
	var result strings.Builder
	for _, step := range expression.Steps() {
		if name, ok := step.(path.ExpressionStepAttributeNameExact); ok {
			if result.Len() > 0 {
				result.WriteString(".")
			}
			result.WriteString(string(name))
		} else {
			result.WriteString("[*]")
		}
	}
	return result.String()
}

func attributePathString(attributePath *tftypes.AttributePath) string {
	var result strings.Builder
	for _, step := range attributePath.Steps() {
		if name, ok := step.(tftypes.AttributeName); ok {
			if result.Len() > 0 {
				result.WriteString(".")
			}
			result.WriteString(string(name))
		} else {
			result.WriteString("[*]")
		}
	}
	return result.String()
}

func sortedKeys[V any](items map[string]V) []string {
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
