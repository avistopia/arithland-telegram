package components

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type KeyboardButton = Renderer[tgbotapi.KeyboardButton]

type keyboardButton struct {
	Text string
}

func NewKeyboardButton(text string) KeyboardButton {
	return &keyboardButton{Text: text}
}

func (b *keyboardButton) Render() tgbotapi.KeyboardButton {
	return tgbotapi.KeyboardButton{Text: b.Text}
}

type InlineKeyboardButton = Renderer[tgbotapi.InlineKeyboardButton]

type inlineKeyboardButton struct {
	Text string

	ActionName string

	// Data can be maximum of 64 bytes
	Data string
}

func NewInlineKeyboardButton(text string, actionName string, data string) InlineKeyboardButton {
	return &inlineKeyboardButton{
		Text:       text,
		ActionName: actionName,
		Data:       data,
	}
}

func (b *inlineKeyboardButton) Render() tgbotapi.InlineKeyboardButton {
	data := fmt.Sprintf("%s:%s", b.ActionName, b.Data)

	return tgbotapi.InlineKeyboardButton{
		Text:         b.Text,
		CallbackData: &data,
	}
}
