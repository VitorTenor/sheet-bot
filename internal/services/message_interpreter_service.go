package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	invalids = map[string]bool{
		"diario":        true,
		"diario-detail": true,
		"zerar":         true,
		"saldo":         true,
	}
	expenseKeywords = [][]string{
		{"comprei"},
		{"paguei"},
		{"fiz", "um", "pix"},
		{"transferi"},
		{"gastei"},
	}
	incomeKeywords = [][]string{
		{"recebi"},
		{"vendi"},
		{"ganhei"},
		{"fiz", "um", "depósito"},
		{"pix", "recebido"},
	}
	preps = map[string]bool{"de": true, "por": true, "para": true, "pra": true}
	units = map[string]bool{"reais": true, "rs": true, "r$": true}
)

type MessageInterpreterService struct {
}

func NewMessageInterpreterService() *MessageInterpreterService {
	return &MessageInterpreterService{}
}

func (mis *MessageInterpreterService) InterpretMessage(message string) interface{} {
	msgLower := strings.ToLower(strings.TrimSpace(message))

	if invalids[msgLower] {
		return false
	}

	tokens := strings.Fields(msgLower)

	containsSequence := func(seq []string) bool {
		for i := 0; i <= len(tokens)-len(seq); i++ {
			match := true
			for j, word := range seq {
				if tokens[i+j] != word {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
		return false
	}

	isExpense := false
	isIncome := false
	for _, seq := range expenseKeywords {
		if containsSequence(seq) {
			isExpense = true
			break
		}
	}
	for _, seq := range incomeKeywords {
		if containsSequence(seq) {
			isIncome = true
			break
		}
	}
	if !isExpense && !isIncome {
		return false
	}

	re := regexp.MustCompile(`(\d+[.,]?\d*)\s*(reais|rs|r\$)?`)
	matches := re.FindStringSubmatch(msgLower)
	if len(matches) < 2 {
		return false
	}

	valStr := strings.ReplaceAll(matches[1], ",", ".")
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return false
	}

	valIndex := -1
	for i, tok := range tokens {
		if strings.Contains(tok, matches[1]) {
			valIndex = i
			break
		}
	}
	if valIndex == -1 {
		return false
	}

	start := valIndex
	if valIndex > 0 && preps[tokens[valIndex-1]] {
		start = valIndex - 1
	}
	end := valIndex + 1

	if end < len(tokens) && units[tokens[end]] {
		end++
	}

	tokens = append(tokens[:start], tokens[end:]...)

	description := strings.Join(tokens, " ")

	if isExpense {
		return fmt.Sprintf("-%v / %s", val, description)
	}
	return fmt.Sprintf("%v / %s", val, description)
}
