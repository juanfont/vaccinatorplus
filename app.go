package vaccinatorplus

import (
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Vaccinator struct {
	dbPath string
	token  string
}

func NewVaccinator(token string, dbPath string, initialYear int) (*Vaccinator, error) {
	v := Vaccinator{
		dbPath: dbPath,
		token:  token,
	}

	err := v.initDB()
	if err != nil {
		return nil, err
	}

	return &v, err
}

func (v *Vaccinator) Run() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		spew.Dump(update.Message)

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		time.Sleep(3 * time.Second)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
