package viper

import (
	"errors"
	"fmt"
	"gochat/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfigFile(baseConfigPath string, overlays ...string) (*config.App, error) {
	v := viper.New()

	// load base config
	if err := readConfigInto(v, baseConfigPath, false); err != nil {
		return nil, err
	}

	// merge optional overlays; later overlays win
	for _, overlay := range overlays {
		if err := mergeIfExists(v, overlay); err != nil {
			return nil, err
		}
	}

	cfg := new(config.App)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("parse config file failed: %w", err)
	}

	return cfg, cfg.Validate()
}

func readConfigInto(v *viper.Viper, path string, merge bool) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file failed: %w", err)
	}

	expandedContent := os.ExpandEnv(string(content))
	cfgType := configTypeFor(path)
	v.SetConfigType(cfgType)

	if merge {
		if err := v.MergeConfig(strings.NewReader(expandedContent)); err != nil {
			return fmt.Errorf("merge config failed: %w", err)
		}
		return nil
	}

	if err := v.ReadConfig(strings.NewReader(expandedContent)); err != nil {
		return fmt.Errorf("read config failed: %w", err)
	}
	return nil
}

func mergeIfExists(v *viper.Viper, path string) error {
	if path == "" {
		return nil
	}

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat config overlay failed: %w", err)
	}
	return readConfigInto(v, path, true)
}

func configTypeFor(path string) string {
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	if ext == "" {
		return "yaml"
	}
	return ext
}
