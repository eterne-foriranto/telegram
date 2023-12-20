package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Response struct {
	Text    string
	Buttons []string
	ChatID  int64
}

func (a App) handle(msg *tgbotapi.Message) Response {
	chatID := int(msg.Chat.ID)
	return response(msg.Text, chatID, &a)
}

func makeMessage(resp Response) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(resp.ChatID, resp.Text)
	buttons := make([]tgbotapi.KeyboardButton, 0)
	for _, v := range resp.Buttons {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(v))
	}

	if resp.Buttons != nil {
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)
	}
	return msg
}
