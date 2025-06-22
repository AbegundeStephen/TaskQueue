package task

import ("encoding/json"
"time")

// Prioority levels for tasks
type Priority int

const (
PriorityLow Priority = iota
PriorityNormal
PriorityHigh
PriorityCritical
)

// Taskstatus represents the current state of a task
type Taskstatus string

const (
	StatusPending   Taskstatus = "pending"
	StatusRunning Taskstatus = "running"
	StatusCompleted Taskstatus = "completed"
	StatusFailed    Taskstatus = "failed"
	StatusRetrying Taskstatus = "retrying"
	StatusCancelled Taskstatus = "cancelled"
)

// Task represents a unit of work to be executed
type Task struct {
	ID		  string    `json:"id" db:"id"`
	Name	string		`json:"name" db:"name"`
	Queue	string		`json:"queue" db:"queue"`
	Priority	Priority	`json:"priority" db:"priority"`
	Status	Taskstatus	`json:"status" db:"status"`
	Payload json.RawMessage `json:"payload" db:"payload"`
	Args []interface{} `json:"args" db:"args"`
	Kwargs map[string]interface{} `json:"kwargs" db:"kwargs"`
	Result json.RawMessage `json:"result,omitempty" db:"result"`
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


//TaskHandler defines the interface for task execution
type TaskHandler interface {
	Handler(task *Task) (interface{}, error)
}

//TaskHandlerFunc is an adapter to allow ordinary functions to be used as TaskHandler
type TaskHandlerFunc func(task *Task) (interface{}, error)