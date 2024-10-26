# mess

A cli tool to clean up your development mess.

## Installation

```bash
go install github.com/olzhasar/mess@latest
```

## Commands

### Git

Find all git repositories in the specified path (useful for piping)

```bash
mess git <path>
```

Filter repositories with uncommitted changes

```bash
mess git <path> --dirty
```

Filter repositories with last commit older than specified number of days

```bash
mess git <path> --older <days>
```

### Clean

Remove temporary files from the specified path

```bash
mess clean <path>
```

Will delete the following patterns:

- '*.pyc'
- '__pycache__'
- '.mypy_cache'
- '.pytest_cache'
- '.ruff_cache'
- '.tox'
- '.nox'
- 'node_modules'

## License
MIT
