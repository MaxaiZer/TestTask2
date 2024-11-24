package integration

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"test-task/internal/dto"
	"testing"
)

func TestAddSong(t *testing.T) {

	request := dto.AddSongRequest{
		Group: "TestGroup",
		Name:  "TestName",
	}
	_, err := AddSong(request)
	assert.NoError(t, err)
}

func TestAddSong_WhenMissingField(t *testing.T) {

	request := "{ \"group\":\"TestGroup\" }"

	req, _ := http.NewRequest("POST", "/api/v1/songs", bytes.NewBuffer([]byte(request)))
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
