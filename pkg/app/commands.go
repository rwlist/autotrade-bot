package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rwlist/autotrade-bot/pkg/trade/chatex"

	log "github.com/sirupsen/logrus"

	"github.com/rwlist/autotrade-bot/pkg/stat"
)

func (h *Handler) handleCommand(chatID int, cmds []string) { //nolint:gocyclo
	if len(cmds) == 0 {
		return
	}

	cmd := cmds[0]
	str := ""
	if len(cmds) > 1 {
		str = cmds[1]
	}

	switch cmd {
	case "/alter":
		h.commandAlter(chatID, str)

	case "/fstat":
		h.commandFstat(chatID, str)

	case "/setScale":
		h.commandSetScale(chatID, str)

	case "/testSwitch":
		h.commandTestModeSwitch(chatID)

	case "/end":
		h.commandEnd(chatID)

	case "/begin":
		h.commandBegin(chatID, str)

	case "/draw":
		h.commandDraw(chatID, str)

	case "/sell":
		h.commandSell(chatID)

	case "/buy":
		h.commandBuy(chatID)

	case "/status":
		h.commandStatus(chatID, str)

	case "/history":
		h.commandHistory(chatID)

	case "/opts":
		h.commandOpts(chatID)

	case "/opt_set":
		h.commandOptSet(chatID, cmds[1:])

	case "/opt_help":
		h.commandOptHelp(chatID)

	case "/opt_auto":
		h.commandOptAuto(chatID, cmds[1:])

	default:
		h.commandNotFound(chatID)
	}
}

func (h *Handler) commandStatus(chatID int, str string) {
	const places = 2

	var st *stat.Service
	if strings.Contains(str, "chatex") {
		st = h.svc.StatusChatex
	} else {
		st = h.svc.Status
	}

	status, err := st.Status()
	if err != nil {
		text := fmt.Sprintf("Error while status:\n\n%s", err)
		h.sendMessage(chatID, text)
		return
	}

	res := fmt.Sprintf("BTC: 1 ≈ %v USDT \n", status.Rate.RoundBank(places))
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
		log.WithError(err).Error("command buy error")
		err = fmt.Errorf("command buy error: %w: ", err)
		h.sendMessage(chatID, err.Error())
		return
	}
	h.sendMessage(chatID, "Command Buy finished")
}

func (h *Handler) commandSell(chatID int) {
	err := h.svc.Logic.Sell(&Sender{h.bot, chatID})
	if err != nil {
		log.WithError(err).Error("command sell error")
		err = fmt.Errorf("command sell error: %w: ", err)
		h.sendMessage(chatID, err.Error())
		return
	}
	h.sendMessage(chatID, "Command Sell finished")
}

func (h *Handler) commandDraw(chatID int, str string) {
	b, err := h.svc.Logic.Draw(str, nil)
	if err != nil {
		log.WithField("str", str).WithError(err).Error("command draw error")
		err = fmt.Errorf("command draw error: %w: ", err)
		h.sendMessage(chatID, err.Error())
		return
	}
	h.sendPhoto(chatID, "graph.png", b)
	h.svc.History.AddFormula(str)
}

func (h *Handler) commandBegin(chatID int, str string) {
	err := h.svc.Logic.Begin(&Sender{h.bot, chatID}, str)
	if err != nil {
		log.WithField("str", str).WithError(err).Error("command begin error")
		err = fmt.Errorf("command begin error: %w: ", err)
		h.sendMessage(chatID, err.Error())
		return
	}
	h.svc.History.AddFormula(str)
}

func (h *Handler) commandEnd(chatID int) {
	err := h.svc.Logic.End(&Sender{h.bot, chatID})
	if err != nil {
		log.WithError(err).Error("command end error")
		err = fmt.Errorf("command end error: %w: ", err)
		h.sendMessage(chatID, err.Error())
		return
	}
}

func (h *Handler) commandFstat(chatID int, str string) {
	status := h.svc.Logic.Fstat(str)
	if status.Err != nil {
		err := status.Err
		log.WithField("str", str).WithError(err).Error("command fstat error")
		err = fmt.Errorf("command fstat error: %w: ", err)
		h.sendMessage(chatID, err.Error())
		return
	}
	h.sendMessage(chatID, status.Txt)
	h.sendPhoto(chatID, "graph.png", status.B)
}

func (h *Handler) commandNotFound(chatID int) {
	h.commandHelp(chatID)
}

func (h *Handler) commandTestModeSwitch(chatID int) {
	cur := h.svc.Logic.SafeTestModeSwitch()
	if cur {
		h.sendMessage(chatID, "Testmode is ON!")
	} else {
		h.sendMessage(chatID, "Testmode is OFF!\nNow, be careful")
	}
}

