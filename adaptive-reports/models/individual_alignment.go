package models

import "github.com/adaptiveteam/adaptive/adaptive-reports/utilities"

type Alignment struct {
	AdvocateStatus string
	CoachStatus string
	Updated string
	Advocate string
	Coach string
	IDOName string
	FocusedOn string
	IsA string
	Driving string
	UpdatedOn string
	CreatedOn string
	CompleteBy string
	FocusedOnID string
	DrivingID string
}

func (a Alignment) GetAdvocate() (rv string) {
	rv = a.Advocate
	return rv
}

func (a Alignment) GetCoach() (rv string) {
	rv = a.Coach
	return rv
}
func (a Alignment) GetIDOName() (rv string) {
	rv = a.IDOName
	return rv
}
func (a Alignment) GetFocusedOn() (rv string) {
	rv = a.FocusedOn
	return rv
}
func (a Alignment) GetFocusedOnID() (rv string) {
	rv = a.FocusedOnID
	return rv
}
func (a Alignment) GetIsA() (rv string) {
	rv = a.IsA
	return rv
}
func (a Alignment) GetDriving() (rv string) {
	rv = a.Driving
	return rv
}
func (a Alignment) GetDrivingID() (rv string) {
	rv = a.DrivingID
	return rv
}
func (a Alignment) GetCreatedOn() (rv string) {
	rv = a.CreatedOn
	return rv
}
func (a Alignment) GetCompleteBy() (rv string) {
	rv = a.CompleteBy
	return rv
}
func (a Alignment) GetUpdated() (rv string) {
	rv = a.Updated
	return rv
}
func (a Alignment) GetUpdatedOn() (rv string) {
	rv = a.UpdatedOn
	return rv
}
func (a Alignment) GetAdvocateStatus() (rv string) {
	rv = a.AdvocateStatus
	return rv
}
func (a Alignment) GetCoachStatus() (rv string) {
	rv = a.CoachStatus
	return rv
}

func CreateStrategyAlignments(table utilities.Table, rows int) (rv []Alignment) {
	rv = make([]Alignment, 0)

	for i := 0; i < rows; i++ {
		newAlignment := Alignment{
			Advocate:       table.GetValue("Advocate",i),
			Coach:          table.GetValue("Coach",i),
			IDOName:        table.GetValue("IDOName",i),
			FocusedOn:      table.GetValue("FocusedOn",i),
			FocusedOnID:    table.GetValue("FocusedOnID",i),
			IsA:            table.GetValue("IsA",i),
			Driving:        table.GetValue("Driving",i),
			DrivingID:      table.GetValue("DrivingID",i),
			CreatedOn:      table.GetValue("CreatedOn",i),
			CompleteBy:     table.GetValue("CompleteBy",i),
			AdvocateStatus: table.GetValue("AdvocateStatus",i),
			CoachStatus:    table.GetValue("CoachStatus",i),
			Updated:        table.GetValue("Updated",i),
			UpdatedOn:        table.GetValue("UpdatedOn",i),
		}
		rv = append(rv, newAlignment)
	}

	return rv
}