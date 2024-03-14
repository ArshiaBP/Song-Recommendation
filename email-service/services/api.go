package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	spotifyBaseURL = "https://spotify23.p.rapidapi.com/recommendations/?limit=5&seed_tracks="
	spotifyHost    = "spotify23.p.rapidapi.com"
	domain         = "sandboxde986d9ae0c445fbab0cd848eb5ad1bd.mailgun.org"
)

func recommendationList(songID string) ([]string, error) {
	type Track struct {
		Name string `json:"name"`
	}
	var res struct {
		Tracks []Track `json:"tracks"`
	}
	songs := make([]string, 5)
	spotifyURL := fmt.Sprintf("%s%s", spotifyBaseURL, songID)
	req, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		return []string{}, err
	}
	req.Header.Add("X-RapidAPI-Key", os.Getenv("API_Key"))
	req.Header.Add("X-RapidAPI-Host", spotifyHost)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("spotify recommender failed")
		return []string{}, err
	}
	resBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("spotify recommender failed")
		return []string{}, err
	}
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		log.Println("spotify recommender failed")
		return []string{}, err
	}
	for i := 0; i < 5; i++ {
		songs[i] = res.Tracks[i].Name
	}
	log.Println("spotify recommender worked successfully")
	return songs, nil
}

func SendMail(subject, body, to string) (string, error) {
	mg := mailgun.NewMailgun(domain, os.Getenv("MailGun_API_Key"))
	m := mg.NewMessage(
		fmt.Sprint("Song Recommender <song-recommender@"+domain+">"),
		subject,
		body,
		to,
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	msg, _, err := mg.Send(ctx, m)
	return msg, err
}
