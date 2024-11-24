package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

		response := details[dto.AddSongRequest{Group: group, Name: song}]
		empty := songDetail{}
		if response == empty {
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

	details[dto.AddSongRequest{Group: "TestGroup", Name: "TestName"}] = songDetail{
		ReleaseDate: "16.07.2006",
		Text:        "lalalala\nlalalalala\nlallalalala\nlalalalala\n\nnanananana\nnanananan\nnananana\nnanananan",
		Link:        "https://www.youtube.com/watch?v=VKTbMpkEM6A",
	}
}
