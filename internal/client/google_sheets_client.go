package client

import (
	"errors"
	"fmt"

	"google.golang.org/api/sheets/v4"
)

type GoogleSheetsClient struct {
	srv *sheets.Service
}

func NewGoogleSheetsClient(srv *sheets.Service) *GoogleSheetsClient {
	return &GoogleSheetsClient{
		srv: srv,
	}
}

func (gsc *GoogleSheetsClient) GetNote(spreadsheetId string, sheetId int64, rowAndColumnRange string) (string, error) {
	resp, err := gsc.srv.Spreadsheets.Get(spreadsheetId).Ranges(rowAndColumnRange).IncludeGridData(true).Do()
	if err != nil {
		return "", err
	}
	for _, sheet := range resp.Sheets {
		if sheet.Properties.SheetId == sheetId {
			if len(sheet.Data) > 0 && len(sheet.Data[0].RowData) > 0 && len(sheet.Data[0].RowData[0].Values) > 0 {
				return sheet.Data[0].RowData[0].Values[0].Note, nil
			}
		}
	}
	return "", errors.New("note not found")
}

func (gsc *GoogleSheetsClient) GetValue(spreadsheetId, rowAndColumnRange string) (*sheets.ValueRange, error) {
	return gsc.srv.Spreadsheets.Values.Get(spreadsheetId, rowAndColumnRange).Do()
}

func (gsc *GoogleSheetsClient) GetSheetId(spreadsheetId, sheetName string) (int64, error) {
	resp, err := gsc.srv.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		return 0, err
	}
	for _, sheet := range resp.Sheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetId, nil
		}
	}
	return 0, fmt.Errorf("sheet with name \"%s\" not found", sheetName)
}

func (gsc *GoogleSheetsClient) BatchUpdate(spreadsheetId string, noteRequest *sheets.BatchUpdateSpreadsheetRequest) error {
	_, err := gsc.srv.Spreadsheets.BatchUpdate(spreadsheetId, noteRequest).Do()
	return err
}

func (gsc *GoogleSheetsClient) UpdateSheet(spreadsheetId, rowAndColumnRange string, newRow []interface{}) error {
	_, err := gsc.srv.Spreadsheets.Values.Update(spreadsheetId, rowAndColumnRange, &sheets.ValueRange{
		Values: [][]interface{}{newRow},
	}).ValueInputOption("USER_ENTERED").Do()
	return err
}
