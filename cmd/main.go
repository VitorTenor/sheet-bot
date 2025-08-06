package main

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/playwright-community/playwright-go"

	"github.com/vitortenor/sheet-bot/internal/client"
	"github.com/vitortenor/sheet-bot/internal/configuration"
	"github.com/vitortenor/sheet-bot/internal/services"
)

func main() {
	log.Info("starting application...")
	ctx := context.Background()

	err := playwright.Install()
	if err != nil {
		log.Error("failed to install playwright: ", err)
	}

	appConfig, err := configuration.InitConfig(ctx, "application.yaml")
	if err != nil {
		log.Fatal("failed to load configuration: ", err)
	}

	googleSrv, err := configuration.BuildGoogleSrv(ctx, appConfig)
	if err != nil {
		log.Fatal("failed to build Google service: ", err)
	}

	oac := client.NewOllamaAIClient(appConfig.Ai.ModelURL)
	oas := services.NewOllamaAIService(appConfig, oac)

	gss := client.NewGoogleSheetsClient(googleSrv)
	gsc := services.NewGoogleSheetsService(appConfig, gss)
	ms := services.NewMessageService(ctx, appConfig, gsc, oas)

	wcs := services.NewWhatsAppCrawlerService(ctx, appConfig, ms)
	wcs.WhatsAppCrawler()

}
