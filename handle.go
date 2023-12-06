package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"slices"
)

type Response struct {
	Text    string
	Buttons []string
	ChatID  int64
}

func (a App) handle(msg *tgbotapi.Message) Response {
	if int(msg.Chat.ID) == a.Owner.ChatID {
		currentService := currentService(int(msg.Chat.ID), a.DB)
		if currentService == "" {
			services := ownerServices(&a)
			if slices.Contains(services, msg.Text) {
				service, ok := a.ServiceMap[msg.Text]
				if ok {
					resp := service.start(a.Owner.ID, a.DB)
				}
			} else {
				return Response{
					Text:    "Привет!",
					Buttons: services,
					ChatID:  msg.Chat.ID,
				}
			}
		} else {
			service, ok := a.ServiceMap[currentService]
			if ok {
				resp := service.next(a.Owner.ID, &a)
			}
		}
		resp.ChatID = msg.Chat.ID
		return resp
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
