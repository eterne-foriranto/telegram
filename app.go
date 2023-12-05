package main

import (
	"github.com/restream/reindexer/v3"
	"strconv"
)

type App struct {
	Owner User
	DB    *reindexer.Reindexer
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

	db := reindexer.NewReindex("cproto://172.19.0.6:6534/fk", reindexer.WithCreateDBIfMissing())
	err = db.OpenNamespace("user", reindexer.DefaultNamespaceOptions(), User{})
	handleError(err)
	err = db.Upsert("user", owner)
	handleError(err)
	err = db.OpenNamespace("service_instance", reindexer.DefaultNamespaceOptions(), ServiceInstance{})
	handleError(err)
	return App{owner, db}
}
