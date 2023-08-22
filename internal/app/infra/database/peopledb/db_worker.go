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
	// tickInsertRateOffset := w.getRandomTickTime(1000, 3000)
	// tickInsertRate := time.Duration(10000*time.Millisecond) + tickInsertRateOffset
	tickInsertRate := time.Duration(10 * time.Second)
	tickInsert := time.Tick(tickInsertRate)

	batchMaxSize := 10000
	batch := make([]Job, 0, batchMaxSize)

	for {
		select {
		case data := <-dataCh:
			batch = append(batch, data)

		case <-tickInsert:
			batchLen := len(batch)
			if batchLen > 0 {
				log.Infof("Tick insert (len=%d)", batchLen)
				insertCh <- batch

				batch = make([]Job, 0, batchMaxSize)
			}
		}
	}
}

func (w Worker) processInsert(insertCh chan []Job) {
	conn, err := w.db.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Failed to acquire connection: %v", err)
	}
	defer conn.Release()

	columns := []string{"id", "nickname", "name", "birthdate", "stack", "search"}
	identifier := pgx.Identifier{"people"}

	for {
		select {
		case payload := <-insertCh:
			_, err := conn.CopyFrom(
				context.Background(),
				identifier,
				columns,
				pgx.CopyFromSlice(len(payload), w.makeCopyFromSlice(payload)),
			)

			if err != nil {
				log.Errorf("Error on insert batch: %v", err)
			}
		}
	}
}

func (Worker) getRandomTickTime(min, max int) time.Duration {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	amount := randomizer.Intn(max-min) + min

	return time.Duration(amount) * time.Millisecond
}

func (Worker) makeCopyFromSlice(batch []Job) func(i int) ([]interface{}, error) {
	return func(i int) ([]interface{}, error) {
		return []interface{}{
			batch[i].Payload.ID,
			batch[i].Payload.Nickname,
			batch[i].Payload.Name,
			batch[i].Payload.Birthdate,
			batch[i].Payload.StackStr(),
			batch[i].Payload.SearchStr(),
		}, nil
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
