package routingpool

import (
	"sync"
	"sync/atomic"
	//"utility"
    "fmt"
)

var (
	//logger = utility.GetLogger()
)

/*var pool *ThreadPool

func init() {
	pool = GetPool(10, 10)
}
*/
type ThreadPool struct {
	NumberOfThread int
	QueueCapacity  int
	ActiveThread   int32
	FreeThread     int32
	TaskWorkQueue  chan Task
	TaskCacheQueue chan Task
	wg             sync.WaitGroup
	StopQueuePool  chan bool
    StopWorkPool   chan bool
	stopCache  bool
}
/*
func SetPoolSize(number int, capacity int) {
	pool.NumberOfThread = number
	pool.QueueCapacity = capacity
}
*/

func GetPool(number int, capacity int) *ThreadPool {
	pool := &ThreadPool {
		NumberOfThread: number,
		QueueCapacity:  capacity,
		ActiveThread:   0,
		FreeThread:     int32(number),
		TaskWorkQueue:  make(chan Task),
		TaskCacheQueue: make(chan Task, capacity),
		StopQueuePool:  make(chan bool),
        StopWorkPool:   make(chan bool),
		wg:             sync.WaitGroup{},
		stopCache:  false,
	}

	return pool
}

// Startup threads
//func Start() {pool.Start()}
func (pool *ThreadPool) Start() {
	if pool.TaskCacheQueue == nil {
		//logger.Debug("Start Thread pool....")
		pool.TaskCacheQueue = make(chan Task, pool.QueueCapacity)
	}

	pool.wg.Add(1)
	go pool.startQueueThread()

	for routine := 0; routine < pool.NumberOfThread; routine++ {
		pool.wg.Add(1)
		go pool.startWorkThread(routine)
	}

	// Waiting for all routines startup.
	pool.wg.Wait()
}

func (pool *ThreadPool) startWorkThread(id int) {
	
    fmt.Println("worker thread %d get ready", id)
    pool.wg.Done()
	//logger.Debug("worker thread %d get ready", id)
	for {
		select {
		case task := <-pool.TaskWorkQueue:
			pool.wg.Add(1)
			//task.SendResponse()

			//logger.Debugf("[Thread id-%d, name-%s] Thread Started! Routing pool status: Active threads-%d, Free threads-%d", id, task.GetTaskName(), pool.ActiveThread, pool.FreeThread)

			task.Run(id)

			atomic.AddInt32(&pool.ActiveThread, -1)
			atomic.AddInt32(&pool.FreeThread, 1)

			pool.wg.Done()

			//logger.Debugf("[Thread id-%d, name-%s] Thread Finished! Routing pool status: Active threads-%d, Free threads-%d", id, task.GetTaskName(), pool.ActiveThread, pool.FreeThread)

		case <-pool.StopWorkPool:
            fmt.Println("Shutdown-ING...WorkThread.")
			//pool.wg.Done()
		}
	}
}

// Start queue thread to collect the requests from client.
func (pool *ThreadPool) startQueueThread() {
	pool.wg.Done()
	//logger.Debug("Queue thread get ready.")
    fmt.Println("Queue thread get ready.")
	for {
		select {
		case task := <-pool.TaskCacheQueue:
			//logger.Debugf("Cache queue tik out task %s", task.GetTaskName())
            fmt.Println("Cache queue tik out task.")
			pool.TaskWorkQueue <- task

			atomic.AddInt32(&pool.ActiveThread, 1)
			atomic.AddInt32(&pool.FreeThread, -1)
		case <-pool.StopQueuePool:
            fmt.Println("Shutdown-ING...QueueThread.")
            pool.stopCacheTask()
            //pool.StopWorkPool <- true
			//pool.wg.Done()
		}
	}
}

//func PutTask(task Task) {pool.PutTask(task)}
func (pool *ThreadPool) PutTask(task Task) bool {
    if pool.stopCache {
        return false
    }
    
	//logger.Debugf("Received task %s. Currently task queue size is %d, capacity is %d", task.GetTaskName(), len(pool.TaskCacheQueue), pool.QueueCapacity)
	if len(pool.TaskCacheQueue) >= pool.QueueCapacity {
		//logger.Errorf("Task queue is full, task %s is aborted.", task.GetTaskName())
		return false
	}
    
    fmt.Println("Received task...........")
	pool.TaskCacheQueue <- task
	//task.WaitForResponse()

	return true
}

//func Wait() {pool.Wait()}
func (pool *ThreadPool) Wait() {
	pool.wg.Wait()
}

func (pool *ThreadPool) Shutdown() {
    fmt.Println("Shutdown()")
    pool.StopQueuePool <- true
}

func (pool *ThreadPool) stopCacheTask(){
    pool.stopCache = true
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



