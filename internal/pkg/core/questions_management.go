package core

import (
	"github.com/avistopia/arithland-telegram/internal/models"
	"github.com/avistopia/arithland-telegram/internal/pkg/components"
	"github.com/avistopia/arithland-telegram/internal/pkg/flows"
	"github.com/avistopia/arithland-telegram/internal/pkg/texts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) questionsManagementFlow() flows.Flow {
	return flows.Flow{KeyboardButtonActions: map[string]components.Action{
		texts.ShowQuestionsManagement: func(user *models.User, message *tgbotapi.Message) error {
			// TODO
			return nil
		},
	},
	}
}
