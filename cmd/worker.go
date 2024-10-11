package gooban

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EndlessUpHill/goakka/core"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkerActor struct {
    *core.BasicActor
    db         *pgxpool.Pool
    queue      string
    workerID   string
    maxRetries int
	functionRegistry *functionRegistry
}

func NewWorkerActor(name string, db *pgxpool.Pool, queue string) *WorkerActor {
    actor := &WorkerActor{
        db:         db,
        queue:      queue,
        workerID:   uuid.New().String(),
        maxRetries: 5,
    }
	actor.BasicActor.ReceiveFunc = func(result *core.ActorResult) *core.ActorResult {
		job, ok := result.Message.(*Job)
		if !ok {
			return &core.ActorResult{Error: fmt.Errorf("invalid message type")}
		}

		// Process the job
		err := actor.processJob(job)
		if err != nil {
			// Handle failure
			job.Attempts++
			if job.Attempts >= job.MaxAttempts {
				job.Status = "failed"
			} else {
				job.Status = "pending"
				job.ScheduledAt = time.Now().Add(time.Minute * time.Duration(job.Attempts))
			}
			UpdateJob(actor.BasicActor.GetContext(), actor.db, job)
			return &core.ActorResult{Error: err, Action: core.ACTOR_RETRY, Message: job}
		}

		// Mark job as completed
		job.Status = "completed"
		 err = UpdateJob(actor.BasicActor.GetContext(),
		  actor.db, job)
		if err != nil {
			return &core.ActorResult{Error: fmt.Errorf("failed to update job status: %w", err)}
		}
		return &core.ActorResult{}
	}

    return actor
}

func (w *WorkerActor) Start() {
    // Start the actor's main loop
    w.BasicActor.Start()

    // Start polling for jobs
    go w.pollForJobs()
}

func (w *WorkerActor) pollForJobs() {
    for {
        select {
		case <-w.BasicActor.GetContext().Done():
            return
        default:
            job, err := FetchJob(w.BasicActor.GetContext(), w.db, w.queue)
            if err != nil {
                time.Sleep(time.Second) // Wait before retrying
                continue
            }
            if job != nil {
                w.SendMessage(job)
            } else {
                time.Sleep(time.Second) // No job found, wait before polling again
            }
        }
    }
}


func (w *WorkerActor) processJob(job *Job) error {
    var args json.RawMessage
    err := json.Unmarshal(job.Args, &args)
    if err != nil {
        return fmt.Errorf("failed to unmarshal job arguments: %w", err)
    }

	result := func(args JobArgs) JobResult {
		// Perform the task
		return JobResult{
			JobID:    job.ID,
			Status:   "completed",
		}
	}

	if (result(JobArgs{}).Status != "completed") {
		return fmt.Errorf("task failed")
	}



	return nil
}
