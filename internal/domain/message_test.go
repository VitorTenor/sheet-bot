package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage_CheckIfIsSystemMessage(t *testing.T) {

	_ = t.Run("empty message", func(t *testing.T) {
		// arrange
		message := Message{}

		// act
		isSystemMessage := message.CheckIfIsSystemMessage()

		// assert
		_ = assert.False(t, isSystemMessage)
	})

	_ = t.Run("invalid message", func(t *testing.T) {
		// arrange
		message := Message{
			Message: "-5,00 / invalid message",
		}

		// act
		isSystemMessage := message.CheckIfIsSystemMessage()

		// assert
		_ = assert.False(t, isSystemMessage)
	})

	_ = t.Run("valid message", func(t *testing.T) {
		// arrange
		message := Message{
			Message: SystemMessagePrefix + "valid message",
		}

		// act
		isSystemMessage := message.CheckIfIsSystemMessage()

		// assert
		_ = assert.True(t, isSystemMessage)
	})
}

func TestMessage_Normalize(t *testing.T) {

	_ = t.Run("'diario' with accent", func(t *testing.T) {
		// arrange
		message := Message{
			Message: "Diário",
		}

		// act
		message.Normalize()

		// assert
		_ = assert.Equal(t, dailyMessage, message.Message)
	})

	_ = t.Run("replace comma with dot", func(t *testing.T) {
		// arrange
		message := Message{
			Message: "100,00 / sold something",
		}

		// act
		message.Normalize()

		// assert
		_ = assert.Equal(t, "100.00 / sold something", message.Message)
	})
}

func TestMessage_IsIncomeOrOutcome(t *testing.T) {

	_ = t.Run("empty message", func(t *testing.T) {
		// arrange
		message := Message{}

		// act
		isIncomeOrOutcome := message.IsIncomeOrOutcome()

		// assert
		_ = assert.False(t, isIncomeOrOutcome)
	})

	_ = t.Run("invalid message", func(t *testing.T) {
		// arrange
		message := Message{
			Message: "invalid message",
		}

		// act
		isIncomeOrOutcome := message.IsIncomeOrOutcome()

		// assert
		_ = assert.False(t, isIncomeOrOutcome)
	})

	_ = t.Run("valid message", func(t *testing.T) {
		// arrange
		message := Message{
			Message: "100.00 / sold something",
		}

		// act
		isIncomeOrOutcome := message.IsIncomeOrOutcome()

		// assert
		_ = assert.True(t, isIncomeOrOutcome)
	})
}

func TestMessage_IsDailyExpense(t *testing.T) {

	_ = t.Run("empty message", func(t *testing.T) {
		// arrange
		message := Message{}

		// act
		isDailyExpense := message.IsDailyExpense()

		// assert
		_ = assert.False(t, isDailyExpense)
	})

	_ = t.Run("invalid message", func(t *testing.T) {
		// arrange
		message := Message{
			Message: "invalid message",
		}

		// act
		isDailyExpense := message.IsDailyExpense()

		// assert
		_ = assert.False(t, isDailyExpense)
	})

	_ = t.Run("valid message", func(t *testing.T) {
		// arrange
		message := Message{
			Message: dailyMessage,
		}

		// act
		isDailyExpense := message.IsDailyExpense()

		// assert
		_ = assert.True(t, isDailyExpense)
	})
}

func TestMessage_IsDailyBalance(t *testing.T) {

	_ = t.Run("empty message", func(t *testing.T) {
		// arrange
		message := Message{}

		// act
		isDailyBalance := message.IsDailyBalance()

		// assert
		_ = assert.False(t, isDailyBalance)
	})

	_ = t.Run("invalid message", func(t *testing.T) {
		// arrange
		message := Message{
			Message: "invalid message",
		}

		// act
		isDailyBalance := message.IsDailyBalance()

		// assert
		_ = assert.False(t, isDailyBalance)
	})

	_ = t.Run("valid message", func(t *testing.T) {
		// arrange
		message := Message{
			Message: dailyBalance,
		}

		// act
		isDailyBalance := message.IsDailyBalance()

		// assert
		_ = assert.True(t, isDailyBalance)
	})
}

func TestMessage_IsSetAsZero(t *testing.T) {

	_ = t.Run("empty message", func(t *testing.T) {
		// arrange
		message := Message{}

		// act
		isSetAsZero := message.IsSetAsZero()

		// assert
		_ = assert.False(t, isSetAsZero)
	})

	_ = t.Run("invalid message", func(t *testing.T) {
		// arrange
		message := Message{
			Message: "invalid message",
		}

		// act
		isSetAsZero := message.IsSetAsZero()

		// assert
		_ = assert.False(t, isSetAsZero)
	})

	_ = t.Run("valid message", func(t *testing.T) {
		// arrange
		message := Message{
			Message: setAsZero,
		}

		// act
		isSetAsZero := message.IsSetAsZero()

		// assert
		_ = assert.True(t, isSetAsZero)
	})
}
