package validators

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const uuid = "0a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d"

func TestYdid(t *testing.T) {
	tests := []struct {
		name      string
		typeCodes []string
		value     types.String
		wantError string
	}{
		{name: "valid", typeCodes: []string{"role"}, value: types.StringValue("ydid:role:" + uuid)},
		{name: "one hex segment", typeCodes: []string{"role"}, value: types.StringValue("ydid:role:000001:" + uuid)},
		{name: "two hex segments", typeCodes: []string{"role"}, value: types.StringValue("ydid:role:000001:ABCDEF:" + uuid)},
		{name: "index suffixes", typeCodes: []string{"task"}, value: types.StringValue("ydid:task:000001:" + uuid + ":1:23")},
		{name: "any type", value: types.StringValue("ydid:group:" + uuid)},
		{name: "one of several types", typeCodes: []string{"imgfam", "imggrp"}, value: types.StringValue("ydid:imggrp:" + uuid)},
		{name: "null", typeCodes: []string{"role"}, value: types.StringNull()},
		{name: "unknown", typeCodes: []string{"role"}, value: types.StringUnknown()},
		{name: "blank, as the platform allows", typeCodes: []string{"role"}, value: types.StringValue("")},
		{
			name:      "wrong type",
			typeCodes: []string{"role"},
			value:     types.StringValue("ydid:group:" + uuid),
			wantError: `must be a role ydid (ydid:role:...), got "ydid:group:` + uuid + `"`,
		},
		{
			name:      "none of several types",
			typeCodes: []string{"imgfam", "imggrp"},
			value:     types.StringValue("ydid:image:" + uuid),
			wantError: "must be a imgfam or imggrp ydid (ydid:imgfam:..., ydid:imggrp:...)",
		},
		{name: "not a ydid", value: types.StringValue("not-a-ydid"), wantError: `must be a YellowDog ID (ydid:...), got "not-a-ydid"`},
		{name: "malformed uuid", value: types.StringValue("ydid:role:0a1b2c3d"), wantError: "must be a YellowDog ID"},
		{name: "lower case hex segment", value: types.StringValue("ydid:role:abcdef:" + uuid), wantError: "must be a YellowDog ID"},
		{name: "upper case uuid", value: types.StringValue("ydid:role:" + strings.ToUpper(uuid)), wantError: "must be a YellowDog ID"},
		{name: "zero index suffix", value: types.StringValue("ydid:task:" + uuid + ":0"), wantError: "must be a YellowDog ID"},
		{name: "trailing text", value: types.StringValue("ydid:role:" + uuid + "x"), wantError: "must be a YellowDog ID"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := validator.StringRequest{Path: path.Root("role_id"), ConfigValue: test.value}
			var resp validator.StringResponse

			Ydid(test.typeCodes...).ValidateString(context.Background(), req, &resp)

			if test.wantError == "" {
				if resp.Diagnostics.HasError() {
					t.Fatalf("unexpected error: %v", resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.ErrorsCount() != 1 {
				t.Fatalf("want one error, got %v", resp.Diagnostics)
			}
			diagnostic := resp.Diagnostics.Errors()[0]
			if !strings.Contains(diagnostic.Detail(), test.wantError) {
				t.Errorf("error %q does not contain %q", diagnostic.Detail(), test.wantError)
			}
		})
	}
}

func TestYdidDescription(t *testing.T) {
	if got := Ydid("role").Description(context.Background()); got != "must be a role ydid (ydid:role:...)" {
		t.Errorf("got %q", got)
	}
}
