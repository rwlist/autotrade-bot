package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/rwlist/autotrade-bot/pkg/trigger"

	"github.com/rwlist/autotrade-bot/pkg/logic"

	"github.com/rwlist/autotrade-bot/pkg/stat"

	"github.com/petuhovskiy/telegram"
	"github.com/petuhovskiy/telegram/updates"
	"github.com/rwlist/autotrade-bot/pkg/app"
	"github.com/rwlist/autotrade-bot/pkg/conf"
	"github.com/rwlist/autotrade-bot/pkg/trade/binance"

	gobinance "github.com/adshao/go-binance"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	cfg, err := conf.ParseEnv()
	if err != nil {
		log.WithError(err).Fatal("in conf.ParseEnv()")
	}

	log.SetFormatter(&log.JSONFormatter{PrettyPrint: cfg.Bot.PrettyPrint})

	bot := telegram.NewBotWithOpts(cfg.Bot.Token, &telegram.Opts{
		Middleware: func(handler telegram.RequestHandler) telegram.RequestHandler {
			return func(methodName string, req interface{}) (message json.RawMessage, err error) {
				res, err := handler(methodName, req)
				if err != nil {
					log.WithError(err).Error("telegram response error")
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
		log.WithError(err).Fatal("in updates.StartPolling()")
	}

	var cli binance.Client
	cli = binance.NewClientDefault(gobinance.NewClient(cfg.Binance.APIKey, cfg.Binance.Secret))
	if cfg.Binance.Debug {
		cli = binance.NewClientLog(cli)
	}

	myBinance := binance.NewBinance(cli)

	tr, err := trigger.NewTrigger(myBinance)
	if err != nil {
		log.WithError(err).Fatal("in trigger.NewTrigger")
	}

	handler := app.NewHandler(
		bot,
		cfg,
		app.Services{
			Logic:  logic.NewLogic(&myBinance, &tr, cfg.Bot.IsTest),
			Status: stat.New(myBinance),
		},
	)

	for upd := range ch {
		upd := upd
		handler.Handle(&upd)
	}
}
