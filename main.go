package main

import (
	"encoding/json"
	"log"

	"github.com/petuhovskiy/telegram"
	"github.com/petuhovskiy/telegram/updates"

	"github.com/rwlist/autotrade-bot/app"
	"github.com/rwlist/autotrade-bot/binance"
	"github.com/rwlist/autotrade-bot/conf"
)

func main() {
	cfg, err := conf.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	bot := telegram.NewBotWithOpts(cfg.Bot.Token, &telegram.Opts{
		Middleware: func(handler telegram.RequestHandler) telegram.RequestHandler {
			return func(methodName string, req interface{}) (message json.RawMessage, err error) {
				res, err := handler(methodName, req)
				if err != nil {
					log.Println("Telegram response error: ", err)
				}

				return res, err
			}
		},
	})

	ch, err := updates.StartPolling(bot, telegram.GetUpdatesRequest{
		Offset:  0,
		Limit:   50,
		Timeout: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	logic := app.NewLogic(binance.NewMyBinance(cfg.Binance, true))
	handler := app.NewHandler(bot, logic, cfg)

	for upd := range ch {
		handler.Handle(&upd)
	}
}
