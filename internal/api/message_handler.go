package api

import (
	"context"
	"net/http"

	"github.com/vitortenor/sheet-bot-api/internal/domain"
	"github.com/vitortenor/sheet-bot-api/internal/services"

	"github.com/danielgtaylor/huma/v2"
)

func InitMessageRoutes(humaApi huma.API, messageHandler *MessageHandler) {
	huma.Register(humaApi, huma.Operation{
		Path:          "/message",
		OperationID:   "handle-message",
		Method:        http.MethodPost,
		DefaultStatus: http.StatusOK,
		Summary:       "Handle a new message",
		Description:   "Handle a new message from the user",
	}, messageHandler.HandleMessage)
}

type MessageHandler struct {
	service *services.MessageService
}

func NewMessageHandler(service *services.MessageService) *MessageHandler {
	return &MessageHandler{
		service: service,
	}
}

type MessageRequest struct {
	Body struct {
		Message string `json:"message" required:"true" description:"The message to be handled"`
	}
}

type MessageResponse struct {
	Body struct {
		Message string `json:"message" description:"The message of the response"`
	}
}

func (mh *MessageHandler) HandleMessage(ctx context.Context, mr *MessageRequest) (*MessageResponse, error) {
	message := mh.service.ProcessAndReply(ctx, mr.toDomain())

	return domainToResponse(message), nil
}

func (mr *MessageRequest) toDomain() *domain.Message {
	return &domain.Message{
		Message: mr.Body.Message,
	}
}

func domainToResponse(message *domain.Message) *MessageResponse {
	return &MessageResponse{
		Body: struct {
			Message string `json:"message" description:"The message of the response"`
		}{
			Message: message.Message,
		},
	}
}
