package testutil

import "github.com/hashicorp/terraform-plugin-testing/helper/acctest"

// ResourcePrefix starts the name of everything acceptance tests create, so the
// sweepers can recognise what a killed run left behind. The shared fixtures,
// e.g. the static namespace and keyring, are named outside it and never swept.
const ResourcePrefix = "tf-acc"

// RandomName returns a unique name for an entity an acceptance test creates.
func RandomName() string {
	return acctest.RandomWithPrefix(ResourcePrefix)
}
