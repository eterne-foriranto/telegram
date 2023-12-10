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

func responseToOwner(msg *tgbotapi.Message, app *App) Response {
	currentService := currentService(int(msg.Chat.ID), app.DB)
	if currentService == "" {
		services := ownerServices(app)
		if slices.Contains(services, msg.Text) {
			service, ok := app.ServiceMap[msg.Text]
			if ok {
				return service.start(app.Owner.ID, app.DB)
			}
		} else {
			return Response{
				Text:    "Привет!",
				Buttons: services,
			}
		}
	} else {
		service, ok := app.ServiceMap[currentService]
		if ok {
			return service.next(msg.Text, app)
		}
	}
	return Response{}
}

func (a App) handle(msg *tgbotapi.Message) Response {
	if int(msg.Chat.ID) == a.Owner.ChatID {
		resp := responseToOwner(msg, &a)
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
