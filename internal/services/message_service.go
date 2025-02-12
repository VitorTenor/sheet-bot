package services

import (
	"context"

	"github.com/vitortenor/sheet-bot-api/internal/domain"
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
	if message.CheckIfIsSystemMessage() {
		return nil
	}

	message.Normalize()

	if message.CheckMessage() {
		return &domain.Message{
			Message: ms.sheetService.ProcessAndUpdateSheet(message.Message),
		}
	}

	if message.IsDailyExpense() {
		return &domain.Message{
			Message: ms.sheetService.GetDailyExpenses(),
		}
	}

	if message.IsDailyBalance() {
		return &domain.Message{
			Message: ms.sheetService.GetBalance(),
		}
	}

	return &domain.Message{
		Message: domain.InvalidMessage + ": " + message.Message,
	}
}
