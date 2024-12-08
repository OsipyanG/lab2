package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const configPath = "configs/server_config.json"

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`

	WriteTimeout string `json:"write_timeout"`
	ReadTimeout  string `json:"read_timeout"`

	LogLevel      string `json:"log_level"`
	ImitationTime string `json:"imitation_time"`

	ReadTimeoutDuration   time.Duration `json:"-"`
	WriteTimeoutDuration  time.Duration `json:"-"`
	ImitationTimeDuration time.Duration `json:"-"`
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

	cfg.ReadTimeoutDuration, err = time.ParseDuration(cfg.ReadTimeout)
	if err != nil {
		log.Fatalf("cannot parse ReadTimeout, err=%s, ReadTimeout=%s", err.Error(), cfg.ReadTimeout)
	}

	cfg.WriteTimeoutDuration, err = time.ParseDuration(cfg.WriteTimeout)
	if err != nil {
		log.Fatalf("cannot parse WriteTimeout, err=%s, WriteTimeout=%s", err.Error(), cfg.WriteTimeout)
	}

	cfg.ImitationTimeDuration, err = time.ParseDuration(cfg.ImitationTime)
	if err != nil {
		log.Fatalf("cannot parse ImitationTime, err=%s, ImitationTime=%s", err.Error(), cfg.ImitationTime)
	}

	return cfg
}