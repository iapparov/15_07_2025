package config

import (
	"log"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
    Env	string	`yaml:"env" env-default:"local"`
    Http_port	int	`yaml:"http_port" env-default:"8080"`
	Max_tasks	int `yaml:"max_tasks" env-default:"3"`
	Allowed_types []string `yaml:"allowed_types" env-default:"jpg,pdf"`
	Max_objects	int	`yaml:"max_objects" env-default:"3"`
}


func MustLoad() *Config{
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == ""{
		log.Fatalf("CONFIG_PATH env is required")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("%s", "failed to read config file:" + err.Error())
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil{
		log.Fatalf("%s", "failed to unmarshal config: " + err.Error())
	}
	return &cfg
}