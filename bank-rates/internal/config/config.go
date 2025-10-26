package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env"`

	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`

	Database struct {
		Host     string `yaml:"host"`
		HostProd string `yaml:"host_prod"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	Migrations struct {
		Path string `yaml:"path"`
	} `yaml:"migrations"`

	NotifyServiceAddr string `yaml:"NotifyServiceAddr"`
}

// MustLoad загружает YAML-конфиг или падает, если не удалось
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty (use --config or CONFIG_PATH)")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	// Пример запуска: ./app --config=./config.yaml
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH") // можно удалить, если хочешь убрать даже эту опцию
	}

	return res
}
