package main

import (
	"context"
	"log"

	"github.com/vitortenor/sheet-bot-api/internal/client"
	"github.com/vitortenor/sheet-bot-api/internal/configs"
	"github.com/vitortenor/sheet-bot-api/internal/services"
)

func main() {
	ctx := context.Background()

	appConfig, err := configs.InitConfig(ctx, "application.yaml")
	if err != nil {
		log.Fatal("failed to load configuration: ", err)
	}

	googleSrv, err := configs.BuildGoogleSrv(ctx, appConfig)
	if err != nil {
		log.Fatal("failed to build Google service: ", err)
	}

	gss := client.NewGoogleSheetsClient(googleSrv)
	gsc := services.NewGoogleSheetsService(appConfig, gss)
	ms := services.NewMessageService(ctx, gsc)

	wcs := services.NewWhatsAppCrawlerService(ctx, appConfig, ms)
	wcs.WhatsAppCrawler()

}
