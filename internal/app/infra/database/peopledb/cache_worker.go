package peopledb

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

var (
	MaxWorker = 2
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
	cache      *PeopleDbCache
}

func NewWorker(workerPool chan chan Job, cache *PeopleDbCache) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		cache:      cache,
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel
			log.Info("Worker registered in worker pool")

			select {
			case job := <-w.JobChannel:
				log.Infof("Worker %d: Received job %s", w.WorkerPool, job.Payload.ID)
				if _, err := w.cache.Set(job.Payload.ID, job.Payload); err != nil {
					log.Errorf("Error inserting person in cache: %v", err)
				}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Dispatcher struct {
	maxWorkers int
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	cache      *PeopleDbCache
	jobQueue   chan Job
}

func NewDispatcher(cache *PeopleDbCache, jobQueue JobQueue) *Dispatcher {
	maxWorkers := MaxWorker

	pool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
		jobQueue:   jobQueue,
		cache:      cache,
	}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.cache)
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
