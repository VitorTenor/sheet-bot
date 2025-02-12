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

	googleSrv, err := configs.BuildGoogleSrv(ctx)
	if err != nil {
		log.Fatal("Failed to build Google service: ", err)
	}

	gss := client.NewGoogleSheetsClient(googleSrv)
	gsc := services.NewGoogleSheetsService(gss)
	ms := services.NewMessageService(ctx, gsc)

	wcs := services.NewWhatsAppCrawlerService(ctx, ms)
	wcs.WhatsAppCrawler()

}
