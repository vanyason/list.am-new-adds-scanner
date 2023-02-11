package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type TgBot struct {
	Token  string `json:"tg_token"`
	ChatId string `json:"chat_id"`
}

func CreateBot(configPath string) (bot TgBot, err error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return bot, fmt.Errorf("error reading config json at path %s : %w", configPath, err)
	}

	err = json.Unmarshal([]byte(file), &bot)
	if err != nil {
		return bot, fmt.Errorf("error parsing config json at path %s : %w", configPath, err)
	}

	fmt.Println(bot)

	return bot, nil
}

func (bot *TgBot) SendMessageSilently(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", bot.Token, bot.ChatId, url.QueryEscape(message))

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error sending tg message (%s) : %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error sending tg message. Reply. Code: %d Body: %s", resp.StatusCode, string(b))
	}

	return nil
}
