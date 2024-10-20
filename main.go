package main

import (
	"flag"
	"log"

	"telegram-bot/clients/telegram"
)

const (
	botHost = "api.telegram.org"
)

func main() {
	client := telegram.New(botHost, mustToken())
}

func mustToken() string {
	token := flag.String("token-bot", "", "token from bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("Error, empty value token")
	}

	return *token
}
