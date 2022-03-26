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
	_, err := GetAppConfigs(testConfigPath)
	assert.EqualError(t, err, fmt.Sprintf("reading app configs yaml file: open %s: no such file or directory", testConfigPath))
}

func Test_CreateConfig(t *testing.T) {
	e := make(map[string]string)
	e["MY_CONFIG"] = "foobar"
	e["HELLO_CONFIG"] = "testing"
	a := []Config{
		{
			RepoName:      "andrewmarklloyd/pi-test",
			ManifestName:  "pi-test-arm",
			HomeDir:       "/home/pi",
			AppUser:       "pi",
			LogForwarding: false,
			EnvVars:       e,
		},
	}

	u, _ := uuid.NewUUID()
	testConfigPath := fmt.Sprintf("/tmp/.pi-app-deployer.app.config.%s", u.String())
	err := WriteAppConfigs(testConfigPath, a)
	assert.NoError(t, err)

	content, err := os.ReadFile(testConfigPath)
	assert.NoError(t, err)
	expectedContent := `- repoName: andrewmarklloyd/pi-test
  manifestName: pi-test-arm
  homeDir: /home/pi
  appUser: pi
  logForwarding: false
  envVars:
    HELLO_CONFIG: testing
    MY_CONFIG: foobar
`
	assert.Equal(t, expectedContent, string(content))

	configs, err := GetAppConfigs(testConfigPath)
	assert.NoError(t, err)

	assert.Equal(t, "pi", configs[0].AppUser)
	assert.Equal(t, "andrewmarklloyd/pi-test", configs[0].RepoName)
	assert.Equal(t, "pi-test-arm", configs[0].ManifestName)
	assert.Equal(t, "/home/pi", configs[0].HomeDir)
	assert.False(t, configs[0].LogForwarding)
	expectedMap := make(map[string]string)
	expectedMap["MY_CONFIG"] = "foobar"
	expectedMap["HELLO_CONFIG"] = "testing"
	assert.Equal(t, expectedMap, configs[0].EnvVars)
	fmt.Println()
}

func Test_CreateMultipleConfigs(t *testing.T) {
	config1Env := make(map[string]string)
	config1Env["MY_CONFIG"] = "foobar"
	config1Env["HELLO_CONFIG"] = "testing"

	config2Env := make(map[string]string)
	config2Env["HELLO_WORLD"] = "hello-world"
	config2Env["CONFIG"] = "config-test"
	a := []Config{
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
	}

	u, _ := uuid.NewUUID()
	testConfigPath := fmt.Sprintf("/tmp/.pi-app-deployer.app.config.%s", u.String())
	err := WriteAppConfigs(testConfigPath, a)
	assert.NoError(t, err)

	content, err := os.ReadFile(testConfigPath)
	assert.NoError(t, err)
	expectedContent := `- repoName: andrewmarklloyd/pi-test
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

	configs, err := GetAppConfigs(testConfigPath)
	assert.NoError(t, err)

	assert.Equal(t, "pi", configs[0].AppUser)
	assert.Equal(t, "andrewmarklloyd/pi-test", configs[0].RepoName)
	assert.Equal(t, "pi-test-arm", configs[0].ManifestName)
	assert.Equal(t, "/home/pi", configs[0].HomeDir)
	assert.False(t, configs[0].LogForwarding)
	expectedMap := make(map[string]string)
	expectedMap["MY_CONFIG"] = "foobar"
	expectedMap["HELLO_CONFIG"] = "testing"
	assert.Equal(t, expectedMap, configs[0].EnvVars)

	assert.Equal(t, "app-runner", configs[1].AppUser)
	assert.Equal(t, "andrewmarklloyd/pi-test-2", configs[1].RepoName)
	assert.Equal(t, "pi-test-amd64", configs[1].ManifestName)
	assert.Equal(t, "/home/app-runner", configs[1].HomeDir)
	assert.True(t, configs[1].LogForwarding)
	expectedMap = make(map[string]string)
	expectedMap["CONFIG"] = "config-test"
	expectedMap["HELLO_WORLD"] = "hello-world"
	assert.Equal(t, expectedMap, configs[1].EnvVars)
}
