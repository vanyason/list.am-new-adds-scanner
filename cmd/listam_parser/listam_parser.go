package main

import (
	"log"

	"github.com/vanyason/list.am-new-adds-scanner/lib/scraper"
	"github.com/vanyason/list.am-new-adds-scanner/lib/util"
)

const appName = "scraper"
const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
const domain = "https://www.list.am/ru"
const path = "/category/56?cmtype=1&pfreq=1&po=1&n=1&price2=300000&crc=0" //< apartments, without agent, monthly, with photo, Yerevan, 300000֏
const path2 = "/category/56?cmtype=1&pfreq=1&po=1&n=1&price2=900&crc=1"   //< apartments, without agent, monthly, with photo, Yerevan, 900$

func main() {
	util.SetupLogging(appName)
	defer util.LogSecondsPass("Total Took")()

	s := scraper.New(userAgent, domain, path)

	ads, err := s.Scrap()
	if err != nil {
		log.Fatal(err)
	}

	s.SetPath(path2)
	_, err = s.Scrap()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Adds parsed: ", len(ads))
}
