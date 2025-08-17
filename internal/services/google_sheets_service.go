package services

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/vitortenor/sheet-bot/internal/client"
	"github.com/vitortenor/sheet-bot/internal/configuration"
	"github.com/vitortenor/sheet-bot/internal/domain"
	"github.com/vitortenor/sheet-bot/internal/utils"
)

const (
	ZeroBalance = domain.ZeroBalanceMessage
	SystemError = domain.SystemErrorMessage
)

type GoogleSheetsService struct {
	appConfig *configuration.ApplicationConfig
	client    *client.GoogleSheetsClient
}

func NewGoogleSheetsService(appConfig *configuration.ApplicationConfig, gsc *client.GoogleSheetsClient) *GoogleSheetsService {
	return &GoogleSheetsService{
		appConfig: appConfig,
		client:    gsc,
	}
}

func (gss *GoogleSheetsService) GetDailyExpenses() string {
	_, err := gss.client.GetSheetId(gss.appConfig.Google.SheetId, strconv.Itoa(utils.GetCurrentYear()))
	if err != nil {
		return SystemError
	}

	valueRange := utils.BuildDailyOutcomeRange()
	response, err := gss.client.GetValue(gss.appConfig.Google.SheetId, valueRange)
	if err != nil {
		return SystemError
	}

	if len(response.Values) > 0 && len(response.Values[0]) > 0 {
		return domain.SystemMessagePrefix + response.Values[0][0].(string)
	}

	return ZeroBalance
}

func (gss *GoogleSheetsService) GetBalance() string {
	_, err := gss.client.GetSheetId(gss.appConfig.Google.SheetId, strconv.Itoa(utils.GetCurrentYear()))
	if err != nil {
		return SystemError
	}

	response, err := gss.client.GetValue(gss.appConfig.Google.SheetId, utils.BuildBalanceRange())
	if err != nil {
		return SystemError
	}

	if len(response.Values) > 0 && len(response.Values[0]) > 0 {
		return domain.SystemMessagePrefix + response.Values[0][0].(string)
	}

	return ZeroBalance
}

func (gss *GoogleSheetsService) SetDailyAsZero() string {
	sheetId, err := gss.client.GetSheetId(gss.appConfig.Google.SheetId, strconv.Itoa(utils.GetCurrentYear()))
	if err != nil {
		return SystemError
	}

	existingNote, err := gss.client.GetNote(gss.appConfig.Google.SheetId, sheetId, utils.BuildDailyOutcomeRange())
	if err != nil {
		return SystemError
	}

	if existingNote != "" {
		return domain.SystemMessagePrefix + "daily value has notes, please remove them before setting as zero"
	}

	err = gss.client.UpdateSheet(gss.appConfig.Google.SheetId, utils.BuildDailyOutcomeRange(), []interface{}{"0"})
	if err != nil {
		return SystemError
	}

	return domain.SystemMessagePrefix + "daily value set as zero"
}

func (gss *GoogleSheetsService) ProcessAndUpdateSheet(inputValue string) string {
	sheetId, err := gss.client.GetSheetId(gss.appConfig.Google.SheetId, strconv.Itoa(utils.GetCurrentYear()))
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
	if value == 0 {
		return nil
	}

	isIncome := value > 0

	row := utils.GetCurrentDayRow()
	column := utils.GetCurrentIncomeColumnNumber()
	rowAndColumnRange := utils.BuildIncomeRange()

	if !isIncome {
		row = utils.GetCurrentDayRow()
		column = utils.GetCurrentDailyOutcomeColumnNumber()
		rowAndColumnRange = utils.BuildDailyOutcomeRange()
	}

	response, err := gss.client.GetValue(gss.appConfig.Google.SheetId, rowAndColumnRange)
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

	existingNote, err := gss.client.GetNote(gss.appConfig.Google.SheetId, sheetId, rowAndColumnRange)
	if err != nil {
		return err
	}

	if existingNote == "" && !isIncome {
		parsedValue = 0
	}

	newValue := strings.Replace(fmt.Sprintf("%.2f", parsedValue+math.Abs(value)), ".", ",", -1)

	err = gss.client.UpdateSheet(gss.appConfig.Google.SheetId, rowAndColumnRange, []interface{}{newValue})
	if err != nil {
		return err
	}

	description := fmt.Sprintf("%.2f - %s", math.Abs(value), strings.Split(inputValue, "/")[1])

	concatenatedNote := description
	if existingNote != "" {
		concatenatedNote = fmt.Sprintf("%s\n%s", existingNote, description)
	}

	noteRequest := utils.BuildNoteRequest(concatenatedNote, sheetId, row, column)
	err = gss.client.BatchUpdate(gss.appConfig.Google.SheetId, noteRequest)
	if err != nil {
		return err
	}

	return nil
}

func (gss *GoogleSheetsService) GetDetailedDailyBalance() string {
	sheetId, err := gss.client.GetSheetId(gss.appConfig.Google.SheetId, strconv.Itoa(utils.GetCurrentYear()))
	if err != nil {
		return SystemError
	}

	existingNote, err := gss.client.GetNote(gss.appConfig.Google.SheetId, sheetId, utils.BuildDailyOutcomeRange())
	if err != nil {
		return SystemError
	}

	if existingNote == "" {
		return domain.SystemMessagePrefix + "\n" + existingNote
	}

	if strings.Contains(existingNote, "\n") {
		notes := strings.Split(existingNote, "\n")
		var formattedNotes []string
		for _, note := range notes {
			formattedNote := domain.SystemMessagePrefix + " " + note
			formattedNotes = append(formattedNotes, formattedNote)
		}

		return strings.Join(formattedNotes, "\n")
	}

	return domain.SystemMessagePrefix + " " + existingNote
}
