package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "configs/async_server_config.json"

type Config struct {
	Host string `env-default:"localhost" json:"host"`
	Port string `env-default:"8080"      json:"port"`

	ReadTimeout  string `env-default:"60s" json:"read_timeout"`
	WriteTimeout string `env-default:"60s" json:"write_timeout"`

	LogLevel      string `env-default:"info" json:"log_level"`
	ImitationTime string `env-default:"5s"   json:"imitation_time"`

	ReadTimeoutDuration   time.Duration `json:"-"`
	WriteTimeoutDuration  time.Duration `json:"-"`
	ImitationTimeDuration time.Duration `json:"-"`
}

func MustLoad() *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	// Convert string durations to time.Duration
	var err error

	cfg.ReadTimeoutDuration, err = time.ParseDuration(cfg.ReadTimeout)
	if err != nil {
		log.Fatalf("invalid read_timeout: %v", err)
	}

	cfg.WriteTimeoutDuration, err = time.ParseDuration(cfg.WriteTimeout)
	if err != nil {
		log.Fatalf("invalid write_timeout: %v", err)
	}

	cfg.ImitationTimeDuration, err = time.ParseDuration(cfg.ImitationTime)
	if err != nil {
		log.Fatalf("invalid imitation_time: %v", err)
	}

	return cfg
}
