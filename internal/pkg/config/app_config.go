package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type AppConfigs []Config

func GetAppConfigs(path string) (AppConfigs, error) {
	var emptyAppConfigs AppConfigs
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return emptyAppConfigs, fmt.Errorf("reading app configs yaml file: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &emptyAppConfigs)
	if err != nil {
		return emptyAppConfigs, fmt.Errorf("unmarshalling app configs %s", err)
	}

	return emptyAppConfigs, nil
}

// WriteAppConfigs will overwrite the whole file
func WriteAppConfigs(path string, appConfigs AppConfigs) error {
	out, err := yaml.Marshal(appConfigs)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, out, 0644)
	if err != nil {
		return fmt.Errorf("writing run script: %s", err)
	}
	return nil
}
