package gooban

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InsertJob inserts a new job into the database.
func InsertJob(ctx context.Context, db *pgxpool.Pool, job *Job) error {
    if job == nil {
        return fmt.Errorf("job cannot be nil")
    }

    // Ensure job has a valid ID
    if job.ID == uuid.Nil {
        job.ID = uuid.New()
    }

    // Set default timestamps
    now := time.Now()
    if job.InsertedAt.IsZero() {
        job.InsertedAt = now
    }
    if job.UpdatedAt.IsZero() {
        job.UpdatedAt = now
    }
    if job.ScheduledAt.IsZero() {
        job.ScheduledAt = now
    }

    // Prepare SQL statement
    sql := `
        INSERT INTO jobs (
            id, queue, status, args, priority, attempts, max_attempts, scheduled_at, inserted_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        )
    `

    // Execute the query
    _, err := db.Exec(ctx, sql,
        job.ID,
        job.Queue,
        job.Status,
        job.Args,
        job.Priority,
        job.Attempts,
        job.MaxAttempts,
        job.ScheduledAt,
        job.InsertedAt,
        job.UpdatedAt,
    )
    if err != nil {
        return fmt.Errorf("failed to insert job into database: %w", err)
    }

    return nil
}

// UpdateJob updates an existing job in the database.
func UpdateJob(ctx context.Context, db *pgxpool.Pool, job *Job) error {
    if job == nil {
        return fmt.Errorf("job cannot be nil")
    }

    // Update the timestamp
    job.UpdatedAt = time.Now()

    // Prepare SQL statement
    sql := `
        UPDATE jobs SET
            status = $1,
            args = $2,
            priority = $3,
            attempts = $4,
            scheduled_at = $5,
            updated_at = $6
        WHERE id = $7
    `

    // Execute the query
    result, err := db.Exec(ctx, sql,
        job.Status,
        job.Args,
        job.Priority,
        job.Attempts,
        job.ScheduledAt,
        job.UpdatedAt,
        job.ID,
    )
    if err != nil {
        return fmt.Errorf("failed to update job in database: %w", err)
    }

    // Check if any row was updated
    if result.RowsAffected() == 0 {
        return fmt.Errorf("no job found with id %s", job.ID)
    }

    return nil
}

// FetchJob fetches the next available job from the specified queue.
func FetchJob(ctx context.Context, db *pgxpool.Pool, queue string) (*Job, error) {
    // Begin a transaction
    tx, err := db.BeginTx(ctx, pgx.TxOptions{
        IsoLevel:   pgx.Serializable, // Ensure data consistency
        AccessMode: pgx.ReadWrite,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer func() {
        if tx != nil {
            tx.Rollback(ctx) // Rollback if not already committed
        }
    }()

    // Prepare SQL statement
    sql := `
        SELECT
            id, queue, status, args, priority, attempts, max_attempts, scheduled_at, inserted_at, updated_at
        FROM jobs
        WHERE
            queue = $1 AND
            status = 'pending' AND
            scheduled_at <= NOW()
        ORDER BY priority DESC, scheduled_at ASC
        LIMIT 1
        FOR UPDATE SKIP LOCKED
    `

    // Query for the job
    row := tx.QueryRow(ctx, sql, queue)

    var job Job

    // Scan the result
    err = row.Scan(
        &job.ID,
        &job.Queue,
        &job.Status,
        &job.Args,
        &job.Priority,
        &job.Attempts,
        &job.MaxAttempts,
        &job.ScheduledAt,
        &job.InsertedAt,
        &job.UpdatedAt,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            // No job available
            err = tx.Commit(ctx) // Commit the transaction
            if err != nil {
                return nil, fmt.Errorf("failed to commit transaction: %w", err)
            }
            tx = nil // Prevent deferred rollback
            return nil, nil
        }
        return nil, fmt.Errorf("failed to fetch job: %w", err)
    }

    // Update the job status to 'processing' and increment attempts
    job.Status = "processing"
    job.Attempts++
    job.UpdatedAt = time.Now()

    // Update the job in the database
    updateSQL := `
        UPDATE jobs SET
            status = $1,
            attempts = $2,
            updated_at = $3
        WHERE id = $4
    `
    _, err = tx.Exec(ctx, updateSQL,
        job.Status,
        job.Attempts,
        job.UpdatedAt,
        job.ID,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to update job status: %w", err)
    }

    // Commit the transaction
    err = tx.Commit(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }
    tx = nil // Prevent deferred rollback

    return &job, nil
}


