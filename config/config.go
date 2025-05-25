package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfigs() Config {
	yamlData, err := os.ReadFile("app_config.yaml")
	if err != nil {
		log.Fatal("Error while reading app config file", err)
	}

	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		log.Fatal("Error while parsing YAML config", err)
	}

	return config
}

func LoadDBConfig() DbConfig {
	yamlData, err := os.ReadFile("app_config.yaml")
	if err != nil {
		log.Fatal("Error while reading app config file", err)
	}

	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		log.Fatal("Error while parsing DB config", err)
	}

	fmt.Println("Loaded DB Config")
	return config.DbConfig
}

//func LoadServerConfig() ServerConfig {
//	yamlData, err := os.ReadFile("app_config.yaml")
//	if err != nil {
//		log.Fatal("Error while reading app config file", err)
//	}
//
//	var config Config
//	if err := yaml.Unmarshal(yamlData, &config); err != nil {
//		log.Fatal("Error while parsing server config", err)
//	}
//
//	fmt.Println("Loaded Server Config")
//	return config.Server
//}
