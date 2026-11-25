// Package sweep_test runs the sweepers, which delete what acceptance tests left
// behind when a run died before its final destroy. See make sweep.
//
// Every sweeper is registered here, one per entity type, so that a single test
// binary runs each once and Dependencies can order them across types. A
// sweeper only deletes names starting with testutil.ResourcePrefix.
package sweep_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	addSweeper("yellowdog_application", sweepApplications())
	addSweeper("yellowdog_group", sweepGroups())
	addSweeper("yellowdog_keyring", sweepKeyrings())
	addSweeper("yellowdog_compute_requirement_template", sweepComputeRequirementTemplates())
	// requirement templates refer to source templates
	addSweeper("yellowdog_compute_source_template", sweepComputeSourceTemplates(),
		"yellowdog_compute_requirement_template")
	addSweeper("yellowdog_machine_image_family", sweepImageFamilies())
	// last, as groups hold roles in namespaces and templates and image families live in them
	addSweeper("yellowdog_namespace", sweepNamespaces(),
		"yellowdog_group", "yellowdog_compute_requirement_template", "yellowdog_compute_source_template",
		"yellowdog_machine_image_family")
}

// addSweeper registers f under name, which plugin-testing also uses to tell
// whether a sweeper has run, so it must be unique.
func addSweeper(name string, f resource.SweeperFunc, dependencies ...string) {
	resource.AddTestSweepers(name, &resource.Sweeper{
		Name:         name,
		F:            f,
		Dependencies: dependencies,
	})
}

// TestMain runs the sweepers when go test is given -sweep.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}
