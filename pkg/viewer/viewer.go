package viewer

import (
	"fmt"
	"io/ioutil"

	"github.com/charmbracelet/glamour"
)

func Launch(filePath string) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	out, err := glamour.Render(string(bytes), "dark")
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}
