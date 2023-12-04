package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Response struct {
	Text   string
	Keys   []tgbotapi.KeyboardButton
	ChatID int64
}

func (a App) handle(msg *tgbotapi.Message) Response {
	if int(msg.Chat.ID) == a.Owner.ChatID {
		return Response{
			Text:   "Привет!",
			Keys:   nil,
			ChatID: msg.Chat.ID,
		}
	}
	return Response{
		Text:   "",
		Keys:   nil,
		ChatID: msg.Chat.ID,
	}
}

func makeMessage(resp Response) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(resp.ChatID, resp.Text)
	if resp.Keys != nil {
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(resp.Keys)
	}
	return msg
}
