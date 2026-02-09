package scheduler

import (
	"myserver/internal/database"
	"time"
)

type SchedulerHandler struct {
	dbHandler *database.DatabaseHandler
	jobsQueue []database.SchedulerJob
}

func (s *SchedulerHandler) getNextTimeSchedule() {
	schedulerJobDb := s.dbHandler.GetSchedulerJobConnection()
	allJobs := schedulerJobDb.GetAllSchedulerJobWithOrder()
	s.jobsQueue = allJobs
}

func executeJob(sh *SchedulerHandler) {
	// var nextJob database.SchedulerJob
	schedulerJobDb := sh.dbHandler.GetSchedulerJobConnection()
	for {
		job := schedulerJobDb.GetNextSchedulerJob()
		// nextJob, sh.jobsQueue = sh.jobsQueue[0], sh.jobsQueue[:1]
		timer := time.NewTimer(job.GetDuration())
		<-timer.C
		schedulerJobDb.MarkJobDoneAndSchedule(&job)
	}
}

func RunScheduler(sh *SchedulerHandler) {
	// scheduler timer to execute
	go executeJob(sh)
}
