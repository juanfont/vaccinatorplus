package vaccinatorplus

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

//go:embed ducks/*
var ducks embed.FS

type Vaccinator struct {
	bot         *tgbotapi.BotAPI
	currentYear int
	dbPath      string
	token       string
}

func NewVaccinator(token string, dbPath string, initialYear int) (*Vaccinator, error) {
	rand.Seed(time.Now().Unix())
	v := Vaccinator{
		currentYear: initialYear,
		dbPath:      dbPath,
		token:       token,
	}

	err := v.initDB()
	if err != nil {
		return nil, err
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	// bot.Debug = true
	v.bot = bot
	log.Printf("Authorized on account %s", v.bot.Self.UserName)
	return &v, err
}

func (v *Vaccinator) Run() error {
	go v.yearChecker()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := v.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		db, err := v.db()
		if err != nil {
			return err
		}

		c := Conversation{}
		result := db.First(&c, "chat_id = ?", update.Message.Chat.ID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("New user %s (%s %s)", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName)
			go v.handleNewUser(update.Message)
		} else {
			go v.handleNewMessage(update.Message, c)
		}
	}

	return nil
}

func (v *Vaccinator) yearChecker() {
	for {
		resp, err := http.Get("https://www.rijksoverheid.nl/onderwerpen/coronavirus-vaccinatie/prikuitnodiging-en-afspraak")
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		text := string(body)
		nextYear := v.currentYear + 1
		if strings.Contains(text, strconv.Itoa(nextYear)) {
			log.Printf("Year %d has appeared!", nextYear)
			v.notifyUsers(nextYear)
			v.currentYear = nextYear
			time.Sleep(5 * time.Minute)
		} else {
			log.Printf("Still in %d", v.currentYear)
			time.Sleep(5 * time.Minute)
		}
	}
}

func (v *Vaccinator) notifyUsers(year int) {
	db, err := v.db()
	if err != nil {
		return
	}
	chats := []Conversation{}
	result := db.Find(&chats)
	log.Printf("Found %d rows", result.RowsAffected)
	for _, c := range chats {
		if c.RequestedYear == year && c.NotifiedYear != year {
			v.handleVaccinationCall(c)
			continue
		}

		if c.NotifyAllYears != nil && *c.NotifyAllYears == true && c.NotifiedYear != year {
			v.handleNotifyAllYears(c, year)
		}
	}
}

func (v *Vaccinator) sendRandomCat(c Conversation) {
	mews, _ := fs.ReadDir(ducks, "ducks")
	random := mews[rand.Intn(len(mews))]
	fname := random.Name()
	data, err := ducks.ReadFile("ducks/" + fname)
	if err != nil {
		log.Fatalln(err)
	}
	photoFileBytes := tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: data,
	}
	v.bot.Send(tgbotapi.NewPhotoUpload(c.ChatID, photoFileBytes))
}
