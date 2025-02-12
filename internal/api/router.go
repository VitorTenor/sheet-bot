package api

import "github.com/danielgtaylor/huma/v2"

func InitRoutes(api huma.API, messageHandler *MessageHandler) {
	InitMessageRoutes(api, messageHandler)
}
