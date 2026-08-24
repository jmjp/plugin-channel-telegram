package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"fluxa/pkg/plugin"
)

// Update contém os campos do Telegram relevantes para normalização.
type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message"`
}
type Message struct {
	MessageID int64   `json:"message_id"`
	Text      string  `json:"text"`
	Chat      Chat    `json:"chat"`
	From      User    `json:"from"`
	Photo     []Photo `json:"photo"`
}
type Chat struct {
	ID int64 `json:"id"`
}
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
type Photo struct {
	FileID string `json:"file_id"`
}

// Normalize converte uma atualização Telegram em contrato agnóstico do core.
func Normalize(update Update) (plugin.InboundMessage, bool) {
	if update.Message == nil || (update.Message.Text == "" && len(update.Message.Photo) == 0) {
		return plugin.InboundMessage{}, false
	}
	metadata := map[string]string{"telegram_message_id": strconv.FormatInt(update.Message.MessageID, 10), "telegram_user_id": strconv.FormatInt(update.Message.From.ID, 10)}
	if update.Message.From.Username != "" {
		metadata["username"] = update.Message.From.Username
	}
	if len(update.Message.Photo) > 0 {
		metadata["telegram_photo_file_id"] = update.Message.Photo[len(update.Message.Photo)-1].FileID
	}
	return plugin.InboundMessage{Channel: "telegram", CustomerRef: strconv.FormatInt(update.Message.Chat.ID, 10), Text: update.Message.Text, Timestamp: time.Now(), Metadata: metadata}, true
}

// Poll busca atualizações e delega somente mensagens normalizadas ao callback.
func Poll(ctx context.Context, client *http.Client, token string, handle func(context.Context, plugin.InboundMessage) error) error {
	if token == "" {
		return fmt.Errorf("telegram: token obrigatório")
	}
	if client == nil {
		client = http.DefaultClient
	}
	var offset int64
	for ctx.Err() == nil {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.telegram.org/bot"+token+"/getUpdates?timeout=25&offset="+strconv.FormatInt(offset, 10), nil)
		if err != nil {
			return fmt.Errorf("telegram: criando poll: %w", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("telegram: getUpdates: %w", err)
		}
		var body struct {
			OK     bool     `json:"ok"`
			Result []Update `json:"result"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("telegram: decodificando updates: %w", err)
		}
		for _, update := range body.Result {
			offset = update.UpdateID + 1
			if msg, ok := Normalize(update); ok {
				if err := handle(ctx, msg); err != nil {
					return fmt.Errorf("telegram: entregando mensagem: %w", err)
				}
			}
		}
	}
	return nil
}
