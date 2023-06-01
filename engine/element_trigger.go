package engine

import (
	log "github.com/sirupsen/logrus"
)

func Trigger() {
	log.Info("Trigger Called")
}

func CheckThreshold(confidences []float32, triggerThreshold float32, classes []string, classIds []int) {
	for _, classId := range classIds {
		//TODO:Change to animal/tiger/leopard
		if classes[classId] == "truck" && confidences[classId] >= triggerThreshold {
			Trigger()
			break
		}
	}
}
