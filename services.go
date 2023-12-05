package main

import (
	"fmt"
	"github.com/restream/reindexer/v3"
)

type ServiceInstance struct {
	ServiceID string `reindex:"service_id"`
	UserID    string `reindex:"user_id"`
	State     string `reindex:"state"`
}

type Invite struct {
	ServiceInstance
}

func (i Invite) start(userId string, db *reindexer.Reindexer) Response {
	i.ServiceID = "invite"
	i.State = "input user_id"
	i.UserID = userId

	err := db.Upsert("service_instance", i)
	handleError(err)

	_, err = db.Update("user", &User{
		ID:               userId,
		CurrentServiceID: "invite",
	})
	handleError(err)
	return Response{"Enter ID of the user to be invited", nil, 0}
}

func (i Invite) next(inp string, app *App) Response {
	key := inviteKey()
	newUserID := inp

	newUser := User{
		ID:            newUserID,
		IsActive:      false,
		IsDeleted:     false,
		InvitationKey: key,
	}

	db := app.DB
	err := db.Upsert("user", newUser)
	handleError(err)

	_, err = db.Update("user", &User{
		ID:               i.UserID,
		CurrentServiceID: "",
	})
	handleError(err)

	err = db.Delete("service_instance", &ServiceInstance{ServiceID: i.ServiceID, UserID: i.UserID})
	handleError(err)
	return Response{fmt.Sprintf("Invite key for user %s is `%s`", newUserID, key), nil, 0}
}
