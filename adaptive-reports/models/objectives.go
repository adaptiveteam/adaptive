package models

import "github.com/adaptiveteam/adaptive-reports/utilities"

type Objectives map[string]*Objective
type Initiatives map[string]*Initiative

type Objective struct {
	ObjectiveName        string
	ObjectiveDescription string
	ObjectiveType        string
	ObjectiveAdvocate    string
	ObjectiveStatus      string
	ObjectiveUpdate      string
	ObjectiveCreatedOn   string
	ObjectiveUpdatedOn   string
	ObjectiveEndDate     string
	PercentTimeLeft      string
	Index                int
	Initiatives          []*Initiative
}

type Initiative struct {
	InitiativeName        string
	InitiativeDescription string
	InitiativeAdvocate    string
	InitiativeStatus      string
	InitiativeUpdate      string
	InitiativeCreatedOn   string
	InitiativeEndDate     string
	InitiativeUpdatedOn   string
	PercentTimeLeft       string
	ObjectiveID           string
	Index int
}

type StrategyComponent interface {
	GetName()            string
	GetDescription()     string
	GetAdvocate()        string
	GetStatus()          string
	GetUpdate()          string
	GetCreatedOn()       string
	GetUpdatedOn()       string
	GetEndDate()         string
	GetPercentTimeLeft() string
	GetIndex()           int
	SetIndex(int)
}

func (c Objective) GetName() string {
	return c.ObjectiveName
}

func (c Objective) GetDescription() string {
	return c.ObjectiveDescription
}

func (c Objective) GetUpdate() string {
	return c.ObjectiveUpdate
}

func (c Objective) GetAdvocate() string {
	return c.ObjectiveAdvocate
}

func (c Objective) GetStatus() string {
	return c.ObjectiveStatus
}

func (c Objective) GetCreatedOn() string {
	return c.ObjectiveCreatedOn
}

func (c Objective) GetUpdatedOn() string {
	return c.ObjectiveUpdatedOn
}

func (c Objective) GetEndDate() string {
	return c.ObjectiveEndDate
}

func (c Objective) GetPercentTimeLeft() string {
	return c.PercentTimeLeft
}

func (c Objective) GetIndex() int {
	return c.Index
}

func (c *Objective) SetIndex(index int) {
	c.Index = index
}

func (c Initiative) GetName() string {
	return c.InitiativeName
}

func (c Initiative) GetDescription() string {
	return c.InitiativeDescription
}

func (c Initiative) GetUpdate() string {
	return c.InitiativeUpdate
}

func (c Initiative) GetAdvocate() string {
	return c.InitiativeAdvocate
}

func (c Initiative) GetStatus() string {
	return c.InitiativeStatus
}

func (c Initiative) GetCreatedOn() string {
	return c.InitiativeCreatedOn
}

func (c Initiative) GetUpdatedOn() string {
	return c.InitiativeUpdatedOn
}

func (c Initiative) GetEndDate() string {
	return c.InitiativeEndDate
}

func (c Initiative) GetPercentTimeLeft() string {
	return c.PercentTimeLeft
}

func (c *Initiative) SetIndex(index int) {
	c.Index = index
}

func (c Initiative) GetIndex() int {
	return c.Index
}

func CreateObjectives(table utilities.Table, rows int) (
	objectivesMap Objectives,
	initiativesMap Initiatives,
) {
	objectivesMap = make(Objectives)
	initiativesMap = make(Initiatives)
	// Load the Objectives
	for i := 0; i < rows; i++ {
		objectiveID := table.GetValue("objective_id", i)
		newInitiative := Initiative{
			InitiativeName:        table.GetValue("Initiative Name", i),
			InitiativeDescription: table.GetValue("Initiative Description", i),
			InitiativeAdvocate:    table.GetValue("Initiative Advocate", i),
			InitiativeStatus:      table.GetValue("Initiative Status", i),
			InitiativeUpdate:      table.GetValue("Initiative Update", i),
			InitiativeCreatedOn:   table.GetValue("Initiative Created On", i),
			InitiativeUpdatedOn:   table.GetValue("Initiative Updated On", i),
			InitiativeEndDate:     table.GetValue("Initiative End Date", i),
			PercentTimeLeft:       table.GetValue("Initiative Time Left", i),
			ObjectiveID:           objectiveID,
		}
		initiativesMap[table.GetValue("initiative_id", i)] = &newInitiative
		o, ok := objectivesMap[objectiveID]
		if !ok {
				o = &Objective{
					ObjectiveName:        table.GetValue("Objective Name", i),
					ObjectiveDescription: table.GetValue("Objective Description", i),
					ObjectiveType:        table.GetValue("Objective Type", i),
					ObjectiveAdvocate:    table.GetValue("Objective Advocate", i),
					ObjectiveStatus:      table.GetValue("Objective Status", i),
					ObjectiveUpdate:      table.GetValue("Objective Update", i),
					ObjectiveCreatedOn:   table.GetValue("Objective Created On", i),
					ObjectiveUpdatedOn:   table.GetValue("Objective Updated On", i),
					ObjectiveEndDate:     table.GetValue("Objective End Date", i),
					PercentTimeLeft:      table.GetValue("Objective Time Left", i),

					Initiatives: []*Initiative {
						&newInitiative,
					},
				}
				objectivesMap[objectiveID] = o
		} else {
			objectivesMap[objectiveID].Initiatives = append(objectivesMap[objectiveID].Initiatives, &newInitiative)
		}
	}
	return objectivesMap, initiativesMap
}
