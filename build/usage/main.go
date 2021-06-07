package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	helpOutput, err := exec.Command("go", "run", "cmd/til/main.go", "--help").Output()
	if err != nil {
		panic(err)
	}

	readme, err := ioutil.ReadFile("README.md")
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile("```shell\\n[^`]+```")
	matches := re.FindStringSubmatch(string(readme))
	replaced := strings.ReplaceAll(string(readme), matches[0],
		fmt.Sprintf("```shell\n%s\n```", helpOutput))

	ioutil.WriteFile("README.md", []byte(replaced), 0644)
}
