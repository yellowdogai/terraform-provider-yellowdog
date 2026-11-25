package testutil

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/yellowdogai/yellowdog-sdk-go/model"
	"github.com/yellowdogai/yellowdog-sdk-go/optional"
)

// SimulatorSourceTemplate creates a simulator compute source template directly
// through the SDK and returns its ID, deleting it when the test finishes.
func SimulatorSourceTemplate(t *testing.T, namespace string) string {
	t.Helper()
	skipUnlessAcceptance(t)

	client := testAccClient(t)

	source := model.SimulatorComputeSource{}
	source.Name = acctest.RandomWithPrefix("tf-acc-sim-src")
	source.Region = "sim-region-1"
	source.InstanceType = "sim.small"
	source.ImageId = "sim-image-1"

	template, err := client.Compute.AddComputeSourceTemplate(context.Background(), &model.ComputeSourceTemplate{
		Namespace:   optional.Of(namespace),
		Description: optional.Of("fixture created by the acceptance suite"),
		Source:      &source,
	})
	if err != nil {
		t.Fatalf("Failed to create simulator source template fixture: %v", err)
	}
	id, ok := template.Id.Get()
	if !ok {
		t.Fatal("Simulator source template fixture was created without an ID")
	}

	t.Cleanup(func() {
		if err := client.Compute.DeleteComputeSourceTemplateById(context.Background(), id); err != nil {
			t.Errorf("Failed to delete simulator source template fixture %s: %v", id, err)
		}
	})
	return id
}
