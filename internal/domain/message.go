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

	dailyMessage         = "diario"
	detailedDailyMessage = "notas"
	dailyBalance         = "saldo"
	setAsZero            = "zerar"
	sysDailyReminder     = "sysdailyreminder"
	regex                = `^-?\d+(?:[.,]\d+)?\s\/\s.+$`
)

func (m *Message) CheckIfIsSystemMessage() bool {
	if m.Message == "" {
		return false
	}

	return strings.HasPrefix(m.Message, SystemMessagePrefix)
}

func (m *Message) Normalize() {
	if m.IsIncomeOrOutcome() {
		m.Message = strings.ReplaceAll(m.Message, ",", ".")
	} else {
		m.Message = strings.ToLower(m.Message)
		// in pt-br the message 'diario' can be written with an accent
		m.Message = strings.ReplaceAll(m.Message, "á", "a")
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

func (m *Message) IsSetAsZero() bool {
	if m.Message == "" {
		return false
	}

	if m.Message == setAsZero {
		return true
	}

	return false
}

func (m *Message) IsDetailedDailyBalance() bool {
	if m.Message == "" {
		return false
	}

	if m.Message == detailedDailyMessage {
		return true
	}

	return false
}

func (m *Message) IsReminderVerification() bool {
	if m.Message == "" {
		return false
	}

	if m.Message == sysDailyReminder {
		return true
	}

	return false
}
