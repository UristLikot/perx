package main

import (
	"flag"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/edwingeng/deque/v2"
	"github.com/gin-gonic/gin"
	"github.com/madflojo/tasks"
)

func dequeStatus(c *gin.Context) {
	if Deck.dq.IsEmpty() {
		c.JSON(http.StatusOK, gin.H{"status": "deque is empty"})
		return
	}
	c.JSON(http.StatusOK, Deck.dq.Dump())
}

func addTask(c *gin.Context) {
	var task newTask
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	timeNow := time.Now().Unix()
	result := taskResult{
		newTask:           newTask{N: task.N, D: task.D, N1: task.N1, L: task.L, TTL: task.TTL},
		TaskPlace:         taskCount,
		Status:            "New",
		IterNum:           0,
		TaskCreatedTime:   timeNow,
		TaskStartedTime:   0,
		TaskCompletedTime: 0,
	}
	taskCount++
	Deck.mu.Lock()
	defer Deck.mu.Unlock()
	Deck.dq.PushBack(result)
	_, err = scheduler.Add(&tasks.Task{
		RunOnce:  true,
		Interval: time.Second,
		TaskFunc: checkTasks,
	})

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "added"})
}
func checkTasks() error {
	var ttl int64
	if !Deck.dq.IsEmpty() {
		Deck.mu.Lock()
		Deck.dq.Range(func(i int, task taskResult) bool {
			if task.Status == "New" {
				task.TaskStartedTime = time.Now().Unix()
				task.Status = "WIP"

				Deck.dq.Replace(i, task)

				for i := 0; i < task.N; i++ {
					task.N1 += task.D
					task.IterNum++
					time.Sleep(time.Second * time.Duration(task.L))
				}

				task.TaskCompletedTime = time.Now().Unix()
				task.Status = "Done"

				Deck.dq.Replace(i, task)

				ttl = task.TTL
			}
			return true
		})
		Deck.mu.Unlock()
	}
	_, _ = scheduler.Add(&tasks.Task{
		RunOnce:  true,
		Interval: time.Second * time.Duration(ttl),
		TaskFunc: removeTasks,
	})

	return nil
}
func removeTasks() error {
	if !Deck.dq.IsEmpty() {
		Deck.mu.Lock()
		Deck.dq.Range(func(i int, task taskResult) bool {
			if task.TaskCompletedTime+task.TTL < time.Now().Unix() {
				Deck.dq.Remove(i)
			}
			return true
		})
		Deck.mu.Unlock()
	}
	return nil
}

func main() {
	scheduler = tasks.New()

	n := flag.Int("n", runtime.NumCPU(), "max goroutines")
	flag.Parse()
	runtime.GOMAXPROCS(*n)

	Deck.dq = deque.NewDeque[taskResult]()
	Deck.mu = sync.Mutex{}

	r := gin.Default()
	r.GET("/deque", dequeStatus)
	r.POST("/deque", addTask)
	_ = r.Run()
}
