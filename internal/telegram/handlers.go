package telegram

import (
	"fmt"
	"log"

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
	// tgbotapi.NewMessage()
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

		fmt.Println(userID)
		// b.handleUser(*userID, chatID)
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

func (b *Bot) handleUser(userID uuid.UUID, chatID int64) (*uuid.UUID, bool) {
	// userID, err := b.repo.User.GetUserByTelegramUsername(username)
	// if err != nil {
	// 	return nil, false
	// }
	// if userID != nil && *userID != uuid.Nil {
	// 	return userID, true
	// }

	return nil, false
}
