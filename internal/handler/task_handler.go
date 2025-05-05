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
	if tasks == nil {
		tasks = []model.Task{}
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

	// service returns the freshly‑created Task, including created_at
	createdTask, err := h.Service.CreateTask(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 201 + full JSON
	c.JSON(http.StatusCreated, createdTask)
}

// PUT /tasks/:id
func (h *TaskHandler) update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	// 1) bind only what client sent
	var payload model.Task
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.ID = id

	// 2) fetch the current record
	existing, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// 3) if text wasn’t provided in JSON, keep the old one
	if payload.Text == "" {
		payload.Text = existing.Text
	}

	// 4) run the update + return the new full object
	updated, err := h.Service.UpdateTask(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
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
