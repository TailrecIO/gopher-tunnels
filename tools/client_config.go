package main

import (
	"flag"
	"github.com/tailrecio/gopher-tunnels/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type ClientFlags struct {
	EnvironmentName string
	ClientConfigFile string
	AppConfigFile string
}

func readEnvironment(flags *ClientFlags) map[string]string {

	flag.StringVar(&flags.EnvironmentName, "env", "dev", "environment name")
	flag.StringVar(&flags.ClientConfigFile, "output", "gopher.yml", "destination of the client config file")
	flag.StringVar(&flags.AppConfigFile, "input", "application.yml", "the application config file")
	flag.Parse()
	log.Printf("Generating a config file for env: `%v` at `%v", flags.EnvironmentName, flags.ClientConfigFile)

	var	environments map[string]map[string]string
	ymlData, err := ioutil.ReadFile(flags.AppConfigFile)
	if err != nil {
		log.Fatalf("Failed to read a file: `%v`", flags.AppConfigFile)
	}
	err = yaml.Unmarshal(ymlData, &environments)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML from a file due to `%v`", err.Error())
	}
	environment := environments[flags.EnvironmentName]
	if environment == nil {
		log.Fatalf("Environment: `%v` not found in the configuration file", flags.EnvironmentName)
	}
	return environment
}

func main() {
	flags := ClientFlags{}
	clientConfig := readEnvironment(&flags)
	if clientConfig[config.BaseQueueEndpoint] != "" {
		// remove account ID from a config map
		delete(clientConfig, config.AccountId)
	} else {
		delete(clientConfig, config.BaseQueueEndpoint)
	}
	ymlData, err := yaml.Marshal(clientConfig)
	if err != nil {
		log.Fatalf("Failed to marshal a config: `%v` to YAML", clientConfig)
	}
	err = ioutil.WriteFile(flags.ClientConfigFile, ymlData, 0644)
	if err != nil {
		log.Fatalf("Failed to write a config file due to: %v", err.Error())
	}
}