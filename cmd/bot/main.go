package main

import (
	"github.com/avistopia/arithland-telegram/internal/models"
	"github.com/avistopia/arithland-telegram/internal/pkg/core"
	"github.com/avistopia/arithland-telegram/internal/pkg/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Panic("failed to load config")
	}

	var connection gorm.Dialector

	switch config.Database.Type {
	case "sqlite":
		connection = sqlite.Open(config.Database.DSN)
	case "postgres":
		connection = postgres.Open(config.Database.DSN)
	}

	db, err := gorm.Open(connection, &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Panic("failed to connect to database")
	}

	userRepo, err := models.NewUserRepo(db)
	if err != nil {
		logrus.WithError(err).Panic("failed to initialize user repo")
	}

	bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		logrus.WithError(err).Panic("failed to initialize telegram bot")
	}

	flow, err := core.NewService(bot, userRepo).Flow()
	if err != nil {
		logrus.WithError(err).Panic("failed to get the core service flows")
	}

	handler.NewHandler(userRepo, bot, flow).Listen()
}
