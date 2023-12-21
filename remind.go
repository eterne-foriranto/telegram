package main

import (
	"fmt"
)

func text(drugName string) string {
	return fmt.Sprintf("Пора принимать %v", drugName)
}

func (j *Job) remind(app *App) {
	db := app.DB
	setUserState(j.ChatID, UnderRemind, db)
	resp := Response{
		Text:    text(j.Name),
		Buttons: []string{"Принял(а)"},
		ChatID:  int64(j.ChatID),
	}
	app.send(resp)
}
