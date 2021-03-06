# TIL-Prompt

TIL-Prompt is an interactive terminal prompt for creating and managing a collection of TIL (Today I Learned) entries.


# Installation
To install from source using golang 1.16

```
go install github.com/roryq/til-prompt/cmd/til@vlatest
```

# Usage

![](demo.gif)

Follow the prompts to save a new TIL entry.
A README.md is regenerated after each save.

<!--usage-shell-->
```
Usage: til <command>

An interactive prompt for managing TIL entries.

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  new
    Create a new TIL entry. (default command)

  config list
    List the current configuration. (default sub-command)

  config edit
    Open the config in your configured $EDITOR

  edit [<keyword>]
    Open today's TIL in your configured editor. You can optionally pass a
    keyword to search for.

  view [<keyword>]
    Display a styled view of today's TIL. You can optionally pass a keyword to
    search for.

Run "til <command> --help" for more information on a command.
```

# License
[MIT](LICENSE)
