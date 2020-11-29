# autotrade-bot
Simple telegram bot for trading semi-automatization

## Basic functionality ideas

Bot should support entering and exiting the position, displaying current status and triggering exit by some fancy stop-loss function.

Of course, the bot must be accessible only by me.

## Config

Place config in `.env` file:

```..env
ADMIN_TELEGRAM_ID=1
BOT_TOKEN=2:A
BINANCE_API_KEY=Jw
BINANCE_SECRET=v
CHATEX_REFRESH_TOKEN=a
REDIS_PASSWORD=abc
REDIS_ADDR=10.10.10.10:6379
```
