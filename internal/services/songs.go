package services

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"test-task/internal/dto"
	"test-task/internal/entities"
	errs "test-task/internal/errors"
	"time"
)

type songsRepository interface {
	Add(ctx context.Context, song entities.Song, verses []string) (int64, error)
	Get(ctx context.Context, request dto.GetSongsRequest) ([]entities.Song, error)
	GetVerses(ctx context.Context, songID int64, page int, pageSize int) ([]entities.Verse, error)
	Update(ctx context.Context, song entities.Song) error
	Delete(ctx context.Context, ID int64) error
}

type SongsService struct {
	songs          songsRepository
	songDetailsAPI *httptest.Server
}

func NewSongsService(songs songsRepository) *SongsService {
	return &SongsService{songs: songs, songDetailsAPI: mockDetailsAPI()}
}

func (s *SongsService) AddSong(ctx context.Context, request dto.AddSongRequest) (int64, error) {

	log.Debugf("AddSong: %v", request)

	url := fmt.Sprintf("%s/info?group=%s&song=%s",
		s.songDetailsAPI.URL,
		url.QueryEscape(request.Group),
		url.QueryEscape(request.Name),
	)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("couldn't send request to song details API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return 0, fmt.Errorf("%w: song doesn't exist", errs.BadRequest)
		}
		return 0, fmt.Errorf("unexpected status from song details API, got: %d", resp.StatusCode)
	}

	var info songDetail
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}
	err = json.Unmarshal(body, &info)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshel song details: %w", err)
	}

	songRelease, err := time.Parse("02.01.2006", info.ReleaseDate)
	if err != nil {
		return 0, fmt.Errorf("failed to parse release date: %w", err)
	}

	song := entities.Song{
		Group:       request.Group,
		Name:        request.Name,
		Link:        info.Link,
		ReleaseDate: songRelease,
	}

	id, err := s.songs.Add(ctx, song, strings.Split(info.Text, "\n\n"))
	return id, err
}

func (s *SongsService) EditSong(ctx context.Context, ID int64, request dto.EditSongRequest) error {

	log.Debugf("EditSong: %v", request)

	song := entities.Song{
		ID:          ID,
		Group:       request.Group,
		Name:        request.Name,
		Link:        request.Link,
		ReleaseDate: request.ReleaseDate,
	}

	return s.songs.Update(ctx, song)
}

func (s *SongsService) DeleteSong(ctx context.Context, id int64) error {
	log.Debugf("DeleteSong: %v", id)
	return s.songs.Delete(ctx, id)
}

func (s *SongsService) GetSongs(ctx context.Context, request dto.GetSongsRequest) ([]entities.Song, error) {
	log.Debugf("GetSongs: %v", request)
	return s.songs.Get(ctx, request)
}

func (s *SongsService) GetSongVerses(ctx context.Context, id int64, limit int, offset int) ([]entities.Verse, error) {
	log.Debugf("GetSongVerses: id:%v limit:%v offset:%v", id, limit, offset)
	return s.songs.GetVerses(ctx, id, limit, offset)
}
