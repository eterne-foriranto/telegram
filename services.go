package main

import (
	"fmt"
	"github.com/restream/reindexer/v3"
)

type ServiceInstance struct {
	ServiceID string   `reindex:"service_id"`
	UserID    string   `reindex:"user_id"`
	State     string   `reindex:"state"`
	_         struct{} `reindex:"service_id+user_id,,composite,pk"`
}

func (si ServiceInstance) assignUser(db *reindexer.Reindexer) {
	err := db.OpenNamespace("service_instance",
		reindexer.DefaultNamespaceOptions(), si)
	handleError(err)
	err = db.Upsert("service_instance", si)
	handleError(err)

	db.Query("user").
		Where("id", reindexer.EQ, si.UserID).
		Set("current_service_id", si.ServiceID).
		Update()
}

func (si ServiceInstance) unlink(db *reindexer.Reindexer) {
	db.Query("user").
		Where("id", reindexer.EQ, si.UserID).
		Set("current_service_id", "").
		Update()

	_, err := db.Query("service_instance").
		Where("service_id", reindexer.EQ, si.ServiceID).
		Where("user_id", reindexer.EQ, si.UserID).
		Delete()
	handleError(err)
}

type ServiceIface interface {
	start(string, *reindexer.Reindexer) Response
	next(string, *App) Response
}

type Invite struct {
	ServiceInstance
}

func (i *Invite) start(userID string, db *reindexer.Reindexer) Response {
	i.ServiceID = "invite"
	i.State = "input user_id"
	i.UserID = userID
	i.assignUser(db)
	return Response{"Enter ID of the user to be invited", nil, 0}
}

func (i *Invite) next(inp string, app *App) Response {
	key := inviteKey()
	newUserID := inp

	newUser := User{
		ID:            newUserID,
		IsActive:      false,
		IsDeleted:     false,
		InvitationKey: key,
	}

	db := app.DB
	i.unlink(db)
	err := db.Upsert("user", newUser)
	handleError(err)
	return Response{fmt.Sprintf("Invite key for user %s is `%s`", newUserID, key), nil, 0}
}
