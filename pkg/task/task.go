package task

import (
	"encoding/json"
	"time"
)

// Priority levels for tasks
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// TaskStatus represents the current state of a task
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
	StatusRetrying  TaskStatus = "retrying"
	StatusCancelled TaskStatus = "cancelled"
)

// Task represents a unit of work to be executed
type Task struct {
	ID          string                 `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Queue       string                 `json:"queue" db:"queue"`
	Priority    Priority               `json:"priority" db:"priority"`
	Status      TaskStatus             `json:"status" db:"status"`
	Payload     json.RawMessage        `json:"payload" db:"payload"`
	Args        []interface{}          `json:"args" db:"args"`
	Kwargs      map[string]interface{} `json:"kwargs" db:"kwargs"`
	Result      json.RawMessage        `json:"result,omitempty" db:"result"`
	Error       string                 `json:"error,omitempty" db:"error"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty" db:"scheduled_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty" db:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty" db:"completed_at"`
	MaxRetries  int                    `json:"max_retries" db:"max_retries"`
	RetryCount  int                    `json:"retry_count" db:"retry_count"`
	RetryDelay  time.Duration          `json:"retry_delay" db:"retry_delay"`
	Timeout     time.Duration          `json:"timeout" db:"timeout"`
	WorkerID    string                 `json:"worker_id,omitempty" db:"worker_id"`
}

// TaskHandler defines the interface for task execution
type TaskHandler interface {
	Handle(task *Task) (interface{}, error)
}

// TaskHandlerFunc is an adapter to allow ordinary functions to be used as TaskHandlers
type TaskHandlerFunc func(task *Task) (interface{}, error)

func (f TaskHandlerFunc) Handle(task *Task) (interface{}, error) {
	return f(task)
}

// TaskRegistry manages registered task handlers
type TaskRegistry struct {
	handlers map[string]TaskHandler
}

func NewTaskRegistry() *TaskRegistry {
	return &TaskRegistry{
		handlers: make(map[string]TaskHandler),
	}
}

func (r *TaskRegistry) Register(name string, handler TaskHandler) {
	r.handlers[name] = handler
}

func (r *TaskRegistry) RegisterFunc(name string, handler func(*Task) (interface{}, error)) {
	r.handlers[name] = TaskHandlerFunc(handler)
}

func (r *TaskRegistry) Get(name string) (TaskHandler, bool) {
	handler, exists := r.handlers[name]
	return handler, exists
}

func (r *TaskRegistry) List() []string {
	names := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		names = append(names, name)
	}
	return names
}