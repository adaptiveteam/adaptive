package strategy

import (
	"fmt"
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive-reports/models"
	"reflect"
	"strconv"
	"strings"
)

func findCompetency(competency string, competencies []models.Competency) (rv int) {
	for i := 0; i < len(competencies) && rv == 0 ; i++ {
		if competencies[i].Name == competency {
			rv = i+1
		}
	}
	return rv
}

func CreateAlignmentSummary(
	f *excel.File,
	strategicAlignment *models.Sheet,
	allAlignments []models.Alignment,
	allObjectives models.Objectives,
	allInitiatives models.Initiatives,
	allCompetencies []models.Competency,
	styles map[string]string,
) {
	const startingRow = 1
	type row struct {
		name string
		width float64
	}
	rows := [] row {
		{"individual status",20},
		{"coach status",20},
		{"weekly update?",20},
		{"this individual",30},
		{"is working with",30},
		{"together they are working on",50},
		{"that aligns with",50},
		{"which is a",15},
		{"that drives",50},
		{"updated on",20},
		{"created on",20},
		{"finish by",20},
	}

	column := 1
	for _,r := range rows {
		strategicAlignment.NewCell(
			column, startingRow,
		).Value(r.name).
			Style(
				models.NewStyle(styles).
					GetStyle("Heading Font").
					GetStyle("Heading Background").
					GetStyle("Horizontal Center").
					GetStyle("White Borders").
					ToNewStyle(f),
			).Width(r.width)
		column++
	}
	strategicAlignment.AddFilter(1, 1, len(rows), 1)

	for r,a := range allAlignments {
		v := reflect.ValueOf(a)
		// minus two because we don't want the ID fields
		for i := 0; i< v.NumField()-2; i++ {
			columnValue := fmt.Sprint(v.Field(i).Interface())
			alignsWith := strategicAlignment.NewCell(
				i+1,
				r+startingRow+1,
			).Value(columnValue)

			if strings.Contains(rows[i].name, "status") || rows[i].name == "weekly update?" {
				alignsWith.
					Style(
						models.NewStyle(styles).
							GetStyle(columnValue).
							GetStyle("Black Borders").
							GetStyle("Italics Font").
							GetStyle("Vertical Center").
							ToNewStyle (f),
					)
			} else {
				alignsWith.
					Style(
						models.NewStyle(styles).
							GetStyle("Black Borders").
							GetStyle("Normal Font").
							GetStyle("Normal Background").
							ToNewStyle (f),
					)
			}
			if rows[i].name == "that aligns with" || rows[i].name == "that drives" {
				alignsWith.Style(
					models.NewStyle(styles).
						GetStyle("Black Borders").
						GetStyle("Vertical Center").
						GetStyle("Link Font").
						GetStyle("Normal Background").
						ToNewStyle(f),
				)
				if rows[i].name == "that aligns with" {
					switch a.IsA {
					case "Initiative":
						componentLink := allInitiatives[a.FocusedOnID].GetIndex()
						alignsWith.SetLink("Initiative "+strconv.Itoa(componentLink),"A",1)
					case "Objective":
						componentLink := allObjectives[a.FocusedOnID].GetIndex()
						alignsWith.SetLink("Objective "+strconv.Itoa(componentLink),"A",1)
					case "Competency":
						competencyIndex := findCompetency(columnValue, allCompetencies)
						alignsWith.SetLink("Competencies","A",competencyIndex+startingRow)
					}
				}
				if rows[i].name == "that drives" {
					if a.IsA == "Initiative" {
						componentLink := allObjectives[allInitiatives[a.FocusedOnID].ObjectiveID].GetIndex()
						alignsWith.SetLink("Objective "+strconv.Itoa(componentLink),"A",1)
					}
					if a.IsA == "Objective" {
						alignsWith.SetLink("Performance Summary","A",2)
					}
					if a.IsA == "Competency" {
						alignsWith.SetLink("Competencies","A",1)
					}
				}
			}
		}
	}
}
