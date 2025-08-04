package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env"`
	StoragePath    string `yaml:"storage_path"`
	MigrationsPath string `yaml:"migrations_path"`
	Port           string `yaml:"port"`
}

func MustLoad() Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) Config {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config:" + err.Error())
	}
	return cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config_path", "", "path to config")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
