package routingpool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"utility"
)

var logger = utility.GetLogger()

type ThreadPool struct {
	numberOfThread int
	queueCapacity  int
	activeThread   int32
	freeThread     int32
	taskWorkQueue  chan Task
	taskCacheQueue chan Task
	wg             sync.WaitGroup
	shutdown       chan bool
}

func GetPool(number int, capacity int) *ThreadPool {
	pool := ThreadPool{
		numberOfThread: number,
		queueCapacity:  capacity,
		activeThread:   0,
		freeThread:     int32(number),
		taskWorkQueue:  make(chan Task),
		taskCacheQueue: make(chan Task),
		shutdown:       make(chan bool),
		wg:             sync.WaitGroup{},
	}

	return &pool
}

// Startup threads
func (pool *ThreadPool) Start() {

	for routine := 0; routine < pool.numberOfThread; routine++ {
		pool.wg.Add(1)
		go pool.startWorkThread(routine)
	}

	pool.wg.Add(1)
	go pool.startQueueThread()

	// Waiting for all routines startup.
	pool.Wait()
}

func (pool *ThreadPool) startWorkThread(id int) {
	pool.wg.Done()

	for {
		select {
		case task := <-pool.taskWorkQueue:
			pool.wg.Add(1)
			task.sendResponse()

			logger.DEBUG(fmt.Sprintf("[Thread id-%d, name-%s] Thread Started! Routing pool status: Active threads-%d, Free threads-%d", id, task.getTaskName(), pool.activeThread, pool.freeThread))

			task.run(id)

			atomic.AddInt32(&pool.activeThread, -1)
			atomic.AddInt32(&pool.freeThread, 1)

			pool.wg.Done()

			logger.DEBUG(fmt.Sprintf("[Thread id-%d, name-%s] Thread Finished! Routing pool status: Active threads-%d, Free threads-%d", id, task.getTaskName(), pool.activeThread, pool.freeThread))

		case <-pool.shutdown:
			//pool.wg.Done()
			return
		}
	}
}

// Start queue thread to collect the requests from client.
func (pool *ThreadPool) startQueueThread() {
	pool.wg.Done()

	for {
		select {
		case task := <-pool.taskCacheQueue:
			pool.taskWorkQueue <- task

			atomic.AddInt32(&pool.activeThread, 1)
			atomic.AddInt32(&pool.freeThread, -1)
		case <-pool.shutdown:

			//pool.wg.Done()
		}
	}
}

func (pool *ThreadPool) PutTask(task Task) {
	pool.taskCacheQueue <- task

	task.waitForResponse()
}

func (pool *ThreadPool) Wait() {
	pool.wg.Wait()
}

// Run the task with what you want, and return the result.
type Task interface {
	run(id int)
	waitForResponse()
	sendResponse()

	getTaskName() string
}

type Caller struct {
	Name string
	Call func(int)
	Data interface{}

	response chan bool
}

func NewCaller(name string, call func(int)) *Caller {
	return &Caller{Name: name, Call: call, response: make(chan bool)}
}

func (c *Caller) run(id int) {
	c.Call(id)
}

func (c *Caller) waitForResponse() {
	<-c.response
}

func (c *Caller) sendResponse() {
	c.response <- true
}

func (c *Caller) getTaskName() string {
	return c.Name
}
