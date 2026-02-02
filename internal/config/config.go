package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ENV    string         `yaml:"env" env-default:"development"`
	DB     DatabaseConfig `yaml:"database"`
	Server HttpServer     `yaml:"http_server"`
}

type HttpServer struct {
	Address string `yaml:"address" default:"localhost"`
	Port    string `yaml:"port" default:"8080"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" default:"localhost"`
	Port     string `yaml:"port" default:"5432"`
	DBName   string `yaml:"dbname" default:"postgres"`
	User     string `yaml:"user" default:"root"`
	Password string `yaml:"password" default:""`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
		//configPath = "./config/config.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("congig file does not exist: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	log.Printf("Config loaded: DB Host=%s, Port=%s", config.DB.Host, config.DB.Port)

	return &config
}
