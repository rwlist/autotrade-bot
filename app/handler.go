package app

import (
	"strings"

	"github.com/rwlist/autotrade-bot/logic"

	"github.com/rwlist/autotrade-bot/app/stat"

	"github.com/rwlist/autotrade-bot/pkg/tostr"

	"github.com/petuhovskiy/telegram"

	"github.com/rwlist/autotrade-bot/pkg/conf"
)

type Handler struct {
	bot    *telegram.Bot
	cfg    *conf.Struct
	svc    Services
	isTest bool
}

type Services struct {
	Logic  *logic.Logic
	Status *stat.Service
}

func NewHandler(bot *telegram.Bot, cfg *conf.Struct, svc Services) *Handler {
	return &Handler{
		bot:    bot,
		cfg:    cfg,
		svc:    svc,
		isTest: true,
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
		ChatID: tostr.Str(chatID),
		Text:   text,
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
