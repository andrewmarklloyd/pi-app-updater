package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type AppConfigs struct {
	Configs []Config `json:"configs"`
}

func GetAppConfigs(path string) (AppConfigs, error) {
	var emptyAppConfigs AppConfigs
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		if err.Error() == fmt.Sprintf("open %s: no such file or directory", path) {
			return AppConfigs{Configs: []Config{}}, nil
		}
		return emptyAppConfigs, fmt.Errorf("reading app configs yaml file: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &emptyAppConfigs)
	if err != nil {
		return emptyAppConfigs, fmt.Errorf("unmarshalling app configs %s", err)
	}

	return emptyAppConfigs, nil
}

// WriteAppConfigs will overwrite the whole file.
// It is up to the caller to make sure all contents are there.
func (a *AppConfigs) WriteAppConfigs(path string) error {
	out, err := yaml.Marshal(a)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, out, 0644)
	if err != nil {
		return fmt.Errorf("writing app configs: %s", err)
	}
	return nil
}

func (a *AppConfigs) ConfigExists(c Config) bool {
	for _, existing := range a.Configs {
		if existing.RepoName == c.RepoName && existing.ManifestName == c.ManifestName {
			return true
		}
	}
	return false
}
