package testutil

import (
	"context"
	"testing"

	"github.com/yellowdogai/yellowdog-sdk-go/account"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// Group creates a group directly through the SDK and returns its ID, deleting it when the test finishes.
func Group(t *testing.T) string {
	t.Helper()
	skipUnlessAcceptance(t)

	client := testAccClient(t)

	group, err := client.Account.AddGroup(context.Background(), &account.AddGroupRequest{
		Name:        RandomName(),
		Description: optional.Of("fixture created by the acceptance suite"),
	})
	if err != nil {
		t.Fatalf("Failed to create group fixture: %v", err)
	}
	id, ok := group.Id.Get()
	if !ok {
		t.Fatal("Group fixture was created without an ID")
	}

	t.Cleanup(func() {
		if err := client.Account.DeleteGroup(context.Background(), id); err != nil {
			t.Errorf("Failed to delete group fixture %s: %v", id, err)
		}
	})
	return id
}

// Application creates an application directly through the SDK and returns its ID, deleting it when the test finishes.
func Application(t *testing.T) string {
	t.Helper()
	skipUnlessAcceptance(t)

	client := testAccClient(t)

	resp, err := client.Account.AddApplication(context.Background(), &account.AddApplicationRequest{
		Name:        RandomName(),
		Description: optional.Of("fixture created by the acceptance suite"),
	})
	if err != nil {
		t.Fatalf("Failed to create application fixture: %v", err)
	}
	application, ok := resp.Application.Get()
	if !ok {
		t.Fatal("Application fixture was created without an application")
	}
	id, ok := application.Id.Get()
	if !ok {
		t.Fatal("Application fixture was created without an ID")
	}

	t.Cleanup(func() {
		if err := client.Account.DeleteApplication(context.Background(), id); err != nil {
			t.Errorf("Failed to delete application fixture %s: %v", id, err)
		}
	})
	return id
}
