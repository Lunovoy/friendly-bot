package main

import (
	"database/sql"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

var token = "6842952380:AAGJyH1ukPjCcP3HD700ikoN-GWNWcMCw2s"

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}
	// Подключение к базе данных
	db, err := sql.Open("pgx", "user:password@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Создание клиента Telegram бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	// Настройка расписания проверки дней рождения
	ticker := time.NewTicker(1 * time.Minute) // Проверка каждые 24 часа
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Запрос к базе данных для получения друзей, у которых сегодня день рождения
			rows, err := db.Query("SELECT name FROM friend WHERE DAY(birth) = DAY(NOW()) AND MONTH(birth) = MONTH(NOW())")
			if err != nil {
				log.Println("Error querying database:", err)
				continue
			}
			defer rows.Close()
			// Отправка сообщений в Telegram
			for rows.Next() {
				var name string
				if err := rows.Scan(&name); err != nil {
					log.Println("Error scanning row:", err)
					continue
				}
				msg := tgbotapi.NewMessage(chatID, "С днём рождения, "+name+"!")
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Error sending message:", err)
				}
			}
			if err := rows.Err(); err != nil {
				log.Println("Error iterating over rows:", err)
			}
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
