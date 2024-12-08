package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const configPath = "configs/sync_server_config.json"

type Config struct {
	Host string `env-default:"localhost" json:"host"`
	Port string `env-default:"8080"      json:"port"`

	ReadWriteTimeout string `env-default:"60s" json:"read_write_timeout"`

	LogLevel      string `env-default:"info" json:"log_level"`
	ImitationTime string `env-default:"5s"   json:"imitation_time"`

	ReadWriteTimeoutDuration time.Duration `json:"-"`
	ImitationTimeDuration    time.Duration `json:"-"`
}

func MustLoad() *Config {
	cfg := &Config{}

	fileBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("cannot open config file: %v", err)
	}

	err = json.Unmarshal(fileBytes, cfg)
	if err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
	}

	cfg.ReadWriteTimeoutDuration, err = time.ParseDuration(cfg.ReadWriteTimeout)
	if err != nil {
		log.Fatalf("cannot parse ReadWriteTimeout, err=%s, ReadWriteTimeout=%s", err.Error(), cfg.ReadWriteTimeout)
	}

	cfg.ImitationTimeDuration, err = time.ParseDuration(cfg.ImitationTime)
	if err != nil {
		log.Fatalf("cannot parse ImitationTime, err=%s, ImitationTime=%s", err.Error(), cfg.ImitationTime)
	}

	return cfg
}
