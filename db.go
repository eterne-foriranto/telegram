package main

import (
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

type User struct {
	ID            string `reindex:"id"`
	ChatID        int    `reindex:"chat_id,,pk"`
	IsActive      bool   `reindex:"is_active"`
	IsDeleted     bool   `reindex:"is_deleted"`
	InvitationKey string `reindex:"invitation_key"`
	State         string `reindex:"state"`
	Jobs          []*Job `reindex:"job,,joined"`
	EditedJobID   int    `reindex:"edited_job_id"`
}

type Job struct {
	ID           int     `reindex:"id,,pk"`
	ChatID       int     `reindex:"chat_id"`
	Name         string  `reindex:"name"`
	Times        []*Time `reindex:"at,,joined"`
	EditedTimeID int     `reindex:"edited_time_id"`
}

type Time struct {
	ID     int `reindex:"id,,pk"`
	JobID  int `reindex:"job_id"`
	Hour   int `reindex:"hour"`
	Minute int `reindex:"minute"`
}

func findUser(chatID int, db *reindexer.Reindexer) (*User, bool) {
	iterator := db.Query("user").
		WhereInt("chat_id", reindexer.EQ, chatID).
		Exec()
	user, err := iterator.FetchOne()
	if err != nil {
		return nil, false
	}
	return user.(*User), true
}

func defaultUpsert(db *reindexer.Reindexer, ns string, item interface{}) {
	err := db.OpenNamespace(ns, reindexer.DefaultNamespaceOptions(), item)
	handleError(err)
	err = db.Upsert(ns, item)
	handleError(err)
}

func (u *User) attachJob(name string, db *reindexer.Reindexer) {
	job := &Job{
		Name:   name,
		ChatID: u.ChatID,
	}
	defaultUpsert(db, "job", job)

	db.Query("user").
		WhereInt("chat_id", reindexer.EQ, u.ChatID).
		Set("edited_job_id", job.ID).
		Update()
}

func (u *User) findEditedJob(db *reindexer.Reindexer) (*Job, bool) {
	iterator := db.Query("job").
		WhereInt("id", reindexer.EQ, u.EditedJobID).
		Exec()
	job, err := iterator.FetchOne()
	if err != nil {
		return nil, false
	}
	return job.(*Job), true
}

func findEditedTime(id int, db *reindexer.Reindexer) (*Time, bool) {
	iterator := db.Query("time").
		WhereInt("id", reindexer.EQ, id).
		Exec()
	time, err := iterator.FetchOne()
	if err != nil {
		return nil, false
	}
	return time.(*Time), true
}

func (j *Job) setEditedTimeID(timeID int, db *reindexer.Reindexer) {
	db.Query("job").
		WhereInt("id", reindexer.EQ, j.ID).
		Set("edited_time_id", timeID).
		Update()
}

func (u *User) setHour(hour int, db *reindexer.Reindexer) {
	job, ok := u.findEditedJob(db)
	if ok {
		if job.EditedTimeID == 0 {
			time := &Time{
				JobID: job.ID,
				Hour:  hour,
			}
			defaultUpsert(db, "time", time)
			job.setEditedTimeID(time.ID, db)
		} else {
			time, ok := findEditedTime(job.EditedTimeID, db)
			if ok {
				time
			}
		}
	}
}

func (u *User) setState(state string, db *reindexer.Reindexer) {
	db.Query("user").
		WhereInt("chat_id", reindexer.EQ, u.ChatID).
		Set("state", state).
		Update()
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
