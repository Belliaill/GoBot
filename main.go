package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

func main() {
	godotenv.Load()

	pref := telebot.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(c telebot.Context) error {
		return c.Send(fmt.Sprintf("Hello, %s!", c.Sender().FirstName))
	})

	b.Start()
}
