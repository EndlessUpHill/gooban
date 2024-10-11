package gooban

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Job struct {
    ID          uuid.UUID       `db:"id"`
    Queue       string          `db:"queue"`
    Status      string          `db:"status"`
    Args        json.RawMessage `db:"args"`
    Priority    int             `db:"priority"`
    Attempts    int             `db:"attempts"`
    MaxAttempts int             `db:"max_attempts"`
    ScheduledAt time.Time       `db:"scheduled_at"`
    InsertedAt  time.Time       `db:"inserted_at"`
    UpdatedAt   time.Time       `db:"updated_at"`
}

type JobArgs struct {
    UserID    int    `json:"user_id"`
    TaskName  string `json:"task_name"`
    ExtraData string `json:"extra_data"`
}

type JobResult struct {
	JobID    uuid.UUID `json:"job_id"`
	Status   string    `json:"status"`
}