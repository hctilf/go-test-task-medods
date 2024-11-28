package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Env        `yaml:"env" mapstructure:"env"`
		HTTPServer `yaml:"http_server" mapstructure:"http_server"`
		Postgres   `yaml:"postgres" mapstructure:"postgres"`
	}

	Env struct {
		Environment string `yaml:"environment" mapstructure:"environment"`
		LogLevel    string `yaml:"log_level" mapstructure:"log_level"`
	}

	HTTPServer struct {
		Address      string        `yaml:"address" mapstructure:"address"`
		Timeout      time.Duration `yaml:"timeout" mapstructure:"timeout"`
		IdleTimeout  time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
		ReadTimeout  time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
	}

	Postgres struct {
		Host     string `yaml:"host" mapstructure:"host"`
		Port     string `yaml:"port" mapstructure:"port"`
		User     string `yaml:"user" mapstructure:"user"`
		Password string `yaml:"password" mapstructure:"password"`
		Database string `yaml:"database" mapstructure:"database"`
		SSLMode  string `yaml:"sslmode" mapstructure:"sslmode"`
	}
)

var (
	ConfigInstance *Config
	once           sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		ConfigInstance = newConfig()
	})

	return ConfigInstance
}

func newConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	configName := os.Getenv("CONFIG_NAME")
	if configPath == "" {
		log.Fatal("CONFIG_NAME is not set")
	}
	configExt := os.Getenv("CONFIG_EXT")
	if configPath == "" {
		log.Fatal("CONFIG_EXT is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var config Config

	viper.SetConfigName(configName)
	viper.SetConfigType(configExt)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.UnmarshalExact(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &config
}
