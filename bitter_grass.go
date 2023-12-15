package main

import (
	"fmt"
	"strconv"
	"strings"
)

var reacts = map[string]string{
	"done":  "Приняла",
	"delay": "Напомнить позже",
}

func reactsFront() []string {
	reactsFront := make([]string, 0)
	for _, v := range reacts {
		reactsFront = append(reactsFront, v)
	}
	return reactsFront
}

type BitterGrass struct {
	ServiceInstance
}

func (bg *BitterGrass) start(userID string, app *App) Response {
	bg.ServiceID = "bitter_grass"
	bg.UserID = userID
	bg.State = "awaiting for react"
	bg.assignUser(app.DB)
	return Response{
		Text:    "Пора принять таблетки",
		Buttons: reactsFront(),
		ChatID:  0,
	}
}

func (bg *BitterGrass) next(inp string, app *App) Response {
	buttons := []string{reacts["done"]}
	if inp == reacts["done"] {
		bg.unlink(app.DB)
		return Response{"Хорошо", nil, 0}
	}

	switch bg.State {
	case "awaiting for react":
		switch inp {
		case reacts["delay"]:
			bg.State = "input minutes to delay"
			return Response{"Через сколько минут?", buttons, 0}
		}
	case "input minutes to delay":
		str := strings.ReplaceAll(inp, ",", ".")
		minutes, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return Response{"Требуется ввести число", buttons, 0}
		}
		bg.unlink(app.DB)
		return Response{fmt.Sprintf("Отложено на %v минуты", minutes), nil, 0}
	}
	return Response{}
}
