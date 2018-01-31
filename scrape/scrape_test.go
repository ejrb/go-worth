package scrape

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func Filter(cards []Card, fs ...func(Card) bool) []Card {
	matches := make([]Card, 0)

	allPass := func(c Card) bool {
		for _, f := range fs {
			if !f(c) {
				return false
			}
		}
		return true
	}

	for _, card := range cards {
		if allPass(card) {
			matches = append(matches, card)
		}
	}
	return matches
}

func getTestDocument() (*goquery.Document, error) {
	f, _ := os.Open("RIXTestData.htm")
	return goquery.NewDocumentFromReader(f)
}

func isRarity(r string) func(c Card) bool {
	return func(c Card) bool { return c.rarity == r }
}

// Tests the makeMTGGoldfishParser function asserting that the parser can
// pull out expected cards from a test HTML page
func TestMTGGoldfishParserUnit(t *testing.T) {
	cs := CardSource{
		getTestDocument,
		makeMTGGoldfishParser("rix", false),
	}

	scraper := Scraper{[]CardSource{cs}}

	cards := scraper.Scrape()

	if n := len(cards); n != 31 {
		t.Error("expected 31 cards, got ", n)
	}
	if n := len(Filter(cards, func(c Card) bool { return c.set == "RIX" })); n != len(cards) {
		t.Error("all cards should be RIX cards, got ", n)
	}

	aDollarOrHigher := func(c Card) bool { return c.price.value > 1. }
	if n := len(Filter(cards, aDollarOrHigher)); n != 16 {
		t.Error("expected 16 cards > $1, got ", n)
	}

	if n := len(Filter(cards, isRarity("Common"))); n != 2 {
		t.Error("expected 2 commons, got ", n)
	}
	if n := len(Filter(cards, isRarity("Uncommon"))); n != 19 {
		t.Error("expected 19 uncommons, got ", n)
	}
	if n := len(Filter(cards, isRarity("Rare"))); n != 9 {
		t.Error("expected 9 rares, got ", n)
	}
	if n := len(Filter(cards, isRarity("Mythic"))); n != 1 {
		t.Error("expected 1 mythics, got ", n)
	}
}

// Download the Rivals of Ixalan prices
// func TestMTGGoldfishScrapingItegration(t *testing.T) {
// 	scraper := NewMTGGoldfishScraper("rix")
//
// 	cards := scraper.Scrape()
//
// 	if n := len(cards); n != 804 {
// 		t.Error("expected 804 RIX cards, got ", n)
// 	}
//
// 	isFoil := func(c Card) bool { return c.foil }
// 	if n := len(Filter(cards, isFoil)); n != 394 {
// 		t.Error("expected 394 foils, got ", n)
// 	}
//
// 	if n := len(Filter(cards, isRarity("Rare"))); n != 199 {
// 		t.Error("expected 199 rares, got ", n)
// 	}
// 	if n := len(Filter(cards, isRarity("Mythic"))); n != 60 {
// 		t.Error("expected 60 mythics, got ", n)
// 	}
// }
