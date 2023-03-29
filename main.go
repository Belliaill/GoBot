package main

import (
	"fmt"
	"gobot/db"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
	"gopkg.in/telebot.v3"
)

const (
	None BotState = iota
	Donate
	Chatting
)

type BotState int

var (
	DB  = db.NewDB("database")
	bot *telebot.Bot

	menu = &telebot.ReplyMarkup{}

	chatBtn   = menu.Data("Chat", "chat")
	donateBtn = menu.Data("Donate", "donate")

	adminChat   telebot.ChatID
	dontateChat telebot.ChatID
)

func main() {
	godotenv.Load()
	ac, _ := strconv.ParseInt(os.Getenv("ADMIN_CHAT_ID"), 10, 64)
	ad, _ := strconv.ParseInt(os.Getenv("DONATE_CHAT_ID"), 10, 64)
	adminChat = telebot.ChatID(ac)
	dontateChat = telebot.ChatID(ad)

	pref := telebot.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	bot = b

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", HandleStart)
	b.Handle(&chatBtn, HandleChatBtn)
	b.Handle(&donateBtn, HandleDonateBtn)
	b.Handle(telebot.OnText, HandleAll)
	b.Handle(telebot.OnPhoto, HandleScreen)

	b.Start()
}

func HandleStart(c telebot.Context) error {
	Id := int(c.Sender().ID)

	e := slices.IndexFunc(DB.GetUsers(), func(u db.User) bool {
		return u.Id == Id
	})

	menu.Inline(
		menu.Row(chatBtn, donateBtn),
	)

	c.Set("state", None)

	if e == -1 {
		Name := c.Sender().FirstName
		DB.AppendUser(db.User{
			Id:   Id,
			Name: Name,
		})
		return c.Send(fmt.Sprintf("Привет, %s!", c.Sender().FirstName), menu)
	} else {
		return c.Send(fmt.Sprintf("Привет снова, %s!", c.Sender().FirstName), menu)
	}
}

func HandleAll(c telebot.Context) error {
	menu.Inline(
		menu.Row(chatBtn, donateBtn),
	)

	// users := DB.GetUsers()
	if adminChat == telebot.ChatID(c.Chat().ID) {
		r := c.Message().ReplyTo
		if r != nil {
			_, err := bot.Forward(r.Sender, c.Message())
			return err
		}
	} else if c.Get("state") == Chatting {
		_, err := bot.Forward(adminChat, c.Message())
		return err
	}
	return c.Reply("Выбирете действие", menu)
}

func HandleChatBtn(c telebot.Context) error {
	c.Set("state", Chatting)
	return c.Send("Отправте сообщения")
}

func HandleDonateBtn(c telebot.Context) error {
	c.Set("state", Donate)
	return c.Send("Отправте скрин")
}

func HandleScreen(c telebot.Context) error {
	if c.Get("state") != Donate {
		return c.Reply("Выбирете действие", menu)
	}
	c.Set("state", None)
	c.Reply("Что дальше?", menu)
	_, err := bot.Forward(dontateChat, c.Message())
	return err
}
