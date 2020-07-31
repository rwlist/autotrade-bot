package app

import (
	"github.com/petuhovskiy/telegram"
	"github.com/rwlist/autotrade-bot/pkg/convert"
)

type Sender struct {
	bot    *telegram.Bot
	chatID int
}

func (s *Sender) Send(text string) {
	_, _ = s.bot.SendMessage(&telegram.SendMessageRequest{
		ChatID: convert.Str(s.chatID),
		Text:   text,
	})
}

func (s *Sender) SendPhoto(name string, b []byte) {
	_, _ = s.bot.SendPhoto(&telegram.SendPhotoRequest{
		ChatID: convert.Str(s.chatID),
		Photo:  NewBytesUploader(name, b),
	})
}
