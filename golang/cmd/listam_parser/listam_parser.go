package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/vanyason/list.am-new-adds-scanner/internal"
)

func run(args internal.CmdArguments, scratchData internal.ScratchData) error {
	/* Measure time */
	start := time.Now()

	/* Get pages from Listam */
	parsedPages, err := internal.ScratchHtmlPages(scratchData)
	if err != nil {
		return fmt.Errorf("error getting html pages : %w", err)
	}

	/* Save map to a json */
	newJson, err := json.Marshal(parsedPages)
	if err != nil {
		return fmt.Errorf("error saving parsed pages to json : %w", err)
	}

	/*
	 * If there is no saved json :
	 *    - save and continue
	 * Else :
	 *    - read old one
	 *    - compare / get the diffs
	 *    - replace old one with the new one
	 *	  - notify
	 */
	_, err = os.Stat(args.DBFileName)
	notExist := os.IsNotExist(err)
	if err != nil && !notExist {
		return fmt.Errorf("error checking file existence (%s) : %w", args.DBFileName, err)
	}

	if notExist {
		err := ioutil.WriteFile(args.DBFileName, newJson, 0644)
		if err != nil {
			return fmt.Errorf("error saving new json : %w", err)
		}

		log.Printf("File (%s) created. Continue", args.DBFileName)
	} else {
		oldJson, err := ioutil.ReadFile(args.DBFileName)
		if err != nil {
			return fmt.Errorf("error reading existing json : %w", err)
		}

		oldPages := make(map[string]string)
		err = json.Unmarshal(oldJson, &oldPages)
		if err != nil {
			return fmt.Errorf("error parsing existing json : %w", err)
		}

		diffs := internal.Compare(oldPages, parsedPages)

		err = ioutil.WriteFile(args.DBFileName, newJson, 0644)
		if err != nil {
			return fmt.Errorf("error rewriting old json with new : %w", err)
		}

		if len(diffs) == 0 {
			log.Println("diffs not found")
		} else {
			log.Printf("diffs : %v", diffs)
			log.Printf("notifying")
		}
	}

	log.Printf("Loop took ~ %f seconds", time.Since(start).Seconds())
	return nil
}

func main() {
	/* Parse cmd args */
	cmdArgs, err := internal.ParseCmdLineArgs()
	if err != nil {
		log.Fatalf("error parsing command line args : %s", err)
	}

	/* Setup loging */
	logFile, err := os.OpenFile(cmdArgs.LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error creating / opening log file : %s", err)
	}
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	/* Generate params for web scratching */
	scratchData := internal.GenerateScratchData(cmdArgs)

	/* Start app */
	log.Printf("\nExecution started...\nParams :\n%+v\nUrls :\n%s\n%s\n%s\n",
		cmdArgs,
		scratchData.ApartmentsIterateUrlDram(1),
		scratchData.HousesIterateUrlDram(1),
		scratchData.TownhousesIterateUrlDram((1)))

	if err := run(cmdArgs, scratchData); err != nil {
		log.Fatalln(err)
	}
}
