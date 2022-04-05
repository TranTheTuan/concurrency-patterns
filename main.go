package main

import (
	"fmt"
	"time"

	"github.com/TranTheTuan/concurrency-patterns/runner"
)

func main() {
	r := runner.NewRunner(time.Second * 2)
	r.AddTasks(CreateTask(), CreateTask(), CreateTask())

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("recover from %v\n", err)
		}
	}()

	err := r.Start()
	if err != nil {
		panic(err)
	}
}

func CreateTask() func(int) {
	return func(i int) {
		fmt.Printf("create task %d\n", i)
		time.Sleep(1*time.Second)
	}
}
