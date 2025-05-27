package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LogLevel     int  `yaml:"logLevel"`
	SendUnparsed bool `yaml:"sendUnparsed"`
	Serial       struct {
		Path     string `yaml:"path"`
		BaudRate int    `yaml:"baudRate"`
		DataBits int    `yaml:"dataBits"`
		StopBits int    `yaml:"stopBits"`
		Parity   string `yaml:"parity"`
	} `yaml:"serial"`
	MQTT struct {
		URL         string `yaml:"url"`
		TopicPrefix string `yaml:"topicPrefix"`
		ClientID    string `yaml:"clientId"`
		Options     struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"options"`
	} `yaml:"mqtt"`
}

func LoadConfig() (*Config, error) {
	file := os.Getenv("HAN_MQTT_CONFIG")
	if file == "" {
		file = "config.yaml"
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing yaml: %w", err)
	}

	parsed, err := url.Parse(cfg.MQTT.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid MQTT URL: %w", err)
	}
	host := parsed.Host
	if !strings.Contains(host, ":") {
		parsed.Host = parsed.Host + ":1883"
		cfg.MQTT.URL = parsed.String()
	}

	return &cfg, nil
}
