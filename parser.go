package main

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/rs/zerolog/log"
)

var obisRegexp = regexp.MustCompile(`([0-9:.-]+)\(([\d.]+)\*(\w+)\)`)

type ObisMessage struct {
	ID          string  `json:"id"`
	Value       float64 `json:"value"`
	Unit        string  `json:"unit"`
	Description string  `json:"description,omitempty"`
}

func ParseObisLine(line string) *ObisMessage {
	match := obisRegexp.FindStringSubmatch(line)
	if len(match) != 4 {
		return nil
	}
	value, err := strconv.ParseFloat(match[2], 64)
	if err != nil {
		log.Error().Err(err).Msgf("failed to convert %s to float", match[2])
		return nil
	}
	return &ObisMessage{
		ID:          match[1],
		Value:       value,
		Unit:        match[3],
		Description: OBISDescriptions[match[1]],
	}
}

func MustJSON(m *ObisMessage) string {
	bytes, _ := json.Marshal(m)
	return string(bytes)
}
