package main

import (
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
	"time"
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
	ID           int           `reindex:"id,,pk"`
	ChatID       int           `reindex:"chat_id"`
	CronID       int           `reindex:"cron_id"`
	Name         string        `reindex:"name"`
	Times        []*Time       `reindex:"at,,joined"`
	EditedTimeID int           `reindex:"edited_time_id"`
	Period       time.Duration `reindex:"period"`
	Count        int           `reindex:"count"`
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
		Count:  0,
	}

	err := db.OpenNamespace("job", reindexer.DefaultNamespaceOptions(), Job{})
	_, err = db.Insert("job", job, "id=serial()")
	handleError(err)

	db.Query("user").
		WhereInt("chat_id", reindexer.EQ, u.ChatID).
		Set("edited_job_id", job.ID).
		Update()
}

func (u *User) setPeriod(hours int, db *reindexer.Reindexer) {
	db.Query("job").
		WhereInt("id", reindexer.EQ, u.EditedJobID).
		Set("period", hours*int(time.Hour)).
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

type NotFound struct{}

func (nf NotFound) Error() string {
	return "not found"
}

func (t Time) value(field string) int {
	values := map[string]int{
		"hour":   t.Hour,
		"minute": t.Minute,
	}
	return values[field]
}

func (u *User) setClockElem(time *Time, field string, db *reindexer.Reindexer) error {
	job, ok := u.findEditedJob(db)
	if !ok {
		return NotFound{}
	}

	if job.EditedTimeID == 0 {
		time.JobID = job.ID
		defaultUpsert(db, "time", time)
		job.setEditedTimeID(time.ID, db)
	} else {
		db.Query("time").
			WhereInt("id", reindexer.EQ, job.EditedTimeID).
			Set(field, time.value(field)).
			Update()
	}
	return nil
}

func (u *User) setHour(hour int, db *reindexer.Reindexer) error {
	return u.setClockElem(&Time{Hour: hour}, "hour", db)
}

func (u *User) setMinute(minute int, db *reindexer.Reindexer) error {
	return u.setClockElem(&Time{Minute: minute}, "minute", db)
}

func clearEdited(ns, field string, ID int, db *reindexer.Reindexer) {
	db.Query(ns).
		WhereInt("id", reindexer.EQ, ID).
		Set(field, 0).
		Update()
}

func (u *User) clearEditedTime(db *reindexer.Reindexer) {
	clearEdited("job", "edited_time_id", u.EditedJobID, db)
}

func (u *User) clearEdited(db *reindexer.Reindexer) {
	db.Query("user").
		WhereInt("chat_id", reindexer.EQ, u.ChatID).
		Set("edited_job_id", 0).
		Update()
}

func setUserState(chatID int, state string, db *reindexer.Reindexer) {
	db.Query("user").
		WhereInt("chat_id", reindexer.EQ, chatID).
		Set("state", state).
		Update()
}
func (u *User) setState(state string, db *reindexer.Reindexer) {
	db.Query("user").
		WhereInt("chat_id", reindexer.EQ, u.ChatID).
		Set("state", state).
		Update()
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
