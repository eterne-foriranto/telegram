package main

import (
	"github.com/restream/reindexer/v3"
	"strconv"
)

type App struct {
	Owner      User
	DB         *reindexer.Reindexer
	ServiceMap map[string]ServiceIface
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

	db := reindexer.NewReindex("cproto://172.20.0.2:6534/fk",
		reindexer.WithCreateDBIfMissing())
	err = db.OpenNamespace("user", reindexer.DefaultNamespaceOptions(), User{})
	handleError(err)
	err = db.Upsert("user", owner)
	handleError(err)
	//err = db.OpenNamespace("service_instance",
	//	reindexer.DefaultNamespaceOptions(), &ServiceInstance{})
	//handleError(err)

	serviceMap := map[string]ServiceIface{
		"invite": &Invite{},
	}
	return App{owner, db, serviceMap}
}

func ownerServices(app *App) []string {
	services := make([]string, 0)
	for k := range app.ServiceMap {
		services = append(services, k)
	}
	return services
}
