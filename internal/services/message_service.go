package services

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/log"

	"github.com/vitortenor/sheet-bot/internal/configuration"
	"github.com/vitortenor/sheet-bot/internal/domain"
)

type MessageService struct {
	context            context.Context
	appConfig          *configuration.ApplicationConfig
	sheetService       *GoogleSheetsService
	aiService          *OllamaAIService
	interpreterService *MessageInterpreterService
}

func NewMessageService(ctx context.Context, appConfig *configuration.ApplicationConfig, gss *GoogleSheetsService,
	oas *OllamaAIService, mis *MessageInterpreterService) *MessageService {
	return &MessageService{
		context:            ctx,
		sheetService:       gss,
		aiService:          oas,
		appConfig:          appConfig,
		interpreterService: mis,
	}
}

func (ms *MessageService) ProcessAndReply(message *domain.Message) *domain.Message {
	log.Info("processing message")

	if message.CheckIfIsSystemMessage() {
		return nil
	}

	if ms.appConfig.Ai.IsEnabled {
		patronizedMessage := ms.aiService.GetOllamaAIResponse(message.Message)
		if patronizedMessage != "false" {
			log.Info("processing income/outcome message")
			message = &domain.Message{
				Message: patronizedMessage,
			}
			message.Normalize()
			return &domain.Message{
				Message: ms.sheetService.ProcessAndUpdateSheet(message.Message),
			}
		}
	}

	patronizedMessage := ms.interpreterService.InterpretMessage(message.Message)
	fmt.Println("Patronized message:", patronizedMessage)
	if patronizedMessage != false {
		message = &domain.Message{
			Message: patronizedMessage.(string),
		}
		if message.IsIncomeOrOutcome() {
			log.Info("processing income/outcome message")
			message.Normalize()
			return &domain.Message{
				Message: ms.sheetService.ProcessAndUpdateSheet(message.Message),
			}
		}
	}

	message.Normalize()

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

	if message.IsSetAsZero() {
		log.Info("processing set as zero message")
		return &domain.Message{
			Message: ms.sheetService.SetDailyAsZero(),
		}
	}

	return &domain.Message{
		Message: domain.InvalidMessage + ": " + message.Message,
	}
}
