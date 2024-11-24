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

func TestGetVerses_Pagination(t *testing.T) {

	pageSize := 1
	url := fmt.Sprintf("/api/v1/songs/1/lyrics?page=2&page_size=%d", pageSize)

	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetSongVersesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, pageSize, len(response.Verses))
	assert.Contains(t, response.Verses[0], "Say live and let die")
}

func TestGetVerses_InvalidPagination(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/v1/songs/1/lyrics?page=1&page_size=-1", nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("GET", "/api/v1/songs/1/lyrics?page=-1&page_size=1", nil)
	w = httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("GET", "/api/v1/songs/1/lyrics?page=1&page_size=100000", nil)
	w = httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestGetVerses_WhenSongDoesntExist(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/songs/123456/lyrics", nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
