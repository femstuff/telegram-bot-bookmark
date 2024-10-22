package main

import (
	"flag"
	"log"

	"telegram-bot/clients/telegram"
	"telegram-bot/consumer/event-consumer"
	"telegram-bot/events/tg"
	"telegram-bot/storage/files"
)

const (
	botHost     = "api.telegram.org"
	pathStorage = "linked-storage"
	batchSize   = 100
)

func main() {
	client := telegram.New(botHost, mustToken())

	eventsProcessor := tg.New(client, files.New(pathStorage))

	log.Print("bot is started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("bot is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("token-bot", "", "token from bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("Error, empty value token")
	}

	return *token
}