func (h *Handler) commandSetScale(chatID int, str string) {
	h.svc.Logic.SetScale(str)
	txt := fmt.Sprintf("Graph scale set to %v!", str)
	h.sendMessage(chatID, txt)
}

func (h *Handler) commandAlter(chatID int, str string) {
	err := h.svc.Logic.Alter(str)
	if err != nil {
		log.WithField("str", str).WithError(err).Error("command alter error")
		err = fmt.Errorf("command alter error: %w: ", err)
		h.sendMessage(chatID, err.Error())
		return
	}
	txt := fmt.Sprintf("Formula set to %v!", str)
	h.sendMessage(chatID, txt)
	h.svc.History.AddFormula(str)
}

func (h *Handler) commandHistory(chatID int) {
	hist := h.svc.History.GetFormulasList()
	txt := "History:\n\n"
	for _, val := range hist {
		txt += val + "\n"
	}
	h.sendMessage(chatID, txt)
}

func (h *Handler) commandHelp(chatID int) {
	str := `Need some help?

/status				displays BTC/USDT rate and your binance wallet balance

/sell				sells all BTC
/buy				buys BTC with all USDT

/draw <formula> (example: rate-10+0.0002*(now-start)^1.2) 		draws graph of given formula
/begin <formula> buys BTC with all USDT, activates trigger 
/end    deactivates trigger and sells all BTC
/fstat 	draws graph and sends status message 

/testSwitch activates/deactivates test mode (only for begin/end commands, trigger must be deactivated). 
			While test mode is active begin/end commands won't place buy/sell orders

/setScale sets the graph scale (can be 1m, 3m, 5m, 15m, 30m, 1H, 2H, 4H, 6H, 8H, 12H, 1D, 3D, 1W, 1M). Default 15m

/alter <formula> sets the formula in the trigger to a new without changing of the start point
/history sends 10 last used formulas
/opts prints all set options
/opt_set <key> <value> allows to set option by key and value
/opt_help get help with std options
/opt_auto set some options by some template
`

	h.sendMessage(chatID, str)
}

func (h *Handler) commandOpts(chatID int) {
	res, err := h.svc.ChatexOpts.GetAll()
	if err != nil {
		log.WithError(err).Error("failed to read opts")
		h.sendMessage(chatID, err.Error())
		return
	}

	var lines []string
	for k, v := range res {
		lines = append(lines, fmt.Sprintf("%s:  %s", k, v))
	}

	sort.Strings(lines)
	lines = append([]string{"All opts:", ""}, lines...)

	h.sendMessage(chatID, strings.Join(lines, "\n"))
}

func (h *Handler) commandOptSet(chatID int, cmd []string) {
	const args = 2
	if len(cmd) != args {
		h.sendMessage(chatID, "Must be exactly 2 arguments. Get /help")
		return
	}

	key := cmd[0]
	value := cmd[1]

	err := h.svc.ChatexOpts.SetOption(key, value)
	if err != nil {
		log.WithError(err).Error("failed to set option")
		h.sendMessage(chatID, err.Error())
		return
	}

	h.sendMessage(chatID, fmt.Sprintf("OK! [%s] => %s\n\n/opts", key, value))
}

func (h *Handler) commandOptHelp(chatID int) {
	h.sendMessage(chatID, strings.TrimSpace(`
Here are some common options:
chatex.collector.state -- if "disable" is set, OrdersCollector will skip collectAndSave
limit.usdt -- contains the maximum available trade amount for usdt
coins.tbtc.disabled -- if "true", then this coin is ignored
chatex.collector.period -- can set to any duration, like "20s"
`))
}

func (h *Handler) commandOptAuto(chatID int, args []string) {
	const helpMsg = "First argument must be valid type. Examples:\n\t- /opt_auto template_from_rate ref_rate_usd.%s"

	if len(args) < 1 {
		h.sendMessage(chatID, helpMsg)
		return
	}

	templateFromRate := func(tmpl string) {
		rates, err := h.svc.Chatex.GetAllRates(chatex.USDT)
		if err != nil {
			log.WithError(err).Error("failed to get all rates")
			h.sendMessage(chatID, err.Error())
			return
		}

		var info []string

		for _, rate := range rates {
			key := fmt.Sprintf(tmpl, rate.Currency)

			err := h.svc.ChatexOpts.SetOption(key, rate.Rate.String())
			if err != nil {
				info = append(info, "failed to set option: "+err.Error())
				continue
			}

			info = append(info, fmt.Sprintf("%s = %s", key, rate.Rate))
		}

		h.sendMessage(chatID, "All set!\n\n"+strings.Join(info, "\n"))
	}

	switch args[0] {
	case "template_from_rate":
		if len(args) != 2 { //nolint:gomnd
			h.sendMessage(chatID, helpMsg)
			return
		}
		templateFromRate(args[1])

	default:
		h.sendMessage(chatID, helpMsg)
		return
	}
}
