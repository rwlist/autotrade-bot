package app

import (
	"fmt"
	"github.com/petuhovskiy/telegram"
	"strings"

	"github.com/petuhovskiy/autotrade-bot/conf"
)

type Handler struct {
	bot *telegram.Bot
	logic *Logic
	cfg *conf.Struct
}

func NewHandler(bot *telegram.Bot, logic *Logic, cfg *conf.Struct) *Handler {
	return &Handler{
		bot: bot,
		logic: logic,
		cfg: cfg,
	}
}

func (h *Handler) Handle(upd telegram.Update) {
	if upd.Message == nil {
		return
	}

	msg := upd.Message
	if msg.From.ID != h.cfg.AdminID {
		return
	}

	h.handleMessage(msg)
}

func (h *Handler) sendMessage(chatID int, text string) {
	_, _ = h.bot.SendMessage(&telegram.SendMessageRequest{
		ChatID:                str(chatID),
		Text:                  text,
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

func (h *Handler) handleCommand(chatID int, cmds []string) {
	if len(cmds) == 0 {
		return
	}

	cmd := cmds[0]
	switch cmd {
	case "/status":
		h.commandStatus(chatID)

	default:
		h.commandNotFound(chatID)
	}
}

func (h *Handler) commandStatus(chatID int) {
	status, err := h.logic.CommandStatus()
	if err != nil {
		text := fmt.Sprintf("Error while status:\n\n%s", err)
		h.sendMessage(chatID, text)
		return
	}
	res := fmt.Sprintf("BTC: 1 ≈ %v USDT\n\nWallet balance:", status.rate)
	if len(status.balances) == 0 {
		res += "\nNo money :^)"
	}
	for _, v := range status.balances {
		res += fmt.Sprintf("\n%v:\nFree: %v\nLocked: %v\n", v.Asset, v.Free, v.Locked)
	}
	h.sendMessage(chatID, res)
}

func (h *Handler) commandNotFound(chatID int) {
	h.commandHelp(chatID)
}

func (h *Handler) commandHelp(chatID int) {
	str := `Need some help?

/status				displays btc/usdt rate and your binance wallet balance`

	h.sendMessage(chatID, str)
}