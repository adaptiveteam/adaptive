package models

import "github.com/adaptiveteam/adaptive-reports/utilities"

type Competency struct {
	Name        string
	Description string
	Type        string
}

func (c Competency) GetName() string {
	return c.Name
}

func (c Competency) GetDescription() string {
	return c.Description
}

func CreateCompetencies(table utilities.Table, rows int) (rv []Competency) {
	rv = make([]Competency,0)
	for i := 0; i < rows; i++ {
		rv = append(
			rv,
			Competency {
				Name:        table.GetValue("name", i),
				Description: table.GetValue("description", i),
				Type: table.GetValue("type", i),
			},
		)
	}

	return rv
}