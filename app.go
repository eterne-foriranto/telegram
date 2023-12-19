package main

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/restream/reindexer/v3"
	"strconv"
)

const (
	StateWelcome = "welcome"
	InpDrugName  = "inp drug name"
	InpPeriod    = "inp period"
	UnderRemind  = "under remind"
	InpHour      = "inp hour"
	InpMinute    = "inp minute"
	ReadyToStart = "ready to start"
)

var reacts = map[string]string{
	"add_drug": "Добавить лекарство",
	"add_time": "Добавить ещё одно время",
	"start":    "Запустить напоминание",
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

func (r *Response) processInt(inp string) (int, bool) {
	ok := true
	value, err := strconv.Atoi(inp)
	if err != nil {
		r.Text = err.Error()
		ok = false
	}
	return value, ok
}

func (r *Response) processClockInp(inp string, min, max int) (int, bool) {
	value, ok := r.processInt(inp)
	if value < min || value > max {
		r.Text = fmt.Sprintf("Число должно быть в пределах %v и %v", min, max)
		ok = false
	}
	return value, ok
}

func (r *Response) processHour(inp string) (int, bool) {
	return r.processClockInp(inp, 0, 23)
}

func (r *Response) processMinute(inp string) (int, bool) {
	return r.processClockInp(inp, 0, 59)
}

func response(inp string, chatID int, app *App) Response {
	db := app.DB
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
		user.setState(InpPeriod, db)
		res.Text = "Каждые сколько часов принимать?"
	case InpPeriod:
		hours, ok := res.processInt(inp)
		if ok {
			user.setPeriod(hours, db)
			job, ok := user.findEditedJob(db)
			if ok {
				task := gocron.NewTask(job.remind, app)
				cronJob := gocron.DurationJob(job.Period)
				//cronJob := gocron.DurationJob(int(job.Period) * time.Hour)
			}
		}
	case InpHour:
		hour, ok := res.processHour(inp)
		if ok {
			err := user.setHour(hour, db)
			if err != nil {
				res.Text = "Cannot set hour"
			}
			user.setState(InpMinute, db)
			res.Text = "Введите минуту приёма"
		}
	case InpMinute:
		minute, ok := res.processMinute(inp)
		if ok {
			err := user.setMinute(minute, db)
			if err != nil {
				res.Text = "cannot set minute"
			}
			user.clearEditedTime(db)
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
