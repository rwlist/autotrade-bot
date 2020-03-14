package app

import (
	"fmt"
	"github.com/adshao/go-binance"
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

func (h *Handler) commandBuy(chatID int) {
	ch := make(chan *OrderInfo)
	go h.logic.CommandBuy(ch)
	for order := range ch {
		if order.Err != nil {
			text := fmt.Sprintf("Error while Buy:\n\n%s", order.Err)
			h.sendMessage(chatID, text)
			return
		}
		if order.InfoType == 1 {
			text := fmt.Sprintf("A %v BTC/USDT order was placed with price = %v.\nWaiting for 2 seconds..", order.Side, order.Price)
			h.sendMessage(chatID, text)
		} else if order.InfoType == 2 {
			text := fmt.Sprintf("Done %v / %v\nStatus: %v", order.ExecutedQuantity, order.OrigQuantity, order.Status)
			h.sendMessage(chatID, text)
			if order.Status == binance.OrderStatusTypeFilled {
				break
			}
		}
	}
	h.sendMessage(chatID, "Command \"/buy\" finished")
}

func (h *Handler) commandTestBuyAll(chatID int) {
	err := h.logic.CommandTestOrderAll()
	if err != nil {
		text := fmt.Sprintf("Error while testBuyAll:\n\n%s", err)
		h.sendMessage(chatID, text)
		return
	}
}

func (h *Handler) commandNotFound(chatID int) {
	h.commandHelp(chatID)
}

func (h *Handler) commandHelp(chatID int) {
	str := `Need some help?

/status				displays btc/usdt rate and your binance wallet balance`

	h.sendMessage(chatID, str)
}
