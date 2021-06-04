package tui

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	gap "github.com/muesli/go-app-paths"
)

func defaultConfig(scope *gap.Scope) Config {
	dd, _ := scope.DataDirs()
	return Config{
		SaveDirectory: dd[0],
	}
}

type Config struct {
	SaveDirectory string
}

func LoadConfig(scope *gap.Scope) (config Config, err error) {
	scope.DataDirs()
	configPath, err := scope.ConfigPath("config.yaml")
	if err != nil {
		return config, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config = defaultConfig(scope)
		config.SaveDirectory = getDataPath(config.SaveDirectory)
		bytes, err := yaml.Marshal(config)
		if err != nil {
			return config, err
		}
		if err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm); err != nil {
			return config, err
		}
		if err := ioutil.WriteFile(configPath, bytes, 0644); err != nil {
			return config, err
		}
	}

	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(bytes, &config)
	return config, err
}

func getDataPath(defaultPath string) string {
	reader := bufio.NewScanner(os.Stdin)
	fmt.Printf("No config found. Set the directory to store your til,"+
		" or press Enter to use: %s\n> ", defaultPath)
	reader.Scan()
	if input := reader.Text(); input != "" {
		return input
	}
	return defaultPath
}
