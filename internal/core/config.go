package core

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/roryq/til-prompt/pkg/editor"

	"github.com/charmbracelet/lipgloss"

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
	Editor        string
}

func (c Config) GetEditor() string {
	return c.Editor
}

func ensureConfigPath(scope *gap.Scope) (string, error) {
	scope.DataDirs()
	configPath, err := scope.ConfigPath("config.yaml")
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := defaultConfig(scope)
		config.SaveDirectory = getDataPath(config.SaveDirectory)
		bytes, err := yaml.Marshal(config)
		if err != nil {
			return "", err
		}
		if err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm); err != nil {
			return "", err
		}
		if err := ioutil.WriteFile(configPath, bytes, 0644); err != nil {
			return "", err
		}
	}

	return configPath, nil
}

func LoadConfig(scope *gap.Scope) (config Config, err error) {
	configPath, err := ensureConfigPath(scope)
	if err != nil {
		return config, err
	}

	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(bytes, &config)
	return config, err
}

func EditConfig(scope *gap.Scope) error {
	configPath, err := ensureConfigPath(scope)
	if err != nil {
		return err
	}

	config, err := LoadConfig(scope)
	if err != nil {
		return err
	}

	return editor.Launch(config, configPath)
}

func (c Config) Formatted() string {
	b := &strings.Builder{}
	headStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	bodyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	b.WriteString(headStyle.Render("SaveDirectory: "))
	b.WriteString(bodyStyle.Render(c.SaveDirectory))
	return b.String()
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

func (c Config) ListCategoryDirectories() []string {
	dirEntries, _ := os.ReadDir(c.SaveDirectory)
	var categories []string
	for _, e := range dirEntries {
		if e.IsDir() {
			categories = append(categories, e.Name())
		}
	}

	return categories
}
