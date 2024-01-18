package main

import (
	"github.com/go-co-op/gocron/v2"
	"time"
)

func pushOneTimeJob(app *App, job *Job) {
	task := gocron.NewTask(job.startFrequentReminder, app)
	jobDefinition := gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Now().Add(job.Period)))
	_, err := app.Scheduler.NewJob(jobDefinition, task)
	handleError(err)
}

func restoreJobs(app *App) {
	for _, job := range allJobs(app.DB) {
		pushOneTimeJob(app, job)
	}
}
