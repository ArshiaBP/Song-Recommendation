package services

import (
	"email-service/configs"
	"fmt"
	"time"
)

func Listen() {
	for {
		go func() {
			email, songID, status, err := checkStatus()
			if err != nil && status != "failure" {
				configs.DB.Unscoped().Where("email = ?", email).Delete("request_infos")
				go sendMail("song recommendation failed", "song recommendation process failed!\nplease try again", []string{email})
			}
			switch status {
			case "failure":
				configs.DB.Unscoped().Where("email = ?", email).Delete("request_infos")
				go sendMail("song recommendation failed", "song recommendation process failed!\nplease try again", []string{email})
			case "ready":
				songs, err := recommendationList(songID)
				if err != nil {
					configs.DB.Unscoped().Where("email = ?", email).Delete("request_infos")
					go sendMail("song recommendation failed", "song recommendation process failed!\nplease try again", []string{email})
				}
				body := "suggested songs according to your request:\n"
				for i := 0; i < 5; i++ {
					body = fmt.Sprintf("%s%d.%s\n", body, i+1, songs[i])
				}
				go sendMail("recommended songs", body, []string{email})
			}
		}()
		time.Sleep(time.Millisecond * 500)
	}
}
