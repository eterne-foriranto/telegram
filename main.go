package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getConfigValue(sectionName string, key string) string {
	cnf, err := config.NewConfig("ini", "config.ini")
	handleError(err)
	section, err := cnf.GetSection(sectionName)
	handleError(err)
	return section[key]
}

func Bot() *tgbotapi.BotAPI {
	token := getConfigValue("telegram", "token")
	bot, err := tgbotapi.NewBotAPI(token)
	handleError(err)
	bot.Debug = false
	return bot
}

func main() {
	bot := Bot()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		fmt.Println(update.Message.Text)
	}
}
