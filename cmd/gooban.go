package gooban

import (
	"context"
	"log"

	"github.com/EndlessUpHill/goakka/core"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
    WorkerCount int
    ProducerCount int
    broker core.MessageBroker
}

type gooban struct {
    registry *core.ActorRegistry
    supervisor *core.Supervisor
    broker core.MessageBroker
    functionRegistry *functionRegistry
}

type Gooban interface {
    AddJob(args JobArgs)
    Start()
    Stop()
}

func NewGoobanInsance(ctx context.Context, config Config) (*gooban, error) {
   if config.broker == nil {
       config.broker = core.NewInMemoryBroker()
   }

   dbpool, err := pgxpool.New(ctx, "postgresql://localhost:5432")

   if err != nil {
        log.Fatal(err)
        return nil, err
   }

   registry := core.NewActorRegistry()
   supervisor := core.NewSupervisor(ctx)
   
   producer := NewProducerActor("producer", dbpool)
   worker := NewWorkerActor("worker", dbpool, "default")
   functionRegistry := NewFunctionRegistry()
   worker.functionRegistry = functionRegistry

   registry.RegisterActor(producer)
   registry.RegisterActor(worker)

   config.broker.Subscribe("producer", producer)

   supervisor.SuperviseActor(producer)
   supervisor.SuperviseActor(worker)

  

   return &gooban{
         registry: registry,
         supervisor: supervisor,
         broker: config.broker,
         functionRegistry: functionRegistry,
    }, nil
}

func (g *gooban) AddJob(taskName string, args JobArgs) {
    g.broker.Publish(taskName, args)
}

func (g *gooban) AddFunction(functionName string, f func(args JobArgs) JobResult) {
    g.functionRegistry.RegisterFunction(functionName, f)
}

func (g *gooban) Start() {
    // Start all actors in the registry
}

func (g *gooban) Stop() {
    // Stop all actors in the registry
}
