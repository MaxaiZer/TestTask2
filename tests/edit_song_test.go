package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"test-task/internal/dto"
	"testing"
	"time"
)

func TestEditSong(t *testing.T) {

	id, err := AddSong(dto.AddSongRequest{
		Group: "TestGroup",
		Name:  "TestName",
	})
	assert.NoError(t, err)

	releaseDate, err := time.Parse("2006-01-02", "2024-11-23")
	assert.NoError(t, err)

	request := dto.EditSongRequest{
		Group:       "NewSuperPuperGroup",
		Name:        "NewSuperPuperSongName",
		Link:        "NewSuperPuperLink",
		ReleaseDate: releaseDate,
	}
	body, _ := json.Marshal(request)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/songs/%d", id), bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET",
		fmt.Sprintf("/api/v1/songs?group=%s&name=%s", request.Group, request.Name), nil)
	w = httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetSongsResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Songs)
}

func TestEditSong_WhenDoesntExist(t *testing.T) {

	releaseDate, err := time.Parse("2006-01-02", "2024-11-23")
	assert.NoError(t, err)

	request := dto.EditSongRequest{
		Group:       "NewSuperPuperGroup",
		Name:        "NewSuperPuperSongName",
		Link:        "NewSuperPuperLink",
		ReleaseDate: releaseDate,
	}
	body, _ := json.Marshal(request)

	req, _ := http.NewRequest("PUT", "/api/v1/songs/123456", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
