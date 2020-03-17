package app

import (
	"fmt"
	"strings"

	"github.com/petuhovskiy/telegram"

	"github.com/rwlist/autotrade-bot/conf"
)

type Handler struct {
	bot   *telegram.Bot
	logic *Logic
	cfg   *conf.Struct
}

func NewHandler(bot *telegram.Bot, logic *Logic, cfg *conf.Struct) *Handler {
	return &Handler{
		bot:   bot,
		logic: logic,
		cfg:   cfg,
	}
}

func (h *Handler) Handle(upd telegram.Update) {
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
		ChatID: str(chatID),
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

func (h *Handler) handleCommand(chatID int, cmds []string) {
	if len(cmds) == 0 {
		return
	}

	cmd := cmds[0]
	switch cmd {
	case "/buy":
		h.commandBuy(chatID)

	case "/status":
		h.commandStatus(chatID)

	//TEST CASES
	case "/testbuy":
		h.commandTestBuyAll(chatID)

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
	res := fmt.Sprintf("BTC: 1 ≈ %v USDT \nTotal in USD ≈ %v $ \n\nWallet balance:", status.rate, status.total)
	if len(status.balances) == 0 {
		res += "\nNo money :^)"
	}
	for _, v := range status.balances {
		res += fmt.Sprintf("\n%v:\nIn USD: %v$\nFree: %v\nLocked: %v\n", v.asset, v.usd, v.free, v.locked)
	}
	h.sendMessage(chatID, res)
}

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

func (h *Handler) commandBuy(chatID int) {
	h.logic.CommandBuy(&Sender{h.bot,chatID})
	h.sendMessage(chatID, "Command \"/buy\" finished")
}

//----------TEST_BUY_COMMAND--------------------------
func (h *Handler) commandTestBuyAll(chatID int) {
	h.logic.TestCommandBuy(&Sender{h.bot,chatID})
	h.sendMessage(chatID, "Command \"/TESTbuy\" finished")
}
//-----------------------------------------------------

func (h *Handler) commandNotFound(chatID int) {
	h.commandHelp(chatID)
}

func (h *Handler) commandHelp(chatID int) {
	str := `Need some help?

/status				displays btc/usdt rate and your binance wallet balance`

	h.sendMessage(chatID, str)
}
