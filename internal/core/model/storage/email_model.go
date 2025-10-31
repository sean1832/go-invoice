package storage

import (
	"encoding/json"
	"fmt"
)

type EmailConfig struct {
	Id         string `json:"id"`
	SMTPConfig struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"smtp"`
	DefaultMessage struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
		Footer  string `json:"footer"`
	} `json:"default_message"`
	Recipients []string `json:"recipients"`
}

func NewEmailConfigFromJSON(data []byte) (*EmailConfig, error) {
	var config EmailConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal email config JSON data: %v", err)
	}
	return &config, nil
}
