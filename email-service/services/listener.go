package services

import (
	"email-service/configs"
	"email-service/models"
	"fmt"
	"log"
	"time"
)

func Listen() {
	for {
		go func() {
			requests, dbErr := CheckStatus()
			for _, req := range requests {
				if dbErr != nil && req.Status != "failure" {
					configs.DB.Unscoped().Where("email = ?", req.Email).Delete(&models.RequestInfo{})
					go func() {
						msg, _ := SendMail("song recommendation failed", "song recommendation process failed!\nplease try again", req.Email)
						log.Println(fmt.Sprint("message: " + msg))
					}()
				}
				switch req.Status {
				case "failure":
					configs.DB.Unscoped().Where("email = ?", req.Email).Delete(&models.RequestInfo{})
					go func() {
						msg, _ := SendMail("song recommendation failed", "song recommendation process failed!\nplease try again", req.Email)
						log.Println(fmt.Sprint("message: " + msg))
					}()
				case "ready":
					songs, err := recommendationList(req.SongID)
					if err != nil {
						configs.DB.Unscoped().Where("email = ?", req.Email).Delete(&models.RequestInfo{})
						go func() {
							msg, _ := SendMail("song recommendation failed", "song recommendation process failed!\nplease try again", req.Email)
							log.Println(fmt.Sprint("message: " + msg))
						}()
					}
					body := "suggested songs according to your request:\n"
					configs.DB.Table("request_infos").Where("email = ?", req.Email).Update("status", "done")
					for i := 0; i < 5; i++ {
						body = fmt.Sprintf("%s%d.%s\n", body, i+1, songs[i])
					}
					go func(mailBody string) {
						msg, _ := SendMail("recommended songs", mailBody, req.Email)
						log.Println(fmt.Sprint("message: " + msg))
					}(body)
				}
			}
		}()
		time.Sleep(time.Second * 10)
	}
}
