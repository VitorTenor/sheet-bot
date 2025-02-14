package services

import (
	"context"

	"github.com/labstack/gommon/log"

	"github.com/vitortenor/sheet-bot/internal/domain"
)

type MessageService struct {
	context      context.Context
	sheetService *GoogleSheetsService
}

func NewMessageService(ctx context.Context, gss *GoogleSheetsService) *MessageService {
	return &MessageService{
		context:      ctx,
		sheetService: gss,
	}
}

func (ms *MessageService) ProcessAndReply(message *domain.Message) *domain.Message {
	log.Info("processing message")
	if message.CheckIfIsSystemMessage() {
		return nil
	}

	message.Normalize()

	if message.IsIncomeOrOutcome() {
		log.Info("processing income/outcome message")
		return &domain.Message{
			Message: ms.sheetService.ProcessAndUpdateSheet(message.Message),
		}
	}

	if message.IsDailyExpense() {
		log.Info("processing daily expenses message")
		return &domain.Message{
			Message: ms.sheetService.GetDailyExpenses(),
		}
	}

	if message.IsDailyBalance() {
		log.Info("processing daily balance message")
		return &domain.Message{
			Message: ms.sheetService.GetBalance(),
		}
	}

	return &domain.Message{
		Message: domain.InvalidMessage + ": " + message.Message,
	}
}
