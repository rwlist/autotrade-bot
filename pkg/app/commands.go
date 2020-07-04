package app

import (
	"fmt"
	"log"
)

func (h *Handler) handleCommand(chatID int, cmds []string) {
	if len(cmds) == 0 {
		return
	}

	cmd := cmds[0]
	switch cmd {
	case "/testSwitch":
		h.commandTestModeSwitch(chatID)

	case "/end":
		h.commandEnd(chatID)

	case "/begin":
		h.commandBegin(chatID, cmds[1])

	case "/draw":
		h.commandDraw(chatID, cmds[1])

	case "/sell":
		h.commandSell(chatID)

	case "/buy":
		h.commandBuy(chatID)

	case "/status":
		h.commandStatus(chatID)

	default:
		h.commandNotFound(chatID)
	}
}

func (h *Handler) commandStatus(chatID int) {
	const places = 2

	status, err := h.svc.Status.Status()
	if err != nil {
		text := fmt.Sprintf("Error while status:\n\n%s", err)
		h.sendMessage(chatID, text)
		return
	}

	res := fmt.Sprintf("BTC: 1 ≈ %v USDT \n", status.Rate)
	res += fmt.Sprintf("Total in USD ≈ %v $ \n\n", status.Total.RoundBank(places))
	res += "Wallet balance:"

	if len(status.Balances) == 0 {
		res += "\nNo money :^)"
	}

	for _, v := range status.Balances {
		res += fmt.Sprintf("\n%v:\n", v.Asset)
		res += fmt.Sprintf("In USD: %v$\n", v.USD.RoundBank(places))
		res += fmt.Sprintf("Free: %v\n", v.Free)
		res += fmt.Sprintf("Locked: %v\n", v.Locked)
	}

	h.sendMessage(chatID, res)
}

func (h *Handler) commandBuy(chatID int) {
	err := h.svc.Logic.Buy(&Sender{h.bot, chatID})
	if err != nil {
		log.Println("Error in commandBuy: ", err)
		txt := fmt.Sprintf("Error in commandBuy:\n%v", err)
		h.sendMessage(chatID, txt)
		return
	}
	h.sendMessage(chatID, "Command Buy finished")
}

func (h *Handler) commandSell(chatID int) {
	err := h.svc.Logic.Sell(&Sender{h.bot, chatID})
	if err != nil {
		log.Println("Error in commandSell: ", err)
		txt := fmt.Sprintf("Error in commandSell:\n%v", err)
		h.sendMessage(chatID, txt)
		return
	}
	h.sendMessage(chatID, "Command Sell finished")
}

func (h *Handler) commandDraw(chatID int, str string) {
	b, err := h.svc.Logic.Draw(str, nil)
	if err != nil {
		log.Println("Error in commandDraw: ", err)
		txt := fmt.Sprintf("Error in commandDraw:\n%v", err)
		h.sendMessage(chatID, txt)
		return
	}
	err = h.sendPhoto(chatID, "graph.png", b)
	if err != nil {
		log.Println("Error in commandDraw: ", err)
		txt := fmt.Sprintf("Error in commandDraw:\n%v", err)
		h.sendMessage(chatID, txt)
		return
	}
}

func (h *Handler) commandBegin(chatID int, str string) {
	err := h.svc.Logic.Begin(&Sender{h.bot, chatID}, str, h.isTest)
	if err != nil {
		log.Println("Error in commandBegin: ", err)
		txt := fmt.Sprintf("Error in commandBegin:\n%v", err)
		h.sendMessage(chatID, txt)
		return
	}
}

func (h *Handler) commandEnd(chatID int) {
	err := h.svc.Logic.End(&Sender{h.bot, chatID}, h.isTest)
	if err != nil {
		log.Println("Error in commandEnd: ", err)
		txt := fmt.Sprintf("Error in commandEnd:\n%v", err)
		h.sendMessage(chatID, txt)
		return
	}
}

func (h *Handler) commandNotFound(chatID int) {
	h.commandHelp(chatID)
}

func (h *Handler) commandTestModeSwitch(chatID int) {
	h.isTest = !h.isTest
	if h.isTest {
		h.sendMessage(chatID, "Testmode is ON!")
	} else {
		h.sendMessage(chatID, "Testmode is OFF!\nNow, be careful")
	}
}

func (h *Handler) commandHelp(chatID int) {
	str := `Need some help?

/status				displays BTC/USDT rate and your binance wallet balance

/sell				sells all BTC
/buy				buys BTC with all USDT

/draw <formula> (example: rate-10+0.0002*(now-start)^1.2) 		draws graph of given formula
/begin <formula> buys BTC with all USDT, activates trigger 
/end    deactivates trigger and sells all BTC
/drawit TBD
`

	h.sendMessage(chatID, str)
}
