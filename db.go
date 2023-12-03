package main

import (
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

type User struct {
	ID            string `reindex:"id,,pk"`
	ChatID        int    `reindex:"chat_id"`
	IsActive      bool   `reindex:"is_active"`
	IsDeleted     bool   `reindex:"is_deleted"`
	InvitationKey string `reindex:"invitation_key"`
}
