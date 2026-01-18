package config

import (
	"errors"
	"os"

	"github.com/jeroendee/ssg/internal/model"
	"gopkg.in/yaml.v3"
)

// yamlConfig represents the YAML file structure.
type yamlConfig struct {
	Site struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		BaseURL     string `yaml:"baseURL"`
		Author      string `yaml:"author"`
		Logo        string `yaml:"logo"`
		Favicon     string `yaml:"favicon"`
	} `yaml:"site"`
	Build struct {
		Content string `yaml:"content"`
		Output  string `yaml:"output"`
		Assets  string `yaml:"assets"`
	} `yaml:"build"`
	Navigation []struct {
		Title string `yaml:"title"`
		URL   string `yaml:"url"`
	} `yaml:"navigation"`
	Analytics struct {
		GoatCounter string `yaml:"goatcounter"`
	} `yaml:"analytics"`
}

// Options provides CLI flag overrides for configuration.
type Options struct {
	ContentDir string
	OutputDir  string
	AssetsDir  string
}

// Load reads configuration from a YAML file.
func Load(path string) (*model.Config, error) {
	return LoadWithOptions(path, Options{})
}

// LoadWithOptions reads configuration with CLI overrides.
func LoadWithOptions(path string, opts Options) (*model.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var yc yamlConfig
	if err := yaml.Unmarshal(data, &yc); err != nil {
		return nil, err
	}

	if yc.Site.Title == "" {
		return nil, errors.New("config: missing required field 'site.title'")
	}
	if yc.Site.BaseURL == "" {
		return nil, errors.New("config: missing required field 'site.baseURL'")
	}

	cfg := &model.Config{
		Title:       yc.Site.Title,
		Description: yc.Site.Description,
		BaseURL:     yc.Site.BaseURL,
		Author:      yc.Site.Author,
		Logo:        yc.Site.Logo,
		Favicon:     yc.Site.Favicon,
		ContentDir:  yc.Build.Content,
		OutputDir:   yc.Build.Output,
		AssetsDir:   yc.Build.Assets,
		Analytics: model.Analytics{
			GoatCounter: yc.Analytics.GoatCounter,
		},
	}

	// Apply defaults
	if cfg.ContentDir == "" {
		cfg.ContentDir = "content"
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = "public"
	}
	if cfg.AssetsDir == "" {
		cfg.AssetsDir = "assets"
	}

	// Apply overrides
	if opts.ContentDir != "" {
		cfg.ContentDir = opts.ContentDir
	}
	if opts.OutputDir != "" {
		cfg.OutputDir = opts.OutputDir
	}
	if opts.AssetsDir != "" {
		cfg.AssetsDir = opts.AssetsDir
	}

	// Convert navigation
	for _, nav := range yc.Navigation {
		cfg.Navigation = append(cfg.Navigation, model.NavItem{
			Title: nav.Title,
			URL:   nav.URL,
		})
	}

	return cfg, nil
}
