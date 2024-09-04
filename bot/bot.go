package bot

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
	"os/signal"
	"tg-backend/config"
)

type TelegramBot struct {
	telegramBot *bot.Bot
	miniAppUrl  string
}

func NewTelegramBot(tgConf config.TgConfig) *TelegramBot {
	opts := []bot.Option{
		bot.WithDefaultHandler(getDefalutHandler(tgConf.TelegramMiniAppUrl)),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, callbackHandler),
	}
	if tgConf.Test {
		opts = append(opts, bot.UseTestEnvironment())
	}

	b, err := bot.New(tgConf.TelegramBotToken, opts...)
	if err != nil {
		panic(err)
	}

	return &TelegramBot{
		telegramBot: b,
		miniAppUrl:  tgConf.TelegramMiniAppUrl,
	}
}

func callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// answering callback query first to let Telegram know that we received the callback query,
	// and we're handling it. Otherwise, Telegram might retry sending the update repetitively
	// as it thinks the callback query doesn't reach to our application. learn more by
	// reading the footnote of the https://core.telegram.org/bots/api#callbackquery type.
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
}

func getDefalutHandler(appUrl string) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "mini app",
						CallbackData: "button_3",
						URL:          appUrl},
				},
			},
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "This is a test mini app",
			ReplyMarkup: kb,
		})
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func (t *TelegramBot) Start() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	t.telegramBot.Start(ctx)
}
