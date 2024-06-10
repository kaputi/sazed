package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sazed/utils"
	"strings"

	"gopkg.in/yaml.v3"
)

// i don't want root and filetype to be exported, so i make 2 structs

type Config struct {
	root     string
	filetype string
}

type yamlConfig struct {
	Root     string `yaml:"root"`
	Filetype string `yaml:"filetype"`
}

func yamlToConfig(yamlConfig yamlConfig) Config {
	return Config{
		root:     yamlConfig.Root,
		filetype: yamlConfig.Filetype,
	}
}

// func configToYaml(config sazedConfig) yamlConfig {
// 	return yamlConfig{
// 		Root:     config.root,
// 		Filetype: config.filetype,
// 	}
// }

var defaultConfig = yamlConfig{
	Root:     "~/.sazed",
	Filetype: "go",
}

func (c *Config) Root() string {
	return c.root
}

func (c *Config) Filetype() string {
	return c.filetype
}

func Load() Config {
	consfigFilePath := fmt.Sprintf("%s/config.yaml", getConfigPath())

	if _, err := os.Stat(consfigFilePath); errors.Is(err, os.ErrNotExist) {
		createConfigFile()
	}

	configFile, err := os.ReadFile(consfigFilePath)
	utils.CheckErr(err)

	var yamlData yamlConfig

	err = yaml.Unmarshal(configFile, &yamlData)
	utils.CheckErr(err)

	rootPathParts := strings.Split(yamlData.Root, string(os.PathSeparator))
	if string(rootPathParts[0]) == "~" {
		rootPathParts[0], _ = os.UserHomeDir()
		yamlData.Root = strings.Join(rootPathParts, string(os.PathSeparator))
	}

	return yamlToConfig(yamlData)
}

func getConfigPath() string {
	if sazedConfigPath, ok := os.LookupEnv("SAZED_CONFIG"); ok {
		return sazedConfigPath
	}

	if xdgConfigHome, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		return fmt.Sprintf("%s/sazed", xdgConfigHome)
	}

	return fmt.Sprintf("%s/.config/sazed", os.Getenv("HOME"))
}

func createConfigFile() {
	configPath := getConfigPath()

	utils.CreateDirIfNotExist(configPath)

	file, err := os.Create(fmt.Sprintf("%s/config.yaml", configPath))
	utils.CheckErr(err)
	defer file.Close()

	yamlData, err := yaml.Marshal(&defaultConfig)
	utils.CheckErr(err)

	_, err = io.WriteString(file, string(yamlData))
	utils.CheckErr(err)
}
