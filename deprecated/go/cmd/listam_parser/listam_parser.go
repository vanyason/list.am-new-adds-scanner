package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	old "github.com/vanyason/list.am-new-adds-scanner/deprecated/go/lib"
	"golang.org/x/exp/maps"
)

func run(args old.CmdArguments, scratchData old.ScratchData, bot old.TgBot) error {
	/* Measure time */
	start := time.Now()
	timer := func() { log.Printf("Round took ~ %f seconds", time.Since(start).Seconds()) }
	defer timer()

	/* Get pages from Listam */
	parsedPages, err := old.ScratchHtmlPages(scratchData)
	if err != nil {
		return fmt.Errorf("error getting html pages : %w", err)
	}

	/*
	 * If there is no saved json :
	 *    - save and continue
	 * Else :
	 *    - read old one
	 *    - compare / get the diffs
	 *    - replace old one with the new one (ADD NEW VALUES TO THE OLD)
	 *	  - notify
	 */
	_, err = os.Stat(args.DBFileName)
	notExist := os.IsNotExist(err)
	if err != nil && !notExist {
		return fmt.Errorf("error checking file existence (%s) : %w", args.DBFileName, err)
	}

	if notExist {
		newJson, err := json.Marshal(parsedPages)
		if err != nil {
			return fmt.Errorf("error saving parsed pages to json : %w", err)
		}

		err = os.WriteFile(args.DBFileName, newJson, 0644)
		if err != nil {
			return fmt.Errorf("error saving new json : %w", err)
		}

		log.Printf("File (%s) created. Continue", args.DBFileName)
	} else {
		oldJson, err := os.ReadFile(args.DBFileName)
		if err != nil {
			return fmt.Errorf("error reading existing json : %w", err)
		}

		oldPages := make(map[string]string)
		err = json.Unmarshal(oldJson, &oldPages)
		if err != nil {
			return fmt.Errorf("error parsing existing json : %w", err)
		}

		diffs := old.Compare(oldPages, parsedPages)

		maps.Copy(oldPages, parsedPages)

		newJson, err := json.Marshal(oldPages)
		if err != nil {
			return fmt.Errorf("error saving parsed pages + old pages to json : %w", err)
		}

		err = os.WriteFile(args.DBFileName, newJson, 0644)
		if err != nil {
			return fmt.Errorf("error rewriting old json with new : %w", err)
		}

		if len(diffs) != 0 {
			for _, diff := range diffs {
				log.Printf("new add : %s", diff)
				if err := bot.SendMessageSilently(diff); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func main() {
	/* create log folder */
	folderPath := "log"
	if _, err := os.Stat(folderPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(folderPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	/* Parse cmd args */
	cmdArgs, err := old.ParseCmdLineArgs()
	if err != nil {
		log.Fatalf("error parsing command line args : %s", err)
	}

	/* Setup logging */
	logFile, err := os.OpenFile(cmdArgs.LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error creating / opening log file : %s", err)
	}
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	/* Generate params for web scratching */
	scratchData := old.GenerateScratchData(cmdArgs)

	/* Create tg bot */
	bot, err := old.CreateBot("config/bot_config.json")
	if err != nil {
		log.Fatalf("error creating tg bot : %s", err)
	}

	/* Start app */
	log.Printf("\nExecution started...\nParams :\n%+v\nUrls :\n%s\n%s\n%s\n",
		cmdArgs,
		scratchData.ApartmentsIterateUrlDram(1),
		scratchData.HousesIterateUrlDram(1),
		scratchData.TownhousesIterateUrlDram((1)))

	/* Loop */
	errCounter := 0
	for {
		if err := run(cmdArgs, scratchData, bot); err != nil {
			log.Println(err)
			errCounter++
		}
		if cmdArgs.ErrorCounter != 0 && errCounter >= int(cmdArgs.ErrorCounter) {
			break
		}
		time.Sleep(time.Duration(cmdArgs.LoopPause) * time.Minute)
	}

	/* Notify that we are broke */
	const finalMsg = "Time to check the logs. We are broken"
	log.Println(finalMsg)
	_ = bot.SendMessageSilently(finalMsg)
}
