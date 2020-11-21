package main

import (
	"encoding/json"

	chatexsdk "github.com/chatex-com/sdk-go"

	"github.com/rwlist/autotrade-bot/pkg/history"
	"github.com/rwlist/autotrade-bot/pkg/trade/chatex"

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

	var binanceCli binance.Client
	binanceCli = binance.NewClientDefault(gobinance.NewClient(cfg.Binance.APIKey, cfg.Binance.Secret))
	if cfg.Binance.Debug {
		binanceCli = binance.NewClientLog(binanceCli)
	}

	myBinance := binance.NewBinance(binanceCli)

	tr := trigger.NewTrigger(myBinance)

	chatexCli := chatexsdk.NewClient("https://api.chatex.com/v1", cfg.Chatex.RefreshToken)
	myChatex := chatex.NewChatex(chatexCli)

	handler := app.NewHandler(
		bot,
		cfg,
		app.Services{
			Logic:        logic.NewLogic(myBinance, &tr, cfg.Bot.IsTest),
			Status:       stat.New(myBinance),
			StatusChatex: stat.New(myChatex),
			History:      history.New(),
		},
	)

	for upd := range ch {
		upd := upd
		handler.Handle(&upd)
	}
}
