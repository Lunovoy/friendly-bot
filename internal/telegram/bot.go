package telegram

import (
	"friendly-bot/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	token string
}

type Bot struct {
	bot  *tgbotapi.BotAPI
	repo *repository.Repository
}

func NewBot(bot *tgbotapi.BotAPI, repo *repository.Repository) *Bot {
	return &Bot{
		bot:  bot,
		repo: repo,
	}
}

func (b *Bot) Start() error {
	b.bot.Debug = true

	updates := b.initUpdatesChannel()

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		b.handleMessage(update.Message)
	}
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	return b.bot.GetUpdatesChan(updateConfig)
}
