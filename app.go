package main

import (
	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/restream/reindexer/v3"
	"strconv"
)

const (
	StateWelcome = "welcome"
	InpDrugName  = "inp drug name"
	InpHour      = "inp hour"
)

var reacts = map[string]string{
	"add_drug": "Добавить лекарство",
}

type App struct {
	DB        *reindexer.Reindexer
	Scheduler gocron.Scheduler
	Bot       *tgbotapi.BotAPI
	Buttons   map[string][]string
}

func getApp() App {
	buttons := map[string][]string{
		StateWelcome: {"Добавить лекарство"},
	}
	ownerChatID, err := strconv.Atoi(getConfigValue("telegram", "owner_chat_id"))
	handleError(err)
	owner := User{
		ID:            "me",
		ChatID:        ownerChatID,
		IsActive:      true,
		IsDeleted:     false,
		InvitationKey: "",
	}

	db := reindexer.NewReindex("cproto://172.19.0.7:6534/fk",
		reindexer.WithCreateDBIfMissing())
	err = db.OpenNamespace("user", reindexer.DefaultNamespaceOptions(), User{})
	handleError(err)
	err = db.Upsert("user", owner)
	handleError(err)

	serviceMap := map[string]ServiceIface{
		"invite":       &Invite{},
		"bitter_grass": &BitterGrass{},
	}

	sch, err := gocron.NewScheduler()
	handleError(err)
	bg := BitterGrass{}
	bot := getBot()
	app := App{db, sch, bot, buttons}
	_, err = sch.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(17, 59, 00),
			),
		),
		gocron.NewTask(bg.start, "me", &app),
	)
	handleError(err)
	sch.Start()
	return app
}

func user(chatID int, db *reindexer.Reindexer) *User {
	user, ok := findUser(chatID, db)

	if ok {
		return user
	}

	user = &User{
		ChatID: chatID,
		State:  StateWelcome,
	}

	defaultUpsert(db, "user", user)

	return user
}

func validateHour(inp string) (Response, bool) {
	resp := Response{}
	ok := true
	hour, err := strconv.Atoi(inp)
	if err != nil {
		resp.Text = err.Error()
		ok = false
	}

	if hour < 0 || hour > 23 {
		resp.Text = "Число должно быть в пределах 0 и 23"
		ok = false
	}
	return resp, ok
}

func responce(inp string, chatID int, db *reindexer.Reindexer) Response {
	user := user(chatID, db)
	res := Response{}
	switch user.State {
	case StateWelcome:
		res.Buttons = []string{reacts["add_drug"]}
		if inp == reacts["add_drug"] {
			user.setState(InpDrugName, db)
			res.Text = "Введите название лекарства"
			res.Buttons = []string{}
		}
	case InpDrugName:
		user.attachJob(inp, db)
		user.setState(InpHour, db)
		res.Text = "Введите час приёма"
	case InpHour:
		res, ok := validateHour(inp)
		if ok {

		}
	}
	res.ChatID = int64(chatID)
	return res
}

func ownerServices(app *App) []string {
	services := make([]string, 0)
	for k := range app.ServiceMap {
		services = append(services, k)
	}
	return services
}
