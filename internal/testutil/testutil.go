package testutil

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/provider"
	"github.com/yellowdogai/terraform-provider-yellowdog/internal/resourceerror"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
)

const (
	// FakeYellowdogProviderHcl can be used in unit tests that leverage the terraform-plugin-framework testing framework
	// but do not need a real environment, for example for PlanOnly steps.
	FakeYellowdogProviderHcl = `
provider "yellowdog" {
  api_key_id     = "test"
  api_key_secret = "test"
  api_url        = "http://127.0.0.1:1"
}`
)

// ProtoV6ProviderFactories are used to instantiate a provider during acceptance testing. The factory function
// will be invoked for every Terraform CLI command executed to create a provider server to which the CLI can reattach.
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"yellowdog": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func PreCheck(t *testing.T) {
	t.Helper()
	for _, name := range []string{provider.EnvAPIURL, provider.EnvAPIKeyID, provider.EnvAPIKeySecret} {
		if os.Getenv(name) == "" {
			t.Fatalf("%s must be set for acceptance tests (see .env.example, or run make env)", name)
		}
	}
}

var (
	testAccClientOnce sync.Once
	testAccClientVal  *yellowdog.PlatformClient
	testAccClientErr  error
)

// platformClient returns a platform client for the assertions Terraform
// cannot make itself -- chiefly CheckDestroy, which has to ask the API whether an
// entity really went away rather than trust the state file.
//
// It reads the same environment the provider does, so it talks to the same
// account, and it is built once for the whole package: tests run in parallel and
// each would otherwise open its own.
//
// It returns an error rather than taking a *testing.T because CheckDestroy is
// handed no *testing.T -- it reports failure by returning an error. testAccClient
// is the wrapper for everywhere else.
func platformClient() (*yellowdog.PlatformClient, error) {
	testAccClientOnce.Do(func() {
		testAccClientVal, testAccClientErr = yellowdog.NewClient(ConfigFromEnv())
	})
	return testAccClientVal, testAccClientErr
}

// ConfigFromEnv returns the SDK config for the environment that PreCheck requires.
func ConfigFromEnv() *yellowdog.Config {
	return yellowdog.DefaultConfig(
		yellowdog.ServicesSchema{DefaultURL: os.Getenv(provider.EnvAPIURL)},
		yellowdog.APIKey{ID: os.Getenv(provider.EnvAPIKeyID), Secret: yellowdog.Secret(os.Getenv(provider.EnvAPIKeySecret))},
	)
}

func testAccClient(t *testing.T) *yellowdog.PlatformClient {
	t.Helper()

	client, err := platformClient()
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	return client
}

// CheckDestroyed builds a CheckDestroy function for a resource type. get
// fetches a single entity by id and is expected to fail with a 404 once the
// entity is gone; every resource of that type in state is checked.
//
// A non-404 error fails rather than counting as proof of deletion. That is
// stricter than it looks, and deliberately so: each resource's Read uses the
// same isNotFound to decide an entity has disappeared out of band, so an entity
// that signals deletion some other way is a live bug in that Read -- one where
// Terraform would report an error instead of planning to recreate. This is
// where it surfaces.
func CheckDestroyed(
	resourceType string,
	get func(context.Context, *yellowdog.PlatformClient, string) error,
) func(*terraform.State) error {
	return CheckDestroyedFromState(resourceType,
		func(ctx context.Context, c *yellowdog.PlatformClient, is *terraform.InstanceState) error {
			return get(ctx, c, is.ID)
		})
}

// CheckDestroyedFromState is CheckDestroyed for a resource that has no ID of
// its own, such as an association between two entities. get is handed the
// resource's pre-destroy state, so it can look the entity up by its attributes.
func CheckDestroyedFromState(
	resourceType string,
	get func(context.Context, *yellowdog.PlatformClient, *terraform.InstanceState) error,
) func(*terraform.State) error {
	return func(state *terraform.State) error {
		client, err := platformClient()
		if err != nil {
			return fmt.Errorf("creating client to verify destroy: %w", err)
		}

		for name, rs := range state.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			switch err := get(context.Background(), client, rs.Primary); {
			case err == nil:
				return fmt.Errorf("%s: %s %s still exists", name, resourceType, rs.Primary.ID)
			case !resourceerror.IsNotFound(err):
				return fmt.Errorf("%s: unexpected error checking %s %s was destroyed: %w", name, resourceType, rs.Primary.ID, err)
			}
		}

		return nil
	}
}

// EnvOrSkip returns the values of the named environment variables, skipping the test if any is unset. It is for the
// cloud secrets a test needs, which the platform checks against the cloud, so that a partial .env runs a partial suite.
func EnvOrSkip(t *testing.T, names ...string) []string {
	t.Helper()
	skipUnlessAcceptance(t)

	values := make([]string, 0, len(names))
	for _, name := range names {
		value := os.Getenv(name)
		if value == "" {
			t.Skipf("Skipping: %s is not set", name)
		}
		values = append(values, value)
	}
	return values
}

// CheckSensitiveResourceAttr is resource.TestCheckResourceAttr for a sensitive attribute, i.e. where the value should
// not be revealed in the error message.
func CheckSensitiveResourceAttr(name, key, value string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(name, key, func(actual string) error {
		if actual != value {
			return fmt.Errorf("%s: Attribute '%s' does not match the expected secret", name, key)
		}
		return nil
	})
}
