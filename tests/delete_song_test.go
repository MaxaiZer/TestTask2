package integration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"test-task/internal/dto"
	"testing"
)

func TestDeleteSong(t *testing.T) {

	id, err := AddSong(dto.AddSongRequest{
		Group: "TestGroup",
		Name:  "TestName",
	})
	assert.NoError(t, err)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/songs/%d", id), nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/songs/%d", id), nil)
	w = httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSong_WhenDoesntExist(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/api/v1/songs/123456", nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
