package components

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message struct {
	Text           string
	InlineKeyboard [][]InlineKeyboardButton
	Keyboard       [][]KeyboardButton
}

func (m Message) Render(chatID int64) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(chatID, m.Text)

	switch {
	case m.Keyboard != nil && m.InlineKeyboard != nil:
		return nil, fmt.Errorf("cannot have both keyboard and inline keyboard defined on a message")
	case m.Keyboard != nil:
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
			ResizeKeyboard: true,
			Keyboard:       renderTable[tgbotapi.KeyboardButton](m.Keyboard),
		}
	case m.InlineKeyboard != nil:
		msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: renderTable[tgbotapi.InlineKeyboardButton](m.InlineKeyboard),
		}
	}

	return &msg, nil
}

func (m Message) RenderEditMessage(chatID int64, messageID int) (*tgbotapi.EditMessageTextConfig, error) {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, m.Text)

	switch {
	case m.Keyboard != nil:
		return nil, fmt.Errorf("cannot update keyboard in edit message")
	case m.InlineKeyboard != nil:
		msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: renderTable[tgbotapi.InlineKeyboardButton](m.InlineKeyboard),
		}
	}

	return &msg, nil
}
