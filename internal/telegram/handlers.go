package telegram

import (
	"fmt"
	"friendly-bot/internal/models"
	"log"
	"sort"
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

type Frequency struct {
	Once        string
	Everyday    string
	Weekdays    string
	Weekly      string
	MonthlyDate string
	MonthlyDay  string
	Annually    string
}

var frequency = Frequency{
	Once:        "once",
	Everyday:    "everyday",
	Weekdays:    "weekdays",
	Weekly:      "weekly",
	MonthlyDate: "monthlyDate",
	MonthlyDay:  "monthlyDay",
	Annually:    "annually",
}

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

func (b *Bot) sendEventsInfo(events []*models.EventWithFriendsAndReminders) {
	for _, event := range events {
		sortedReminders := sortRemindersByMinutesUntilEvent(event.Reminders)
		for _, reminder := range sortedReminders {
			if reminder.IsActive {
				b.repo.Event.UpdateReminderStatus(reminder.ID)
				break
			}
		}
		var friendsMessage string
		for _, event := range events {
			for _, friend := range event.Friends {
				friendsMessage += fmt.Sprintf("- Имя: %s %s; Способ связи: %s %s\n", friend.FirstName, friend.LastName, friend.CommunicationMethod, friend.Messenger)
			}
		}
		// Формируем сообщение с информацией о событии
		message := fmt.Sprintf("Событие: %s\n Описание: %s\n Начало: %v\n Окончание: %v\n Участники:\n %s\n", event.Event.Title, event.Event.Description, event.Event.StartDate.Time.Format(time.RFC1123), event.Event.EndDate.Time.Format(time.RFC1123), friendsMessage)

		// Отправляем сообщение в Telegram каждому пользователю
		tgChat, err := b.repo.GetTgChatByUserID(event.Event.UserID)
		if err != nil {
			log.Printf("error getting chat: %s\n", err.Error())
			continue
		}

		// Отправляем сообщение в чат
		_, err = b.bot.Send(tgbotapi.NewMessage(tgChat.ChatID, message))
		if err != nil {
			log.Printf("error sending message: %s;\nTo User: %v\n", err.Error(), tgChat.UserID)
			continue
		}

		switch *event.Event.Frequency {
		// Если повторение однократное, то после отправки сообщения,
		// обновляем статус события, чтобы не отправлять уведомления повторно
		case frequency.Once:
			err := b.repo.Event.UpdateActiveStatus(event.Event.ID, event.Event.UserID)
			if err != nil {
				log.Printf("error updating event status: %s", err.Error())
				continue
			}

		case frequency.Everyday:
			startDate := event.Event.StartDate.Time.AddDate(0, 0, 1)
			endDate := event.Event.EndDate.Time.AddDate(0, 0, 1)
			err := b.repo.Event.UpdateStartAndEndDate(event.Event.ID, event.Event.UserID, startDate, endDate)
			if err != nil {
				log.Printf("error updating event status: %s", err.Error())
				continue
			}
		case frequency.Weekdays:
			// Вычисляем следующий день
			nextDay := event.Event.StartDate.Time.AddDate(0, 0, 1)
			// Находим день недели
			nextWeekday := nextDay.Weekday()
			var daysToAdd int
			if nextWeekday == time.Saturday {
				daysToAdd = 3
			} else if nextWeekday == time.Sunday {
				daysToAdd = 2
			} else {
				daysToAdd = 1
			}
			startDate := event.Event.StartDate.Time.AddDate(0, 0, daysToAdd)
			endDate := event.Event.EndDate.Time.AddDate(0, 0, daysToAdd)
			err := b.repo.Event.UpdateStartAndEndDate(event.Event.ID, event.Event.UserID, startDate, endDate)
			if err != nil {
				log.Printf("error updating event status: %s", err.Error())
				continue
			}

		case frequency.Weekly:
			startDate := event.Event.StartDate.Time.AddDate(0, 0, 7)
			endDate := event.Event.EndDate.Time.AddDate(0, 0, 7)
			err := b.repo.Event.UpdateStartAndEndDate(event.Event.ID, event.Event.UserID, startDate, endDate)
			if err != nil {
				log.Printf("error updating event status: %s", err.Error())
				continue
			}
		case frequency.MonthlyDate:
			eventStartDate := event.Event.StartDate.Time
			eventEndDate := event.Event.StartDate.Time

			startDate := eventStartDate.AddDate(0, 1, 0)
			if startDate.Day() == 1 {
				startDate = startDate.AddDate(0, 0, -1)
			}
			// fmt.Println("Next date ", startDate)
			endDate := eventEndDate.AddDate(0, 1, 0)
			if endDate.Day() == 1 {
				endDate = endDate.AddDate(0, 0, -1)
			}
			err := b.repo.Event.UpdateStartAndEndDate(event.Event.ID, event.Event.UserID, startDate, endDate)
			if err != nil {
				log.Printf("error updating event status: %s", err.Error())
				continue
			}

		case frequency.MonthlyDay:
			startDate := event.Event.StartDate.Time
			endDate := event.Event.EndDate.Time
			sub := endDate.Sub(startDate)

			newStartDate := getNextDayNumberInMonth(startDate)
			newEndDate := newStartDate.Add(sub)

			err := b.repo.Event.UpdateStartAndEndDate(event.Event.ID, event.Event.UserID, newStartDate, newEndDate)
			if err != nil {
				log.Printf("error updating event status: %s", err.Error())
				continue
			}

		case frequency.Annualy:
			startDate := event.Event.StartDate.Time
			endDate := event.Event.EndDate.Time
			newStartDate := startDate.AddDate(1, 0, 0)
			newEndDate := endDate.AddDate(1, 0, 0)
			err := b.repo.Event.UpdateStartAndEndDate(event.Event.ID, event.Event.UserID, newStartDate, newEndDate)
			if err != nil {
				log.Printf("error updating event status: %s", err.Error())
				continue
			}

		default:
			log.Printf("error: invalid frequency in event %s", event.Event.ID)
			continue
		}

	}
}

func sortRemindersByMinutesUntilEvent(reminders []models.Reminder) []models.Reminder {
	sortedReminders := make([]models.Reminder, len(reminders))
	copy(sortedReminders, reminders)

	sort.SliceStable(sortedReminders, func(i, j int) bool {
		return sortedReminders[i].MinutesUntilEvent > sortedReminders[j].MinutesUntilEvent
	})

	return sortedReminders
}

func getNextDayNumberInMonth(date time.Time) time.Time {
	year, month, day := date.Date()
	loc := date.Location()
	hour, minute := date.Hour(), date.Minute()
	fmt.Println("next Month ", month+1)
	// Инициализируем счетчик для подсчета вторых вторников
	count := 0
	weekday := date.Weekday()
	// Начинаем с первого дня месяца
	currentDate := time.Date(year, month+1, 1, hour, minute, 0, 0, loc)
	// Находим порядковый номер дня недели в месяце
	ordinal := (day-1)/7 + 1
	// Перебираем все дни месяца
	for {
		fmt.Println("Перебор: ", currentDate)
		// Если текущий день - вторник, увеличиваем счетчик
		if currentDate.Weekday() == weekday {
			count++
			fmt.Println("Неделя номер ", count)
			// Если это второй вторник, возвращаем эту дату
			if count == ordinal {
				return currentDate
			}
		}

		// Переходим к следующему дню
		currentDate = currentDate.AddDate(0, 0, 1)
		fmt.Println("Следующий день ", currentDate)
		// Если достигли следующего месяца, выходим из цикла
		if currentDate.Month() != month+1 {
			break
		}
	}

	// Если не удалось найти второй вторник, возвращаем нулевую дату
	return time.Time{}
}
