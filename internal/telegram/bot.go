package telegram

import (
	"fmt"
	"friendly-bot/internal/repository"
	"log"
	"time"

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

	// fmt.Println("CHECK!!!")
	// go b.checkEventsPeriodically()

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C: // Когда таймер срабатывает
				fmt.Println("Минута!!!")
				currentTime := time.Now()
				events, err := b.repo.Event.GetEvents(currentTime)
				if err != nil {
					log.Printf("error getting events: %v", err)
					continue
				}

				b.sendEventsInfo(events)
			}
		}
	}()

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
