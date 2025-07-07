package core

import (
	"fmt"
	"github.com/avistopia/arithland-telegram/internal/models"
	"github.com/avistopia/arithland-telegram/internal/pkg/flows"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	bot      *tgbotapi.BotAPI
	userRepo *models.UserRepo
}

func NewService(bot *tgbotapi.BotAPI, userRepo *models.UserRepo) *Service {
	return &Service{bot: bot, userRepo: userRepo}
}

func (s *Service) Flow() (*flows.Flow, error) {
	flow, err := flows.MergeFlows([]flows.Flow{
		s.mainMenuFlow(),
		s.arithlandConstitutionFlow(),
		s.profileManagementFlow(),
		s.questionsManagementFlow(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to merge flows: %w", err)
	}

	return flow, nil
}
