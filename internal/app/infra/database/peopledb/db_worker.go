package peopledb

import (
	"context"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

var (
	MaxWorker = 1
	MaxQueue  = 1
)

type JobQueue chan Job

// Job represents the job to be run
type Job struct {
	Payload *people.Person
}

// A buffered channel that we can send work requests on.

func NewJobQueue() JobQueue {
	return make(JobQueue, MaxQueue)
}

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
	db         *pgxpool.Pool
}

func (w Worker) Start() {
	dataCh := make(chan Job)
	insertCh := make(chan []Job)

	go w.bootstrap(dataCh)

	go w.processData(dataCh, insertCh)

	go w.processInsert(insertCh)
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func (w Worker) bootstrap(dataCh chan Job) {
	for {
		w.WorkerPool <- w.JobChannel

		select {
		case job := <-w.JobChannel:
			dataCh <- job

		case <-w.quit:
			return
		}
	}
}

func (w Worker) processData(dataCh chan Job, insertCh chan []Job) {
	batchMaxSize := 5000
	batch := make([]Job, 0, batchMaxSize)

	min := 5
	max := 15

	randomTickTime := time.Duration(rand.Intn(max-min) + min)

	tick := time.Tick(randomTickTime * time.Second)

	for {
		select {
		case data := <-dataCh:
			batch = append(batch, data)

		case <-tick:
			log.Infof("Insert tick, current batch length is %d", len(batch))
			if len(batch) > 0 {
				insertCh <- batch
				batch = make([]Job, 0, batchMaxSize)
				tick = time.Tick(randomTickTime * time.Second)
			}
		}
	}
}

func (w Worker) processInsert(insertCh chan []Job) {
	for {
		select {
		case batch := <-insertCh:
			_, err := w.db.CopyFrom(
				context.Background(),
				pgx.Identifier{"people"},
				[]string{"id", "nickname", "name", "birthdate", "stack", "search"},
				pgx.CopyFromSlice(len(batch), func(i int) ([]interface{}, error) {
					return []interface{}{
						batch[i].Payload.ID,
						batch[i].Payload.Nickname,
						batch[i].Payload.Name,
						batch[i].Payload.Birthdate,
						batch[i].Payload.StackStr(),
						batch[i].Payload.SearchStr(),
					}, nil
				}))

			if err != nil {
				log.Errorf("Error on insert batch: %v", err)
			}
		}
	}
}

func NewWorker(workerPool chan chan Job, db *pgxpool.Pool) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		db:         db,
	}
}

type Dispatcher struct {
	maxWorkers int
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	jobQueue   chan Job
	db         *pgxpool.Pool
}

func NewDispatcher(db *pgxpool.Pool, jobQueue JobQueue) *Dispatcher {
	maxWorkers := MaxWorker

	pool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
		jobQueue:   jobQueue,
		db:         db,
	}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.db)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
