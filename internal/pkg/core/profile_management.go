package core

import (
	"fmt"
	"github.com/avistopia/arithland-telegram/internal/models"
	"github.com/avistopia/arithland-telegram/internal/pkg/clean"
	"github.com/avistopia/arithland-telegram/internal/pkg/components"
	"github.com/avistopia/arithland-telegram/internal/pkg/flows"
	"github.com/avistopia/arithland-telegram/internal/pkg/texts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	changeDisplayNameOnClick           = "ChangeDisplayNameOnClick"
	backToShowProfileManagementOnClick = "backToShowProfileManagementOnClick"
)

func (s *Service) profileManagementFlow() flows.Flow {
	return flows.Flow{
		KeyboardButtonActions: map[string]components.Action{
			texts.ShowProfileManagement: func(user *models.User, message *tgbotapi.Message) error {
				response, err := s.profileMessage(user).Render(message.Chat.ID)
				if err != nil {
					return fmt.Errorf("failed to render message: %w", err)
				}

				if _, err = s.bot.Send(response); err != nil {
					return fmt.Errorf("failed to send message: %w", err)
				}

				return nil
			},
		},
		InlineButtonActions: map[string]components.InlineButtonAction{
			changeDisplayNameOnClick: func(
				user *models.User, query *tgbotapi.CallbackQuery, data string,
			) (string, error) {
				user.State = models.NewWaitingForUserFieldState(models.UserFieldNameDisplayName)

				err := s.userRepo.Save(user)
				if err != nil {
					return "", fmt.Errorf("failed to save user state: %w", err)
				}

				editMsg, err := s.enterDisplayNameMessage().RenderEditMessage(
					query.Message.Chat.ID, query.Message.MessageID,
				)
				if err != nil {
					return "", fmt.Errorf("failed to render edit message: %w", err)
				}

				if _, err := s.bot.Send(editMsg); err != nil {
					log.Printf("failed to edit message: %v", err)
				}

				return texts.ChangeDisplayName, nil
			},
			backToShowProfileManagementOnClick: func(
				user *models.User, query *tgbotapi.CallbackQuery, data string,
			) (string, error) {
				user.State = models.NewDefaultState()

				err := s.userRepo.Save(user)
				if err != nil {
					return "", fmt.Errorf("failed to save user state: %w", err)
				}

				editMsg, err := s.profileMessage(user).RenderEditMessage(query.Message.Chat.ID, query.Message.MessageID)
				if err != nil {
					return "", fmt.Errorf("failed to render edit message: %w", err)
				}

				if _, err := s.bot.Send(editMsg); err != nil {
					log.Printf("failed to edit message: %v", err)
				}

				return texts.Cancelled, nil
			},
		},
		MessageActions: map[models.StateName]components.Action{
			models.StateName_WaitingForUserField: func(user *models.User, message *tgbotapi.Message) error {
				displayName, validationError := clean.UserDisplayName(message.Text)

				if validationError != "" {
					response, err := components.Message{Text: validationError}.Render(message.Chat.ID)
					if err != nil {
						return fmt.Errorf("failed to render validation error message: %w", err)
					}

					if _, err = s.bot.Send(response); err != nil {
						return fmt.Errorf("failed to send validation error message: %w", err)
					}

					return nil
				}

				user.State = models.NewDefaultState()
				user.DisplayName = displayName

				err := s.userRepo.Save(user)
				if err != nil {
					return fmt.Errorf("failed to save user state: %w", err)
				}

				response, err := s.profileMessage(user).Render(message.Chat.ID)
				if err != nil {
					return fmt.Errorf("failed to render message: %w", err)
				}

				if _, err = s.bot.Send(response); err != nil {
					return fmt.Errorf("failed to send message: %w", err)
				}

				return nil
			},
		},
	}
}

func (s *Service) profileMessage(user *models.User) components.Message {
	return components.Message{
		Text: texts.Format(texts.ProfileManagement, map[string]string{
			"displayName": user.DisplayName,
			"balance":     texts.FormatFloat(user.Balance),
		}),
		InlineKeyboard: [][]components.InlineKeyboardButton{
			{
				components.NewInlineKeyboardButton(texts.ChangeDisplayName, changeDisplayNameOnClick, ""),
			},
		},
	}
}

func (s *Service) enterDisplayNameMessage() components.Message {
	return components.Message{
		Text: texts.EnterDisplayName,
		InlineKeyboard: [][]components.InlineKeyboardButton{
			{
				components.NewInlineKeyboardButton(texts.Cancel, backToShowProfileManagementOnClick, ""),
			},
		},
	}
}
