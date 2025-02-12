package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

const (
	ColumnPerMonth   = 6
	RowColumnPattern = "%s!%s%d"
)

func CleanMoneyValue(value string) string {
	newValue := strings.ReplaceAll(value, "R$", "")
	newValue = strings.ReplaceAll(newValue, ".", "")
	newValue = strings.ReplaceAll(newValue, ",", ".")

	return strings.TrimSpace(newValue)
}

func convertToXlsxColumn(columnNumber int) string {
	title := ""
	for columnNumber > 0 {
		columnNumber--
		title = string(rune(columnNumber%26+65)) + title
		columnNumber = columnNumber / 26
	}
	return title
}

func BuildNoteRequest(concatenatedNote string, sheetId int64, row, column int) *sheets.BatchUpdateSpreadsheetRequest {
	return &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				UpdateCells: &sheets.UpdateCellsRequest{
					Rows: []*sheets.RowData{
						{
							Values: []*sheets.CellData{
								{
									Note: concatenatedNote,
								},
							},
						},
					},
					Fields: "note",
					Range: &sheets.GridRange{
						SheetId:          sheetId,
						StartRowIndex:    int64(row - 1),
						EndRowIndex:      int64(row),
						StartColumnIndex: int64(column),
						EndColumnIndex:   int64(column + 1),
					},
				},
			},
		},
	}
}

/* range methods */

// saldo
func BuildBalanceRange() string {
	column := convertToXlsxColumn(GetCurrentBalanceColumnNumber())
	return buildRowColumnPattern(strconv.Itoa(GetCurrentYear()), column)
}

// diario
func BuildDailyOutcomeRange() string {
	column := convertToXlsxColumn(GetCurrentDailyOutcomeColumnNumber() + 1)
	return buildRowColumnPattern(strconv.Itoa(GetCurrentYear()), column)
}

// entrada
func BuildIncomeRange() string {
	column := convertToXlsxColumn(GetCurrentIncomeColumnNumber() + 1)
	return buildRowColumnPattern(strconv.Itoa(GetCurrentYear()), column)
}

func buildRowColumnPattern(sheetName string, column string) string {
	return fmt.Sprintf(RowColumnPattern, sheetName, column, GetCurrentDayRow())
}

/* row and column methods */

func GetCurrentIncomeColumnNumber() int {
	// 5 is the difference between the end of month range and the income 'entrada' column
	return getCurrentMonthColumn() - 5
}

func GetCurrentDailyOutcomeColumnNumber() int {
	// 3 is the difference between the end of month range and the outcome 'diario' column
	return getCurrentMonthColumn() - 3
}

func GetCurrentBalanceColumnNumber() int {
	// 1 is the difference between the end of month range and the balance 'saldo' column
	return getCurrentMonthColumn() - 1
}

/* date methods */

func GetCurrentYear() int {
	return time.Now().Year()
}

func getCurrentMonthColumn() int {
	// x6 because each month has 6 columns
	return int(time.Now().Month()) * ColumnPerMonth
}

func GetCurrentDayRow() int {
	// +2 because the first and second rows are reserved for the header
	return time.Now().Day() + 2
}
