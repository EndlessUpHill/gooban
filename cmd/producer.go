package gooban

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/EndlessUpHill/goakka/core"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProducerActor struct {
    *core.BasicActor
    db *pgxpool.Pool
}

func NewProducerActor(name string, db *pgxpool.Pool) *ProducerActor {
     actor := &ProducerActor{
		BasicActor: core.NewBasicActor(name),
        db: db,
    }
	actor.BasicActor.ReceiveFunc = func(result *core.ActorResult) *core.ActorResult {
		// Handle job submission
		jobArgs, ok := result.Message.(JobArgs)

		if !ok {
			return &core.ActorResult{Error: fmt.Errorf("invalid message type")}
		}

		// Serialize jobArgs to JSON
		argsJSON, err := json.Marshal(jobArgs)
		if err != nil {
			return &core.ActorResult{Error: fmt.Errorf("failed to marshal job arguments: %w", err)}
		}

		job := &Job{
			Queue:       "default",
			Status:      "pending",
			Args:        argsJSON,
			Priority:    0,
			Attempts:    0,
			MaxAttempts: 5,
			// Other fields as necessary
		}

		// Insert job into the database
		err = InsertJob(actor.BasicActor.GetContext(), actor.db, job)
		if err != nil {
			return &core.ActorResult{Error: err, Action: core.ACTOR_FAIL}
		}

		return &core.ActorResult{}

	}
    return actor

}

// Implementing the Actor interface by delegating to basicActor
func (p *ProducerActor) Start() {
    p.BasicActor.Start()
}

func (p *ProducerActor) Stop() {
    p.BasicActor.Stop()
}

func (p *ProducerActor) SendMessage(msg interface{}) {
    p.BasicActor.SendMessage(msg)
}

func (p *ProducerActor) GetID() uuid.UUID {
    return p.BasicActor.GetID()
}

func (p *ProducerActor) GetName() string {
    return p.BasicActor.GetName()
}

func (p *ProducerActor) SetWaitGroup(wg *sync.WaitGroup) {
    p.BasicActor.SetWaitGroup(wg)
}

func (p *ProducerActor) SetContext(ctx context.Context) {
    p.BasicActor.SetContext(ctx)
}

func (p *ProducerActor) SetFailureChannel(ch chan *core.ActorResult) {
    p.BasicActor.SetFailureChannel(ch)
}
