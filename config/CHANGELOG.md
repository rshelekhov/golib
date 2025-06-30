# Changelog

## v1.1.0

### Changes

- Replaced viper with cristalhq/aconfig for better performance and simpler configuration
- Changed struct tags from `mapstructure` to `yaml`
- Removed viper-specific options (`WithEnvPrefix`, `WithConfigName`, `WithConfigType`)

### New Features

- **Auto-discovery**: Automatically finds config files in standard locations
- **Multi-format support**: YAML and .env files
- **Multi-file merging**: Combine multiple config files with `WithMergeFiles()`
- **Custom search paths**: Configure custom discovery paths with `WithSearchPaths()`
- **Enhanced options**: New configuration options for fine-tuned control

## v1.0.0

- Initial release of the config package
