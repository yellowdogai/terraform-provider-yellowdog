package testutil

import (
	"context"
	"testing"

	"github.com/yellowdogai/yellowdog-sdk-go/model"
)

// MachineImageFamily creates a private machine image family in the shared namespace directly through the SDK and
// returns its ID, deleting it, and anything left in it, when the test finishes.
func MachineImageFamily(t *testing.T, osType model.ImageOsType) string {
	t.Helper()
	skipUnlessAcceptance(t)

	client := testAccClient(t)

	family, err := client.Images.AddImageFamily(context.Background(), &model.MachineImageFamily{
		Namespace: SharedNamespace(t),
		Name:      RandomName(),
		OsType:    osType,
		Access:    model.ImageAccessPrivate,
	})
	if err != nil {
		t.Fatalf("Failed to create machine image family fixture: %v", err)
	}
	if _, ok := family.Id.Get(); !ok {
		t.Fatal("Machine image family fixture was created without an ID")
	}

	t.Cleanup(func() {
		if err := client.Images.DeleteImageFamily(context.Background(), family); err != nil {
			t.Errorf("Failed to delete machine image family fixture %s: %v", family.Id.Or(""), err)
		}
	})
	return family.Id.Or("")
}

// MachineImageGroup creates a machine image group in a new MachineImageFamily directly through the SDK and returns its
// ID, deleting it when the test finishes.
func MachineImageGroup(t *testing.T, osType model.ImageOsType) string {
	t.Helper()
	skipUnlessAcceptance(t)

	client := testAccClient(t)

	family, err := client.Images.GetImageFamily(context.Background(), MachineImageFamily(t, osType))
	if err != nil {
		t.Fatalf("Failed to get machine image family fixture: %v", err)
	}
	group, err := client.Images.AddImageGroup(context.Background(), family, &model.MachineImageGroup{
		Name:   RandomName(),
		OsType: osType,
	})
	if err != nil {
		t.Fatalf("Failed to create machine image group fixture: %v", err)
	}
	if _, ok := group.Id.Get(); !ok {
		t.Fatal("Machine image group fixture was created without an ID")
	}

	t.Cleanup(func() {
		if err := client.Images.DeleteImageGroup(context.Background(), group); err != nil {
			t.Errorf("Failed to delete machine image group fixture %s: %v", group.Id.Or(""), err)
		}
	})
	return group.Id.Or("")
}
