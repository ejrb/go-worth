package scrape

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Price struct {
	value float64
	ccy   string
}

type Card struct {
	name   string
	set    string
	rarity string
	price  Price
	foil   bool
}

func ChanToSlice(ch interface{}) interface{} {
	chv := reflect.ValueOf(ch)
	slv := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(ch).Elem()), 0, 0)
	for {
		v, ok := chv.Recv()
		if !ok {
			return slv.Interface()
		}
		slv = reflect.Append(slv, v)
	}
}

// A scraper downloads all relevant pages and parses them into a list
// of cards with prices

type CardSource struct {
	download func() (*goquery.Document, error)
	parse    func(*goquery.Document, chan Card)
}

type Scraper struct {
	sources []CardSource
}

func (scraper Scraper) Scrape() []Card {
	var wg sync.WaitGroup
	ch := make(chan Card)

	for _, source := range scraper.sources {
		wg.Add(1)
		go func(cs CardSource) {
			doc, err := cs.download()
			if err == nil {
				cs.parse(doc, ch)
			}
			wg.Done()
		}(source)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ChanToSlice(ch).([]Card)
}

func makeMTGGoldfishDownloadFn(set string, foil bool) func() (*goquery.Document, error) {
	set = strings.ToUpper(set)
	if foil {
		set = fmt.Sprintf("%s_F", set)
	}
	url := fmt.Sprintf("https://www.mtggoldfish.com/index/%s#paper", set)

	return func() (*goquery.Document, error) {
		log.Printf("Downloading: %s", url)
		return goquery.NewDocument(url)
	}
}

func makeMTGGoldfishParser(set string, foil bool) func(doc *goquery.Document, c chan Card) {

	return func(doc *goquery.Document, c chan Card) {
		parseCard := func(i int, s *goquery.Selection) {
			name := strings.TrimSpace(s.Find(".card a").Text())
			rarity := strings.TrimSpace(s.Find("td:nth-child(3)").Text())
			p := strings.TrimSpace(s.Find("td:nth-child(4)").Text())

			if name != "" {
				if price, err := strconv.ParseFloat(p, 64); err == nil {
					c <- Card{name, set, rarity, Price{price, "USD"}, foil}
				}
			}
		}
		doc.Find(".index-price-table tr").Each(parseCard)
	}
}

func NewMTGGoldfishScraper(set string) Scraper {
	reg_downloader := CardSource{
		makeMTGGoldfishDownloadFn(set, false),
		makeMTGGoldfishParser(set, false),
	}
	foil_downloader := CardSource{
		makeMTGGoldfishDownloadFn(set, true),
		makeMTGGoldfishParser(set, true),
	}
	return Scraper{[]CardSource{reg_downloader, foil_downloader}}
}
