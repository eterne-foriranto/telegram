package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Response struct {
	Text    string
	Buttons []string
	ChatID  int64
}

func (a App) handle(msg *tgbotapi.Message) Response {
	if int(msg.Chat.ID) == a.Owner.ChatID {
		return Response{
			Text:    "Привет!",
			Buttons: []string{"invite"},
			ChatID:  msg.Chat.ID,
		}
	}
	return Response{
		Text:    "",
		Buttons: nil,
		ChatID:  msg.Chat.ID,
	}
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
