package main

import (
	"sync"

	"github.com/edwingeng/deque/v2"
	"github.com/madflojo/tasks"
)

var scheduler *tasks.Scheduler
var taskCount int
var Deck struct {
	dq *deque.Deque[taskResult]
	mu sync.Mutex
}

type newTask struct {
	N   int     `json:"n"`
	D   float64 `json:"d"`
	N1  float64 `json:"n1"`
	L   float64 `json:"l"`
	TTL int64   `json:"ttl"`
}
type taskResult struct {
	newTask           `json:"task"`
	TaskPlace         int    `json:"task_place,omitempty"`
	Status            string `json:"status,omitempty"`
	IterNum           int    `json:"iter_num,omitempty"`
	TaskCreatedTime   int64  `json:"task_created_time,omitempty"`
	TaskStartedTime   int64  `json:"task_started_time,omitempty"`
	TaskCompletedTime int64  `json:"task_completed_time,omitempty"`
}
