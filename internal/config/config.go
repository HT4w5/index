package config

import (
	"fmt"

	"github.com/HT4w5/index/internal/meta"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Log        LogConfig        `mapstructure:"log"`
	Filesystem FileSystemConfig `mapstructure:"filesystem"`
	HTTP       HTTPConfig       `mapstructure:"http"`
	Cache      CacheConfig      `mapstructure:"cache"`
}

type FileSystemConfig struct {
	Root string `mapstructure:"root" validator:"dirpath"`
}

type HTTPConfig struct {
	Addr string `mapstructure:"addr" validator:"ip"`
	Port int    `mapstructure:"port" validator:"port"`
}

type CacheConfig struct {
	// Max cache size in bytes
	MaxSize int   `mapstructure:"max_size"`
	TTL     int64 `mapstructure:"ttl"`
}

type LogConfig struct {
	Level string `mapstructure:"level" validator:"oneof=debug warn info error none"`
}

func (cfg *Config) Load() error {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath(fmt.Sprintf("/etc/%s/", meta.Name))
	vp.AddConfigPath(".")

	err := vp.ReadInConfig()
	if err != nil {
		return err
	}

	return vp.Unmarshal(cfg)
}

func (cfg *Config) LoadFromPath(path string) error {
	vp := viper.New()
	vp.SetConfigFile(path)

	err := vp.ReadInConfig()
	if err != nil {
		return err
	}

	return vp.Unmarshal(cfg)
}

func (cfg *Config) Validate() ([]string, bool) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cfg)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		msgs := make([]string, 0, len(validationErrors))
		for _, e := range validationErrors {
			msgs = append(msgs, fmt.Sprintf(
				"Name: %s. Got: %s. Expected: %s. Reason: %s.",
				e.Field(),
				e.Value(),
				e.Param(),
				e.Tag(),
			))
		}
	}
	return nil, true
}
