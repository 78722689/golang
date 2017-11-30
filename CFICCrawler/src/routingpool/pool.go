package routingpool

import (
"fmt"
"sync/atomic"
"sync"
	"time"
)

type ThreadPool struct {
	numberOfThread int
	queueCapacity int
	activeThread int32
	freeThread int32
	taskWorkQueue chan Task
	taskCacheQueue chan Task
	wg sync.WaitGroup
	shutdown chan bool
}

// Startup threads
func (pool *ThreadPool)Start() {
	for routine := 0; routine < pool.numberOfThread; routine++ {
		go pool.startWorkThread(routine)

	}

	go pool.startQueueThread()
}

func (pool *ThreadPool)startWorkThread(id int) {
	for {
		select {
		case task := <-pool.taskWorkQueue:
			pool.wg.Add(1)

			task.run(id)

			atomic.AddInt32(&pool.activeThread, -1)
			atomic.AddInt32(&pool.freeThread, 1)
			fmt.Println(fmt.Sprintf("active thread - %d, free thread - %d, released thread id - %d", pool.activeThread, pool.freeThread, id))

			pool.wg.Done()

		case <-pool.shutdown:
			return
		}
	}
}

// Start queue thread to collect the requests from client.
func (pool *ThreadPool)startQueueThread() {
	for {
		select {
		case task := <-pool.taskCacheQueue:
			//pool.wg.Add(1)
			pool.taskWorkQueue <- task

			atomic.AddInt32(&pool.activeThread, 1)
			atomic.AddInt32(&pool.freeThread, -1)
		}
	}
}

func(pool *ThreadPool)PutTask(task Task) {
	pool.taskCacheQueue <- task
}

func GetPool(number int, capacity int) *ThreadPool {
	pool := ThreadPool {
		numberOfThread : number,
		queueCapacity  : capacity,
		activeThread   : 0,
		freeThread     : int32(number),
		taskWorkQueue  : make(chan Task),
		taskCacheQueue : make(chan Task),
		shutdown       : make(chan bool),
	}

	return &pool
}
func(pool *ThreadPool)Wait() {
	pool.wg.Wait()
}

// Run the task with what you want, and return the result.
type Task interface {
	run(id int)
}

type MyTask struct {
	Name string
	Call func(id int)
}

func (task *MyTask)run(id int)  {
	//fmt.Println(fmt.Sprintf("Thread - %d is running with task - %s", id, task.Name))
	task.Call(id)

	time.Sleep(time.Second*2)
}

/*
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//sync.WaitGroup{runtime.NumCPU()}

	myPool := GetPool(100, 100)
	myPool.Start()

	for i := 0; i<=20; i++ {
		go func(id int) {
			for j := 0; j<=20; j++ {
				task := &MyTask{name : fmt.Sprintf(" Task{id - %d, queue-%d}", id, j)}
				myPool.PutTask(task)
			}
		}(i)
	}

	// Waiting for all threads finish and exit
	myPool.wg.Wait()
	fmt.Println("waiting main end.")
}
*/