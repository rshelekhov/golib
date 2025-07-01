package config

import (
	"flag"
	"log"
	"os"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigdotenv"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

const CONFIG_PATH = "CONFIG_PATH"

type LoaderConfig struct {
	Files              []string
	AllowUnknownFields bool
	SkipFlags          bool
	MergeFiles         bool
	SearchPaths        []string
}

type Option func(*LoaderConfig)

func WithFiles(files []string) Option {
	return func(cfg *LoaderConfig) {
		cfg.Files = files
	}
}

func WithAllowUnknownFields(allow bool) Option {
	return func(cfg *LoaderConfig) {
		cfg.AllowUnknownFields = allow
	}
}

func WithSkipFlags(skip bool) Option {
	return func(cfg *LoaderConfig) {
		cfg.SkipFlags = skip
	}
}

func WithMergeFiles(merge bool) Option {
	return func(cfg *LoaderConfig) {
		cfg.MergeFiles = merge
	}
}

func WithSearchPaths(paths []string) Option {
	return func(cfg *LoaderConfig) {
		cfg.SearchPaths = paths
	}
}

func MustLoad[T any](opts ...Option) *T {
	cfg := new(T)

	// Default loader config
	loaderCfg := &LoaderConfig{
		AllowUnknownFields: true,
		SkipFlags:          true,
		MergeFiles:         true,
		SearchPaths:        getDefaultSearchPaths(),
	}

	// Apply options
	for _, opt := range opts {
		opt(loaderCfg)
	}

	configPath := fetchConfigPath(loaderCfg.SkipFlags)

	var files []string

	// If a path to the config is specified, use it
	if configPath != "" {
		files = []string{configPath}
	} else if len(loaderCfg.Files) > 0 {
		// Use explicitly provided files
		files = loaderCfg.Files
	} else {
		// Auto-discover config files
		files = discoverConfigFiles(loaderCfg.SearchPaths)
		if len(files) == 0 {
			log.Fatalf("no config files found in search paths: %v", loaderCfg.SearchPaths)
		}
	}

	loader := aconfig.LoaderFor(cfg, aconfig.Config{
		Files:              files,
		AllowUnknownFields: loaderCfg.AllowUnknownFields,
		SkipFlags:          loaderCfg.SkipFlags,
		MergeFiles:         loaderCfg.MergeFiles,
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".yml":  aconfigyaml.New(),
			".env":  aconfigdotenv.New(),
		},
	})

	if err := loader.Load(); err != nil {
		log.Fatalf("failed to load config from files %v: %v", files, err)
	}

	return cfg
}

func fetchConfigPath(skipFlags bool) string {
	var v string

	if !skipFlags {
		// Check if flag is already defined
		if flag.Lookup("config") == nil {
			flag.StringVar(&v, "config", "", "path to config file")
			if !flag.Parsed() {
				flag.Parse()
			}
		}
	}

	// If flag exists, get its value
	if configFlag := flag.Lookup("config"); configFlag != nil {
		v = configFlag.Value.String()
	}

	// Fallback to environment variable
	if v == "" {
		v = os.Getenv(CONFIG_PATH)
	}

	return v
}

func getDefaultSearchPaths() []string {
	return []string{
		// Current directory
		"config.yaml",
		"config.yml",
		".env",

		// Config subdirectory
		"config/config.yaml",
		"config/config.yml",
		"config/.env",

		// Parent directory
		"../config/config.yaml",
		"../config/config.yml",
		"../config/.env",
	}
}

func discoverConfigFiles(searchPaths []string) []string {
	var existingFiles []string
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			existingFiles = append(existingFiles, path)
		}
	}
	return existingFiles
}
