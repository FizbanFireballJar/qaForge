# QA Forge

AI-powered testing tool for QA teams.

## Installation

### macOS
```bash
# Coming soon: brew install qaforge

# Manual (dla developmentu)
git clone https://github.com/yourcompany/qaforge.git
cd qaforge
make build
sudo mv qaforge /usr/local/bin/
```

### Linux
```bash
# Manual
git clone https://github.com/yourcompany/qaforge.git
cd qaforge
make build
sudo mv qaforge /usr/local/bin/
```

### Windows
```powershell
# Coming soon: scoop install qaforge

# Manual
git clone https://github.com/yourcompany/qaforge.git
cd qaforge
go build -o qaforge.exe cmd/qaforge/*.go
```

## Quick Start
```bash
# Wyświetl pomoc
qaforge --help

# Moduły (w budowie)
qaforge gen create
qaforge run tree
```

## Development
```bash
# Build
make build

# Build dla wszystkich platform
make build-all

# Run tests
make test

# Run bez build
make run
```

## Project Status

🚧 **MVP w budowie** - tydzień 1/6

- [x] Podstawowa struktura CLI
- [ ] qaforge gen create
- [ ] qaforge gen review
- [ ] qaforge run tree
- [ ] qaforge run plan
- [ ] qaforge run start