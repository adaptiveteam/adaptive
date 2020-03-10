package model

import "time"

type DBModel struct {
	DBCreatedAt time.Time
	DBUpdatedAt time.Time
	DBDeletedAt *time.Time `sql:"index"`
}

func (d DBModel) AsAdd() (op DBModel) {
	op = d
	currentTime := time.Now()
	op.DBCreatedAt = currentTime
	op.DBUpdatedAt = currentTime
	return
}

func (d DBModel) AsUpdate() (op DBModel) {
	op = d
	currentTime := time.Now()
	op.DBUpdatedAt = currentTime
	return
}

func (d DBModel) AsDelete() (op DBModel) {
	op = d
	currentTime := time.Now()
	op.DBDeletedAt = &currentTime
	return
}
