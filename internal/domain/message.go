package domain

import (
	"regexp"
	"strings"
)

type Message struct {
	Message string
}

const (
	SystemMessagePrefix = "sys: "

	SystemErrorMessage = SystemMessagePrefix + "system error"
	InvalidMessage     = SystemMessagePrefix + "invalid message"
	ZeroBalanceMessage = SystemMessagePrefix + "R$ 0,00"

	dailyMessage = "diario"
	dailyBalance = "saldo"
	regex        = `^-?\d+\s\/\s.+(?:\s\d+)?$`
)

func (m *Message) CheckIfIsSystemMessage() bool {
	if m.Message == "" {
		return false
	}

	return strings.HasPrefix(m.Message, SystemMessagePrefix)
}

func (m *Message) Normalize() {
	if !m.IsIncomeOrOutcome() {
		m.Message = strings.ToLower(m.Message)
	}
}

func (m *Message) IsIncomeOrOutcome() bool {
	if m.Message == "" {
		return false
	}

	return regexp.MustCompile(regex).MatchString(m.Message)
}

func (m *Message) IsDailyExpense() bool {
	if m.Message == "" {
		return false
	}

	if m.Message == dailyMessage {
		return true
	}

	return false
}

func (m *Message) IsDailyBalance() bool {
	if m.Message == "" {
		return false
	}

	if m.Message == dailyBalance {
		return true
	}

	return false
}
