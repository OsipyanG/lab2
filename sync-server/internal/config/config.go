package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "configs/sync_server_config.json"

type Config struct {
	Host string `json:"host" env-default:"localhost"`
	Port string `json:"port" env-default:"8080"`

	ReadTimeout  string `json:"read_timeout" env-default:"60s"`
	WriteTimeout string `json:"write_timeout" env-default:"60s"`

	LogLevel      string `json:"log_level" env-default:"info"`
	ImitationTime string `json:"imitation_time" env-default:"5s"`

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
