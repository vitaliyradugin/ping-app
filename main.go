package main

import (
	"fmt"
	"os"
	"os/signal"
	"ping-app/workerpool"
	"syscall"
	"time"
)

const (
	REQ_INTERVAL  = 5 * time.Second
	REQ_TIMEOUT   = 1 * time.Second
	WORKERS_COUNT = 5
)

var urlPool = []string{
	"https://pkg.go.dev/std",
	"https://metanit.com/go/tutorial/2.10.php",
	"https://leetcode.com/problemset/",
	"https://edu.postgrespro.ru/postgresql_internals-17.pdf",
	"https://psv4.userapi.com/s/v1/d/PvXfMHhYp2wqPZv_EoX00UAzTXeS__70NaPVcuQlOGGqXuMVGvlCFi4Ba6quaB2v7mGtqD6nJ3gEiJqFIMkjRbv31sQkZVJrirt5sKx5TCB5ve60g4WZvg/Grokaem_algoritmy_2.pdf",
}

/*
TODO Добавить сохранение ошибок в бд
подумать над тем как прикрутить сюда что-то еще вроде grafana prometheus
или сохранение логов
или трасировку
*/
func main() {
	results := make(chan workerpool.Result)
	workerPool := workerpool.New(WORKERS_COUNT, REQ_TIMEOUT, results)

	workerPool.Init()

	go generateJobs(workerPool)
	go processResults(results)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	workerPool.Stop()
}

func processResults(results chan workerpool.Result) {
	go func() {
		for result := range results {
			fmt.Println(result.Info())
		}
	}()
}

func generateJobs(wp *workerpool.Pool) {
	for {
		for _, url := range urlPool {
			wp.Push(workerpool.Job{URL: url})
		}
		time.Sleep(REQ_INTERVAL)
	}
}
