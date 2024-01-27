package main

import (
	"github.com/go-co-op/gocron/v2"
	"time"
)

func pushOneTimeJob(app *App, job *Job) {
	task := gocron.NewTask(job.startFrequentReminder, app)
	jobDefinition := gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Now().Add(job.Period)))
	cronJob, err := app.Scheduler.NewJob(jobDefinition, task)
	job.setCronID(cronJob.ID(), app.DB)
	handleError(err)
}

func restoreJobs(app *App) {
	for _, job := range allJobs(app.DB) {
		pushOneTimeJob(app, job)
	}
}

func (u *User) cancelJob(name string, app *App) {
	db := app.DB
	job := u.activeJobByName(name, db)
	job.setInactive(db)
	err := app.Scheduler.RemoveJob(decodeCronID(job.CronID))
	handleError(err)
}
