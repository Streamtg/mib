package commands

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/cache"
	"EverythingSuckz/fsb/internal/utils"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
	"github.com/celestix/gotgproto/types"
	"github.com/gotd/td/tg"
)

func (m *command) LoadStream(dispatcher dispatcher.Dispatcher) {
	defer m.log.Sugar().Info("Loaded")
	dispatcher.AddHandler(
		handlers.NewMessage(nil, sendLink),
	)
}

func supportedMediaFilter(m *types.Message) (bool, error) {
	if m.Media == nil {
		return false, dispatcher.EndGroups
	}
	switch m.Media.(type) {
	case *tg.MessageMediaDocument:
		return true, nil
	case *tg.MessageMediaPhoto:
		return true, nil
	case tg.MessageMediaClass:
		return false, dispatcher.EndGroups
	default:
		return false, nil
	}
}

// Convierte bytes a tamaño legible
func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	default:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	}
}

// Emoji según tipo de archivo
func fileTypeEmoji(mime string) string {
	lowerMime := strings.ToLower(mime)
	switch {
	case strings.Contains(lowerMime, "video"):
		return "🎬"
	case strings.Contains(lowerMime, "image"):
		return "🖼️"
	case strings.Contains(lowerMime, "audio"):
		return "🎵"
	case strings.Contains(lowerMime, "pdf"):
		return "📕"
	case strings.Contains(lowerMime, "zip"), strings.Contains(lowerMime, "rar"):
		return "🗜️"
	case strings.Contains(lowerMime, "text"):
		return "📝"
	case strings.Contains(lowerMime, "application"):
		return "📄"
	default:
		return "📄"
	}
}

