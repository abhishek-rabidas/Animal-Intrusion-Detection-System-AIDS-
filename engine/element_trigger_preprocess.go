package engine

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Trigger(class string) {
	msg := DetectionMessage{
		Timestamp: time.Now().Format("02-01-2006 15:04:05"),
		Class:     class,
	}

	payload, _ := json.Marshal(msg)

	post, err := http.Post("http://localhost:55555", "application/json", bytes.NewReader(payload))
	if err != nil || post.StatusCode != 200 {
		log.Error(err)
	}
}

func CheckThreshold(confidences []float32, triggerThreshold float32, classes []string, classIds []int) {
	for _, classId := range classIds {
		//TODO:Change to animal/tiger/leopard
		if classes[classId] == "truck" && confidences[classId] >= triggerThreshold {
			Trigger("truck")
			break
		}
	}
}
