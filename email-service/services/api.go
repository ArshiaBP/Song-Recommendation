package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	spotifyBaseURL = "https://spotify23.p.rapidapi.com/recommendations/?limit=5&seed_tracks="
	spotifyHost    = "spotify23.p.rapidapi.com"
)

func recommendationList(songID string) ([]string, error) {
	type track struct {
		Name string `json:"name"`
	}
	tracks := make([]track, 5)
	songs := make([]string, 5)
	spotifyURL := fmt.Sprintf("%s%s", spotifyBaseURL, songID)
	req, _ := http.NewRequest("GET", spotifyURL, nil)
	req.Header.Add("X-RapidAPI-Key", os.Getenv("API_Key"))
	req.Header.Add("X-RapidAPI-Host", spotifyHost)
	response, _ := http.DefaultClient.Do(req)
	defer response.Body.Close()
	for i := 0; i < 5; i++ {
		err := json.NewDecoder(response.Body).Decode(&(tracks[i]))
		if err != nil {
			return []string{}, err
		}
		songs[i] = tracks[i].Name
	}
	return songs, nil
}
