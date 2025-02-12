package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"

	"github.com/vitortenor/sheet-bot-api/internal/domain"
)

const (
	interval = 2 * time.Second // Message check interval
)

var (
	groupName   = "sheet-bot"
	whatsappURL = "https://web.whatsapp.com/"
	userDataDir = "./user_data"
)

func launchBrowser() (playwright.BrowserContext, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	browserContext, err := pw.Chromium.LaunchPersistentContext(userDataDir, playwright.BrowserTypeLaunchPersistentContextOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	return browserContext, nil
}

func openWhatsAppPage(browser playwright.BrowserContext) (playwright.Page, error) {
	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}
	_, err = page.Goto(whatsappURL, playwright.PageGotoOptions{
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

func openGroupChat(page playwright.Page) error {
	sheetBot, err := page.QuerySelector(fmt.Sprintf(`span[title="%s"]`, groupName))
	if err != nil {
		return err
	}
	if err := sheetBot.Click(); err != nil {
		return err
	}
	_, err = page.WaitForSelector(".x10l6tqk")
	return err
}

func processMessages(page playwright.Page, messagesText []string, ms *MessageService) error {
	messageTextSize := len(messagesText)
	if messageTextSize == 0 {
		return nil
	}

	if !checkIfMessageStartsWithIgnoredValues(messagesText[messageTextSize-1]) {
		log.Println("---------- Processing messages ----------")
		var messagesToSave []string
		counter := 1

		for messageTextSize-counter >= 0 && !checkIfMessageStartsWithIgnoredValues(messagesText[messageTextSize-counter]) {
			messagesToSave = append(messagesToSave, messagesText[messageTextSize-counter])
			counter++
		}

		for _, message := range messagesToSave {
			// http request to save message
			domainMessage := &domain.Message{
				Message: message,
			}
			response := ms.ProcessAndReply(nil, domainMessage)
			log.Printf("Message: %s", message)
			err := typeAndSend(page, response.Message)
			if err != nil {
				return err
			}

		}
		log.Println("----------------------------------------")
	}
	return nil
}

func checkMessages(page playwright.Page, ms *MessageService) {
	for {
		messagesText, err := getMessagesText(page)
		if err != nil {
			log.Printf("Error processing messages: %v", err)
		} else {
			if err := processMessages(page, messagesText, ms); err != nil {
				log.Printf("Error processing messages: %v", err)
			}
		}
		time.Sleep(interval)
	}
}

func WhatsAppCrawler(ms *MessageService) {
	browser, err := launchBrowser()
	if err != nil {
		log.Fatalf("Error launching browser: %v", err)
	}
	defer browser.Close()

	page, err := openWhatsAppPage(browser)
	if err != nil {
		log.Fatalf("Error opening WhatsApp page: %v", err)
	}

	if err := openGroupChat(page); err != nil {
		log.Fatalf("Error opening group chat: %v", err)
	}

	log.Println("----------------------------------------")
	checkMessages(page, ms)
}

func getMessagesText(page playwright.Page) ([]string, error) {
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

func typeAndSend(page playwright.Page, message string) error {
	messageBox, err := page.QuerySelector(`div[aria-label="Digite uma mensagem"]`)
	if err != nil {
		return err
	}
	if err := messageBox.Type(message); err != nil {
		return err
	}
	return page.Keyboard().Press("Enter")
}

func checkIfMessageStartsWithIgnoredValues(message string) bool {
	return strings.HasPrefix(message, "sys:")
}
