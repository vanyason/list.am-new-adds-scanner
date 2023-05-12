package scraper

import (
	"log"
	"strings"

	"github.com/gocolly/colly"

	"github.com/vanyason/list.am-new-adds-scanner/lib/ads"
	"github.com/vanyason/list.am-new-adds-scanner/lib/util"
)

type Ad = ads.Ad

// Entity that collects data from list.am
type Scraper struct {
	domain  string
	path    string
	pageURL string
	colly   *colly.Collector
}

// Constructor
func New(userAgent, domain, path string) Scraper {
	return Scraper{
		domain:  domain,
		path:    path,
		pageURL: domain + path,
		colly: colly.NewCollector(
			colly.UserAgent(userAgent),
		),
	}
}

// Function to override path
// Note : domain is set already
func (s *Scraper) SetPath(path string) {
	s.path = path
	s.pageURL = s.domain + path
}

// Goes to the list.am according to the Scraper.domain + Scraper.path and parses ads on the site.
// When parsing besides ads looks for the button with the next page. If found - overrides url and
// continues parsing
func (s *Scraper) Scrap() (ads []Ad, err error) {
	defer util.LogSecondsPass("Scrap Took")()

	const cssSelectorAd = "div.gl a"
	const cssSelectorNextPage = "div.dlf > a"

	// Reset callbacks
	s.colly = s.colly.Clone()

	// Setup callbacks
	s.colly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	s.colly.OnError(func(_ *colly.Response, err error) {
		log.Printf("colly.OnError : %v", err)
	})

	s.colly.OnHTML(cssSelectorAd, func(el *colly.HTMLElement) { //< Parse ad
		ad := Ad{
			Link:  el.Attr("href"),
			Price: el.ChildText("div.p"),
			Dscr:  strings.Join(strings.Fields(el.ChildText("div.l")), " "),
			At:    strings.Join(strings.Fields(el.ChildText("div.at")), " "),
		}

		if ad.Link == "" || ad.Price == "" || ad.Dscr == "" || ad.At == "" {
			log.Printf("colly.OnHtml : error parsing ads (%v)", ad)
		}

		ads = append(ads, ad)
	})

	s.colly.OnHTML(cssSelectorNextPage, func(el *colly.HTMLElement) { //< Find "Next" button
		path := el.Attr("href")
		if path == "" {
			log.Printf("colly.OnHtml : error parsing next page url")
			return
		}

		newUrl := s.domain + path
		s.colly.Visit(newUrl)
	})

	return ads, s.colly.Visit(s.pageURL)
}
