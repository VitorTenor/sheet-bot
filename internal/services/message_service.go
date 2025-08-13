package services

import (
	"context"

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
		if resp := ms.aiService.GetOllamaAIResponse(message.Message); resp != "false" {
			return ms.processIncomeOutcome(&domain.Message{
				Message: resp,
			})
		}
	}

	if message.IsIncomeOrOutcome() {
		return ms.processIncomeOutcome(message)
	}

	if resp := ms.interpreterService.InterpretMessage(message.Message); resp != false {
		if msg, ok := resp.(string); ok {
			interpreted := &domain.Message{
				Message: msg,
			}
			if interpreted.IsIncomeOrOutcome() {
				return ms.processIncomeOutcome(interpreted)
			}
		}
	}

	message.Normalize()

	switch {
	case message.IsDailyExpense():
		log.Info("processing daily expenses message")
		return ms.newReply(ms.sheetService.GetDailyExpenses())

	case message.IsDailyBalance():
		log.Info("processing daily balance message")
		return ms.newReply(ms.sheetService.GetBalance())

	case message.IsSetAsZero():
		log.Info("processing set as zero message")
		return ms.newReply(ms.sheetService.SetDailyAsZero())

	default:
		return ms.newReply(domain.InvalidMessage + ": " + message.Message)
	}
}

func (ms *MessageService) processIncomeOutcome(msg *domain.Message) *domain.Message {
	log.Info("processing income/outcome message")
	msg.Normalize()

	return ms.newReply(ms.sheetService.ProcessAndUpdateSheet(msg.Message))
}

func (ms *MessageService) newReply(content string) *domain.Message {
	return &domain.Message{
		Message: content,
	}
}
