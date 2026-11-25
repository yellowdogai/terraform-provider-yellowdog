package provider

import (
	"reflect"
	"testing"

	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
)

func TestApplicationProducts(t *testing.T) {
	tests := map[string]struct {
		terraformVersion string
		want             []yellowdog.Product
	}{
		"with Terraform version": {"1.9.0", []yellowdog.Product{
			{Name: "terraform-provider-yellowdog", Version: "0.2.0"},
			{Name: "terraform", Version: "1.9.0"},
		}},
		"without Terraform version": {"", []yellowdog.Product{
			{Name: "terraform-provider-yellowdog", Version: "0.2.0"},
		}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := applicationProducts("0.2.0", tt.terraformVersion); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applicationProducts() = %v, want %v", got, tt.want)
			}
		})
	}
}
