package router

import (
	"github.com/gin-gonic/gin"
	"github.com/siddharthray/task-manager-api/internal/handler"
)

func NewRouter(taskHandler *handler.TaskHandler) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())
	taskHandler.Register(r)
	return r
}
