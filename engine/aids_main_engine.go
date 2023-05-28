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

func (aids *AIDSEngine) Start() error {
	err := aids.Detector.Load()
	if err != nil {
		log.Error("Loading Detection Failed, Process will not start")
		return err
	}
	aids.Detector.Process()
	return nil
}

func (aids *AIDSEngine) Stop() error {
	aids.Detector.Close()
	return nil
}
