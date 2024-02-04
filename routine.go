package main

import (
	"fmt"
	"github.com/restream/reindexer/v4"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Year struct {
	Validator
	Out int
}

func validateMax(value, max int) (bool, string) {
	if value <= max {
		return true, ""
	}
	return false, fmt.Sprintf("Число должно быть не больше %v", max)
}

func makeYear(inp string) *Year {
	year := &Year{}
	year.Inp = inp
	year.OK = true
	return year
}

func (y *Year) validate() {
	ok, value, err := validateInt(y.Inp)
	if y.checkOK(ok, err) {
		y.Out = value
	}

	y.checkOK(validateMin(y.Out, time.Now().Year()))
}

func (r *Response) handleFirstYear(inp string, user *User, db *reindexer.Reindexer) {
	year := makeYear(inp)
	year.validate()
	if year.OK {
		job, ok := user.findEditedJob(db)
		if ok {
			job.setNextYear(year.Out, db)
			user.setState(InpFirstStartMonth, db)
			r.Text = "Введите месяц первого приёма"
			r.Buttons = months()
		}
	} else {
		r.Text = strings.Join(year.Errors, ". ")
	}
}

type Month struct {
	Validator
	Out string
}

func makeMonth(inp string) *Month {
	month := &Month{}
	month.Inp = inp
	month.OK = true
	return month
}

func (m *Month) validate() {
	if slices.Contains(months(), m.Inp) {
		m.Out = m.Inp
	} else {
		m.OK = false
		m.Errors = []string{"Месяц не распознан"}
	}
}

func (r *Response) handleFirstMonth(inp string, user *User, db *reindexer.Reindexer) {
	month := makeMonth(inp)
	month.validate()
	if month.OK {
		job, ok := user.findEditedJob(db)
		if ok {
			job.setNextMonth(month.Out, db)
			user.setState(InpFirstStartDay, db)
			r.Text = "Введите первое число приёма"
		}
	} else {
		r.Text = strings.Join(month.Errors, ". ")
	}
}

type Day struct {
	Validator
	Out int
}

func makeDay(inp string) *Day {
	day := &Day{}
	day.Inp = inp
	day.OK = true
	return day
}

func (d *Day) validate() {
	ok, value, err := validateInt(d.Inp)
	if d.checkOK(ok, err) {
		d.Out = value
	}

	d.checkOK(validateMin(d.Out, 1))
	d.checkOK(validateMax(d.Out, 31))
}

func (r *Response) handleFirstDay(inp string, user *User, db *reindexer.Reindexer) {
	day := makeDay(inp)
	day.validate()
	if day.OK {
		job, ok := user.findEditedJob(db)
		if ok {
			job.setNextDay(day.Out, db)
			user.setState(InpFirstStartHour, db)
			r.Text = "Введите первый час приёма"
		}
	} else {
		r.Text = strings.Join(day.Errors, ". ")
	}
}

type Hour struct {
	Validator
	Out int
}

func makeHour(inp string) *Hour {
	hour := &Hour{}
	hour.Inp = inp
	hour.OK = true
	return hour
}

func (h *Hour) validate() {
	ok, value, err := validateInt(h.Inp)
	if h.checkOK(ok, err) {
		h.Out = value
	}
	h.checkOK(validateMin(h.Out, 0))
	h.checkOK(validateMax(h.Out, 23))
}

func (r *Response) handleFirstHour(inp string, user *User,
	db *reindexer.Reindexer) {
	hour := makeHour(inp)
	hour.validate()
	if hour.OK {
		job, ok := user.findEditedJob(db)
		if ok {
			job.setNextHour(hour.Out, db)
			user.setState(InpFirstStartMinute, db)
			r.Text = "Введите первую минуту приёма"
		}
	} else {
		r.Text = strings.Join(hour.Errors, ". ")
	}
}

type Minute struct {
	Validator
	Out int
}

func makeMinute(inp string) *Minute {
	minute := &Minute{}
	minute.Inp = inp
	minute.OK = true
	return minute
}

func (m *Minute) validate() {
	ok, value, err := validateInt(m.Inp)
	if m.checkOK(ok, err) {
		m.Out = value
	}

	m.checkOK(validateMin(m.Out, 0))
	m.checkOK(validateMax(m.Out, 60))
}

func (r *Response) handleFirstMinute(inp string, user *User,
	db *reindexer.Reindexer) {
	minute := makeMinute(inp)
	minute.validate()
	if minute.OK {
		job, ok := user.findEditedJob(db)
		if ok {
			job.setNextMinute(minute.Out, db)
			if job.NextTime.After(time.Now()) {
				user.setState(InpPeriod, db)
				r.Text = "Каждые сколько часов принимать?"
			} else {
				user.setState(InpFirstStartYear, db)
				r.Text = "Дата напоминания должна быть в будущем"
				r.Buttons = []string{strconv.Itoa(time.Now().Year())}
			}
		}
	} else {
		r.Text = strings.Join(minute.Errors, ". ")
	}
}
