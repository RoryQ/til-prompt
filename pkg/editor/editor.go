package editor

import (
	"errors"
	"fmt"
	"os"

	"github.com/roryq/til-prompt/pkg/exec"
)

type GetEditorer interface {
	GetEditor() string
}

// Launch will open the filePath in either the configured editor or one set by the VISUAL / EDITOR envars
func Launch(config GetEditorer, filePath string) error {
	editor := coalesceString(config.GetEditor(), os.Getenv("VISUAL"), os.Getenv("EDITOR"))
	if editor == "" {
		return errors.New("no editor found: please configure an editor in the config, or set your $EDITOR env")
	}

	cmd := exec.CommandFromString(fmt.Sprintf("%s %s </dev/tty", editor, filePath))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func coalesceString(stringSlice ...string) string {
	for i := range stringSlice {
		if stringSlice[i] != "" {
			return stringSlice[i]
		}
	}

	return ""
}
