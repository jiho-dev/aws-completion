package config

import (
	//"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

/////////////////////

const API_TERMINATED = "api-terminated"

type ApiOption struct {
	OutputField string   `yaml:"OutputField"`
	Required    []string `yaml:"Required"`
	Args        []string `yaml:"Args"`
}

type AwscConfig struct {
	Version   string   `yaml:"Version"`
	ApiFilter []string `yaml:"ApiFilter"`
	// key: full-api, value: api-option
	ApiOptions map[string]ApiOption `yaml:"ApiOptions"`
}

func (ao *ApiOption) GetFlags(cmdName string) *flag.FlagSet {
	flag := flag.NewFlagSet(cmdName, flag.ExitOnError)

	for _, c := range ao.Required {
		flag.String(c, "", c)
	}

	for _, c := range ao.Args {
		flag.String(c, "", c)
	}

	return flag
}

func ParseConfig(fileName string) (*AwscConfig, error) {
	fileName, _ = filepath.Abs(fileName)
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var cc AwscConfig
	err = yaml.Unmarshal(yamlFile, &cc)
	if err != nil {
		return nil, err
	}

	return &cc, nil
}

func WriteConfig(conf *AwscConfig, fileName string) error {
	yamlData, err := yaml.Marshal(conf)

	fileName, _ = filepath.Abs(fileName)
	err = ioutil.WriteFile(fileName, yamlData, 644)
	if err != nil {
		return err
	}

	return nil
}

func YamlTest() {
	c := AwscConfig{
		Version: "1",

		ApiFilter: []string{
			/*
				"describe-network",
				"describe-no",
			*/
		},

		/*
			ApiPrefix: map[string]map[string][]string{
				"describe": {
					"network": []string{"acls", "interfaces"},
					"volumes": []string{API_TERMINATED, "modifications"},
				},
				"import": {
					"image":    []string{API_TERMINATED},
					"key":      []string{"pair"},
					"snapshot": []string{API_TERMINATED},
				},
				"modify": {
					"hosts":    []string{API_TERMINATED},
					"instance": []string{"attribute", "capacity-reservation-attributes", "event-window"},
					"snapshot": []string{API_TERMINATED},
				},
				"wait": {
					API_TERMINATED: []string{},
				},
			},
		*/

		ApiOptions: map[string]ApiOption{
			"describe-network-interfaces": {
				OutputField: "NetworkInterfaces",
				Args:        []string{"host-ip", "network-interface-id", "vpc-id"},
			},

			"describe-network-acls": {
				OutputField: "NetworkInterfaces",
				Args:        []string{"host-ip", "nacl-id", "vpc-id"},
			},
		},
	}

	sort.Strings(c.ApiFilter)

	yamlData, err := yaml.Marshal(&c)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	fmt.Println(" --- YAML ---")
	fmt.Printf("%s \n", string(yamlData))
}
