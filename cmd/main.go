package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/labstack/gommon/log"
	"github.com/playwright-community/playwright-go"

	"github.com/vitortenor/sheet-bot/internal/client"
	"github.com/vitortenor/sheet-bot/internal/configuration"
	"github.com/vitortenor/sheet-bot/internal/domain"
	"github.com/vitortenor/sheet-bot/internal/services"
)

func main() {
	ctx := context.Background()

	appConfig, err := configuration.InitConfig(ctx, "application.yaml")
	if err != nil {
		log.Fatal("failed to load configuration: ", err)
	}

	log.SetOutput(&configuration.LogInterceptor{Writer: os.Stdout, AppConfig: appConfig})
	log.Info("starting application...")

	filePath := appConfig.UserDataFile
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("failed to read user configuration file: ", err)
	}

	var whatsappUsers *[]domain.WhatsappUser
	if err := json.Unmarshal(data, &whatsappUsers); err != nil {
		log.Fatal("failed to unmarshal user configuration file: ", err)
	}

	err = playwright.Install()
	if err != nil {
		log.Error("failed to install playwright: ", err)
	}

	googleSrv, err := configuration.BuildGoogleSrv(ctx, appConfig)
	if err != nil {
		log.Fatal("failed to build Google service: ", err)
	}

	oac := client.NewOllamaAIClient(appConfig.Ai.ModelURL)
	oas := services.NewOllamaAIService(appConfig, oac)

	gss := client.NewGoogleSheetsClient(googleSrv)
	gsc := services.NewGoogleSheetsService(appConfig, gss)
	mis := services.NewMessageInterpreterService()
	ms := services.NewMessageService(ctx, appConfig, gsc, oas, mis)

	wcs := services.NewWhatsAppCrawlerService(ctx, appConfig, ms, whatsappUsers)

	wcs.WhatsAppCrawler()
}
