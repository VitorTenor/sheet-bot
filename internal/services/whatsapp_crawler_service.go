package services

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"

	"github.com/vitortenor/sheet-bot/internal/configuration"
	"github.com/vitortenor/sheet-bot/internal/domain"
)

type WhatsAppCrawlerService struct {
	context        context.Context
	appConfig      *configuration.ApplicationConfig
	messageService *MessageService
}

func NewWhatsAppCrawlerService(ctx context.Context, appConfig *configuration.ApplicationConfig, ms *MessageService) *WhatsAppCrawlerService {
	return &WhatsAppCrawlerService{
		context:        ctx,
		appConfig:      appConfig,
		messageService: ms,
	}
}

const interval = 2 * time.Second

var playwrightTimeout = playwright.Float(3600000) // 1 hour timeout

func (wcs *WhatsAppCrawlerService) WhatsAppCrawler() {
	// Launch browser
	browser, err := wcs.launchBrowser()
	if err != nil {
		log.Fatalf("error launching browser: %v", err)
	}
	defer browser.Close()

	// Open WhatsApp Web page
	page, err := wcs.openWhatsAppPage(browser)
	if err != nil {
		log.Fatalf("error opening WhatsApp page: %v", err)
	}

	// Open archived chats if configured
	if wcs.appConfig.WhatsApp.IsArchived {
		if err = wcs.openArchivedChats(page); err != nil {
			log.Fatalf("error opening archived chats: %v", err)
		}
	}

	// Open group chat
	if err = wcs.openGroupChat(page); err != nil {
		log.Fatalf("error opening group chat: %v", err)
	}

	log.Println("WhatsApp crawler started successfully")

	// Choose method based on headless config
	if wcs.appConfig.Crawler.Headless {
		wcs.checkMessagesHeadless(page)
	} else {
		wcs.checkMessagesNonHeadless(page)
	}
}

func (wcs *WhatsAppCrawlerService) launchBrowser() (playwright.BrowserContext, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	browserContext, err := pw.Chromium.LaunchPersistentContext(
		wcs.appConfig.Crawler.UserDataDir,
		playwright.BrowserTypeLaunchPersistentContextOptions{
			Channel:  playwright.String("chrome"),
			Headless: playwright.Bool(wcs.appConfig.Crawler.Headless),
			Timeout:  playwrightTimeout,
		},
	)
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

	return page, nil
}

func (wcs *WhatsAppCrawlerService) openArchivedChats(page playwright.Page) error {
	// Wait for "Archived" button and click it
	_, err := page.WaitForSelector("text='Arquivadas'", playwright.PageWaitForSelectorOptions{Timeout: playwrightTimeout})
	if err != nil {
		return err
	}

	archivedButton, err := page.QuerySelector("text='Arquivadas'")
	if err != nil {
		return err
	}

	return archivedButton.Click()
}

func (wcs *WhatsAppCrawlerService) openGroupChat(page playwright.Page) error {
	// Wait for group name and click
	_, err := page.WaitForSelector("text='"+wcs.appConfig.WhatsApp.GroupName+"'", playwright.PageWaitForSelectorOptions{Timeout: playwrightTimeout})
	if err != nil {
		return err
	}

	chat, err := page.QuerySelector(`span[title="` + wcs.appConfig.WhatsApp.GroupName + `"]`)
	if err != nil {
		return err
	}

	if err := chat.Click(); err != nil {
		return err
	}

	// Wait for messages container
	_, err = page.WaitForSelector(".x10l6tqk", playwright.PageWaitForSelectorOptions{Timeout: playwrightTimeout})
	return err
}

// ---------- HEADLESS: real-time message capture ----------
func (wcs *WhatsAppCrawlerService) checkMessagesHeadless(page playwright.Page) {
	// Set up MutationObserver to track only the last message
	_, err := page.Evaluate(`
        window.lastMessage = null;
        const targetNode = document.querySelector('#main');
        const config = { childList: true, subtree: true };
        const callback = (mutationsList) => {
            mutationsList.forEach((mutation) => {
                mutation.addedNodes.forEach((node) => {
                    if(!node.querySelector) return;
                    const span = node.querySelector('.selectable-text.copyable-text span');
                    if(span){
                        window.lastMessage = span.textContent;
                    }
                });
            });
        };
        const observer = new MutationObserver(callback);
        observer.observe(targetNode, config);
    `)
	if err != nil {
		log.Fatalf("failed to set up mutation observer: %v", err)
	}

	// Infinite loop to process new messages
	for {
		lastMsgRaw, err := page.Evaluate(`window.lastMessage`)
		if err != nil {
			log.Printf("error evaluating last message: %v", err)
			time.Sleep(interval)
			continue
		}

		if lastMsg, ok := lastMsgRaw.(string); ok && lastMsg != "" {
			// Skip messages starting with "sys:"
			if !strings.HasPrefix(lastMsg, "sys:") {
				domainMsg := &domain.Message{Message: lastMsg}
				response := wcs.messageService.ProcessAndReply(domainMsg)
				if err := wcs.typeAndSend(page, response.Message); err != nil {
					log.Printf("error sending message: %v", err)
				}
			}
			// Reset lastMessage to avoid reprocessing
			_, _ = page.Evaluate(`window.lastMessage = null`)
		}

		time.Sleep(interval)
	}
}

// ---------- NON-HEADLESS: polling the last message ----------
func (wcs *WhatsAppCrawlerService) checkMessagesNonHeadless(page playwright.Page) {
	for {
		// Get the last message from the chat
		lastMsg, err := page.Evaluate(`(() => {
			const spans = document.querySelectorAll('#main .selectable-text.copyable-text span');
			if(spans.length === 0) return '';
			return spans[spans.length -1].textContent;
		})()`)
		if err != nil {
			log.Printf("error getting last message: %v", err)
			time.Sleep(interval)
			continue
		}

		// Skip messages starting with "sys:"
		if msgStr, ok := lastMsg.(string); ok && msgStr != "" && !strings.HasPrefix(msgStr, "sys:") {
			domainMsg := &domain.Message{Message: msgStr}
			response := wcs.messageService.ProcessAndReply(domainMsg)
			if err := wcs.typeAndSend(page, response.Message); err != nil {
				log.Printf("error sending message: %v", err)
			}
		}

		time.Sleep(interval)
	}
}

func (wcs *WhatsAppCrawlerService) typeAndSend(page playwright.Page, message string) error {
	msgBox, err := page.QuerySelector(`div[aria-label="Digite uma mensagem"]`)
	if err != nil {
		return err
	}
	if err := msgBox.Type(message); err != nil {
		return err
	}
	return page.Keyboard().Press("Enter")
}
