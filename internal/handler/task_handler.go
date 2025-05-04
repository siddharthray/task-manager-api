package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/siddharthray/task-manager-api/internal/model"
	"github.com/siddharthray/task-manager-api/internal/service"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	Service service.TaskService
}

func NewTaskHandler(s service.TaskService) *TaskHandler {
	return &TaskHandler{Service: s}
}

func (h *TaskHandler) Register(r gin.IRouter) {
	tasks := r.Group("/tasks")
	tasks.GET("", h.list)
	tasks.GET("/:id", h.get)
	tasks.POST("", h.create)
	tasks.PUT("/:id", h.update)
	tasks.DELETE("/:id", h.delete)
}

// GET /tasks/:id
func (h *TaskHandler) get(c *gin.Context) {
	// 1) parse and validate the path param
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	// 2) call your service
	task, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3) return JSON
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) list(c *gin.Context) {
	tasks, err := h.Service.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// POST /tasks
func (h *TaskHandler) create(c *gin.Context) {
	var payload model.Task
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newID, err := h.Service.CreateTask(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": newID})
}

// PUT /tasks/:id
func (h *TaskHandler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var payload model.Task
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.ID = id

	if err := h.Service.UpdateTask(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// DELETE /tasks/:id
func (h *TaskHandler) delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.Service.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
