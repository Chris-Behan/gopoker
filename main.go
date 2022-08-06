package main

import (
	"fmt"

	"github.com/Chris-Behan/gopoker/cards"
)

func main() {
	deck1 := cards.GenerateDeck()
	fmt.Printf("Deck1: %v\n", deck1)
	c := deck1.GetCards()
	fmt.Printf("Number of cards in deck: %v", len(c))
	// for n := 0; n < 52; n++ {
	// 	card, err := deck1.Draw()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("Draw %v: %v", n, card)
	// }
}
