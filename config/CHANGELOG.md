# Changelog

All notable changes to the Config package will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2025-07-01

### Changed

- Improved flag handling with automatic flag definition and parsing
- Enhanced configuration path resolution with flag and environment variable support
- Better error handling with detailed error messages including file paths
- Optimized config file discovery with existence checks

### Added

- **Smart flag handling**: Automatically defines `-config` flag if not already defined
- **Fallback configuration**: Environment variable `CONFIG_PATH` as fallback when flag is not provided
- **Enhanced error messages**: More detailed error reporting with file paths and search paths
- **Robust file discovery**: Improved file existence checking in discovery process

### Fixed

- Fixed potential flag redefinition issues by checking if flag already exists
- Improved flag parsing logic to handle pre-parsed flags correctly

## [1.1.0] - 2025-06-30

### Changed

- Replaced viper with cristalhq/aconfig for better performance and simpler configuration
- Changed struct tags from `mapstructure` to `yaml`
- Removed viper-specific options (`WithEnvPrefix`, `WithConfigName`, `WithConfigType`)

### Added

- **Auto-discovery**: Automatically finds config files in standard locations
- **Multi-format support**: YAML and .env files
- **Multi-file merging**: Combine multiple config files with `WithMergeFiles()`
- **Custom search paths**: Configure custom discovery paths with `WithSearchPaths()`
- **Enhanced options**: New configuration options for fine-tuned control

## [1.0.0] - 2025-05-16

### Added

- Initial release of the config package
