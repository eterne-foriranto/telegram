package main

import (
	"fmt"
)

func text(drugName string) string {
	return fmt.Sprintf("Пора принимать %v", drugName)
}

func (j *Job) remind(app *App) {
	db := app.DB
	bot := app.Bot
	setUserState(j.ChatID, UnderRemind, db)
	resp := Response{
		Text:    text(j.Name),
		Buttons: []string{"Принял(а)"},
		ChatID:  int64(j.ChatID),
	}
	msg := makeMessage(resp)
	_, err := bot.Send(msg)
	handleError(err)
}
