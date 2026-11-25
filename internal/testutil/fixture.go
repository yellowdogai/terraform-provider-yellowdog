package testutil

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	yellowdog "github.com/yellowdogai/yellowdog-sdk-go"
	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// defaultTestAccNamespace is the namespace acceptance tests share when YELLOWDOG_NAMESPACE is not set.
const defaultTestAccNamespace = "tf-provider-test-static"

// sharedNamespaces resolves each shared namespace once per test binary, keyed by name.
var sharedNamespaces sync.Map

type sharedNamespace struct {
	once sync.Once
	err  error
}

// SharedNamespace returns a shared namespace for tests to use. The namespace will be created if it doesn't already
// exist. The namespace is not automatically deleted by the tests.
func SharedNamespace(t *testing.T) string {
	t.Helper()
	skipUnlessAcceptance(t)
	return ensureSharedNamespace(t, sharedNamespaceName())
}

// SharedAlternateNamespace returns a second shared namespace, for tests that move a resource between namespaces.
func SharedAlternateNamespace(t *testing.T) string {
	t.Helper()
	skipUnlessAcceptance(t)
	return ensureSharedNamespace(t, sharedNamespaceName()+"-alternative")
}

func sharedNamespaceName() string {
	return cmp.Or(os.Getenv("YELLOWDOG_NAMESPACE"), defaultTestAccNamespace)
}

func ensureSharedNamespace(t *testing.T, name string) string {
	t.Helper()

	v, _ := sharedNamespaces.LoadOrStore(name, &sharedNamespace{})
	ns := v.(*sharedNamespace)

	ns.once.Do(func() {
		client, err := platformClient()
		if err != nil {
			ns.err = err
			return
		}

		// GetNamespace takes an ID, not a name, so look the namespace up by name.
		namespaces, err := client.Namespaces.GetNamespaces(&account.NamespaceSearch{
			Namespace: optional.Of(name),
		}).ListAll(context.Background())
		if err != nil {
			ns.err = fmt.Errorf("could not get namespace %s: %w", name, err)
			return
		}
		if len(namespaces) == 0 {
			if _, err = client.Namespaces.CreateNamespace(context.Background(), &account.CreateNamespaceRequest{
				Namespace: name,
			}); err != nil {
				var apiErr *yellowdog.APIError
				if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusConflict {
					// Another package's test binary may have created it since the lookup. The sync.Once only guards
					// against other tests in this binary.
					return
				}
				ns.err = fmt.Errorf("could not create namespace %s: %w", name, err)
				return
			}
		}
	})

	if ns.err != nil {
		t.Fatalf("Failed to resolve test namespace: %v", ns.err)
	}
	return name
}

// skipUnlessAcceptance skips a test that does live work before resource.Test is
// reached. resource.Test performs that check itself, but only once called, so a
// helper that talks to the API while building a config -- SharedNamespace is
// the one -- would otherwise fail the credential-free `make test` run rather
// than skip it.
func skipUnlessAcceptance(t *testing.T) {
	t.Helper()

	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless %s is set", resource.EnvTfAcc)
	}
}
