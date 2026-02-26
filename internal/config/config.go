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
	Root string `mapstructure:"root" validate:"dirpath"`
}

type HTTPConfig struct {
	Addr string `mapstructure:"addr" validate:"ip"`
	Port uint   `mapstructure:"port" validate:"port"`
}

type CacheConfig struct {
	// Max cache size in bytes
	MaxSize string `mapstructure:"max_size" validate:"byte_size"`
	TTL     string `mapstructure:"ttl" validate:"duration"`
}

type LogConfig struct {
	Level string `mapstructure:"level" validate:"oneof=debug warn info error none"`
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

func (cfg *Config) Validate() (validator.ValidationErrors, bool) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("byte_size", validateByteSize)
	validate.RegisterValidation("duration", validateDuration)
	err := validate.Struct(cfg)
	if err != nil {
		return err.(validator.ValidationErrors), false
	}
	return nil, true
}
