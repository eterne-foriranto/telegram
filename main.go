package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
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

func getBot() *tgbotapi.BotAPI {
	token := getConfigValue("telegram", "token")
	bot, err := tgbotapi.NewBotAPI(token)
	handleError(err)
	bot.Debug = false
	return bot
}

func main() {
	bot := getBot()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	ownerChatID, err := strconv.Atoi(getConfigValue("telegram", "owner_chat_id"))
	handleError(err)
	state := State{int64(ownerChatID)}

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		msg := makeMessage(state.handle(update.Message))
		_, err := bot.Send(msg)
		handleError(err)
	}
}
