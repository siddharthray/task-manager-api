package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/siddharthray/task-manager-api/internal/handler"
	"time"
)

func NewRouter(taskHandler *handler.TaskHandler) *gin.Engine {
	r := gin.Default()
	// 1) CORS middleware: allow only http://localhost:5173
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 2) Global middleware
	r.Use(gin.Logger(), gin.Recovery())
	// 3) Register your task routes
	taskHandler.Register(r)
	return r
}
