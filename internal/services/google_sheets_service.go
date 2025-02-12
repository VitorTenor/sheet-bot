package services

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/vitortenor/sheet-bot-api/internal/client"
	"github.com/vitortenor/sheet-bot-api/internal/domain"
	"github.com/vitortenor/sheet-bot-api/internal/utils"
)

const (
	SpreadsheetID = "1M8vonVq4defB0LfdqQbU26vyyJb3s2I0R2Z0YfAKodQ"
	ZeroBalance   = domain.ZeroBalanceMessage
	SystemError   = domain.SystemErrorMessage
)

type GoogleSheetsService struct {
	client *client.GoogleSheetsClient
}

func NewGoogleSheetsService(gsc *client.GoogleSheetsClient) *GoogleSheetsService {
	return &GoogleSheetsService{
		client: gsc,
	}
}

func (gss *GoogleSheetsService) GetDailyExpenses() string {
	_, err := gss.client.GetSheetId(SpreadsheetID, strconv.Itoa(utils.GetCurrentYear()))
	if err != nil {
		return SystemError
	}

	response, err := gss.client.GetValue(SpreadsheetID, utils.BuildMoneyDailyOutput())
	if err != nil {
		return SystemError
	}

	if len(response.Values) > 0 && len(response.Values[0]) > 0 {
		return domain.SystemMessagePrefix + response.Values[0][0].(string)
	}

	return ZeroBalance
}

func (gss *GoogleSheetsService) GetBalance() string {
	_, err := gss.client.GetSheetId(SpreadsheetID, strconv.Itoa(utils.GetCurrentYear()))
	if err != nil {
		return SystemError
	}

	response, err := gss.client.GetValue(SpreadsheetID, utils.BuildMoneyBalance())
	if err != nil {
		return SystemError
	}

	if len(response.Values) > 0 && len(response.Values[0]) > 0 {
		return domain.SystemMessagePrefix + response.Values[0][0].(string)
	}

	return ZeroBalance
}

func (gss *GoogleSheetsService) ProcessAndUpdateSheet(inputValue string) string {
	sheetId, err := gss.client.GetSheetId(SpreadsheetID, strconv.Itoa(utils.GetCurrentYear()))
	if err != nil {
		return SystemError
	}

	err = gss.updateSheetValuesAndNotes(sheetId, inputValue)
	if err != nil {
		return SystemError
	}

	return domain.SystemMessagePrefix + "processed " + inputValue
}

func (gss *GoogleSheetsService) updateSheetValuesAndNotes(sheetId int64, inputValue string) error {
	value, err := strconv.ParseFloat(strings.TrimSpace(strings.Split(inputValue, "/")[0]), 64)
	if err != nil {
		return err
	}

	if value != 0 {
		isIncome := value > 0

		rowIndex := utils.LetterToIndex(utils.GetMoneyInputColumn()) + 2
		noteColumnIndex := utils.LetterToIndex(utils.GetMoneyInputColumn())

		column := utils.BuildMoneyInput()

		if !isIncome {
			column = utils.BuildMoneyDailyOutput()
			rowIndex = utils.LetterToIndex(utils.GetMoneyOutputDailyColumn())
			noteColumnIndex = utils.LetterToIndex(utils.GetMoneyOutputDailyColumn())
		}

		response, err := gss.client.GetValue(SpreadsheetID, column)
		if err != nil {
			return err
		}

		currentValue := "0"
		if len(response.Values) > 0 && len(response.Values[0]) > 0 {
			currentValue = utils.CleanMoneyValue(response.Values[0][0].(string))
		}

		parsedValue, err := strconv.ParseFloat(currentValue, 64)
		if err != nil {
			return err
		}

		existingNote, err := gss.client.GetNote(SpreadsheetID, sheetId, column)
		if err != nil {
			return err
		}

		if existingNote == "" {
			parsedValue = 0
		}

		newValue := strconv.FormatFloat(parsedValue+math.Abs(value), 'f', -1, 64)

		err = gss.client.UpdateSheet(SpreadsheetID, column, []interface{}{newValue})
		if err != nil {
			return err
		}

		description := fmt.Sprintf("%f - %s", math.Abs(value), strings.Split(inputValue, "/")[1])

		concatenatedNote := description
		if existingNote != "" {
			concatenatedNote = fmt.Sprintf("%s\n%s", existingNote, description)
		}

		noteRequest := utils.BuildNoteRequest(concatenatedNote, sheetId, rowIndex, noteColumnIndex)
		err = gss.client.BatchUpdate(SpreadsheetID, noteRequest)
		if err != nil {
			return err
		}

	}
	return nil
}
