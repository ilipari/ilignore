# ilignore

A CLI tool to check files against `.gitignore`-style patterns, helping prevent commits of files that are already tracked in git repositories.

## Features

- Check files against multiple ignore files
- Multiple file source options:
  - Command line arguments
  - Shell command output (`git ls-files`, etc)
  - Standard input
  - Git diff integration
- Concurrent file processing
- Configurable output format
- YAML configuration support

## Installation

```bash
go install github.com/ilipari/ilignore@latest
```

## Configuration
Create `.ilignorerc.yaml` in your home directory or project root. You can specify a different config file using the `--config` flag.

Here is an axample of such file:

```yaml
check:
  input: git diff --cached --name-only --diff-filter=ACMD
  ignore: 
    - .gitignore
    - .iliignore
  name-only: true
  concurrency: 4
```

## Usage
### Check specific files
```bash
ilignore check file1.txt file2.txt
```

### Check git staged files
```bash
ilignore check --input "git diff --cached --name-only"
```

### Use multiple ignore files
```bash
ilignore check -g .gitignore -g .iliignore
```

### Show only filenames
```bash
ilignore check --name-only files.txt
```

## Command Line Flags
- `-i, --input`: Command to get files to check
- `-g, --ignore`: Ignore file paths (default: .ilignore)
- `--name-only`: Show only file names
- `-c, --concurrency`: Number of worker goroutines
- `-v, --verbose`: Enable verbose output (same as --log info)
- `--log`: Set log level (debug, info, warn, error) (default is WARN)
- `--config`: Config file path (default is $HOME/.ilignorerc.yaml)


## License
GNU General Public License v3.0
