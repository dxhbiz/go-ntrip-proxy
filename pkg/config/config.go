package config

import (
	"path"

	"github.com/dxhbiz/go-ntrip-proxy/pkg/kit/exe"
	"github.com/spf13/viper"
)

var (
	exePath = exe.Path()
)

// Config
type Config struct {
	Server  ServerConfig   `mapstructure:"server", json:"server"`
	Casters []CasterConfig `mapstructure:"casters", json:"casters"`
	Log     LogConfig      `mapstructure:"log", json:"log"`
}

// ServerConfig
type ServerConfig struct {
	Host string `mapstructure:"host", json:"host"`
	Port uint16 `mapstructure:"port", json:"port"`
}

// CasterConfig
type CasterConfig struct {
	Name       string `mapstructure:"name", json:"name"`
	Host       string `mapstructure:"host", json:"host"`
	Port       uint16 `mapstructure:"port", json:"port"`
	Username   string `mapstructure:"username", json:"username"`
	Password   string `mapstructure:"password", json:"password"`
	Mountpoint string `mapstructure:"mountpoint", json:"mountpoint"`
}

// LogConfig
type LogConfig struct {
	Development bool   `mapstructure:"development"`
	Level       string `mapstructure:"level"`
	Filename    string `mapstructure:"filename"`
	MaxSize     int    `mapstructure:"maxSize"`
	MaxBackups  int    `mapstructure:"maxBackups"`
	MaxAge      int    `mapstructure:"maxAge"`
	Compress    bool   `mapstructure:"compress"`
}

var cfg Config

// InitConfig
func InitConfig(configFile string) error {
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	isAbs := path.IsAbs(cfg.Log.Filename)
	if !isAbs {
		cfg.Log.Filename = path.Join(exePath, cfg.Log.Filename)
	}

	return nil
}

// GetConfig
func GetConfig() Config {
	return cfg
}
