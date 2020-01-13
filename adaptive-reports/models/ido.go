package models

import (
	"github.com/adaptiveteam/adaptive/adaptive-reports/utilities"
	"strconv"
)

type IDOs []*IDO
type IDOUpdates []*IDOUpdate

type IDO struct {
	name                 string
	description          string
	advocate             string
	coach                string
	completed            bool
	createdDate          string
	focusedOnName        string
	focusedOnDescription string
	isA                  string
	driving              string
	updates              []*IDOUpdate
}

func NewIDO(
	name string,
	description string,
	advocate string,
	coach string,
	completed bool,
	createdDate string,
	focusedOnName string,
	focusedOnDescription string,
	isA string,
	driving string,
) *IDO {
	return &IDO{
		name:          name,
		description:   description,
		advocate:      advocate,
		coach:         coach,
		completed:     completed,
		createdDate:   createdDate,
		focusedOnName: focusedOnName,
		focusedOnDescription: focusedOnDescription,
		isA:           isA,
		driving:       driving,
	}
}

func (I *IDO) Driving() string {
	return I.driving
}

func (I *IDO) SetDriving(driving string) {
	I.driving = driving
}

func (I *IDO) IsA() string {
	return I.isA
}

func (I *IDO) SetIsA(isA string) {
	I.isA = isA
}

func (I *IDO) FocusedOnName() string {
	return I.focusedOnName
}

func (I *IDO) FocusedOnDescription() string {
	return I.focusedOnDescription
}

func (I *IDO) Updates() []*IDOUpdate {
	return I.updates
}

func (I *IDO) CreatedDate() string {
	return I.createdDate
}

func (I *IDO) Completed() bool {
	return I.completed
}

func (I *IDO) Coach() string {
	return I.coach
}

func (I *IDO) Advocate() string {
	return I.advocate
}

func (I *IDO) Description() string {
	return I.description
}

func (I *IDO) Name() string {
	return I.name
}

type IDOUpdate struct {
	updateDate       string
	advocateStatus   string
	advocateComments string
	coachStatus      string
	coachComments    string
}

func (I IDOUpdate) CoachComments() string {
	return I.coachComments
}

func (I IDOUpdate) CoachStatus() string {
	return I.coachStatus
}

func (I IDOUpdate) AdvocateComments() string {
	return I.advocateComments
}

func (I IDOUpdate) AdvocateStatus() string {
	return I.advocateStatus
}

func (I IDOUpdate) UpdateDate() string {
	return I.updateDate
}

func CreateIDOUpdateAndInsertIntoIDO(
	ido *IDO,
	updateDate string,
	advocateStatus string,
	advocateComments string,
	coachStatus string,
	coachComments string,
) (rv *IDOUpdate) {
	rv = &IDOUpdate{
		updateDate:       updateDate,
		advocateStatus:   advocateStatus,
		advocateComments: advocateComments,
		coachStatus:      coachStatus,
		coachComments:    coachComments,
	}
	ido.updates = append(ido.updates, rv)
	return rv
}

func CreateIDOs(table utilities.Table, rows int) (
	idoList IDOs,
) {
	idoList = make(IDOs, 0)
	var currentIDO *IDO

	for i := 0; i < rows; i++ {
		idoName := table.GetValue("ido_name", i)
		if currentIDO == nil || currentIDO.Name() != idoName {
			completed, _ := strconv.ParseBool(table.GetValue("completed", i))
			currentIDO = NewIDO(
				table.GetValue("ido_name", i),
				table.GetValue("ido_description", i),
				table.GetValue("advocate", i),
				table.GetValue("coach", i),
				completed,
				table.GetValue("ido_created_at", i),
				table.GetValue("focused_on_name", i),
				table.GetValue("focused_on_description", i),
				table.GetValue("is_a", i),
				table.GetValue("driving", i),
			)

			idoList = append(idoList, currentIDO)
		}

		CreateIDOUpdateAndInsertIntoIDO(
			currentIDO,
			table.GetValue("updated_at", i),
			table.GetValue("advocate_status", i),
			table.GetValue("advocate_comments", i),
			table.GetValue("coach_status", i),
			table.GetValue("coach_comments", i),
		)
	}

	return idoList
}
