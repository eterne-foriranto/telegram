package main

import (
	"github.com/restream/reindexer/v3"
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

func (bg BitterGrass) start(userID string, db *reindexer.Reindexer) Response {
	bg.ServiceID = "bitter_grass"
	bg.UserID = userID
	bg.State = "awaiting for react"
	bg.assignUser(db)
	return Response{
		Text:    "Пора принять таблетки",
		Buttons: reactsFront(),
		ChatID:  0,
	}
}

func (bg BitterGrass) next(inp string, app *App) Response {
	switch bg.State {
	case "awaiting for react":
		switch inp {
		case reacts["done"]:
			bg.unlink(app.DB)
			return Response{"Хорошо", nil, 0}
		case reacts["delay"]:
			bg.State = "input minutes to delay"
			return Response{"Через сколько минут?", nil, 0}
		}
	case "input minutes to delay":
		str := strings.ReplaceAll(inp, ",", ".")
		minutes, err := strconv.ParseFloat(str, 32)

	}
	return Response{}
}
