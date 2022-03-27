package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_NoConfig(t *testing.T) {
	u, _ := uuid.NewUUID()
	testConfigPath := fmt.Sprintf("/tmp/.pi-app-deployer.app.config.%s", u.String())
	appConfigs, err := GetAppConfigs(testConfigPath)
	assert.NoError(t, err)
	assert.NotNil(t, appConfigs)
	assert.NotNil(t, appConfigs.Configs)
	assert.Equal(t, []Config{}, appConfigs.Configs)
}

func Test_CreateConfig(t *testing.T) {
	e := make(map[string]string)
	e["MY_CONFIG"] = "foobar"
	e["HELLO_CONFIG"] = "testing"
	a := AppConfigs{[]Config{
		{
			RepoName:      "andrewmarklloyd/pi-test",
			ManifestName:  "pi-test-arm",
			HomeDir:       "/home/pi",
			AppUser:       "pi",
			LogForwarding: false,
			EnvVars:       e,
		},
	}}

	u, _ := uuid.NewUUID()
	testConfigPath := fmt.Sprintf("/tmp/.pi-app-deployer.app.config.%s", u.String())

	err := a.WriteAppConfigs(testConfigPath)
	assert.NoError(t, err)

	content, err := os.ReadFile(testConfigPath)
	assert.NoError(t, err)
	expectedContent := `configs:
- repoName: andrewmarklloyd/pi-test
  manifestName: pi-test-arm
  homeDir: /home/pi
  appUser: pi
  logForwarding: false
  envVars:
    HELLO_CONFIG: testing
    MY_CONFIG: foobar
`
	assert.Equal(t, expectedContent, string(content))

	aConf, err := GetAppConfigs(testConfigPath)
	assert.NoError(t, err)

	assert.Equal(t, "pi", aConf.Configs[0].AppUser)
	assert.Equal(t, "andrewmarklloyd/pi-test", aConf.Configs[0].RepoName)
	assert.Equal(t, "pi-test-arm", aConf.Configs[0].ManifestName)
	assert.Equal(t, "/home/pi", aConf.Configs[0].HomeDir)
	assert.False(t, aConf.Configs[0].LogForwarding)
	expectedMap := make(map[string]string)
	expectedMap["MY_CONFIG"] = "foobar"
	expectedMap["HELLO_CONFIG"] = "testing"
	assert.Equal(t, expectedMap, aConf.Configs[0].EnvVars)
	fmt.Println()
}

func Test_CreateMultipleConfigs(t *testing.T) {
	config1Env := make(map[string]string)
	config1Env["MY_CONFIG"] = "foobar"
	config1Env["HELLO_CONFIG"] = "testing"

	config2Env := make(map[string]string)
	config2Env["HELLO_WORLD"] = "hello-world"
	config2Env["CONFIG"] = "config-test"
	a := AppConfigs{[]Config{
		{
			RepoName:      "andrewmarklloyd/pi-test",
			ManifestName:  "pi-test-arm",
			HomeDir:       "/home/pi",
			AppUser:       "pi",
			LogForwarding: false,
			EnvVars:       config1Env,
		},
		{
			RepoName:      "andrewmarklloyd/pi-test-2",
			ManifestName:  "pi-test-amd64",
			HomeDir:       "/home/app-runner",
			AppUser:       "app-runner",
			LogForwarding: true,
			EnvVars:       config2Env,
		},
	}}

	u, _ := uuid.NewUUID()
	testConfigPath := fmt.Sprintf("/tmp/.pi-app-deployer.app.config.%s", u.String())
	err := a.WriteAppConfigs(testConfigPath)
	assert.NoError(t, err)

	content, err := os.ReadFile(testConfigPath)
	assert.NoError(t, err)
	expectedContent := `configs:
- repoName: andrewmarklloyd/pi-test
  manifestName: pi-test-arm
  homeDir: /home/pi
  appUser: pi
  logForwarding: false
  envVars:
    HELLO_CONFIG: testing
    MY_CONFIG: foobar
- repoName: andrewmarklloyd/pi-test-2
  manifestName: pi-test-amd64
  homeDir: /home/app-runner
  appUser: app-runner
  logForwarding: true
  envVars:
    CONFIG: config-test
    HELLO_WORLD: hello-world
`
	assert.Equal(t, expectedContent, string(content))

	a, err = GetAppConfigs(testConfigPath)
	assert.NoError(t, err)

	assert.Equal(t, "pi", a.Configs[0].AppUser)
	assert.Equal(t, "andrewmarklloyd/pi-test", a.Configs[0].RepoName)
	assert.Equal(t, "pi-test-arm", a.Configs[0].ManifestName)
	assert.Equal(t, "/home/pi", a.Configs[0].HomeDir)
	assert.False(t, a.Configs[0].LogForwarding)
	expectedMap := make(map[string]string)
	expectedMap["MY_CONFIG"] = "foobar"
	expectedMap["HELLO_CONFIG"] = "testing"
	assert.Equal(t, expectedMap, a.Configs[0].EnvVars)

	assert.Equal(t, "app-runner", a.Configs[1].AppUser)
	assert.Equal(t, "andrewmarklloyd/pi-test-2", a.Configs[1].RepoName)
	assert.Equal(t, "pi-test-amd64", a.Configs[1].ManifestName)
	assert.Equal(t, "/home/app-runner", a.Configs[1].HomeDir)
	assert.True(t, a.Configs[1].LogForwarding)
	expectedMap = make(map[string]string)
	expectedMap["CONFIG"] = "config-test"
	expectedMap["HELLO_WORLD"] = "hello-world"
	assert.Equal(t, expectedMap, a.Configs[1].EnvVars)
}

func Test_ConfigExists(t *testing.T) {
	c1 := Config{
		RepoName:      "andrewmarklloyd/pi-test",
		ManifestName:  "pi-test-arm",
		HomeDir:       "/home/pi",
		AppUser:       "pi",
		LogForwarding: false,
		EnvVars:       map[string]string{"MY_CONFIG": "foobar", "HELLO_CONFIG": "testing"},
	}
	c2 := Config{
		RepoName:      "andrewmarklloyd/pi-test-2",
		ManifestName:  "pi-test-amd64",
		HomeDir:       "/home/app-runner",
		AppUser:       "app-runner",
		LogForwarding: true,
		EnvVars:       map[string]string{"HELLO_WORLD": "hello-world", "CONFIG": "config-test"},
	}
	c3 := Config{
		RepoName:      "andrewmarklloyd/pi-test",
		ManifestName:  "pi-agent-arm",
		HomeDir:       "/home/app-runner",
		AppUser:       "app-runner",
		LogForwarding: true,
		EnvVars:       map[string]string{"HELLO_WORLD": "hello-world", "CONFIG": "config-test"},
	}
	appConfigs := AppConfigs{[]Config{c1, c2}}
	exists := appConfigs.ConfigExists(c1)
	assert.True(t, exists, "Config should exist in the app configs struct")

	exists = appConfigs.ConfigExists(c3)
	assert.False(t, exists, "Config should NOT exist in the app configs struct")
}
