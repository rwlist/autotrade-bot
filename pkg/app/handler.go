package app

import (
	"strings"

	"github.com/rwlist/autotrade-bot/pkg/history"

	"github.com/rwlist/autotrade-bot/pkg/convert"

	"github.com/rwlist/autotrade-bot/pkg/logic"

	"github.com/rwlist/autotrade-bot/pkg/stat"

	"github.com/petuhovskiy/telegram"

	"github.com/rwlist/autotrade-bot/pkg/conf"
)

type Handler struct {
	bot *telegram.Bot
	cfg *conf.Struct
	svc Services
}

type Services struct {
	Logic        *logic.Logic
	Status       *stat.Service
	StatusChatex *stat.Service
	History      *history.History
}

func NewHandler(bot *telegram.Bot, cfg *conf.Struct, svc Services) *Handler {
	return &Handler{
		bot: bot,
		cfg: cfg,
		svc: svc,
	}
}

func (h *Handler) Handle(upd *telegram.Update) {
	if upd.Message == nil {
		return
	}

	msg := upd.Message
	if msg.From.ID != h.cfg.Bot.AdminID {
		return
	}

	h.handleMessage(msg)
}

func (h *Handler) sendMessage(chatID int, text string) {
	_, _ = h.bot.SendMessage(&telegram.SendMessageRequest{
		ChatID: convert.Str(chatID),
		Text:   text,
	})
}

func (h *Handler) sendPhoto(chatID int, name string, b []byte) {
	_, _ = h.bot.SendPhoto(&telegram.SendPhotoRequest{
		ChatID: convert.Str(chatID),
		Photo:  NewBytesUploader(name, b),
	})
}

func (h *Handler) handleMessage(msg *telegram.Message) {
	text := msg.Text
	if !strings.HasPrefix(text, "/") {
		return
	}

	cmds := strings.Split(text, " ")
	h.handleCommand(msg.Chat.ID, cmds)
}
