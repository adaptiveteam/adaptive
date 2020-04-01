package ido

import (
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/adaptiveteam/adaptive/adaptive-reports/models"
)

// CreateIDOWorksheet creates a new sheet and sets Landscape orientation
func CreateIDOWorksheet(f *excel.File, title string) *models.Sheet {
	return models.NewSheet(f, title).Landscape()
}
