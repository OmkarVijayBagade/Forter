// Package config handles configuration management for forter
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Category represents a file category with its extensions
type Category struct {
	Name        string   `mapstructure:"name"`
	Extensions  []string `mapstructure:"extensions"`
	Description string   `mapstructure:"description"`
}

// Config holds the application configuration
type Config struct {
	Categories      []Category        `mapstructure:"categories"`
	CustomMappings  map[string]string `mapstructure:"custom_mappings"`
	DefaultCategory string            `mapstructure:"default_category"`
	SkipHidden      bool              `mapstructure:"skip_hidden"`
	SkipDirs        bool              `mapstructure:"skip_dirs"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Categories: []Category{
			{
				Name:        "Documents",
				Extensions:  []string{"pdf", "doc", "docx", "txt", "rtf", "odt", "xls", "xlsx", "ppt", "pptx", "csv", "md"},
				Description: "Document files",
			},
			{
				Name:        "Images",
				Extensions:  []string{"jpg", "jpeg", "png", "gif", "bmp", "svg", "webp", "ico", "tiff", "raw", "psd"},
				Description: "Image files",
			},
			{
				Name:        "Videos",
				Extensions:  []string{"mp4", "avi", "mkv", "mov", "wmv", "flv", "webm", "m4v", "mpg", "mpeg", "3gp"},
				Description: "Video files",
			},
			{
				Name:        "Audio",
				Extensions:  []string{"mp3", "wav", "flac", "aac", "ogg", "m4a", "wma", "opus"},
				Description: "Audio files",
			},
			{
				Name:        "Archives",
				Extensions:  []string{"zip", "rar", "7z", "tar", "gz", "bz2", "xz", "tgz", "tbz", "iso"},
				Description: "Archive files",
			},
			{
				Name:        "Code",
				Extensions:  []string{"go", "py", "js", "ts", "jsx", "tsx", "html", "css", "scss", "sass", "java", "c", "cpp", "h", "hpp", "rs", "rb", "php", "swift", "kt", "json", "xml", "yaml", "yml", "sql", "sh", "bash", "zsh"},
				Description: "Source code files",
			},
		},
		CustomMappings:  make(map[string]string),
		DefaultCategory: "Others",
		SkipHidden:      true,
		SkipDirs:        true,
	}
}

// Load loads configuration from file or creates default
func Load() (*Config, error) {
	cfg := DefaultConfig()

	viper.SetConfigName(".forter")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("default_category", cfg.DefaultCategory)
	viper.SetDefault("skip_hidden", cfg.SkipHidden)
	viper.SetDefault("skip_dirs", cfg.SkipDirs)
	viper.SetDefault("categories", cfg.Categories)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
		// Config file not found, use defaults
	} else {
		if err := viper.Unmarshal(cfg); err != nil {
			return nil, fmt.Errorf("error unmarshaling config: %w", err)
		}
	}

	return cfg, nil
}

// Save saves the current configuration to file
func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".forter.yaml")
	viper.Set("categories", c.Categories)
	viper.Set("custom_mappings", c.CustomMappings)
	viper.Set("default_category", c.DefaultCategory)
	viper.Set("skip_hidden", c.SkipHidden)
	viper.Set("skip_dirs", c.SkipDirs)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetCategoryForExtension returns the category name for a given extension
func (c *Config) GetCategoryForExtension(ext string) string {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))

	// Check custom mappings first
	if category, ok := c.CustomMappings[ext]; ok {
		return category
	}

	// Check category definitions
	for _, cat := range c.Categories {
		for _, e := range cat.Extensions {
			if strings.EqualFold(e, ext) {
				return cat.Name
			}
		}
	}

	return c.DefaultCategory
}

// GetAllCategoryNames returns all category names
func (c *Config) GetAllCategoryNames() []string {
	names := make([]string, 0, len(c.Categories)+1)
	for _, cat := range c.Categories {
		names = append(names, cat.Name)
	}
	names = append(names, c.DefaultCategory)
	return names
}

// IsHiddenFile checks if a file is hidden (starts with .)
func IsHiddenFile(filename string) bool {
	base := filepath.Base(filename)
	return strings.HasPrefix(base, ".")
}
