# Changelog

## v2.0.0

### Breaking Changes

- **BREAKING**: Replaced viper with cristalhq/aconfig for better performance and simpler configuration
- **BREAKING**: Changed struct tags from `mapstructure` to `yaml`
- **BREAKING**: Removed viper-specific options (`WithEnvPrefix`, `WithConfigName`, `WithConfigType`)

### New Features

- **Auto-discovery**: Automatically finds config files in standard locations
- **Multi-format support**: YAML and .env files
- **Multi-file merging**: Combine multiple config files with `WithMergeFiles()`
- **Custom search paths**: Configure custom discovery paths with `WithSearchPaths()`
- **Enhanced options**: New configuration options for fine-tuned control

## v1.0.0

- Initial release of the config package
