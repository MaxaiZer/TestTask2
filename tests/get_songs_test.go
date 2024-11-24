package integration

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"test-task/internal/dto"
	"testing"
)

func TestGetSongs_ByGroup(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/songs?group=Wings", nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetSongsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Songs)
	assert.Equal(t, "Wings", response.Songs[0].Group)
	assert.Equal(t, "Live and Let Die", response.Songs[0].Name)
}

func TestGetSongs_ByGroup_CaseInsensitive(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/songs?group=wings", nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetSongsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Songs)
	assert.Equal(t, "Wings", response.Songs[0].Group)
	assert.Equal(t, "Live and Let Die", response.Songs[0].Name)
}

func TestGetSongs_ByName(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/songs?name=Live and Let Die", nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetSongsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Songs)
	assert.Equal(t, "Wings", response.Songs[0].Group)
	assert.Equal(t, "Live and Let Die", response.Songs[0].Name)
}

func TestGetSongs_ByName_CaseInsensitive(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/songs?name=live and Let die", nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetSongsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Songs)
	assert.Equal(t, "Wings", response.Songs[0].Group)
	assert.Equal(t, "Live and Let Die", response.Songs[0].Name)
}

func TestGetSongs_ByReleaseAfter(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/songs?release_after=2013-01-01", nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetSongsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Songs)
	assert.Equal(t, "Swans", response.Songs[0].Group)
}

func TestGetSongs_ByReleaseBefore(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/songs?release_before=2000-01-01", nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetSongsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Songs)
	assert.Equal(t, "Wings", response.Songs[0].Group)
}

func TestGetSongs_ByLink(t *testing.T) {

	link := "https://www.youtube.com/watch?v=wYQZHNwIUq8"

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/songs?link=%s", link), nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetSongsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Songs)
	for _, song := range response.Songs {
		assert.Equal(t, link, song.Link)
	}
}

func TestGetSongs_Pagination(t *testing.T) {

	pageSize := 1
	url := fmt.Sprintf("/api/v1/songs?page=2&page_size=%d", pageSize)

	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetSongsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, pageSize, len(response.Songs))
	assert.Equal(t, "Swans", response.Songs[0].Group)
}

func TestGetSongs_InvalidPagination(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/songs?page=1&page_size=-1", nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("GET", "/api/v1/songs?page=-1&page_size=1", nil)
	w = httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("GET", "/api/v1/songs?page=1&page_size=100000", nil)
	w = httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
