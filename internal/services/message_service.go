package services

import (
	"context"

	"github.com/vitortenor/sheet-bot-api/internal/domain"
)

type MessageService struct {
	sheetService *GoogleSheetsService
}

func NewMessageService(gss *GoogleSheetsService) *MessageService {
	return &MessageService{
		sheetService: gss,
	}
}

func (ms *MessageService) ProcessAndReply(_ context.Context, message *domain.Message) *domain.Message {
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
