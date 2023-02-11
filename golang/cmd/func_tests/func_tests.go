package main

import (
	"log"

	"github.com/vanyason/list.am-new-adds-scanner/internal"
)

func testTgBot() error {
	log.Println("Testing Tg bot started ...")

	bot, err := internal.CreateBot("config/bot_config.json")
	if err != nil {
		return err
	}

	err = bot.SendMessageSilently("test tg bot")
	if err != nil {
		return err
	}

	log.Println("Testing Tg bot succeeded")
	return nil
}

func main() {
	if err := testTgBot(); err != nil {
		log.Printf("test tg bot failed with : %s", err)
	}
}
