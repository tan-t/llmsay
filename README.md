# llmsay

`llmsay` is a command-line interface (CLI) tool that provides access to major language models such as GPT, Claude, and Gemini. It allows users to interact with these models directly from the terminal.

## Installation

[Provide installation instructions here]

## Usage

The basic syntax for using `llmsay` is:

```
llmsay [flags]
llmsay [command]
```

### Available Commands

- `completion`: Generate the autocompletion script for the specified shell
- `configure`: Configure the tool settings
- `help`: Get help about any command

### Flags

- `-f, --file string`: Specify the config file path
- `-h, --help`: Display help for llmsay
- `-m, --model string`: Specify the model name (default: "gpt-4o")

## Examples

[Provide some usage examples here]

## Configuration

The configuration file for `llmsay` is named `config.toml`.

Typical location: `$HOME/.config/llmsay/config.toml`

This location uses the `$HOME` environment variable, which expands to your home directory. For example:
- On Linux or macOS: `/home/username/.config/llmsay/config.toml`
- On Windows: `C:\Users\username\.config\llmsay\config.toml`

You can specify a different configuration file using the `-f` or `--file` flag:

```
llmsay -f /path/to/your/config.toml
```

To manage your configuration:

1. Run `llmsay` without any arguments to see information about the config file location.
2. Use the `configure` command to set up your configuration: `llmsay configure`

## Models

`llmsay` supports multiple language models. The default model is "gpt-4o", but you can specify a different model using the `-m` or `--model` flag.

## Getting Help

To get more information about a specific command, use:

```
llmsay [command] --help
```
