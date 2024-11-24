package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
	"test-task/internal/config"
	"test-task/internal/handlers"
	"test-task/internal/repositories"
	"test-task/internal/router"
	"test-task/internal/services"
)

// @title Songs API
// @version 1.0
// @description API for test task.
// @host localhost:8080
// @BasePath /api/v1
func main() {

	cfg := config.Get()

	if cfg.Mode == config.DebugMode {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	dbContext, err := repositories.NewDbContext(cfg.DbConnectionString)
	if err != nil {
		log.Fatalf("error create dbContext: %v", err)
		return
	}
	defer dbContext.Close()

	err = dbContext.Migrate()
	if err != nil {
		log.Fatalf("error run database migration: %v", err)
		return
	}
	log.Infof("database migration complete")

	songsRepository := repositories.NewSongsRepository(dbContext.DB)
	songsService := services.NewSongsService(songsRepository)
	songsHandler := handlers.NewSongsHandler(songsService)

	gin.SetMode(cfg.Mode)
	ginEngine := gin.New()

	router.Setup(ginEngine, songsHandler)

	log.Errorf(ginEngine.Run(":" + strconv.Itoa(cfg.Port)).Error())
}
