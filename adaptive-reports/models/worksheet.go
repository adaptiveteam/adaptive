package models

import (
	"fmt"
	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
	"log"
	"regexp"
	"strconv"
)

type Style struct {
	f *excel.File
	style string
	styles map[string]string
}

func GetCoordinates(column string, row int) (rv string) {
	c, err := excel.ColumnNameToNumber(column)
	if err != nil {
		log.Panic("invalid column name - "+column)
	}
	rv, err = excel.CoordinatesToCellName(c, row)
	if err != nil {

	}
	return rv
}

func NewStyle(style map[string]string) (rv Style) {
	rv.styles = style
	rv.style = `{}`
	return rv
}

func (style Style) GetStyle(desiredStyle string) (rv Style){
	rv.styles = style.styles
	if s, ok := style.styles[desiredStyle]; ok {
		rv.style = style.style[:len(style.style)-1]
		if len(style.style) > 2 {
			rv.style = rv.style + ","
		}
		rv.style = rv.style + s+"}"
	} else {
		rv = style
	}
	return rv
}

func (style Style) ToNewStyle (f *excel.File) (rv int) {
	rv, err := f.NewStyle(style.style)
	if err != nil {
		log.Panic("invalid style - "+style.style)
	}

	return rv
}

type Sheet struct {
	f         *excel.File
	sheetName string
	sheetNum int
}

type Cell struct {
	sheet *Sheet
	startingCell string
	endingCell string
}

func NewSheet(
	f *excel.File,
	sheetName string,
) *Sheet {
	sheetNum := f.NewSheet(sheetName)

	// Delete "Sheet1" if it exists. This sheet is the default on file creation
	// f.DeleteSheet("Sheet1")
	return &Sheet{
		f,
		// the excel limit on worksheet names is 31
		sheetName,
		sheetNum,
	}
}

func DeleteSheet(
	f *excel.File,
	sheetName string,
) {
	if f.GetSheetIndex(sheetName) != 0 {
		f.DeleteSheet(sheetName)
	}
}

func (s *Sheet) Landscape() *Sheet {
	err := s.f.SetPageLayout(
		s.sheetName,
		excel.PageLayoutOrientation(excel.OrientationLandscape),
	)
	if err != nil {
		log.Panic("Error setting Landscape for Sheet "+s.sheetName+" with error "+fmt.Sprint(err))
	}
	return s
}

func (s *Sheet) AddFilter(
	startFilterColumn int,
	startFilterRow int,
	endFilterColumn int,
	endFilterRow int,
) *Sheet {
	startingColumn, err := excel.ColumnNumberToName(startFilterColumn)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Error getting starting column for Sheet "+s.sheetName))
	}
	endingColumn, err := excel.ColumnNumberToName(endFilterColumn)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Error getting ending column for Sheet "+s.sheetName))
	}
	start := GetCoordinates(
		startingColumn,
		startFilterRow,
	)
	end := GetCoordinates(
		endingColumn,
		endFilterRow,
	)
	err = s.f.AutoFilter(s.sheetName, start, end, "")
	if err != nil {
		fmt.Println(errors.Wrap(err, "Error setting filter for Sheet "+s.sheetName+" with error "))
	}
	return s
}

func (s *Sheet) Activate() *Sheet {
	s.f.SetActiveSheet(s.sheetNum)
	return s
}

func (s *Sheet) NewCell(
	startMergeColumn int,
	startMergeRow int,
) (rv *Cell) {
	return s.NewMergedCell(
		startMergeColumn,
		startMergeRow,
		startMergeColumn,
		startMergeRow,
	)
}

func (s *Sheet) NewMergedCell(
	startMergeColumn int,
	startMergeRow int,
	endMergeColumn int,
	endMergeRow int,
) (rv *Cell) {
	startMergeColumnName,_ := excel.ColumnNumberToName(startMergeColumn)
	endMergeColumnName,_ := excel.ColumnNumberToName(endMergeColumn)
	rv = &Cell{
		s,
		GetCoordinates(startMergeColumnName, startMergeRow),
		GetCoordinates(endMergeColumnName, endMergeRow),
	}
	err := s.f.MergeCell(s.sheetName, rv.startingCell, rv.endingCell)
	if err != nil {
		log.Panic("Error merging cells from "+rv.startingCell+" to "+rv.endingCell+" with error "+fmt.Sprint(err))
	}
	return rv
}

func (m *Cell) Value(v string) *Cell {
	err := m.sheet.f.SetCellValue(m.sheet.sheetName,m.startingCell, v)
	if err != nil {
		log.Panic("Error Setting value "+v+" at location "+m.startingCell+" with error "+fmt.Sprint(err))
	}
	return m
}

func (m *Cell) FloatValue(v float64) *Cell {
	err := m.sheet.f.SetCellValue(m.sheet.sheetName,m.startingCell, v)
	if err != nil {
		log.Panic("Error Setting value "+fmt.Sprintf("%f",v)+" at location "+m.startingCell+" with error "+fmt.Sprint(err))
	}
	return m
}

func (m *Cell) IntValue(v int) *Cell {
	err := m.sheet.f.SetCellValue(m.sheet.sheetName,m.startingCell, v)
	if err != nil {
		log.Panic("Error Setting value "+fmt.Sprintf("%d",v)+" at location "+m.startingCell+" with error "+fmt.Sprint(err))
	}
	return m
}

func (m *Cell) Style(s int) *Cell {
	err := m.sheet.f.SetCellStyle(
		m.sheet.sheetName,
		m.startingCell,
		m.endingCell,
		s,
	)
	if err != nil {
		log.Panic("Error setting column style "+strconv.Itoa(s)+" at location "+m.startingCell+" with error "+fmt.Sprint(err))
	}
	return m
}

func (m *Cell) Width(w float64) *Cell {
	err := m.sheet.f.SetColWidth(m.sheet.sheetName,string(m.startingCell[0]),string(m.startingCell[0]), w)
	if err != nil {
		log.Panic("Error setting column width "+fmt.Sprintf("%f", w)+" at location "+m.startingCell+" with error "+fmt.Sprint(err))
	}
	return m
}

func (m *Cell) Height(h float64) *Cell {
	re := regexp.MustCompile(`^([A-Z]+)([1-9]\d*)$`)
	components := re.FindStringSubmatch(m.startingCell)
	rowNum,_ := strconv.Atoi(components[2])
	err := m.sheet.f.SetRowHeight(m.sheet.sheetName, rowNum, h)
	if err != nil {
		log.Panic("Error setting column width "+fmt.Sprintf("%f", h)+" at location "+m.startingCell+" with error "+fmt.Sprint(err))
	}
	return m
}

func (m *Cell) SetLink(sheetName, column string, row int) *Cell {
	cell := GetCoordinates(column, row)
	// err := f.SetCellHyperLink("Sheet1", "A3", "Sheet1!A40", "Location")
	location := `'`+sheetName+`'!`+cell
	err := m.sheet.f.SetCellHyperLink(m.sheet.sheetName, m.startingCell,location, "Location" )
	if err != nil {
		log.Panic("Error setting link to Cell at location "+cell+" with error "+fmt.Sprint(err))
	}
	return m
}
