package main

import (
	"encoding/json"
	"log"

	"github.com/rwlist/autotrade-bot/logic"

	"github.com/rwlist/autotrade-bot/app/stat"

	"github.com/petuhovskiy/telegram"
	"github.com/petuhovskiy/telegram/updates"
	"github.com/rwlist/autotrade-bot/app"
	"github.com/rwlist/autotrade-bot/pkg/conf"
	"github.com/rwlist/autotrade-bot/trade/binance"
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

	myBinance := binance.NewBinance(cfg.Binance, cfg.Binance.Debug)

	handler := app.NewHandler(
		bot,
		cfg,
		app.Services{
			Logic:  logic.NewLogic(myBinance),
			Status: stat.New(myBinance),
		},
	)

	for upd := range ch {
		upd := upd
		handler.Handle(&upd)
	}
}
