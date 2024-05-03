package client

import (
	"fmt"

	"github.com/dcwk/metrics/internal/logger"
)

type WorkerPool struct {
	jobs chan func()
}

func NewWorkerPool(limit int64) *WorkerPool {
	workerPool := &WorkerPool{
		jobs: make(chan func()),
	}

	for i := 0; i < int(limit); i++ {
		logger.Log.Info(fmt.Sprintf("Started worker %d", i))
		go workerPool.Work()
	}

	return workerPool
}

func (wp *WorkerPool) Produce(job func()) {
	wp.jobs <- job
}

func (wp *WorkerPool) Work() {
	for job := range wp.jobs {
		job()
	}
}
