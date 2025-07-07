package handler

import (
	"encoding/json"
	"fmt"
	"github.com/avistopia/arithland-telegram/internal/pkg/components"
	"github.com/avistopia/arithland-telegram/internal/pkg/flows"
	"github.com/avistopia/arithland-telegram/internal/pkg/texts"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"

	"github.com/avistopia/arithland-telegram/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	userRepo *models.UserRepo
	bot      *tgbotapi.BotAPI
	flow     *flows.Flow
}

func NewHandler(userRepo *models.UserRepo, bot *tgbotapi.BotAPI, flow *flows.Flow) *Handler {
	return &Handler{userRepo: userRepo, bot: bot, flow: flow}
}

func (h *Handler) Listen() {
	logrus.Info("Started listening...")

	updates := h.bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 5})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		logrus.Info("Shutting down gracefully...")
		h.bot.StopReceivingUpdates()
	}()

	for update := range updates {
		logrus.WithField("update", tryToMap(update)).Info("Received an update. Started handling...")
		h.handle(update)
		logrus.Info("Handled the update.")
	}
}

func (h *Handler) handle(update tgbotapi.Update) {
	switch {
	case update.CallbackQuery != nil:
		response, err := h.handleCallbackQuery(update.CallbackQuery)
		if err != nil {
			logrus.WithError(err).Error("failed to handle callback query")

			if _, sendErr := h.bot.Send(&tgbotapi.CallbackConfig{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            texts.InternalError,
			}); sendErr != nil {
				logrus.WithError(sendErr).Error("failed to send error of handle callback query to the user")
			}

			return
		}

		if _, err := h.bot.Send(&tgbotapi.CallbackConfig{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            response,
		}); err != nil {
			logrus.WithError(err).Error("failed to send response of handle callback query to the user")
		}
	case update.Message != nil:
		err := h.handleMessage(update.Message)
		if err != nil {
			logrus.WithError(err).Error("failed to handle message")

			errResponse, err := components.Message{Text: texts.InternalError}.Render(update.Message.Chat.ID)
			if err != nil {
				logrus.WithError(err).Error("failed to send error response of handle message to the user")
			}

			if _, sendErr := h.bot.Send(errResponse); sendErr != nil {
				logrus.WithError(sendErr).Error("failed to send error response of handle message to the user")
			}
		}
	default:
		logrus.Info("Update type not supported, ignoring.")
	}
}

func (h *Handler) handleCallbackQuery(query *tgbotapi.CallbackQuery) (string, error) {
	user, err := h.getOrCreateUserByTelegramUserID(query.From.ID)
	if err != nil {
		return "", fmt.Errorf("failed to get or create user: %w", err)
	}

	values := strings.Split(query.Data, ":")
	if len(values) < 2 {
		return "", fmt.Errorf("received data should be of format 'actionName:data'")
	}

	actionKey, data := values[0], values[1]

	action, ok := h.flow.InlineButtonActions[actionKey]
	if !ok {
		return "", fmt.Errorf("action %q not found", actionKey)
	}

	response, err := action(user, query, data)
	if err != nil {
		return "", fmt.Errorf("failed to run action for %q: %w", actionKey, err)
	}

	return response, nil
}

func (h *Handler) handleMessage(message *tgbotapi.Message) error {
	if !message.Chat.IsPrivate() {
		return fmt.Errorf("cannot handle update, the chat is not private")
	}

	user, err := h.getOrCreateUserByTelegramUserID(message.From.ID)
	if err != nil {
		return fmt.Errorf("failed to get or create user: %w", err)
	}

	if message.IsCommand() {
		if err := h.handleCommand(user, message); err != nil {
			return fmt.Errorf("failed to handle command: %w", err)
		}

		return nil
	}

	if h.isKeyboardButton(message) {
		if err := h.handleKeyboardButton(user, message); err != nil {
			return fmt.Errorf("failed to handle keyboard button: %w", err)
		}

		return nil
	}

	if err := h.handlePlainMessage(user, message); err != nil {
		return fmt.Errorf("failed to handle plain message: %w", err)
	}

	return nil
}

func (h *Handler) handleCommand(user *models.User, message *tgbotapi.Message) error {
	actionKey := message.Command()

	action, ok := h.flow.CommandActions[actionKey]
	if !ok {
		return fmt.Errorf("action %q not found", actionKey)
	}

	err := action(user, message)
	if err != nil {
		return fmt.Errorf("failed to run action for %q: %w", actionKey, err)
	}

	return nil
}

func (h *Handler) isKeyboardButton(message *tgbotapi.Message) bool {
	_, ok := h.flow.KeyboardButtonActions[message.Text]
	return ok
}

func (h *Handler) handleKeyboardButton(user *models.User, message *tgbotapi.Message) error {
	actionKey := message.Text

	action, ok := h.flow.KeyboardButtonActions[actionKey]
	if !ok {
		return fmt.Errorf("action %q not found", actionKey)
	}

	err := action(user, message)
	if err != nil {
		return fmt.Errorf("failed to run action for %q: %w", actionKey, err)
	}

	return nil
}

func (h *Handler) handlePlainMessage(user *models.User, message *tgbotapi.Message) error {
	actionKey := user.State.Name

	action, ok := h.flow.MessageActions[actionKey]
	if !ok {
		return fmt.Errorf("action %q not found", actionKey)
	}

	err := action(user, message)
	if err != nil {
		return fmt.Errorf("failed to run action for %q: %w", actionKey, err)
	}

	return nil
}

func (h *Handler) getOrCreateUserByTelegramUserID(telegramUserID int64) (*models.User, error) {
	user, err := h.userRepo.GetOrCreateUserByTelegramUserID(telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create user by telegram user ID")
	}

	return user, nil
}

func tryToMap(v any) map[string]any {
	marshalled, err := json.Marshal(v)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("failed to marshal to json: %s", err.Error()),
		}
	}

	var result map[string]any

	err = json.Unmarshal(marshalled, &result)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("failed to unmarshal from json: %s", err.Error()),
		}
	}

	return result
}
