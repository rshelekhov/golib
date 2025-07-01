# Config Loader

Flexible configuration loader for Go with auto-discovery and multiple formats.

## Features

- Generic type for any config struct
- Auto-discovery of config files
- YAML and .env file support
- Multiple config files merging
- Flag-based config path override with smart flag handling
- Environment variables integration with `CONFIG_PATH` support
- Robust error handling with detailed error messages

## Installation

```bash
go get github.com/rshelekhov/golib/config
```

## Usage

### Basic Usage

```go
type AppConfig struct {
    AppEnv string `yaml:"app_env"`
    Port   int    `yaml:"port"`
}

cfg := config.MustLoad[AppConfig]()
```

Auto-discovers config files:

- `config.yaml|yml`, `.env` (current directory)
- `config/config.yaml|yml`, `config/.env` (config subdirectory)
- `../config/config.yaml|yml`, `../config/.env` (parent directory)

### Config Files

**YAML:**

```yaml
app_env: production
port: 8080
database:
  host: localhost
  port: 5432
```

**.env:**

```env
APP_ENV=production
PORT=8080
DATABASE_HOST=localhost
DATABASE_PORT=5432
```

### Configuration Priority

1. Command-line flag: `./app -config ./config.yaml` (automatically defined if not exists)
2. Environment variable: `CONFIG_PATH=/path/to/config.yaml`
3. Explicit files via `WithFiles()`
4. Auto-discovered files in search paths

### Smart Flag Handling

The loader automatically handles the `-config` flag:

- Defines the flag if it doesn't already exist
- Safely works with pre-existing flag definitions
- Parses flags only when necessary
- Falls back to `CONFIG_PATH` environment variable

### Options

```go
cfg := config.MustLoad[AppConfig](
    config.WithFiles([]string{"base.yaml", "override.env"}),
    config.WithMergeFiles(true),
    config.WithSearchPaths([]string{"./custom/*.yaml"}),
    config.WithAllowUnknownFields(false),
    config.WithSkipFlags(false), // Enable flag parsing
)
```

**Available Options:**

- `WithFiles([]string)` - explicit config files
- `WithMergeFiles(bool)` - merge multiple files (default: true)
- `WithSearchPaths([]string)` - custom search paths for auto-discovery
- `WithAllowUnknownFields(bool)` - allow unknown fields (default: true)
- `WithSkipFlags(bool)` - skip CLI flags parsing (default: true)

### Error Handling

The loader provides detailed error messages:

- Shows which files were attempted to load
- Lists search paths when no files are found
- Reports configuration loading failures with context

### Struct Tags

```go
type Config struct {
    DatabaseURL string `yaml:"database_url"`
    Debug       bool   `yaml:"debug"`

    Database struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"database"`
}
```

### Environment Variables

Set configuration path via environment:

```bash
export CONFIG_PATH=/path/to/config.yaml
./app
```

Or use command-line flag:

```bash
./app -config /path/to/config.yaml
```

### Kubernetes

**Local development:**

```go
cfg := config.MustLoad[MyConfig]() // finds .env automatically
```

**Kubernetes with mounted ConfigMap:**

```go
cfg := config.MustLoad[MyConfig](
    config.WithSearchPaths([]string{
        "/etc/config/config.yaml", // mounted ConfigMap
        "/etc/secrets/.env",       // mounted Secret
    }),
)
```

## License

MIT
