package app

import "github.com/petuhovskiy/telegram"

type Sender struct {
	bot    *telegram.Bot
	chatID int
}

func (s *Sender) Send(text string) {
	_, _ = s.bot.SendMessage(&telegram.SendMessageRequest{
		ChatID: str(s.chatID),
		Text:   text,
	})
}

func (s *Sender) SendPhoto(name string, b []byte) error {
	_, err := s.bot.SendPhoto(&telegram.SendPhotoRequest{
		ChatID: str(s.chatID),
		Photo: NewBytesUploader("Graph.png", b),
	})
	return err
}
