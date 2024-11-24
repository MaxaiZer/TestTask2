package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
	"test-task/internal/dto"
	"test-task/internal/entities"
)

type songsService interface {
	AddSong(ctx context.Context, request dto.AddSongRequest) (int64, error)
	EditSong(ctx context.Context, ID int64, request dto.EditSongRequest) error
	DeleteSong(ctx context.Context, ID int64) error
	GetSongVerses(ctx context.Context, ID int64, limit int, offset int) ([]entities.Verse, error)
	GetSongs(ctx context.Context, request dto.GetSongsRequest) ([]entities.Song, error)
}

type SongsHandler struct {
	service songsService
}

const defaultPage = 1
const defaultPageSize = 10
const maxPageSize = 100

func NewSongsHandler(service songsService) *SongsHandler {
	return &SongsHandler{service: service}
}

// AddSong godoc
// @Summary Add a new song
// @Description Add a new song to the database
// @Tags songs
// @Accept  json
// @Produce  json
// @Param request body dto.AddSongRequest true "Song Data"
// @Success 201 {object} dto.AddSongResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /songs [post]
func (h *SongsHandler) AddSong(ctx *gin.Context) {

	var request dto.AddSongRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, dto.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := h.service.AddSong(ctx, request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	log.Infof("Song was added successfully, id: %d", id)
	ctx.JSON(201, dto.AddSongResponse{ID: id})
}

// EditSong godoc
// @Summary Edit an existing song
// @Description Edit the details of an existing song in the database
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int64 true "Song ID"
// @Param request body dto.EditSongRequest true "Updated Song Data"
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /songs/{id} [put]
func (h *SongsHandler) EditSong(ctx *gin.Context) {

	var request dto.EditSongRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, dto.ErrorResponse{Error: err.Error()})
		return
	}

	songID, ok := getIdFromRoute(ctx)
	if !ok {
		return
	}

	err := h.service.EditSong(ctx, songID, request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	log.Infof("Song was edited successfully, id: %d", songID)
}

// DeleteSong godoc
// @Summary Delete a song
// @Description Delete a song from the database
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int64 true "Song ID"
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /songs/{id} [delete]
func (h *SongsHandler) DeleteSong(ctx *gin.Context) {

	songID, ok := getIdFromRoute(ctx)
	if !ok {
		return
	}

	err := h.service.DeleteSong(ctx, songID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	log.Infof("Song was deleted successfully, id: %d", songID)
}

// GetSongs godoc
// @Summary Get all songs
// @Description Get a list of all songs, optionally filtered by query parameters
// @Tags songs
// @Accept  json
// @Produce  json
// @Param request query dto.GetSongsRequest false "Song Request Parameters"
// @Success 200 {object} dto.GetSongsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /songs [get]
func (h *SongsHandler) GetSongs(ctx *gin.Context) {

	var request dto.GetSongsRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(400, dto.ErrorResponse{Error: err.Error()})
		return
	}

	empty := dto.GetSongsRequest{}
	if request == empty {
		ctx.JSON(400, dto.ErrorResponse{Error: "request is empty"})
		return
	}

	if !validatePagination(&request.Page, &request.PageSize, ctx) {
		return
	}

	log.Debugf("Get songs request: %v", request)
	songs, err := h.service.GetSongs(ctx, request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(200, dto.GetSongsResponse{Songs: songs})
}

// GetSongVerses godoc
// @Summary Get verses of a song
// @Description Get a list of verses for a specific song
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int64 true "Song ID"
// @Param request query dto.GetSongVersesRequest false "Verses Request Parameters"
// @Success 200 {object} dto.GetSongVersesResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /songs/{id}/lyrics [get]
func (h *SongsHandler) GetSongVerses(ctx *gin.Context) {

	var request dto.GetSongVersesRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(400, dto.ErrorResponse{Error: err.Error()})
		return
	}

	songID, ok := getIdFromRoute(ctx)
	if !ok {
		return
	}

	if !validatePagination(&request.Page, &request.PageSize, ctx) {
		return
	}

	verses, err := h.service.GetSongVerses(ctx, songID, request.Page, request.PageSize)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	versesAsStr := make([]string, 0, len(verses))
	for i := 0; i < len(verses); i++ {
		versesAsStr = append(versesAsStr, verses[i].Text)
	}

	log.Infof("Song verses with number from %d to %d were returned successfully",
		verses[0].VerseNumber, verses[len(verses)-1].VerseNumber)
	ctx.JSON(200, dto.GetSongVersesResponse{Verses: versesAsStr})
}

func getIdFromRoute(ctx *gin.Context) (int64, bool) {
	ID := ctx.Param("id")

	if ID == "" {
		ctx.JSON(400, dto.ErrorResponse{Error: "ID is required in route"})
		return 0, false
	}

	IDint, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		ctx.JSON(400, dto.ErrorResponse{Error: err.Error()})
		return 0, false
	}

	return IDint, true
}

func validatePagination(page *int, pageSize *int, ctx *gin.Context) bool {

	if *page == 0 {
		*page = defaultPage
	}

	if *pageSize == 0 {
		*pageSize = defaultPageSize
	}

	if *page < 0 {
		ctx.JSON(400, dto.ErrorResponse{Error: "page must be greater than 0"})
		return false
	}

	if *pageSize < 0 {
		ctx.JSON(400, dto.ErrorResponse{Error: "page size must be greater than 0"})
		return false
	}

	if *pageSize > maxPageSize {
		ctx.JSON(400, dto.ErrorResponse{Error: fmt.Sprintf("page size must be less than %d", maxPageSize)})
		return false
	}

	return true
}
