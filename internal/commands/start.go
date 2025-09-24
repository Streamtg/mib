package commands

import (
	"fmt"
	"strings"

	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/utils"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
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
	if len(config.ValueOf.AllowedUsers) != 0 && !utils.Contains(config.ValueOf.AllowedUsers, chatId) {
		ctx.Reply(u, "You are not allowed to use this bot.", nil)
		return dispatcher.EndGroups
	}

	// FORZAR SUSCRIPCIÓN A MULTIPLES CANALES
	if config.ValueOf.ForceSubChannels != "" {
		channels := strings.Split(config.ValueOf.ForceSubChannels, ",")
		notJoined := []string{}

		for _, ch := range channels {
			ch = strings.TrimSpace(ch)
			isSubscribed, err := utils.IsUserSubscribed(ctx, ctx.Raw, ctx.PeerStorage, chatId)
			if err != nil || !isSubscribed {
				notJoined = append(notJoined, ch)
			}
		}

		if len(notJoined) > 0 {
			rows := []tg.KeyboardButtonRow{}
			for _, ch := range notJoined {
				row := tg.KeyboardButtonRow{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonURL{
							Text: fmt.Sprintf("Join @%s", ch),
							URL:  fmt.Sprintf("https://t.me/%s", ch),
						},
					},
				}
				rows = append(rows, row)
			}
			markup := &tg.ReplyInlineMarkup{Rows: rows}
			ctx.Reply(u, "Please join all required channels to use this bot:", &ext.ReplyOpts{Markup: markup})
			return dispatcher.EndGroups
		}
	}

	// MENSAJE PRINCIPAL
	message := `Hello! 👋 I'm your file-sharing assistant.

📂 Send or forward me any file (in any format!) and I'll instantly give you a direct link to download or view online. ⚡

💡 You can also use this bot as a *host* for movie and series channels, etc. 🎬

How to get started?

1️⃣ Send or forward me a file
2️⃣ Wait a few seconds ⏱️
3️⃣ Receive your link 🚀

🎬 Follow our movies and series channels

🇺🇸 English Movies
https://t.me/moviegxg

🇲🇽 Películas en español Latino
https://t.me/peligxg

Official channel: @yoelbots

💡 To view bot statistics, type /stats 📊`

	ctx.Reply(u, message, nil)
	return dispatcher.EndGroups
}
