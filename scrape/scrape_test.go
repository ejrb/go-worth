package scrape

import (
	"fmt"
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

func TestMTGGoldfishScraping(t *testing.T) {
	scraper := NewMTGGoldfishScraper("rix")

	cards := scraper.Scrape()

	for _, card := range cards {
		fmt.Printf("%s $%.2f", card.name, card.price.value)
		if card.foil {
			fmt.Printf(" (FOIL)")
		}
		fmt.Printf("\n")
	}
}
