package main

import (
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

type DB *reindexer.Reindexer

type User struct {
	ID               string `reindex:"id,,pk"`
	ChatID           int    `reindex:"chat_id"`
	IsActive         bool   `reindex:"is_active"`
	IsDeleted        bool   `reindex:"is_deleted"`
	InvitationKey    string `reindex:"invitation_key"`
	CurrentServiceID string `reindex:"current_service_id"`
}

type Service struct {
	ID string `reindex:"id,,pk"`
}

func currentService(chatID int, db *reindexer.Reindexer) string {
	query := db.Query("user").
		WhereInt("chat_id", reindexer.EQ, chatID)
	iterator := query.Exec()
	defer iterator.Close()
	for iterator.Next() {
		elem := iterator.Object().(*User)
		return elem.CurrentServiceID
	}
	return ""
}

type userNotFound struct{}

func (unf userNotFound) Error() string {
	return "User not found"
}

func chatIDByUserID(userID string, db *reindexer.Reindexer) (int64, error) {
	query := db.Query("user").
		WhereString("id", reindexer.EQ, userID)
	iterator := query.Exec()
	defer iterator.Close()
	for iterator.Next() {
		elem := iterator.Object().(*User)
		return int64(elem.ChatID), nil
	}
	return 0, userNotFound{}
}
