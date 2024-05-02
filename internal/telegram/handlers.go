package telegram

import (
	"fmt"
	"friendly-bot/internal/models"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

const (
	commandStart   = "start"
	welcomeText    = "Поздравляю, вы подключились к боту. Теперь я буду вам отправлять уведомления по различным событиям из приложения"
	unknownUser    = "Вы не указали имя telegram аккаунта в приложении либо имя указано с ошибкой!"
	unknownCommand = "Такой команды не существует"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	chatID := message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, unknownCommand)

	switch message.Command() {
	case commandStart:

		fmt.Println(message.From.UserName)
		userID, ok := b.checkUser(message.From.UserName)
		if !ok {
			msg.Text = unknownUser
		} else {
			msg.Text = welcomeText
		}
		_, err := b.bot.Send(msg)
		if err != nil {
			return err
		}

		fmt.Println("UserID: ", userID)
		if userID != nil {
			tgChat, err := b.handleUser(userID, chatID)
			if err != nil {
				return err
			}
			fmt.Println(*tgChat)
		}
		return err
	default:
		_, err := b.bot.Send(msg)
		return err

	}

}

func (b *Bot) handleMessage(message *tgbotapi.Message) {

	chatID := message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, "Функционал ответа на сообщения отсутствует")

	if _, err := b.bot.Send(msg); err != nil {
		log.Fatalf("error sending message: %s", err.Error())
	}
}

func (b *Bot) checkUser(username string) (*uuid.UUID, bool) {
	userID, err := b.repo.User.GetUserByTelegramUsername(username)
	if err != nil {
		return nil, false
	}
	if userID != nil {
		return userID, true
	}

	return nil, false
}

func (b *Bot) checkChatAlreadyExists(username string) (*uuid.UUID, bool) {
	userID, err := b.repo.User.GetUserByTelegramUsername(username)
	if err != nil {
		return nil, false
	}
	if userID != nil {
		return userID, true
	}

	return nil, false
}

func (b *Bot) handleUser(userID *uuid.UUID, chatID int64) (*models.TgChat, error) {

	fmt.Println("Enter handle user:", *userID)
	tgChat, err := b.repo.GetTgChatByUserID(*userID)
	fmt.Println("tgChat: ", tgChat)
	if tgChat == nil {
		fmt.Println("tgChat if")
		newTgChatID, err := b.repo.TgChat.CreateTgChat(chatID, *userID)
		if err != nil {
			return nil, err
		}
		tgChat, err := b.repo.GetTgChatByID(newTgChatID)
		return tgChat, err
	}

	fmt.Println("tgChat get")

	if tgChat.ChatID != chatID {
		err := b.repo.TgChat.UpdateTgChat(chatID, *userID)
		if err != nil {
			return nil, err
		}
	}

	return tgChat, err
}

func (b *Bot) checkEventsPeriodically() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

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
}

func (b *Bot) sendEventsInfo(events []*models.EventWithFriends) {
	for _, event := range events {
		// Формируем сообщение с информацией о событии
		message := fmt.Sprintf("Событие: %s\nНачало: %v\nОкончание: %v", event.Event.Title, event.Event.StartDate, event.Event.EndDate)

		// Отправляем сообщение в Telegram каждому пользователю
		tgChat, err := b.repo.GetTgChatByUserID(event.Event.UserID)
		if err != nil {
			log.Printf("error getting chat: %s\n", err.Error())
			continue
		}
		// Для примера, отправим сообщение в чат, из которого была получена команда /start
		_, err = b.bot.Send(tgbotapi.NewMessage(tgChat.ChatID, message))
		if err != nil {
			log.Printf("error sending message: %s;\nTo User: %v\n", err.Error(), tgChat.UserID)
			continue
		}

		// После отправки сообщения, обновляем статус события, чтобы не отправлять уведомления повторно
		err = b.repo.Event.UpdateStartEventStatus(event.Event.ID, event.Event.UserID)
		if err != nil {
			log.Printf("error updating event status: %s", err.Error())
			continue
		}
	}
}
