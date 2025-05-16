# Config Loader

Simple and flexible configuration loader for Go applications using Viper.

## Features

- Generic type support
- Environment variables integration
- Config file support
- Flexible Viper configuration through options
- Kubernetes-friendly (works with ConfigMaps and Secrets)

## Installation

```bash
go get github.com/rshelekhov/golib/config
```

## Usage

### Basic Usage

```go
package main

import (
    "github.com/rshelekhov/golib/config"
)

type AppConfig struct {
    AppEnv string `mapstructure:"APP_ENV"`
    AppID  string `mapstructure:"APP_ID"`
    Port   int    `mapstructure:"PORT"`
}

func main() {
    cfg := config.MustLoad[AppConfig]()
    // Use cfg...
}
```

### With Options

```go
cfg := config.MustLoad[AppConfig](
    config.WithEnvPrefix("MYAPP"),     // Environment variables will be prefixed with MYAPP_
    config.WithConfigType("yaml"),     // Set config file type
    config.WithConfigName("config"),   // Set config file name
)
```

### Configuration Sources

1. Environment Variables:

```bash
export APP_ENV=production
export APP_ID=myapp
export PORT=8080
```

2. Config File (e.g., config.yaml):

```yaml
app_env: production
app_id: myapp
port: 8080
```

Load config file by:

- Setting CONFIG_PATH environment variable
- Using -config flag: `./app -config ./config.yaml`

### Field Tags

Use `mapstructure` tag to map config values to struct fields:

```go
type Config struct {
    DatabaseURL string `mapstructure:"DATABASE_URL"`
    RedisHost   string `mapstructure:"REDIS_HOST"`
    Debug       bool   `mapstructure:"DEBUG"`
}
```

### Kubernetes Integration

In Kubernetes, all settings come through environment variables from ConfigMap and Secret, so config files are only needed for local development.

## License

MIT
