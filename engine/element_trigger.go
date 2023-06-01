package engine

import (
	log "github.com/sirupsen/logrus"
)

func Trigger() {
	log.Info("Trigger Called")
}

func CheckThreshold(confidences []float32, triggerThreshold float32) {
	for _, confidence := range confidences {
		if confidence >= triggerThreshold {
			Trigger()
			break
		}
	}
}
