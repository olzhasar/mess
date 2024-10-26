# mess

A cli tool to spot and clean up mess in your projects.

## Features

- Find unfinished work
- Find outdated repositories
- Clean up temporary files

## Installation

```bash
go install github.com/olzhasar/mess@latest
```

## Commands

### Git

Find all git repositories in the specified path. Mainly useful in pipelines

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

**Deletes** the following files and directories:

- `*.pyc` — Python compiled files
- `__pycache__` — Python cache directory
- `.mypy_cache` — MyPy cache directory
- `.pytest_cache` — Pytest cache directory
- `.ruff_cache` — Ruff cache directory
- `.tox` — Tox virtual environment files
- `.nox` — Nox virtual environment files
- `node_modules` — Node.js dependencies

## Performance

`mess git` is significantly faster for finding Git repositories compared to the `find` + `test` alternative:

```bash
find <path> -type d -execdir test -d {}/.git \; -prune -print
```

Scanning a 300GB+ home directory on an M1 Macbook Pro:

```bash
mess git ~  1.90s user 7.22s system 39% cpu 23.096 total
find ~ -type d -execdir test -d {}/.git \; -print -prune  60.76s user 184.72s system 70% cpu 5:46.94 total
```

`mess` did 15X better in this case.


## License
MIT
