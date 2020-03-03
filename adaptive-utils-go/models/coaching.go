package models

import (
	"github.com/adaptiveteam/adaptive/daos/common"
)
type CoachingRelationship struct {
	CoachQuarterYear   string     `json:"coach_quarter_year"`
	CoacheeQuarterYear string     `json:"coachee_quarter_year"`
	Coachee            string     `json:"coachee"`
	Quarter            int        `json:"quarter"`
	Year               int        `json:"year"`
	CoachRequested     bool       `json:"coach_requested"`
	CoacheeRequested   bool       `json:"coachee_requested"`
	PlatformID         common.PlatformID `json:"platform_id"`
}

type TargetQY struct {
	Target  string `json:"target"`
	Quarter int    `json:"quarter"`
	Year    int    `json:"year"`
}

type CoachingRejection struct {
	Id               string     `json:"id"`
	CoachRequested   bool       `json:"coach_requested"`
	CoacheeRequested bool       `json:"coachee_requested"`
	CoachRejected    bool       `json:"coach_rejected"`
	CoacheeRejected  bool       `json:"coachee_rejected"`
	Comments         string     `json:"comments"`
	PlatformID       common.PlatformID `json:"platform_id"`
}
