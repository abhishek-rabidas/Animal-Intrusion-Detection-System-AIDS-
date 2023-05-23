package engine

import (
	"AIDS/config"
	log "github.com/sirupsen/logrus"
)

type AIDSEngine struct {
	Detector *Detector
	Config   config.Config
}

func Initialize() *AIDSEngine {

	AIDSConfig, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	log.Info("Initializing AIDS Engine")

	return &AIDSEngine{
		Detector: nil,
		Config:   *AIDSConfig,
	}
}

func (aids *AIDSEngine) start() error {

	return nil
}

func (aids *AIDSEngine) stop() error {

	return nil
}
