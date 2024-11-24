package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"test-task/internal/dto"
)

func AddSong(request dto.AddSongRequest) (int64, error) {

	body, _ := json.Marshal(request)

	req, _ := http.NewRequest("POST", "/api/v1/songs", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	ginEngine.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("couldn't add song: status %v error %s", w.Code, response.Error)
	}

	var response dto.AddSongResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return 0, err
	}
	return response.ID, nil
}
