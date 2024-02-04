package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getConfigValue(sectionName string, key string) string {
	cnf, err := config.NewConfig("ini", "/config/config.ini")
	handleError(err)
	section, err := cnf.GetSection(sectionName)
	handleError(err)
	return section[key]
}

func getBot() *tgbotapi.BotAPI {
	token := getConfigValue("telegram", "token")
	bot, err := tgbotapi.NewBotAPI(token)
	handleError(err)
	bot.Debug = false
	return bot
}

func shutDown(s gocron.Scheduler) {
	err := s.Shutdown()
	handleError(err)
}

func (a App) send(r *Response) {
	msg := makeMessage(r)
	_, err := a.Bot.Send(msg)
	handleError(err)
}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	app := getApp()
	defer shutDown(app.Scheduler)
	bot := app.Bot
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		msg := makeMessage(app.handle(update.Message))
		_, err := bot.Send(msg)
		handleError(err)
	}
}
