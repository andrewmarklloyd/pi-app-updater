package config

import (
	"testing"

	"github.com/andrewmarklloyd/pi-app-deployer/api/v1/manifest"
	"github.com/stretchr/testify/assert"
)

func Test_ValidEnv(t *testing.T) {
	m := manifest.Manifest{
		Env: []string{"MY_CONFIG", "HELLO_CONFIG"},
	}

	v := make(map[string]string)
	v["MY_CONFIG"] = "foobar"
	v["HELLO_CONFIG"] = "testing"

	cfg := Config{
		EnvVars: v,
	}

	err := ValidateEnvVars(m, cfg)
	assert.NoError(t, err)
}

func Test_InvalidEnv(t *testing.T) {
	m := manifest.Manifest{
		Env: []string{"MY_CONFIG", "HELLO_CONFIG"},
	}

	v := make(map[string]string)
	v["MY_CONFIG"] = "foobar"

	cfg := Config{
		EnvVars: v,
	}

	err := ValidateEnvVars(m, cfg)
	assert.EqualError(t, err, "manifest defined env vars should exactly match env vars configured in agent install command")

	m = manifest.Manifest{
		Env: []string{"MY_CONFIG"},
	}

	v = make(map[string]string)
	v["MY_CONFIG"] = "foobar"
	v["HELLO_CONFIG"] = "testing"

	cfg = Config{
		EnvVars: v,
	}

	err = ValidateEnvVars(m, cfg)
	assert.EqualError(t, err, "manifest defined env vars should exactly match env vars configured in agent install command")
}
