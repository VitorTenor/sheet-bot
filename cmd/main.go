package main

import (
	"context"

	"github.com/labstack/gommon/log"

	"github.com/vitortenor/sheet-bot/internal/client"
	"github.com/vitortenor/sheet-bot/internal/configs"
	"github.com/vitortenor/sheet-bot/internal/services"
)

func main() {
	log.Info("starting application...")
	ctx := context.Background()

	appConfig, err := configs.InitConfig(ctx, "application.yaml")
	if err != nil {
		log.Error("failed to load configuration: ", err)
		log.Fatal("failed to load configuration: ", err)
	}

	googleSrv, err := configs.BuildGoogleSrv(ctx, appConfig)
	if err != nil {
		log.Error("failed to build Google service: ", err)
		log.Fatal("failed to build Google service: ", err)
	}

	gss := client.NewGoogleSheetsClient(googleSrv)
	gsc := services.NewGoogleSheetsService(appConfig, gss)
	ms := services.NewMessageService(ctx, gsc)

	wcs := services.NewWhatsAppCrawlerService(ctx, appConfig, ms)
	wcs.WhatsAppCrawler()

}
