package vaccinatorplus

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (v *Vaccinator) handleNewMessage(m *tgbotapi.Message, c Conversation) {
	log.Printf("New message from %s: %s", c.ToHumanName(), m.Text)

	if c.RequestedYear == 0 {
		v.handleSetup1(m, c)
		return
	}

	if c.NotifyAllYears == nil {
		v.handleSetup2(m, c)
		return
	}

	if m.Text == "/setup" {
		db, err := v.db()
		if err != nil {
			return
		}
		c.RequestedYear = 0
		c.NotifyAllYears = nil
		db.Save(&c)
		v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
		time.Sleep(1 * time.Second)
		msg := tgbotapi.NewMessage(c.ChatID, "Let's start! In what year were you born?")
		v.bot.Send(msg)
		return
	}

	v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
	time.Sleep(time.Millisecond * 500)
	msg := tgbotapi.NewMessage(c.ChatID, "Sadly, my creator has not teached me how to respond to that")
	v.bot.Send(msg)
	v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
	time.Sleep(time.Millisecond * 500)
	msg = tgbotapi.NewMessage(c.ChatID, "But I can send you a duckpic!")
	v.bot.Send(msg)
	v.sendRandomCat(c)

	v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
	time.Sleep(time.Millisecond * 500)
	msg = tgbotapi.NewMessage(c.ChatID, fmt.Sprintf("Btw, they are already calling people born in %d. I will keep you posted!", v.currentYear))
	v.bot.Send(msg)
}

func (v *Vaccinator) handleSetup1(m *tgbotapi.Message, c Conversation) {
	year, err := strconv.Atoi(m.Text)
	if err != nil {
		v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
		time.Sleep(1 * time.Second)
		msg := tgbotapi.NewMessage(c.ChatID, "That's not a year...")
		v.bot.Send(msg)
		return
	}

	v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
	time.Sleep(time.Millisecond * 500)
	msg := tgbotapi.NewMessage(c.ChatID, strconv.Itoa(year))
	v.bot.Send(msg)
	v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
	time.Sleep(time.Millisecond * 500)
	msg = tgbotapi.NewMessage(c.ChatID, "Ok!")
	v.bot.Send(msg)

	db, err := v.db()
	if err != nil {
		return
	}
	c.RequestedYear = year
	db.Save(&c)
	v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
	time.Sleep(time.Millisecond * 500)
	msg = tgbotapi.NewMessage(c.ChatID, "Would you like to be notified when any other year is called for registration?")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
		{
			Text: "Yes",
		},
		{
			Text: "No",
		}})

	v.bot.Send(msg)
}

func (v *Vaccinator) handleSetup2(m *tgbotapi.Message, c Conversation) {
	if m.Text != "Yes" && m.Text != "No" {
		v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
		time.Sleep(time.Millisecond * 500)
		msg := tgbotapi.NewMessage(c.ChatID, "That's not a valid answer. Please use the buttons 'Yes' or 'No'")
		v.bot.Send(msg)
		return
	}

	yes := m.Text == "Yes"
	c.NotifyAllYears = &yes
	db, err := v.db()
	if err != nil {
		return
	}
	db.Save(&c)
	v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
	time.Sleep(time.Millisecond * 500)
	msg := tgbotapi.NewMessage(c.ChatID, "Duly noted!")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	v.bot.Send(msg)
	v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
	time.Sleep(time.Millisecond * 500)

	var finalSetup string
	if *c.NotifyAllYears {
		finalSetup = fmt.Sprintf("You will be notified when %d pops in, and when any other year is called", c.RequestedYear)
	} else {
		finalSetup = fmt.Sprintf("You will only be notified when %d is called", c.RequestedYear)
	}
	msg = tgbotapi.NewMessage(c.ChatID, finalSetup)
	v.bot.Send(msg)

	v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
	time.Sleep(time.Millisecond * 500)
	msg = tgbotapi.NewMessage(c.ChatID, "Btw, you can always change these settings by typing '/setup'")
	v.bot.Send(msg)
}

func (v *Vaccinator) handleNewUser(m *tgbotapi.Message) error {
	db, err := v.db()
	if err != nil {
		return err
	}
	chatID := m.Chat.ID
	var name string
	if m.From.FirstName != "" {
		name = m.From.FirstName
	} else {
		name = m.From.UserName
	}
	c := Conversation{
		ChatID:    chatID,
		FirstName: name,
		LastName:  m.From.LastName,
		Username:  m.From.UserName,
	}
	db.Create(&c)

	_, _ = v.bot.Send(tgbotapi.NewChatAction(chatID, "typing"))
	time.Sleep(1 * time.Second)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Hello there, %s", name))
	v.bot.Send(msg)
	v.bot.Send(tgbotapi.NewChatAction(chatID, "typing"))
	time.Sleep(1 * time.Second)
	msg = tgbotapi.NewMessage(chatID, "💉 I am VaccinatorPlus, an extremely GDPR-compliant bot developed by Juan Font (juanfontalonso@gmail.com)")
	v.bot.Send(msg)
	v.bot.Send(tgbotapi.NewChatAction(chatID, "typing"))
	time.Sleep(1 * time.Second)
	msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("I will let you know as soon as you can be registered for vaccination in NL 🇳🇱 (currently calling people born in %d)", v.currentYear))
	v.bot.Send(msg)
	v.bot.Send(tgbotapi.NewChatAction(chatID, "typing"))
	time.Sleep(1 * time.Second)
	msg = tgbotapi.NewMessage(chatID, "Let's start! In what year were you born?")
	v.bot.Send(msg)
	return nil
}

func (v *Vaccinator) handleVaccinationCall(c Conversation) error {
	msg := tgbotapi.NewMessage(c.ChatID, fmt.Sprintf("%s!! Your year (%d) has been called for registration @%s", c.FirstName, c.RequestedYear, c.Username))
	v.bot.Send(msg)

	v.bot.Send(tgbotapi.NewChatAction(c.ChatID, "typing"))
	time.Sleep(time.Millisecond * 500)

	msg = tgbotapi.NewMessage(c.ChatID, fmt.Sprintf("💉💉💉💉 Run to https://coronatest.nl and get your appointment!"))
	v.bot.Send(msg)

	c.NotifiedYear = c.RequestedYear
	db, err := v.db()
	if err != nil {
		return err
	}
	db.Save(&c)
	return nil
}

func (v *Vaccinator) handleNotifyAllYears(c Conversation, year int) error {
	msg := tgbotapi.NewMessage(c.ChatID, fmt.Sprintf("%s, the cohort born in %d is now being called to make their appointments", c.FirstName, year))
	v.bot.Send(msg)

	c.NotifiedYear = year
	db, err := v.db()
	if err != nil {
		return err
	}
	db.Save(&c)
	return nil
}
