package main

import (
	"fmt"
	"scheduler/internal/scheduler"
	"time"
)

func testCallback() {
	fmt.Println("Job done")
}

func runJobTest() {
	schedulerHandler := scheduler.NewSchedulerHandler()

	job := scheduler.NewDailyJob(
		"TestJob", 0, 49, testCallback,
	)

	schedulerHandler.AddJob(job)

	go schedulerHandler.Run()

	schedulerHandler.Stop(time.Second * 100)
}

func main() {
	runJobTest()

	// opts := config.GetOptions()

	// httpService := servicehandler.NewHttpService(opts)
	// serviceHandler := servicehandler.NewServiceHandler(httpService)

	// serviceHandler.RunService()
}
