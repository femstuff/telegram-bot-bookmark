package tg

import (
	"errors"

	"telegram-bot/clients/telegram"
	"telegram-bot/events"
	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

var errMsgProcess = errors.New("unknown event type")
var errMsgMeta = errors.New("unknown meta type")

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("cant get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, upd := range updates {
		res = append(res, event(upd))
	}

	p.offset = updates[len(updates)-1].UpdateID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("cant create process", errMsgProcess)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("cant process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return e.Wrap("cant process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("cant get meta", errMsgMeta)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if fetchType(upd) == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.UserName,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
