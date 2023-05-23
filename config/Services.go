package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type Config struct {
	Model          string  `json:"model"`
	Cfg            string  `json:"cfg"`
	Feed           string  `json:"feed"`
	Classnames     string  `json:"classnames"`
	GpuEnabled     bool    `json:"gpuEnabled"`
	ScoreThreshold float32 `json:"scoreThreshold"`
	NmsThreshold   float32 `json:"nmsThreshold"`
}

func LoadConfig() (*Config, error) {
	configFile, err := os.OpenFile("./config/aids.cfg", os.O_RDONLY, 0666)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	configBytes, _ := io.ReadAll(configFile)

	var config Config = Config{}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {

		return nil, err
	}

	log.Printf("\nLoaded Config: \t%+v", config)

	return &config, nil
}
