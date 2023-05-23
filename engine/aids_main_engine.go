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
		Detector: InitializeDetector(AIDSConfig),
		Config:   *AIDSConfig,
	}
}

func (aids *AIDSEngine) start() error {
	aids.Detector.Load()
	aids.Detector.Process()
	return nil
}

func (aids *AIDSEngine) stop() error {
	aids.Detector.Close()
	return nil
}
