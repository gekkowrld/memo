# Memo

This project aims to create a simple memo app.

It is inspired by:
[mattn memo](https://github.com/mattn/memo)

## Installation

```bash
go install github.com/gekkowrld/memo@v0.0.2
```

:warning: I can't gurantee that you'll get the latest release by using "@latest".
This may be due to how go looks for releases when asked for "@latest" instead of tags (idk how go does it)
Look for the latest tag and use it instead

## Usage

```txt
Document your life in style

Usage:
  memo [command]

Available Commands:
  config      Configure your environment
  delete      Delete a memo
  help        Help about any command
  list        List the memos already created
  new         Add a new memo
  view        View Your Memo

Flags:
  -h, --help   help for memo

Use "memo [command] --help" for more information about a command.
```

## License

[GNU GPL](./LICENSE)
