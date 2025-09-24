package commands

import (
	"fmt"
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/utils"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/tg"
)

func (m *command) LoadStart(dispatcher dispatcher.Dispatcher) {
	log := m.log.Named("start")
	defer log.Sugar().Info("Loaded")
	dispatcher.AddHandler(handlers.NewCommand("start", start))
}

func start(ctx *ext.Context, u *ext.Update) error {
	chatId := u.EffectiveChat().GetID()
	peerChatId := ctx.PeerStorage.GetPeerById(chatId)
	if peerChatId.Type != int(storage.TypeUser) {
		return dispatcher.EndGroups
	}

	// Verificar AllowedUsers
	if len(config.ValueOf.AllowedUsers) != 0 && !utils.Contains(config.ValueOf.AllowedUsers, chatId) {
		ctx.Reply(u, "You are not allowed to use this bot.", nil)
		return dispatcher.EndGroups
	}

	// Forzar suscripción
	if config.ValueOf.ForceSubChannel != "" {
		isSubscribed, err := utils.IsUserSubscribed(ctx, ctx.Raw, ctx.PeerStorage, chatId)
		if err != nil || !isSubscribed {
			row := tg.KeyboardButtonRow{
				Buttons: []tg.KeyboardButtonClass{
					&tg.KeyboardButtonURL{
						Text: "Join @yoelbots",
						URL:  fmt.Sprintf("https://t.me/%s", config.ValueOf.ForceSubChannel),
					},
				},
			}
			markup := &tg.ReplyInlineMarkup{Rows: []tg.KeyboardButtonRow{row}}
			ctx.Reply(u, "Please join our official channel @yoelbots to use the bot.", &ext.ReplyOpts{Markup: markup})
			return dispatcher.EndGroups
		}
	}

	// Mensaje principal
	message := `Hello! 👋 I'm your file-sharing assistant.

📂 Send or forward me any file (in any format!) and I'll instantly give you a direct link to download or view online. ⚡

💡 You can also use this bot as a *host* for movie and series channels, etc. 🎬

How to get started?

1️⃣ Send or forward me a file
2️⃣ Wait a few seconds ⏱️
3️⃣ Receive your link 🚀

🎬 Follow our movies and series channels

Official channel: @yoelbots

💡 To view bot statistics, type /stats 📊`

	// Botones de canales
	row1 := tg.KeyboardButtonRow{
		Buttons: []tg.KeyboardButtonClass{
			&tg.KeyboardButtonURL{
				Text: "🇺🇸 English Movies",
				URL:  "https://t.me/moviegxg",
			},
			&tg.KeyboardButtonURL{
				Text: "🇲🇽 Películas en español Latino",
				URL:  "https://t.me/peligxg",
			},
		},
	}
	row2 := tg.KeyboardButtonRow{
		Buttons: []tg.KeyboardButtonClass{
			&tg.KeyboardButtonURL{
				Text: "Official channel @yoelbots",
				URL:  "https://t.me/yoelbots",
			},
		},
	}
	markup := &tg.ReplyInlineMarkup{Rows: []tg.KeyboardButtonRow{row1, row2}}

	ctx.Reply(u, message, &ext.ReplyOpts{Markup: markup})
	return dispatcher.EndGroups
}