func sendLink(ctx *ext.Context, u *ext.Update) error {
	chatId := u.EffectiveChat().GetID()
	peerChatId := ctx.PeerStorage.GetPeerById(chatId)
	if peerChatId.Type != int(storage.TypeUser) {
		return dispatcher.EndGroups
	}

	if len(config.ValueOf.AllowedUsers) != 0 && !utils.Contains(config.ValueOf.AllowedUsers, chatId) {
		ctx.Reply(u, "You are not allowed to use this bot.", nil)
		return dispatcher.EndGroups
	}

	// ---- FORZAR SUSCRIPCIÓN REAL A VARIOS CANALES ----
	forceChannels := []string{"yoelbots", "pelisgxg"}
	subscribed := true

	for _, ch := range forceChannels {
		// Verificar si el usuario es participante del canal
		_, err := ctx.Raw.ChannelsGetParticipant(context.Background(), &tg.ChannelsGetParticipantRequest{
			Channel: &tg.InputChannelUsername{Username: ch},
			UserId:  &tg.InputUser{UserID: chatId},
		})
		if err != nil {
			subscribed = false
			break
		}
	}

	// --- Construir botones de canales ---
	var rows []tg.KeyboardButtonRow
	for _, ch := range forceChannels {
		text := ch
		if ch == "pelisgxg" {
			text = "🎬 Películas y Series en Español"
		} else {
			text = "📢 @" + ch
		}
		row := tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				&tg.KeyboardButtonURL{
					Text: text,
					URL:  fmt.Sprintf("https://t.me/%s", ch),
				},
			},
		}
		rows = append(rows, row)
	}

	if !subscribed {
		ctx.Reply(u,
			"🚨 Para usar este bot debes unirte a nuestros canales obligatorios:\n\n"+
				"Después de unirte, reenvía tu archivo otra vez. ✅",
			&ext.ReplyOpts{Markup: &tg.ReplyInlineMarkup{Rows: rows}},
		)
		return dispatcher.EndGroups
	}

	// --- VALIDAR TIPO DE ARCHIVO ---
	supported, err := supportedMediaFilter(u.EffectiveMessage)
	if err != nil {
		return err
	}
	if !supported {
		ctx.Reply(u, "Sorry, this message type is unsupported.", nil)
		return dispatcher.EndGroups
	}

	// --- REENVÍO AL CANAL DE LOGS ---
	update, err := utils.ForwardMessages(ctx, chatId, config.ValueOf.LogChannelID, u.EffectiveMessage.ID)
	if err != nil {
		ctx.Reply(u, fmt.Sprintf("Error - %s", err.Error()), nil)
		return dispatcher.EndGroups
	}

	messageID := update.Updates[0].(*tg.UpdateMessageID).ID
	doc := update.Updates[1].(*tg.UpdateNewChannelMessage).Message.(*tg.Message).Media
	file, err := utils.FileFromMedia(doc)
	if err != nil {
		ctx.Reply(u, fmt.Sprintf("Error - %s", err.Error()), nil)
		return dispatcher.EndGroups
	}

	// --- DETECTAR NOMBRE Y EXTENSIÓN ---
	if file.FileName == "" {
		var ext string
		lowerMime := strings.ToLower(file.MimeType)
		switch {
		case strings.Contains(lowerMime, "image/jpeg"):
			ext = ".jpg"
			file.FileName = "photo" + ext
		case strings.Contains(lowerMime, "image/png"):
			ext = ".png"
			file.FileName = "photo" + ext
		case strings.Contains(lowerMime, "image/gif"):
			ext = ".gif"
			file.FileName = "animation" + ext
		case strings.Contains(lowerMime, "video"):
			ext = ".mp4"
			file.FileName = "video" + ext
		case strings.Contains(lowerMime, "audio"):
			ext = ".mp3"
			file.FileName = "audio" + ext
		case strings.Contains(lowerMime, "pdf"):
			ext = ".pdf"
			file.FileName = "document" + ext
		case strings.Contains(lowerMime, "zip"):
			ext = ".zip"
			file.FileName = "archive" + ext
		case strings.Contains(lowerMime, "rar"):
			ext = ".rar"
			file.FileName = "archive" + ext
		case strings.Contains(lowerMime, "text"):
			ext = ".txt"
			file.FileName = "text" + ext
		case strings.Contains(lowerMime, "application"):
			ext = ".bin"
			file.FileName = "file" + ext
		default:
			file.FileName = "unknown"
		}
	}

	// --- MENSAJE VISUAL ---
	emoji := fileTypeEmoji(file.MimeType)
	size := formatFileSize(file.FileSize)
	message := fmt.Sprintf(
		"%s File Name: %s\n\n%s File Type: %s\n\n💾 Size: %s\n\n⏳ @yoelbots",
		emoji, file.FileName,
		emoji, file.MimeType,
		size,
	)

	fullHash := utils.PackFile(file.FileName, file.FileSize, file.MimeType, file.ID)
	hash := utils.GetShortHash(fullHash)

	statsCache := cache.GetStatsCache()
	if statsCache != nil {
		_ = statsCache.RecordFileProcessed(file.FileSize)
	}

	// --- BOTÓN STREAMING / DOWNLOAD ---
	videoParam := fmt.Sprintf("%d?hash=%s", messageID, hash)
	encodedVideoParam := url.QueryEscape(videoParam)
	encodedFilename := url.QueryEscape(file.FileName)
	streamURL := fmt.Sprintf("https://file.streamgramm.workers.dev/?video=%s&filename=%s", encodedVideoParam, encodedFilename)

	streamRow := tg.KeyboardButtonRow{
		Buttons: []tg.KeyboardButtonClass{
			&tg.KeyboardButtonURL{
				Text: "Streaming / Download",
				URL:  streamURL,
			},
		},
	}
	rows = append(rows, streamRow)

	markup := &tg.ReplyInlineMarkup{Rows: rows}

	_, err = ctx.Reply(u, message, &ext.ReplyOpts{
		Markup:           markup,
		NoWebpage:        false,
		ReplyToMessageId: u.EffectiveMessage.ID,
	})
	if err != nil {
		ctx.Reply(u, fmt.Sprintf("Error - %s", err.Error()), nil)
	}

	return dispatcher.EndGroups
}
