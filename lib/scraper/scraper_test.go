package scraper

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
)

const userAgent = "Test User Agent"
const domain = "https://www.test.com"
const path = "/test-path"

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func loadFile(t *testing.T, filePath string) (dat []byte) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Can not read file \"%s\"", filePath)
	}

	return dat
}

func (s *Scraper) resetTransport(dat []byte, statusCode int) {
	s.colly.WithTransport(RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		header := make(http.Header)
		header.Set("Content-Type", "html")

		return &http.Response{
			StatusCode: statusCode,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader(dat)),
		}, nil
	}))
}

func makeScraper(dat []byte, statusCode int) Scraper {
	s := New(userAgent, domain, path)
	s.resetTransport(dat, statusCode)
	return s
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Smth returned error (but shouldn`t) \"%s\"", err)
	}
}

// Test constructor
func TestNew(t *testing.T) {
	s := New(userAgent, domain, path)

	if s.domain != domain {
		t.Errorf("Expected domain to be %s but got %s", domain, s.domain)
	}

	if s.path != path {
		t.Errorf("Expected path to be %s but got %s", path, s.path)
	}

	if s.pageURL != domain+path {
		t.Errorf("Expected pageURL to be %s but got %s", domain+path, s.pageURL)
	}

	if s.colly == nil {
		t.Errorf("Expected colly to be initialized")
	}
}

// Test "SetPath" method
func TestSetPath(t *testing.T) {
	s := New(userAgent, domain, path)

	newPath := "/new-test-path"
	s.SetPath("/new-test-path")

	if s.path != newPath {
		t.Errorf("Expected path to be %s but got %s", newPath, s.path)
	}

	if s.pageURL != domain+newPath {
		t.Errorf("Expected pageURL to be %s but got %s", domain+newPath, s.pageURL)
	}
}

// Test behavior if 404 returned
func TestScrapResCode404(t *testing.T) {
	s := makeScraper(nil, http.StatusNotFound)

	ads, err := s.Scrap()

	if len := len(ads); len != 0 {
		t.Errorf("Expected ads len to be %d but got %d", 0, len)
	}

	expectedErrMsg := "Not Found"
	if errMsg := err.Error(); errMsg != expectedErrMsg {
		t.Errorf("Expected error message to be \"%s\" but got \"%s\"", expectedErrMsg, errMsg)
	}
}

// Test parsing of a typical page. 105 ads should be returned
func TestScrapDramFull(t *testing.T) {
	dat := loadFile(t, "testdata/ads_dram_full.html")

	s := makeScraper(dat, http.StatusOK)

	ads, err := s.Scrap()
	if err != nil {
		t.Fatalf("Expected Scrap result not to return err but got \"%s\"", err)
	}

	const expectedLen = 105
	if len := len(ads); len != expectedLen {
		t.Fatalf("Expected ads len to be %d but got %d", expectedLen, len)
	}
}

// Test parsing typical page (shorter version). Check results here
func TestScrapDramShort(t *testing.T) {
	dat := loadFile(t, "testdata/ads_dram_short.html")

	s := makeScraper(dat, http.StatusOK)

	ads, err := s.Scrap()
	checkErr(t, err)

	const expectedLen = 2
	expectedAds := [expectedLen]Ad{
		{Link: "/ru/item/19279834", Price: "250,000 ֏ в месяц", Dscr: "1-комн. квартира на ул. Мамиконянца, 25 кв.м., 5/9 этаж, дизайнерский ремонт", At: "Арабкир, 1 ком., 25 кв.м., 5/9 этаж"},
		{Link: "/ru/item/18789181", Price: "250,000 ֏ в месяц", Dscr: "2-комн. квартира, 2-й переулок улицы Багяна, 60 кв.м., 1/5 этаж, евроремонт, каменное здание", At: "Нор Норк, 2 ком., 60 кв.м., 1/5 этаж"},
	}

	if len := len(ads); len != expectedLen {
		t.Fatalf("Expected ads len to be %d but got %d", expectedLen, len)
	}

	for i := range ads {
		if ads[i] != expectedAds[i] {
			t.Fatalf("Expected ad %v but got %v", ads[i], expectedAds[i])
		}
	}
}

// If page has a link to the next one - it will parse next one. Because I use RoundTripper,
// urls are ignored and same page will be returned
func TestScrapDramShortLinked(t *testing.T) {
	dat := loadFile(t, "testdata/ads_dram_short_linked_0.html")

	s := makeScraper(dat, http.StatusOK)

	ads, err := s.Scrap()
	checkErr(t, err)

	const expectedLen = 4
	expectedAds := [expectedLen]Ad{
		{Link: "/ru/item/19279834", Price: "250,000 ֏ в месяц", Dscr: "1-комн. квартира на ул. Мамиконянца, 25 кв.м., 5/9 этаж, дизайнерский ремонт", At: "Арабкир, 1 ком., 25 кв.м., 5/9 этаж"},
		{Link: "/ru/item/18789181", Price: "250,000 ֏ в месяц", Dscr: "2-комн. квартира, 2-й переулок улицы Багяна, 60 кв.м., 1/5 этаж, евроремонт, каменное здание", At: "Нор Норк, 2 ком., 60 кв.м., 1/5 этаж"},
		{Link: "/ru/item/19279834", Price: "250,000 ֏ в месяц", Dscr: "1-комн. квартира на ул. Мамиконянца, 25 кв.м., 5/9 этаж, дизайнерский ремонт", At: "Арабкир, 1 ком., 25 кв.м., 5/9 этаж"},
		{Link: "/ru/item/18789181", Price: "250,000 ֏ в месяц", Dscr: "2-комн. квартира, 2-й переулок улицы Багяна, 60 кв.м., 1/5 этаж, евроремонт, каменное здание", At: "Нор Норк, 2 ком., 60 кв.м., 1/5 этаж"},
	}

	if len := len(ads); len != expectedLen {
		t.Fatalf("Expected ads len to be %d but got %d", expectedLen, len)
	}

	for i := range ads {
		if ads[i] != expectedAds[i] {
			t.Fatalf("Expected ad %v but got %v", ads[i], expectedAds[i])
		}
	}
}

// Test invalid html
func TestScrapInvalidHtml(t *testing.T) {
	dat := loadFile(t, "testdata/invalid.html")

	s := makeScraper(dat, http.StatusOK)

	ads, err := s.Scrap()
	checkErr(t, err)

	const expectedLen = 0
	if len := len(ads); len != expectedLen {
		t.Fatalf("Expected ads len to be %d but got %d", expectedLen, len)
	}
}

// Test broken html
func TestScrapIncorrectHtml(t *testing.T) {
	dat := loadFile(t, "testdata/incorrect.html")

	s := makeScraper(dat, http.StatusOK)

	ads, err := s.Scrap()
	checkErr(t, err)

	const expectedLen = 0
	if len := len(ads); len != expectedLen {
		t.Fatalf("Expected ads len to be %d but got %d", expectedLen, len)
	}
}

// Test that changing paths and visiting different url works fine
func TestScrapDifferentHtml(t *testing.T) {
	dat := loadFile(t, "testdata/ads_dram_short.html")
	s := makeScraper(dat, http.StatusOK)

	ads, err := s.Scrap()
	checkErr(t, err)

	dat = loadFile(t, "testdata/ads_dram_short_2.html")
	s.resetTransport(dat, http.StatusOK)
	s.SetPath("/newPath")

	ads2, err := s.Scrap()
	checkErr(t, err)

	ads = append(ads, ads2...)

	const expectedLen = 4
	expectedAds := [expectedLen]Ad{
		{Link: "/ru/item/19279834", Price: "250,000 ֏ в месяц", Dscr: "1-комн. квартира на ул. Мамиконянца, 25 кв.м., 5/9 этаж, дизайнерский ремонт", At: "Арабкир, 1 ком., 25 кв.м., 5/9 этаж"},
		{Link: "/ru/item/18789181", Price: "250,000 ֏ в месяц", Dscr: "2-комн. квартира, 2-й переулок улицы Багяна, 60 кв.м., 1/5 этаж, евроремонт, каменное здание", At: "Нор Норк, 2 ком., 60 кв.м., 1/5 этаж"},
		{Link: "/ru/item/1234", Price: "250,000 ֏ в месяц", Dscr: "aa", At: "Арабкир, 1 ком., 25 кв.м., 5/9 этаж"},
		{Link: "/ru/item/1234", Price: "250,000 ֏ в месяц", Dscr: "bb", At: "Нор Норк, 2 ком., 60 кв.м., 1/5 этаж"},
	}

	if len := len(ads); len != expectedLen {
		t.Fatalf("Expected ads len to be %d but got %d", expectedLen, len)
	}

	for i := range ads {
		if ads[i] != expectedAds[i] {
			t.Fatalf("Expected ad %v but got %v", ads[i], expectedAds[i])
		}
	}
}

// TODO
// Test houses when implemented
