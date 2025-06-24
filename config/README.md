# Config Loader

Flexible configuration loader for Go with auto-discovery and multiple formats.

## Features

- Generic type for any config struct
- Auto-discovery of config files
- YAML and .env file support
- Multiple config files merging
- Flag-based config path override
- Environment variables integration

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

1. Command-line flag: `./app -config ./config.yaml`
2. Environment variable: `CONFIG_PATH=/path/to/config.yaml`
3. Explicit files via `WithFiles()`
4. Auto-discovered files

### Options

```go
cfg := config.MustLoad[AppConfig](
    config.WithFiles([]string{"base.yaml", "override.env"}),
    config.WithMergeFiles(true),
    config.WithSearchPaths([]string{"./custom/*.yaml"}),
    config.WithAllowUnknownFields(false),
)
```

**Available Options:**

- `WithFiles([]string)` - explicit config files
- `WithMergeFiles(bool)` - merge multiple files (default: true)
- `WithSearchPaths([]string)` - custom search paths
- `WithAllowUnknownFields(bool)` - allow unknown fields (default: true)
- `WithSkipFlags(bool)` - skip CLI flags parsing (default: true)

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
