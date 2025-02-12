package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

const (
	RowColumnPattern = "%s!%s%d"
)

func GetCurrentYear() int {
	return time.Now().Year()
}

func GetMoneyInputColumn() string {
	return convertToTitle((-5 + 6*getCurrentMonthNumber()) + 1)
}

func GetMoneyOutputDailyColumn() string {
	return convertToTitle(-5 + 6*getCurrentMonthNumber() + 3)
}

func GetMoneyOutputBalanceColumn() string {
	return convertToTitle(-5 + 6*getCurrentMonthNumber() + 4)
}

func GetCurrentDayColumn() int {
	return time.Now().Day() + 2
}

func LetterToIndex(letter string) int {
	return int(letter[0]) - 65
}

func CleanMoneyValue(value string) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, "R$", ""), ".", ""), ",", "."))
}

func BuildNoteRequest(concatenatedNote string, sheetId int64, rowIndex, noteColumnIndex int) *sheets.BatchUpdateSpreadsheetRequest {
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
						StartRowIndex:    int64(rowIndex + 3),
						EndRowIndex:      int64(rowIndex + 4),
						StartColumnIndex: int64(noteColumnIndex),
						EndColumnIndex:   int64(noteColumnIndex + 1),
					},
				},
			},
		},
	}
}

func BuildMoneyBalance() string {
	return buildRowColumnPattern(strconv.Itoa(GetCurrentYear()), GetMoneyOutputBalanceColumn(), GetCurrentDayColumn())
}

func BuildMoneyDailyOutput() string {
	return buildRowColumnPattern(strconv.Itoa(GetCurrentYear()), GetMoneyOutputDailyColumn(), GetCurrentDayColumn())
}

func BuildMoneyInput() string {
	return buildRowColumnPattern(strconv.Itoa(GetCurrentYear()), GetMoneyInputColumn(), GetCurrentDayColumn())
}

func buildRowColumnPattern(sheetName string, column string, row int) string {
	return fmt.Sprintf(RowColumnPattern, sheetName, column, row)
}

func convertToTitle(columnNumber int) string {
	title := ""
	for columnNumber > 0 {
		columnNumber--
		title = string(rune(columnNumber%26+65)) + title
		columnNumber = columnNumber / 26
	}
	return title
}

func getCurrentMonthNumber() int {
	return int(time.Now().Month())
}
