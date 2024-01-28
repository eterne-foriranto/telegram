package main

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/restream/reindexer/v4"
	"slices"
	"strconv"
	"strings"
)

const (
	StateWelcome     = "welcome"
	EnterDrugName    = "Введите название лекарства"
	InpDrugName      = "inp drug name"
	InpPeriod        = "inp period"
	UnderRemind      = "under remind"
	InpHour          = "inp hour"
	InpMinute        = "inp minute"
	Cancel           = "cancel_job"
	InpJobToCancel   = "inp job to cancel"
	EnterJobToCancel = "Какое напоминание отменить?"
)

var reacts = map[string]string{
	"add_drug":   "Добавить лекарство",
	"add_time":   "Добавить ещё одно время",
	"start":      "Запустить напоминание",
	"have_taken": "Принял(а)",
	"cancel_job": "Отменить напоминание",
}

type App struct {
	DB        *reindexer.Reindexer
	Scheduler gocron.Scheduler
	Bot       *tgbotapi.BotAPI
	Buttons   []string
}

func getApp() App {
	buttons := []string{reacts["add_drug"]}
	ownerChatID, err := strconv.Atoi(getConfigValue("telegram", "owner_chat_id"))
	handleError(err)
	owner := User{
		ID:        "me",
		ChatID:    ownerChatID,
		IsActive:  true,
		IsDeleted: false,
		State:     StateWelcome,
	}

	db := reindexer.NewReindex("cproto://reindexer:6534/fk",
		reindexer.WithCreateDBIfMissing())
	err = db.OpenNamespace("user", reindexer.DefaultNamespaceOptions(), User{})
	handleError(err)
	err = db.Upsert("user", owner)
	handleError(err)

	sch, err := gocron.NewScheduler()
	handleError(err)
	bot := getBot()
	app := App{db, sch, bot, buttons}
	restoreJobs(&app)
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

func (u *User) activeJobNames(db *reindexer.Reindexer) []string {
	names := make([]string, 0)
	for _, job := range u.activeJobs(db) {
		names = append(names, job.Name)
	}
	return names
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

func validateInt(inp string) (bool, int, string) {
	value, err := strconv.Atoi(inp)
	ok := err == nil
	return ok, value, fmt.Sprintf("%v", err)
}

type Validator struct {
	Inp    string
	Errors []string
	OK     bool
}

func (v *Validator) checkOK(ok bool, err string) bool {
	if !ok {
		v.OK = false
		v.Errors = append(v.Errors, err)
		return false
	}
	return true
}

type Period struct {
	Validator
	Out int
}

func (p *Period) validate() {
	ok, value, err := validateInt(p.Inp)
	if p.checkOK(ok, err) {
		p.Out = value
	}

	p.checkOK(validateMin(p.Out, 1))
}

func validateMin(value, min int) (bool, string) {
	if value >= min {
		return true, ""
	}
	return false, fmt.Sprintf("Число должно быть не меньше %v", min)
}

func (r *Response) processClockInp(inp string, min, max int) (int, bool) {
	value, ok := r.processInt(inp)
	if !ok {
		return value, false
	}

	if value < min || value > max {
		r.Text = fmt.Sprintf("Число должно быть в пределах %v и %v", min, max)
		ok = false
	}
	return value, ok
}

func (r *Response) processPeriod(inp string) (int, bool) {
	return r.processClockInp(inp, 1, 23)
}

func (r *Response) processHour(inp string) (int, bool) {
	return r.processClockInp(inp, 0, 23)
}

func (r *Response) processMinute(inp string) (int, bool) {
	return r.processClockInp(inp, 0, 59)
}

func (r *Response) validateJobToCancel(inp string, user *User, db *reindexer.Reindexer) bool {
	if slices.Contains(user.activeJobNames(db), inp) {
		return true
	}
	r.Text = "Нужно выбрать одно из имеющихся напоминаний"
	return false
}

func response(inp string, chatID int, app *App) Response {
	db := app.DB
	user := user(chatID, db)
	res := Response{}
	switch user.State {
	case StateWelcome:
		res.Buttons = app.Buttons
		if len(user.activeJobs(db)) != 0 {
			res.Buttons = append(res.Buttons, reacts["cancel_job"])
		}
		res.Text = "Готов к установке напоминаний"
		switch inp {
		case reacts["add_drug"]:
			user.setState(InpDrugName, db)
			res.Text = EnterDrugName
			res.Buttons = []string{}
		case reacts["cancel_job"]:
			if len(user.activeJobs(db)) == 0 {
				res.Text = "Напоминаний нет"
				res.Buttons = app.Buttons
			} else {
				user.setState(InpJobToCancel, db)
				res.Text = EnterJobToCancel
				res.Buttons = user.activeJobNames(db)
			}
		}
	case InpJobToCancel:
		if res.validateJobToCancel(inp, user, db) {
			user.cancelJob(inp, app)
			user.setState(StateWelcome, db)
			res.Text = "Успешно отменено"
			res.Buttons = app.Buttons
		}
	case InpDrugName:
		user.attachJob(inp, db)
		user.setState(InpPeriod, db)
		res.Text = "Каждые сколько часов принимать?"
	case InpPeriod:
		period := &Period{}
		period.OK = true
		period.Inp = inp
		period.validate()
		if period.OK {
			user.setPeriod(period.Out, db)
			job, ok := user.findEditedJob(db)
			if ok {
				pushOneTimeJob(app, job)
				user.clearEdited(db)
				res.Text = "Напоминание установлено"
				res.Buttons = []string{reacts["add_drug"]}
				user.setState(StateWelcome, db)
			}
		} else {
			res.Text = strings.Join(period.Errors, ". ")
		}
	case UnderRemind:
		if inp == reacts["have_taken"] {
			user.setState(StateWelcome, db)
			job, ok := user.findEditedJob(db)
			if ok {
				user.stopFrequentReminder(app)
				pushOneTimeJob(app, job)
				job.resetCount(db)
			}
			res.Text = "Хорошо"
			res.Buttons = app.Buttons
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
