package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

type CmdArguments struct {
	Price        uint `short:"p" long:"price" description:"price in drams" required:"true"`
	Rooms        uint `short:"r" long:"rooms" description:"amount of rooms. 0 means do not care" default:"0"`
	ErrorCounter uint `short:"e" long:"errorcounter" description:"amount of errors to stop execution. 0 means never" default:"15"`
}

func parseCmdLineArgs() (CmdArguments, error) {
	var args CmdArguments
	if _, err := flags.NewParser(&args, flags.HelpFlag|flags.PassDoubleDash).Parse(); err != nil {
		return args, err
	}

	const maxPrice uint = 10000000
	const maxRooms uint = 10
	const maxErrors uint = 50

	if args.Price > maxPrice {
		return CmdArguments{}, fmt.Errorf("invalid price: %d. max: %d", args.Price, maxPrice)
	}
	if args.Rooms > maxRooms {
		return CmdArguments{}, fmt.Errorf("invalid rooms amount counter: %d. max: %d", args.Rooms, maxRooms)
	}
	if args.ErrorCounter > maxErrors {
		return CmdArguments{}, fmt.Errorf("invalid error counter: %d. max: %d", args.ErrorCounter, maxErrors)
	}

	return args, nil
}

type ScratchData struct {
	Header                   string
	ApartmentsIterateUrlDram func(page int) string
	ApartmentsStopWord       string
	TownhousesIterateUrlDram func(page int) string
	HousesIterateUrlDram     func(page int) string
	HousesStopWord           string
}

func generateScratchData(args CmdArguments) ScratchData {
	sd := ScratchData{
		Header:             "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
		ApartmentsStopWord: "/category/56",
		HousesStopWord:     "/category/63",
	}

	const listamUrl = "https://www.list.am"
	const apartmentsPath = "/category/56/"
	const housesPath = "/category/63/"

	if args.Rooms != 0 {
		sd.ApartmentsIterateUrlDram = func(page int) string {
			return fmt.Sprintf("%s%s%d?cmtype=1&pfreq=1&po=1&n=1&price2=%d&crc=0&_a4=%d", listamUrl, apartmentsPath, page, args.Price, args.Rooms)
		}
		sd.TownhousesIterateUrlDram = func(page int) string {
			return fmt.Sprintf("%s%s%d?cmtype=1&pfreq=1&po=1&n=1&price2=%d&crc=0&_a4=%d&sid=366", listamUrl, housesPath, page, args.Price, args.Rooms)
		}
		sd.HousesIterateUrlDram = func(page int) string {
			return fmt.Sprintf("%s%s%d?cmtype=1&pfreq=1&po=1&n=1&price2=%d&crc=0&_a4=%d&sid=365", listamUrl, housesPath, page, args.Price, args.Rooms)
		}
	} else {
		sd.ApartmentsIterateUrlDram = func(page int) string {
			return fmt.Sprintf("%s%s%d?cmtype=1&pfreq=1&po=1&n=1&price2=%d&crc=0", listamUrl, apartmentsPath, page, args.Price)
		}
		sd.TownhousesIterateUrlDram = func(page int) string {
			return fmt.Sprintf("%s%s%d?cmtype=1&pfreq=1&po=1&n=1&price2=%d&crc=0&&sid=366", listamUrl, housesPath, page, args.Price)
		}
		sd.HousesIterateUrlDram = func(page int) string {
			return fmt.Sprintf("%s%s%d?cmtype=1&pfreq=1&po=1&n=1&price2=%d&crc=0&sid=365", listamUrl, housesPath, page, args.Price)
		}
	}
	return sd
}

func getHtmlPages(sd ScratchData) ([]string, error) {
	getPagesPerCategory := func(urlGenerator func(page int) string, stopWord string) ([]string, error) {
		client := http.Client{Timeout: time.Duration(10) * time.Second}
		var pages []string

		for page := 1; true; page++ {
			req, err := http.NewRequest("GET", urlGenerator(page), nil)
			if err != nil {
				return nil, fmt.Errorf("error generating request for http.client : %w", err)
			}

			req.Header.Set("User-Agent", sd.Header)

			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("error executing request for http.client : %w ; status code : %d", err, resp.StatusCode)
			}
			defer resp.Body.Close()

			if resp.Request.URL.Path == stopWord {
				break
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			pages = append(pages, buf.String())
		}

		return pages, nil
	}

	var allPages []string
	pagesPerCategory, err := getPagesPerCategory(sd.ApartmentsIterateUrlDram, sd.ApartmentsStopWord)
	if err != nil {
		return nil, err
	}
	allPages = append(allPages, pagesPerCategory...)

	pagesPerCategory, err = getPagesPerCategory(sd.TownhousesIterateUrlDram, sd.HousesStopWord)
	if err != nil {
		return nil, err
	}
	allPages = append(allPages, pagesPerCategory...)

	pagesPerCategory, err = getPagesPerCategory(sd.HousesIterateUrlDram, sd.HousesStopWord)
	if err != nil {
		return nil, err
	}
	allPages = append(allPages, pagesPerCategory...)

	return allPages, nil
}

func run(args []string) error {
	cmdArgs, err := parseCmdLineArgs()
	if err != nil {
		return fmt.Errorf("error parsing command line args : %w", err)
	}

	scratchData := generateScratchData(cmdArgs)

	pages, err := getHtmlPages(scratchData)
	if err != nil {
		return fmt.Errorf("error getting html pages : %w", err)
	}

	fmt.Println(scratchData.ApartmentsIterateUrlDram(1))
	fmt.Println(scratchData.ApartmentsStopWord)
	fmt.Println(scratchData.TownhousesIterateUrlDram(1))
	fmt.Println(scratchData.ApartmentsStopWord)
	fmt.Println(scratchData.HousesIterateUrlDram(1))
	fmt.Println(scratchData.HousesStopWord)
	fmt.Println(len(pages))

	/* Parse htmls and save them to jsons with the corresponding name. Use descriptions` strings as a unique id */

	/* Compare with the old entries. Diffs are the new adds */

	/* Notify */

	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
