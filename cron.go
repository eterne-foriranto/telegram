package main

import (
	"github.com/go-co-op/gocron/v2"
	"time"
)

func pushOneTimeJob(app *App, job *Job, firstTime bool) {
	task := gocron.NewTask(job.startFrequentReminder, app)
	if !firstTime {
		job.NextTime = job.NextTime.Add(job.Period)
		defaultUpsert(app.DB, "job", job)
	}
	jobDefinition := gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(job.NextTime))
	cronJob, err := app.Scheduler.NewJob(jobDefinition, task)
	job.setCronID(cronJob.ID(), app.DB)
	handleError(err)
}

func restoreJobs(app *App) {
	for _, job := range allJobs(app.DB) {
		if job.NextTime.After(time.Now()) {
			pushOneTimeJob(app, job, true)
		}
	}
}

func (u *User) cancelJob(name string, app *App) {
	db := app.DB
	job := u.activeJobByName(name, db)
	job.setInactive(db)
	err := app.Scheduler.RemoveJob(decodeCronID(job.CronID))
	handleError(err)
}
