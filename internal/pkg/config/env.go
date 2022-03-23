package config

import (
	"fmt"
	"reflect"

	"github.com/andrewmarklloyd/pi-app-deployer/api/v1/manifest"
)

func ValidateEnvVars(m manifest.Manifest, cfg Config) error {
	keys := make([]string, len(cfg.EnvVars))

	i := 0
	for k := range cfg.EnvVars {
		keys[i] = k
		i++
	}

	if !reflect.DeepEqual(keys, m.Env) {
		return fmt.Errorf("manifest defined env vars should exactly match env vars configured in agent install command")
	}

	return nil
}
