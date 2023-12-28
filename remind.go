package main

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"time"
)

func text(drugName string) string {
	return fmt.Sprintf("Пора принимать %v", drugName)
}

func (j *Job) remind(app *App) {
	db := app.DB
	resp := Response{
		Text:    text(j.Name),
		Buttons: []string{"Принял(а)"},
		ChatID:  int64(j.ChatID),
	}
	app.send(resp)
	j.increaseCount(db)
}

func (j *Job) startFrequentReminder(app *App) {
	db := app.DB
	setUserState(j.ChatID, UnderRemind, db)
	j.remind(app)
	task := gocron.NewTask(j.remind, app)
	jobDef := gocron.DurationJob(10 * time.Second)
	cronJob, err := app.Scheduler.NewJob(jobDef, task)
	handleError(err)
	j.setCronID(cronJob.ID(), db)
}
