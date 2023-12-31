package server

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/praveenvoonna/weather-app/backend/config"
	"github.com/praveenvoonna/weather-app/backend/handlers"
	"go.uber.org/zap"
)

func StartServer() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	config.LoadEnv(logger)

	router := gin.Default()

	corsConfig := config.GetCorsConfig()
	router.Use(corsConfig)

	databaseConfig := config.GetDatabaseConfig(logger)
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		databaseConfig.Host, databaseConfig.Port, databaseConfig.User, databaseConfig.Password, databaseConfig.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Error connecting to the database", zap.Error(err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Error("Error pinging the database", zap.Error(err))
	}

	logger.Info("Successfully connected to the PostgreSQL database!")

	router.POST("/login", func(c *gin.Context) {
		handlers.Login(c, db, logger)
	})

	router.POST("/register", func(c *gin.Context) {
		handlers.Register(c, db, logger)
	})

	router.GET("/weather", func(c *gin.Context) {
		handlers.GetCurrentWeather(c, db, logger)
	})

	router.GET("/history", func(c *gin.Context) {
		handlers.GetSearchHistory(c, db, logger)
	})

	router.DELETE("/history", func(c *gin.Context) {
		handlers.ClearSearchHistory(c, db, logger)
	})

	router.Run(":8080")
}
