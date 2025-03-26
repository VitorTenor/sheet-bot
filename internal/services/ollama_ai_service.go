package services

import (
	"encoding/json"
	"fmt"
	"strings" // Add this import

	"github.com/vitortenor/sheet-bot/internal/client"
	"github.com/vitortenor/sheet-bot/internal/configs"
)

type OllamaAIService struct {
	appConfig *configs.ApplicationConfig
	client    *client.OllamaAIClient
}

var PROMPT = "Your task is to format a text message based on whether it describes an income or expense transaction. Follow these rules carefully:\n\n1. **Expense:** If the message describes a purchase or spending action (e.g., 'comprei uma água por 5 reais'), output the amount as a negative number, followed by the item or description. \n   - Example: 'comprei uma cerveja 30 reais' → `-30 / cerveja`\n\n2. **Income:** If the message describes earning or receiving money (e.g., 'vendi um produto por 200 reais'), output the amount as a positive number, followed by the description. \n   - Example: 'vendi um produto por 200 reais' → `200 / vendi um produto`\n\n3. **Invalid message:** If the input does not clearly describe a valid transaction with an amount, or if the input is one of the following: \"diario\", \"diario-detail\", \"zerar\", or \"saldo\", return `false`. \n   - Example: 'comrpe aaa' → `false`\n\n4. **Specific conditions:** If the input message is exactly \"diario\", \"diario-detail\", \"zerar\", or \"saldo\", return `false`.\n\nProcess the following input message: **'%s'**\n\nReturn only the formatted result, without explanations or additional text."

func NewOllamaAIService(appConfig *configs.ApplicationConfig, oac *client.OllamaAIClient) *OllamaAIService {
	return &OllamaAIService{
		appConfig: appConfig,
		client:    oac,
	}
}

func (oas *OllamaAIService) GetOllamaAIResponse(message string) string {
	promptMessage := fmt.Sprintf(PROMPT, message)
	response, err := oas.client.GetOllamaAIResponse(oas.appConfig.Ai.ModelName, promptMessage)
	if err != nil {
		return fmt.Sprintf("failed to generate text: %s", err)
	}

	type responseObj struct {
		Response string `json:"response"`
	}

	var r responseObj
	err = json.NewDecoder(strings.NewReader(response)).Decode(&r) // Convert the response to a reader
	if err != nil {
		return fmt.Sprintf("failed to decode response: %s", err)
	}

	return r.Response
}
