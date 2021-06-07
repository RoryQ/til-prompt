# TIL-Prompt

TIL-Prompt is an interactive terminal prompt for creating and managing a collection of TIL (Today I Learned) entries.


# Installation
To install from source using golang 2.16

```
go install github.com/roryq/til-prompt/cmd/til@vlatest
```

# Usage

![](demo.gif)

Follow the prompts to save a new TIL entry.
A README.md is regenerated after each save.

```shell
Usage: til <command>

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  config
    Displays the current configuration.

Run "til <command> --help" for more information on a command.

```

# License
[MIT](LICENSE)
