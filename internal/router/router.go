package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
	"test-task/docs"
	"test-task/internal/config"
	"test-task/internal/dto"
	errs "test-task/internal/errors"
	"test-task/internal/handlers"
)

func Setup(engine *gin.Engine, walletHandler *handlers.SongsHandler) {

	engine.Use(gin.Recovery())
	engine.Use(errorHandler)

	engine.GET("/api/v1/songs", walletHandler.GetSongs)
	engine.GET("/api/v1/songs/:id/lyrics", walletHandler.GetSongVerses)
	engine.POST("/api/v1/songs", walletHandler.AddSong)
	engine.PUT("/api/v1/songs/:id", walletHandler.EditSong)
	engine.DELETE("/api/v1/songs/:id", walletHandler.DeleteSong)

	if config.Get().Mode == config.DebugMode {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		docs.SwaggerInfo.Host = "localhost:" + strconv.Itoa(config.Get().Port)
	}
}

func errorHandler(ctx *gin.Context) {

	ctx.Next()

	if len(ctx.Errors) > 0 {
		err := ctx.Errors.Last()

		if errors.Is(err, errs.NotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		} else if errors.Is(err, errs.BadRequest) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Internal server error"})
			logError(ctx, err)
		}
	}
}

func logError(ctx *gin.Context, err error) {
	log.Errorf("%s %s error: %s", ctx.Request.Method, ctx.Request.URL.Path, err)
}
