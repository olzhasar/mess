# mess

`mess` deletes common temporary development files in your projects. It ships with common cleanup patterns for popular programming languages like caches, build artifacts, bytecode files, linters output, etc.

Why a separate tool instead of a bash script?
- Reports reclaimed disk space
- Easy to configure and override patterns
- Simple and fast
- Reduces the number of bash scripts in the world

## Installation

### Pre-built binaries

Download a ready-to-use pre-built binary for your platform (Linux and macOS are currently supported) in 
[GitHub releases](https://github.com/olzhasar/mess/releases).

### Install with go

```bash
go install github.com/olzhasar/mess@latest
```

### Building from source

```bash
make install
```

By default, this installs `mess` to `~/.local/bin`. To install elsewhere, set
`PREFIX`:

```bash
make install PREFIX=/usr/local
```

## Usage

**CAUTION:** *this tool deletes files and directories from your filesystem!*

*While the pre-configured patterns should be safe for most users, please, read the patterns list first and ensure it fits your needs. See [Patterns](#patterns)*

### Clean a directory

Clean a specific path (scans only direct children by default):

```sh
mess clean ~/my_projects
```

Scan subdirectories too:

```sh
mess clean -r ~/my_projects
```

Print each removed path:

```sh
mess clean -v ~/my_projects
```

Use custom patterns for one run. These replace the configured patterns:

```sh
mess clean -r -p "node_modules,*.pyc" ~/dev/lots-of-js/
```

### Patterns

To display the currently configured patterns, run:

```sh
mess patterns
```

By default, `mess` uses its built-in [PATTERNS](./lib/PATTERNS). To override them you can use either use a `-p` flag or you can create one of these files:

- `~/.mess_patterns`
- `<user-config-dir>/mess/patterns`

The user config directory is provided by your operating system. On Linux it's `$XDG_CONFIG_HOME` (defaults to `~/.config`), on Darwin - `~/Library/Application Support/`

Pattern files use one pattern per line. Empty lines and lines starting with `#`
are ignored.

Patterns use a shell glob-style syntax.

Example:

```sh
# Python
*.pyc
__pycache__

# JavaScript
node_modules

# Specific path
foo/bar/
```

## License
MIT
