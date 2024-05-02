package main

import (
	"friendly-bot/internal/repository"
	"friendly-bot/internal/telegram"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	// Подключение к базе данных
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("PG_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("error connecting postgres DB: %s", err.Error())
	}

	repo := repository.NewRepository(db)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	telegramBot := telegram.NewBot(bot, repo)
	if err := telegramBot.Start(); err != nil {
		log.Fatalf("error starting th bot: %s", err.Error())
	}

	ticker := time.NewTicker(1 * time.Minute) // Проверка каждые 24 часа
	defer ticker.Stop()
	// currentTime := time.Now()
	// queryGet := fmt.Sprintf(`SELECT e.*, f.*
	// 						FROM %s e
	// 						JOIN %s fe ON e.id = fe.event_id
	// 						JOIN %s f ON fe.friend_id = fe.friend_id
	// 						WHERE e.start_date <= $1 AND e.end_date > $2`, eventTable, friendsEventsTable, friendTable)
	// queryUpdateStartNotify := fmt.Sprintf("UPDATE %s e SET start_notify_sent = true", eventTable)

	// for {
	// 	select {
	// 	case <-ticker.C:
	// 		// Запрос к базе данных для получения друзей, у которых сегодня день рождения
	// 		rows, err := db.Query("SELECT name FROM friend WHERE DAY(birth) = DAY(NOW()) AND MONTH(birth) = MONTH(NOW())")
	// 		if err != nil {
	// 			log.Println("Error querying database:", err)
	// 			continue
	// 		}
	// 		defer rows.Close()
	// 		// Отправка сообщений в Telegram
	// 		for rows.Next() {
	// 			var name string
	// 			if err := rows.Scan(&name); err != nil {
	// 				log.Println("Error scanning row:", err)
	// 				continue
	// 			}
	// 			msg := tgbotapi.NewMessage(chatID, "С днём рождения, "+name+"!")
	// 			_, err := bot.Send(msg)
	// 			if err != nil {
	// 				log.Println("Error sending message:", err)
	// 			}
	// 		}
	// 		if err := rows.Err(); err != nil {
	// 			log.Println("Error iterating over rows:", err)
	// 		}
	// 	}
	// }

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Bot shutting down...")

	if err := db.Close(); err != nil {
		log.Fatalf("error while closing database connection: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
