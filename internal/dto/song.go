package dto

import (
	"test-task/internal/entities"
	"time"
)

type AddSongRequest struct {
	Group string `json:"group" binding:"required" example:"The Beatles"`
	Name  string `json:"name" binding:"required" example:"Let It Be"`
}

type AddSongResponse struct {
	ID int64 `json:"id"`
}

type GetSongsRequest struct {
	Group         string    `form:"group" example:"The Beatles"`
	Name          string    `form:"name" example:"Let It Be"`
	Link          string    `form:"link" example:"https://www.youtube.com/watch?v=dQw4w9WgXcQ"`
	ReleaseAfter  time.Time `form:"release_after" time_format:"2006-01-02" example:"2006-01-30"`
	ReleaseBefore time.Time `form:"release_before" time_format:"2006-01-02" example:"2006-01-30"`
	pagination
}

type GetSongsResponse struct {
	Songs []entities.Song `json:"songs"`
}

type EditSongRequest struct {
	Group       string    `json:"group" binding:"required" example:"The Beatles"`
	Name        string    `json:"name" binding:"required" example:"Let It Be"`
	Link        string    `json:"link" binding:"required" example:"https://www.youtube.com/watch?v=dQw4w9WgXcQ"`
	ReleaseDate time.Time `json:"release_date" binding:"required" time_format:"2006-01-02" example:"2006-01-30"`
}

type GetSongVersesRequest struct {
	pagination
}

type GetSongVersesResponse struct {
	Verses []string `json:"verses"`
}

type pagination struct {
	Page     int `form:"page" default:"1"`
	PageSize int `form:"page_size" default:"10"`
}
