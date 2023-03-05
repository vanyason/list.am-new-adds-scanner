package main

import (
	"log"

	old "github.com/vanyason/list.am-new-adds-scanner/deprecated/go/lib"
)

func testTgBot() error {
	log.Println("Testing Tg bot started ...")

	bot, err := old.CreateBot("config/testbot_config.json")
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
