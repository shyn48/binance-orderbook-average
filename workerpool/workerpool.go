package workerpool

import (
	"sync"

	"binance-orderbook-average/common"
)

type Task interface{}

type WorkerPool struct {
	tasks       chan Task
	workerCount int
	wg          sync.WaitGroup
	stopChan    chan struct{}
}

func New(workerCount int, taskBuffer int) *WorkerPool {
	return &WorkerPool{
		tasks:       make(chan Task, taskBuffer),
		workerCount: workerCount,
		stopChan:    make(chan struct{}),
	}
}

func (wp *WorkerPool) Start(workerFunc func(Task)) {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go func(id int) {
			defer wp.wg.Done()
			for {
				select {
				case task, ok := <-wp.tasks:
					if !ok {
						return
					}
					workerFunc(task)
				case <-wp.stopChan:
					return
				}
			}
		}(i)
	}
}

func (wp *WorkerPool) Submit(task Task) {
	select {
	case wp.tasks <- task:
	default:
		common.Logger().Warn("WorkerPool: Task buffer full. Dropping task.")
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.stopChan)
	close(wp.tasks)
	wp.wg.Wait()
}
