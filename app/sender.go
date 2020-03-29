package app

import (
	"github.com/petuhovskiy/telegram"
	"github.com/rwlist/autotrade-bot/to_str"
)

type Sender struct {
	bot    *telegram.Bot
	chatID int
}

func (s *Sender) Send(text string) {
	_, _ = s.bot.SendMessage(&telegram.SendMessageRequest{
		ChatID: to_str.Str(s.chatID),
		Text:   text,
	})
}

func (s *Sender) SendPhoto(name string, b []byte) error {
	_, err := s.bot.SendPhoto(&telegram.SendPhotoRequest{
		ChatID: to_str.Str(s.chatID),
		Photo:  NewBytesUploader(name, b),
	})
	return err
}
