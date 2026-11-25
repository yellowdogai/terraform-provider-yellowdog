package resourceerror

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
)

func TestIsNotFound(t *testing.T) {
	tests := map[string]struct {
		err  error
		want bool
	}{
		"404":               {&yellowdog.APIError{StatusCode: http.StatusNotFound}, true},
		"wrapped 404":       {fmt.Errorf("get: %w", &yellowdog.APIError{StatusCode: http.StatusNotFound}), true},
		"500":               {&yellowdog.APIError{StatusCode: http.StatusInternalServerError}, false},
		"403":               {&yellowdog.APIError{StatusCode: http.StatusForbidden}, false},
		"transport failure": {errors.New("connection reset by peer"), false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestIsNotFound_WrappedErrNotFound_IsNotNotFound(t *testing.T) {

	// Given
	err := NewNotFoundError("aws-prod", "AwsCredential")

	// When
	isNotFound := IsNotFound(err)

	// Then
	if !isNotFound {
		t.Errorf("IsNotFound(%v) = false, want true", err)
	}
}
