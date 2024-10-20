package tg

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

const (
	CmdHelp   = "/help"
	CmdRandom = "/random"
	CmdStart  = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("insert new command %s from %s", text, username)

	if isLink(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case CmdStart:
		return p.sendHello(chatID)
	case CmdRandom:
		return p.sendRandom(chatID, username)
	case CmdHelp:
		return p.sendHelp(chatID)
	default:
		return p.tg.SendMessage(chatID, unknownMsg)
	}
}

func (p *Processor) savePage(chatID int, pageURL, username string) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	IsExist, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if IsExist {
		return p.tg.SendMessage(chatID, alreadyExistsMsg)
	}

	if err := p.storage.Save(page); err != nil {
		return e.Wrap("cant save page", err)
	}

	if err := p.tg.SendMessage(chatID, successfulMsg); err != nil {
		return e.Wrap("cant send message", err)
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) error {
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, errors.New("no saved page")) {
		return e.Wrap("cant found page", err)
	}

	if errors.Is(err, errors.New("no saved page")) {
		return p.tg.SendMessage(chatID, enoughPagesMsg)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return e.Wrap("cant send link", err)
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, helpMsg)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, helloMsg)
}

func isLink(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
