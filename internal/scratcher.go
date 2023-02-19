package internal

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/net/html"
)

type ScratchData struct {
	Header                   string
	ApartmentsIterateUrlDram func(page int) string
	ApartmentsStopWord       string
	TownhousesIterateUrlDram func(page int) string
	HousesIterateUrlDram     func(page int) string
	HousesStopWord           string
}

func GenerateScratchData(args CmdArguments) ScratchData {
	sd := ScratchData{
		Header:             "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
		ApartmentsStopWord: "/category/56",
		HousesStopWord:     "/category/63",
	}

	const listamUrl = "https://www.list.am/ru"
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

func ScratchHtmlPages(sd ScratchData) (map[string]string, error) {
	getPagesPerCategory := func(urlGenerator func(page int) string, stopWord string) (map[string]string, error) {
		/* If listam behaves strangely (it can) break */
		const emergencyTimeoutSec = 90 //< Is it enough to iterate over everything ?
		start := time.Now()

		client := http.Client{}
		parsedPages := make(map[string]string)

		for page := 1; true; page++ {
			url := urlGenerator(page)

			req, err := http.NewRequest("GET", url, nil)
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

			parsedPages, err = parseHtml(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("error parsing page (%s) : %w", url, err)
			}

			if len(parsedPages) == 0 {
				return parsedPages, nil
			}

			/* Protection from strange listam behaviour */
			if time.Since(start).Seconds() > emergencyTimeoutSec {
				return nil, fmt.Errorf("error scratching. timeout. check for the infinite loop : %s", url)
			}
		}

		return parsedPages, nil
	}

	allPages := make(map[string]string)

	pagesPerCategory, err := getPagesPerCategory(sd.ApartmentsIterateUrlDram, sd.ApartmentsStopWord)
	if err != nil {
		return nil, err
	}
	maps.Copy(allPages, pagesPerCategory)

	pagesPerCategory, err = getPagesPerCategory(sd.TownhousesIterateUrlDram, sd.HousesStopWord)
	if err != nil {
		return nil, err
	}
	maps.Copy(allPages, pagesPerCategory)

	pagesPerCategory, err = getPagesPerCategory(sd.HousesIterateUrlDram, sd.HousesStopWord)
	if err != nil {
		return nil, err
	}
	maps.Copy(allPages, pagesPerCategory)

	return allPages, nil
}

func parseHtml(r io.Reader) (parsedPage map[string]string, err error) {
	re, err := regexp.Compile(`[\n\s]`)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex : %w", err)
	}

	parsedPage = make(map[string]string)

	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("error parsing page : %w", err)
	}

	var parseLink func(*html.Node)
	parseLink = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.Contains(a.Val, "/ru/item/") {
					var description string
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type != html.ElementNode || c.Data != "div" {
							continue
						}
						for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
							description += (cc.Data + ";")
						}
					}
					description = re.ReplaceAllString(description, "")
					parsedPage[description] = a.Val
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseLink(c)
		}
	}

	parseLink(doc)

	return parsedPage, nil
}
