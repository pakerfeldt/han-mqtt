package main

import (
	"bufio"

	"os"
	"os/signal"
	"strings"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tarm/serial"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04",
	})

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal().Msgf("Error loading config: %v", err)
	}

	zerolog.SetGlobalLevel(zerolog.Level(cfg.LogLevel))

	log.Info().Msg("Starting han-mqtt")
	opts := mqtt.NewClientOptions().AddBroker(cfg.MQTT.URL)
	if cfg.MQTT.ClientID == "" {
		cfg.MQTT.ClientID = "han-mqtt-client-001"
	}
	opts.SetClientID(cfg.MQTT.ClientID)

	if cfg.MQTT.Options.Username != "" {
		opts.SetUsername(cfg.MQTT.Options.Username)
	}
	if cfg.MQTT.Options.Password != "" {
		opts.SetPassword(cfg.MQTT.Options.Password)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Msgf("Error connecting to MQTT: %v", token.Error())
	}

	serialConfig := &serial.Config{Name: cfg.Serial.Path, Baud: cfg.Serial.BaudRate}
	s, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatal().Msgf("Error opening serial port: %v", err)
	}
	defer s.Close()

	prefix := cfg.MQTT.TopicPrefix
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		client.Disconnect(250)
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		line := scanner.Text()
		msg := ParseObisLine(line)
		if msg == nil {
			continue
		}

		topic := prefix + msg.ID
		var payload string
		if cfg.SendUnparsed {
			payload = line
		} else {
			msg.Description = OBISDescriptions[msg.ID]
			payload = MustJSON(msg)
		}

		client.Publish(topic, 0, false, payload)
		log.Debug().Msgf("(%s): %s", topic, payload)
	}

	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("Error reading serial")
	}
}
