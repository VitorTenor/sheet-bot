package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/playwright-community/playwright-go"

	"github.com/vitortenor/sheet-bot/internal/configs"
	"github.com/vitortenor/sheet-bot/internal/domain"
)

type WhatsAppCrawlerService struct {
	context        context.Context
	appConfig      *configs.ApplicationConfig
	messageService *MessageService
}

func NewWhatsAppCrawlerService(ctx context.Context, appConfig *configs.ApplicationConfig, ms *MessageService) *WhatsAppCrawlerService {
	return &WhatsAppCrawlerService{
		context:        ctx,
		appConfig:      appConfig,
		messageService: ms,
	}
}

const (
	interval = 2 * time.Second
)

func (wcs *WhatsAppCrawlerService) WhatsAppCrawler() {
	browser, err := wcs.launchBrowser()
	if err != nil {
		log.Fatalf("error launching browser: %v", err)
	}
	defer browser.Close()

	page, err := wcs.openWhatsAppPage(browser)
	if err != nil {
		log.Fatalf("error opening WhatsApp page: %v", err)
	}

	if err = wcs.openGroupChat(page); err != nil {
		log.Fatalf("error opening group chat: %v", err)
	}

	wcs.checkMessages(page)
}

func (wcs *WhatsAppCrawlerService) launchBrowser() (playwright.BrowserContext, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	browserContext, err := pw.Chromium.LaunchPersistentContext(wcs.appConfig.Crawler.UserDataDir, playwright.BrowserTypeLaunchPersistentContextOptions{
		Channel:  playwright.String("chrome"),
		Headless: playwright.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	return browserContext, nil
}

func (wcs *WhatsAppCrawlerService) openWhatsAppPage(browser playwright.BrowserContext) (playwright.Page, error) {
	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}

	_, err = page.Goto(wcs.appConfig.WhatsApp.WebURL, playwright.PageGotoOptions{
		WaitUntil: (*playwright.WaitUntilState)(playwright.String("networkidle")),
	})
	if err != nil {
		return nil, err
	}

	_, err = page.WaitForSelector(".x1qlqyl8")
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (wcs *WhatsAppCrawlerService) openGroupChat(page playwright.Page) error {
	sheetBot, err := page.QuerySelector(fmt.Sprintf(`span[title="%s"]`, wcs.appConfig.WhatsApp.GroupName))
	if err != nil {
		return err
	}

	if err := sheetBot.Click(); err != nil {
		return err
	}

	_, err = page.WaitForSelector(".x10l6tqk")
	return err
}

func (wcs *WhatsAppCrawlerService) checkMessages(page playwright.Page) {
	for {
		if err := wcs.handleMessages(page); err != nil {
			log.Printf("error handling messages: %v", err)
		}
		time.Sleep(interval)
	}
}

func (wcs *WhatsAppCrawlerService) handleMessages(page playwright.Page) error {
	messagesText, err := wcs.getMessagesText(page)
	if err != nil {
		return fmt.Errorf("error getting messages text: %w", err)
	}

	if err := wcs.processMessages(page, messagesText); err != nil {
		return fmt.Errorf("error processing messages: %w", err)
	}

	return nil
}

func (wcs *WhatsAppCrawlerService) processMessages(page playwright.Page, messagesText []string) error {
	messageTextSize := len(messagesText)
	if messageTextSize == 0 {
		return nil
	}

	if !wcs.checkIfIsSystemMessage(messagesText[messageTextSize-1]) {
		log.Info("processing messages...")
		var messagesToSave []string
		counter := 1

		for messageTextSize-counter >= 0 && !wcs.checkIfIsSystemMessage(messagesText[messageTextSize-counter]) {
			messagesToSave = append(messagesToSave, messagesText[messageTextSize-counter])
			counter++
		}

		for _, message := range messagesToSave {
			domainMessage := &domain.Message{
				Message: message,
			}

			log.Info("processing message: ", domainMessage.Message)
			response := wcs.messageService.ProcessAndReply(domainMessage)
			err := wcs.typeAndSend(page, response.Message)
			log.Info("message processed: ", response.Message)

			if err != nil {
				return err
			}
		}
		log.Info("messages processed")
	}
	return nil
}

func (wcs *WhatsAppCrawlerService) getMessagesText(page playwright.Page) ([]string, error) {
	mainDiv, err := page.QuerySelector(`div[role="application"]`)
	if err != nil {
		return nil, err
	}
	children, err := mainDiv.QuerySelectorAll(".selectable-text")
	if err != nil {
		return nil, err
	}
	var messagesText []string
	for _, child := range children {
		messageText, err := child.TextContent()
		if err != nil {
			return nil, err
		}
		messagesText = append(messagesText, messageText)
	}
	return messagesText, nil
}

func (wcs *WhatsAppCrawlerService) typeAndSend(page playwright.Page, message string) error {
	messageBox, err := page.QuerySelector(`div[aria-label="Digite uma mensagem"]`)
	if err != nil {
		return err
	}
	if err := messageBox.Type(message); err != nil {
		return err
	}
	return page.Keyboard().Press("Enter")
}

func (wcs *WhatsAppCrawlerService) checkIfIsSystemMessage(message string) bool {
	return strings.HasPrefix(message, domain.SystemMessagePrefix)
}
