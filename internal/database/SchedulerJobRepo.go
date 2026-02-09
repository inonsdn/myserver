package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const (
	SCHEDULER_JOB_TABLE_NAME = "schedulerJob"
	ONE_TIME_JOB_TYPE        = 0
	EVERY_MINUTE_JOB_TYPE    = 1
	EVERY_HOUR_JOB_TYPE      = 2
	EVERY_DAY_JOB_TYPE       = 3
	EVERY_MONTH_JOB_TYPE     = 4
	EVERY_YEAR_JOB_TYPE      = 5
)

type SchedulerJobRepo struct {
	*BaseRepo
}

type ScheduleInfo struct {
	periodTime int
}

type SchedulerJob struct {
	id                uuid.UUID
	name              string
	jobType           int
	scheduleTimestamp time.Time
	scheduleInfo      ScheduleInfo
}

func (s *SchedulerJob) GetDuration() time.Duration {
	return time.Now().Sub(s.scheduleTimestamp)
}

func (s *SchedulerJob) GetNextScheduleTime() time.Time {
	if s.jobType == EVERY_MINUTE_JOB_TYPE {
		return time.Now().Add(time.Minute * time.Duration(s.scheduleInfo.periodTime))
	}
	return time.Now().Add(time.Minute * time.Duration(s.scheduleInfo.periodTime))
}

func (s *SchedulerJobRepo) MarkJobDoneAndSchedule(sj *SchedulerJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := fmt.Sprintf("UPDATE %s SET (scheduleTimestamp) = ($1) WHERE id = $2", SCHEDULER_JOB_TABLE_NAME)
	_, err := s.pool.Exec(ctx, statement, sj.GetNextScheduleTime(), sj.id)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (s *SchedulerJobRepo) CreateMinuteJob(name string, minute int, scheduleTime time.Time) uuid.UUID {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := fmt.Sprintf("INSERT INTO %s (name, jobType, scheduleTimestamp) VALUES ($1, $2, $3) RETURNING id", SCHEDULER_JOB_TABLE_NAME)
	var jobId uuid.UUID
	err := s.pool.QueryRow(ctx, statement, name, EVERY_MINUTE_JOB_TYPE, scheduleTime).Scan(&jobId)
	if err != nil {
		slog.Error(err.Error())
		return uuid.Nil
	}
	return jobId
}

func (s *SchedulerJobRepo) GetAllSchedulerJobWithOrder() []SchedulerJob {
	allSchedulers := []SchedulerJob{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := fmt.Sprintf("SELECT id, name, jobType, scheduleTimestamp, scheduleInfo FROM %s ORDER BY scheduleTimestamp ASC", SCHEDULER_JOB_TABLE_NAME)
	rows, err := s.pool.Query(ctx, statement)
	if err != nil {
		slog.Error(err.Error())
		return allSchedulers
	}
	defer rows.Close()
	for rows.Next() {
		var schedulerJob SchedulerJob
		if err := rows.Scan(&schedulerJob.id, &schedulerJob.name, &schedulerJob.jobType, &schedulerJob.scheduleTimestamp, &schedulerJob.scheduleInfo); err != nil {
			slog.Error(err.Error())
			return allSchedulers
		}
		allSchedulers = append(allSchedulers, schedulerJob)
	}
	return allSchedulers
}

func (s *SchedulerJobRepo) GetNextSchedulerJob() SchedulerJob {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := fmt.Sprintf("SELECT id, name, jobType, scheduleTimestamp, scheduleInfo FROM %s ORDER BY scheduleTimestamp ASC LIMIT 1", SCHEDULER_JOB_TABLE_NAME)
	var schedulerJob SchedulerJob
	err := s.pool.QueryRow(ctx, statement).Scan(&schedulerJob.id, &schedulerJob.name, &schedulerJob.jobType, &schedulerJob.scheduleTimestamp, &schedulerJob.scheduleInfo)
	if err != nil {
		slog.Error(err.Error())
		return schedulerJob
	}
	return schedulerJob
}
