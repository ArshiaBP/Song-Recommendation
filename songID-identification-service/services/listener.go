package services

import (
	"songID-identification-service/configs"
	"strconv"
)

func Listen() {
	messages, err := configs.Ch.Consume(configs.Queue.Name, "", false, false, false, false, nil)
	if err != nil {
		Listen()
	}
	quit := make(chan struct{})
	go func() {
		for msg := range messages {
			requestID, _ := strconv.Atoi(string(msg.Body))
			msg.Ack(true)
			go func() {
				err = saveSpotifyID(requestID)
				if err != nil {
					configs.DB.Table("request_infos").Where("id = ?", requestID).Update("status", "failure")
				}
			}()
		}
	}()
	<-quit
}
