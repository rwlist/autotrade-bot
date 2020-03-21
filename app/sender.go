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
