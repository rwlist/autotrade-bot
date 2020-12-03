package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	chatexsdk "github.com/chatex-com/sdk-go"

	"github.com/rwlist/autotrade-bot/pkg/exproc"
	"github.com/rwlist/autotrade-bot/pkg/history"
	"github.com/rwlist/autotrade-bot/pkg/store/redisdb"
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

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":2112", mux)
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("prometheus server error")
		}
	}()

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

	adminSender := app.NewSender(bot, cfg.Bot.AdminID)

	ch, err := updates.StartPolling(bot, telegram.GetUpdatesRequest{
		Offset:  0,
		Limit:   50,
		Timeout: 10,
	})
	if err != nil {
		log.WithError(err).Fatal("in updates.StartPolling()")
	}

	redisDB := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})

	var binanceCli binance.Client
	binanceCli = binance.NewClientDefault(gobinance.NewClient(cfg.Binance.APIKey, cfg.Binance.Secret))
	if cfg.Binance.Debug {
		binanceCli = binance.NewClientLog(binanceCli)
	}

	myBinance := binance.NewBinance(binanceCli)

	tr := trigger.NewTrigger(myBinance)

	chatexOpts := chatex.NewTradeOpts(
		redisdb.NewHash("trade_opts:chatex", redisDB),
	)
	chatexCli := chatexsdk.NewClient("https://api.chatex.com/v1", cfg.Chatex.RefreshToken)
	chatexSnapshotList := redisdb.NewList("chatex_order_snapshots", redisDB)
	ordersCollector := chatex.NewOrdersCollector(chatexCli, chatexSnapshotList, chatexOpts)

	go func() {
		err := ordersCollector.CollectInf(context.Background())
		if err != nil && err != context.Canceled {
			log.WithError(err).Fatal("collect inf finished")
		}
	}()

	exFinder := exproc.NewFinder(chatexCli, ordersCollector, chatexOpts, adminSender)
	ordersCollector.RegisterCallback(exFinder.OnSnapshot)

	myChatex := chatex.NewChatex(chatexCli, ordersCollector)

	handler := app.NewHandler(
		bot,
		cfg,
		app.Services{
			Logic:        logic.NewLogic(myBinance, &tr, cfg.Bot.IsTest),
			Status:       stat.New(myBinance),
			StatusChatex: stat.New(myChatex),
			History:      history.New(),
			ChatexOpts:   chatexOpts,
		},
	)

	for upd := range ch {
		upd := upd
		handler.Handle(&upd)
	}
}
