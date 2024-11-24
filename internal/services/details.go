package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"test-task/internal/dto"
)

type songDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

var details = make(map[dto.AddSongRequest]songDetail)

func mockDetailsAPI() *httptest.Server {

	fillDetails()

	handler := http.NewServeMux()
	handler.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {

		group := r.URL.Query().Get("group")
		song := r.URL.Query().Get("song")

		if group == "" || song == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		group = strings.ToLower(group)
		song = strings.ToLower(song)

		response, ok := details[dto.AddSongRequest{Group: group, Name: song}]
		if !ok {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	return httptest.NewServer(handler)
}

func fillDetails() {

	details[dto.AddSongRequest{Group: "testgroup", Name: "testname"}] = songDetail{
		ReleaseDate: "16.07.2006",
		Text:        "lalalala\nlalalalala\nlallalalala\nlalalalala\n\nnanananana\nnanananan\nnananana\nnanananan",
		Link:        "https://www.youtube.com/watch?v=VKTbMpkEM6A",
	}

	details[dto.AddSongRequest{Group: "the beatles", Name: "let it be"}] = songDetail{
		ReleaseDate: "06.03.1970",
		Text: "When I find myself in times of trouble\nMother Mary comes to me\nSpeaking words of wisdom\n" +
			"Let it be\nAnd in my hour of darkness\nShe is standing right in front of me\nSpeaking words of wisdom\n" +
			"Let it be\n\nLet it be, let it be\nLet it be, let it be\nWhisper words of wisdom\nLet it be\n\n" +
			"And when the broken-hearted people\nLiving in the world agree\nThere will be an answer\nLet it be\n" +
			"For though they may be parted\nThere is still a chance that they will see\nThere will be an answer\n" +
			"Let it be\n\nLet it be, let it be\nLet it be, let it be\nYeah, there will be an answer\nLet it be\n" +
			"Let it be, let it be\nLet it be, let it be\nWhisper words of wisdom\nLet it be\n\nAnd when the night is cloudy\n" +
			"There is still a light that shines on me\nShine until tomorrow\nLet it be\nI wake up to the sound of music\n" +
			"Mother Mary comes to me\nSpeaking words of wisdom\nLet it be, yeah\n\nLet it be, let it be\nLet it be, yeah, let it be\n" +
			"Oh, there will be an answer\nLet it be\nLet it be, let it be\nLet it be, yeah, let it be\nOh, there will be an answer\n" +
			"Let it be\nLet it be, let it be\nLet it be, yeah, let it be\nWhisper words of wisdom\nLet it be",
		Link: "https://www.youtube.com/watch?v=QDYfEBY9NM4",
	}
}
