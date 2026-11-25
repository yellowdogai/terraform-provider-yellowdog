package sweep_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/resourceerror"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/testutil"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
)

// sweepPrefix is what a name made by testutil.RandomName starts with.
const sweepPrefix = testutil.ResourcePrefix + "-"

// platformClient reads the same environment as the provider, so it sweeps the
// account the acceptance tests ran against.
var platformClient = sync.OnceValues(func() (*yellowdog.PlatformClient, error) {
	return yellowdog.NewClient(testutil.ConfigFromEnv())
})

func createSweeperFunc[T any](
	getAll func(ctx context.Context, client *yellowdog.PlatformClient) ([]T, error),
	getName func(v T) string,
	delete func(ctx context.Context, client *yellowdog.PlatformClient, v T) error,
) resource.SweeperFunc {
	return func(_ string) error {
		client, err := platformClient()
		if err != nil {
			return err
		}
		ctx := context.Background()

		// the API search narrows the candidates, the prefix check decides
		candidates, err := getAll(ctx, client)
		if err != nil {
			return fmt.Errorf("listing %T: %w", *new(T), err)
		}

		var errs []error
		for _, candidate := range candidates {
			name := getName(candidate)
			if !strings.HasPrefix(name, sweepPrefix) {
				continue
			}
			if err := delete(ctx, client, candidate); err != nil && !resourceerror.IsNotFound(err) {
				errs = append(errs, fmt.Errorf("deleting %T %s: %w", candidate, name, err))
			}
		}
		return errors.Join(errs...)
	}
}
