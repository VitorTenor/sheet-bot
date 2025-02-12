package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"

	"github.com/vitortenor/sheet-bot-api/internal/api"
	"github.com/vitortenor/sheet-bot-api/internal/client"
	"github.com/vitortenor/sheet-bot-api/internal/configuration"
	"github.com/vitortenor/sheet-bot-api/internal/services"
)

func main() {
	ctx := context.Background()

	e := echo.New()
	humaApi := humaecho.New(e, huma.DefaultConfig("Sheet Bot API", "1.0.0"))

	address := fmt.Sprintf("%s:%d", "localhost", 8080)
	log.Println("Server started on " + address)

	googleSrv, err := configuration.BuildGoogleSrv(ctx)
	if err != nil {
		log.Fatal("Failed to build Google service: ", err)
	}

	gss := client.NewGoogleSheetsClient(googleSrv)
	gsc := services.NewGoogleSheetsService(gss)
	ms := services.NewMessageService(gsc)
	mh := api.NewMessageHandler(ms)

	api.InitRoutes(humaApi, mh)

	err = http.ListenAndServe(address, e)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
