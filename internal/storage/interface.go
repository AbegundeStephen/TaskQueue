package storage

import (
	"context"
	"encoding/json"
	"internal/syscall/windows/registry"
	"net/http"
	"time"

	"github.com/AbegundeStephen/taskqueue/internal/storage"
	"github.com/taskqueue/pkg/task"
)

type Server struct {
	storage storage.Storage
	registry *task.TaskRegistry
	router *gin.Engine
	upgrader websocket.Upgrader
	cLients map[*websocket.Conn]bool
	mu      sync.RWMutext
}



func NewServer(storage storage.Storage, registry *task.TaskRegistry) *Server {
	server := &Server{
		storage: storage,
		registry: registry,
		router: gin.Default(),
		upgrader: websocket.upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			}
		},
		cLients: make(map[*websocket.Conn]bool),
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")

	{
		// Task management endpoints
		api.POST("/tasks", s.createTask)
		api.GET("/tasks/:id", s.getTask)
		api.PUT("/tasks/:id", s.updateTask)
		api.DELETE("/tasks/:id", s.deleteTask)
		api.POST("/tasks/:id/retry", s.retryTask)
		api.POST("/tasks/:id/complete", s.completeTask)
		api.POST("/tasks/:id/cancel", s.cancelTask)

		// Queue management endpoints
		api.GET("/queues", s.ListQueues)
		api.GET("/queues/:name/stats",s.getQueueStats)
		api.POST("/queues/:name/purge", s.purgeQueue)

		// Worker management endpoints
		api.GET("/workers", s.ListWorkers)
		api.GET("/workers/:id", s.getWorker)

		//Metrics and monitoring endpoints
		api.GET("/metrics", s.getMetrics)
		api.GET("/health", s.healthCheck)
	}

	// WebSocket endpoint for real-time updates
	s.router.GET("/ws", s.handleWebSocket)

	//Serve statis dashboard files
	s.router.Static("/static", "./web/static")
	s.router.LoadHTMLGlob("./web/templates/*")
	s.router.GET("/", s.serveDashboard)
}

func (s *Server) createTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create task with defaults
	t := &task.Task{
		ID:		  uuid.New().String(),
		Name:      req.Name,
		Queue:     req.Queue,
		Priority:  req.Priority,
		Status:    task.TaskStatusPending,
		Args:      req.Args,
		Kwargs:    req.Kwargs,
		CreatedAt: time.Now(),
		MaxRetries: req.MaxRetries,
		RetryDelay: req.RetryDelay,
		Timeout:   req.Timeout,
	}

	if req.Payload != nil {
		payLoadJSON, err := json
	}
}