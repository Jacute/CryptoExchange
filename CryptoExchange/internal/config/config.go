package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env" json:"env" env:"ENV" env-default:"local"`
	AppConfig      `yaml:"app" json:"app_config" env-required:"true"`
	DatabaseConfig `yaml:"database" json:"db_config" env-required:"true"`
}

type AppConfig struct {
	IP           string        `yaml:"ip" json:"ip" env:"APP_IP" env-default:"127.0.0.1"`
	Port         int           `yaml:"port" json:"port" env:"APP_PORT" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout" env-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" json:"idle_timeout" env-default:"60s"`
	Lots         []string      `yaml:"lots" json:"lots" env-required:"true"`
	TokenLen     int           `yaml:"token_length" json:"token_length" env-default:"32"`
}

type DatabaseConfig struct {
	IP   string `yaml:"ip" json:"ip" env:"DATABASE_IP" env-default:"127.0.0.1"`
	Port int    `yaml:"port" json:"port" env:"DATABASE_PORT" env-default:"7432"`
}

func MustLoad() Config {
	path := getConfigPath()
	if path == "" {
		panic("config path not provided")
	}
	cfg := MustLoadByPath(path)

	return cfg
}

func MustLoadByPath(path string) Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config not found: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Can't load config: " + err.Error())
	}

	return cfg
}

func getConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to the config file")
	flag.Parse()
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
