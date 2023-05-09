package agentworkerspool

import (
	"context"
	"fmt"
	"sync"
)

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, errCh chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			// fan-in job execution multiplexing errCh into the errCh channel
			errCh <- job.execute(ctx)
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			errCh <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}

type WorkerPool struct {
	workersCount int
	jobs         chan Job
	errCh        chan Result
	Done         chan struct{}
}

func New(wcount int) WorkerPool {
	return WorkerPool{
		workersCount: wcount,
		jobs:         make(chan Job, wcount),
		errCh:        make(chan Result, wcount),
		Done:         make(chan struct{}),
	}
}

func (wp WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, wp.jobs, wp.errCh)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.errCh)
}

func (wp WorkerPool) ErrCh() <-chan Result {
	return wp.errCh
}

func (wp WorkerPool) GenerateFrom(jobsBulk []Job) {
	for i := range jobsBulk {
		wp.jobs <- jobsBulk[i]
	}
	close(wp.jobs)
}
