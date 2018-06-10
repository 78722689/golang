package routingpool

import (
	"sync"
	"sync/atomic"
	"utility"
)

var (
	logger = utility.GetLogger()
)

var pool *ThreadPool

func init() {
	pool = GetPool(10, 10)
}

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

func SetPoolSize(number int, capacity int) {
	pool.numberOfThread = number
	pool.queueCapacity = capacity
}

func GetPool(number int, capacity int) *ThreadPool {
	pool := ThreadPool{
		numberOfThread: number,
		queueCapacity:  capacity,
		activeThread:   0,
		freeThread:     int32(number),
		taskWorkQueue:  make(chan Task),
		//taskCacheQueue: make(chan Task, capacity),
		shutdown:       make(chan bool),
		wg:             sync.WaitGroup{},
	}

	return &pool
}

// Startup threads
func Start() {pool.Start()}
func (pool *ThreadPool) Start() {
	if pool.taskCacheQueue == nil {
		//logger.Debug("Start Thread pool....")
		pool.taskCacheQueue = make(chan Task, pool.queueCapacity)
	}

	pool.wg.Add(1)
	go pool.startQueueThread()

	for routine := 0; routine < pool.numberOfThread; routine++ {
		pool.wg.Add(1)
		go pool.startWorkThread(routine)
	}



	// Waiting for all routines startup.
	pool.Wait()
}

func (pool *ThreadPool) startWorkThread(id int) {
	pool.wg.Done()
	//logger.Debug("worker thread %d get ready", id)
	for {
		select {
		case task := <-pool.taskWorkQueue:
			pool.wg.Add(1)
			//task.SendResponse()

			logger.Debugf("[Thread id-%d, name-%s] Thread Started! Routing pool status: Active threads-%d, Free threads-%d", id, task.GetTaskName(), pool.activeThread, pool.freeThread)

			task.Run(id)

			atomic.AddInt32(&pool.activeThread, -1)
			atomic.AddInt32(&pool.freeThread, 1)

			pool.wg.Done()

			logger.Debugf("[Thread id-%d, name-%s] Thread Finished! Routing pool status: Active threads-%d, Free threads-%d", id, task.GetTaskName(), pool.activeThread, pool.freeThread)

		case <-pool.shutdown:
			//pool.wg.Done()
			return
		}
	}
}

// Start queue thread to collect the requests from client.
func (pool *ThreadPool) startQueueThread() {
	pool.wg.Done()
	//logger.Debug("Queue thread get ready.")
	for {
		select {
		case task := <-pool.taskCacheQueue:
			//logger.Debugf("Cache queue tik out task %s", task.GetTaskName())
			pool.taskWorkQueue <- task

			atomic.AddInt32(&pool.activeThread, 1)
			atomic.AddInt32(&pool.freeThread, -1)
		case <-pool.shutdown:

			//pool.wg.Done()
		}
	}
}

func PutTask(task Task) {pool.PutTask(task)}
func (pool *ThreadPool) PutTask(task Task) bool {
	//logger.Debugf("Received task %s. Currently task queue size is %d, capacity is %d", task.GetTaskName(), len(pool.taskCacheQueue), pool.queueCapacity)
	if len(pool.taskCacheQueue) >= pool.queueCapacity {
		logger.Errorf("Task queue is full, task %s is aborted.", task.GetTaskName())
		return false
	}

	pool.taskCacheQueue <- task
	//task.WaitForResponse()

	return true
}

func Wait() {pool.Wait()}
func (pool *ThreadPool) Wait() {
	pool.wg.Wait()
}

// Run the task with what you want, and return the result.
type Task interface {
	Run(id int)
	WaitForResponse()
	SendResponse()

	GetTaskName() string
}

type Base struct {
	Name string
	Call func(int)
	//Data interface{}

	Response chan bool
}

func NewCaller(name string, call func(int)) *Base {
	return &Base{Name: name, Call: call, Response: make(chan bool)}
}

func (c *Base) Run(id int) {
	c.Call(id)
}

func (c *Base) WaitForResponse() {
	<-c.Response
}

func (c *Base) SendResponse() {
	c.Response <- true
}

func (c *Base) GetTaskName() string {
	return c.Name
}
