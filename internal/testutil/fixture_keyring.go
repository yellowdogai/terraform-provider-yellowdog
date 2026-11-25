package testutil

import (
	"context"
	"testing"
)

// Keyring creates a keyring directly through the SDK and returns its name, deleting it, and the credentials in it,
// when the test finishes.
func Keyring(t *testing.T) string {
	t.Helper()
	skipUnlessAcceptance(t)

	client := testAccClient(t)

	name := RandomName()
	if _, err := client.Keyring.AddKeyring(context.Background(), name, "fixture created by the acceptance suite"); err != nil {
		t.Fatalf("Failed to create keyring fixture: %v", err)
	}

	t.Cleanup(func() {
		if err := client.Keyring.DeleteKeyringByName(context.Background(), name); err != nil {
			t.Errorf("Failed to delete keyring fixture %s: %v", name, err)
		}
	})
	return name
}
