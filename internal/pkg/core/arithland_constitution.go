package core

import (
	"fmt"
	"github.com/avistopia/arithland-telegram/internal/models"
	"github.com/avistopia/arithland-telegram/internal/pkg/components"
	"github.com/avistopia/arithland-telegram/internal/pkg/flows"
	"github.com/avistopia/arithland-telegram/internal/pkg/texts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) arithlandConstitutionFlow() flows.Flow {
	return flows.Flow{
		KeyboardButtonActions: map[string]components.Action{
			texts.ShowArithlandConstitution: func(user *models.User, incomingMessage *tgbotapi.Message) error {
				response, err := components.Message{Text: texts.ArithlandConstitution}.Render(incomingMessage.Chat.ID)
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
