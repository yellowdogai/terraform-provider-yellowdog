package resourceerror

import (
	"errors"
	"fmt"
	"net/http"

	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
)

// ErrNotFound indicates a resource could not be found, see IsNotFound
var ErrNotFound = errors.New("resource not found")

// NewNotFoundError creates an error that will match IsNotFound
func NewNotFoundError(id string, resourceType string) error {
	return fmt.Errorf("%s %s: %w", resourceType, id, ErrNotFound)
}

// IsNotFound whether err is a yellowdog.APIError 404, or a ErrNotFound.
func IsNotFound(err error) bool {
	var apiErr *yellowdog.APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return errors.Is(err, ErrNotFound)
}
