package main

import (
	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/restream/reindexer/v3"
	"strconv"
)

type App struct {
	Owner      User
	DB         *reindexer.Reindexer
	ServiceMap map[string]ServiceIface
	Scheduler  gocron.Scheduler
	Bot        *tgbotapi.BotAPI
}

func getApp() App {
	ownerChatID, err := strconv.Atoi(getConfigValue("telegram", "owner_chat_id"))
	handleError(err)
	owner := User{
		ID:               "me",
		ChatID:           ownerChatID,
		IsActive:         true,
		IsDeleted:        false,
		InvitationKey:    "",
		CurrentServiceID: "",
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
	_, err = sch.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(16, 7, 0),
			),
		),
		gocron.NewTask(bg.start, "me", db),
	)
	handleError(err)
	sch.Start()
	bot := getBot()
	return App{owner, db, serviceMap, sch, bot}
}

func ownerServices(app *App) []string {
	services := make([]string, 0)
	for k := range app.ServiceMap {
		services = append(services, k)
	}
	return services
}
