package testutil

import (
	"maps"

	"github.com/hashicorp/terraform-plugin-testing/config"
)

func CloneWithout(m config.Variables, key string) config.Variables {
	clone := maps.Clone(m)
	delete(clone, key)
	return clone
}

func CloneWith(m config.Variables, key string, value config.Variable) config.Variables {
	clone := maps.Clone(m)
	clone[key] = value
	return clone
}
