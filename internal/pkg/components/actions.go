package components

import (
	"github.com/avistopia/arithland-telegram/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Action func(user *models.User, message *tgbotapi.Message) error

type InlineButtonAction func(user *models.User, query *tgbotapi.CallbackQuery, data string) (string, error)
