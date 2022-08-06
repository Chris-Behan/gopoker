package cards

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Suit represents a cards suit. Ex Spade
type Suit string

// Suits, which in the game of poker, all have the same value.
const (
	Spade   Suit = "Spades"
	Club    Suit = "Clubs"
	Heart   Suit = "Hearts"
	Diamond Suit = "Diamonds"
)

var suits = []Suit{Spade, Club, Heart, Diamond}

// Rank represents a cards value. Ex. Jack
type Rank int8

// Card ranks, two is the lowest, Ace is the highest.
const (
	Two   Rank = 2
	Three Rank = 3
	Four  Rank = 4
	Five  Rank = 5
	Six   Rank = 6
	Seven Rank = 7
	Eight Rank = 8
	Nine  Rank = 9
	Ten   Rank = 10
	Jack  Rank = 11
	Queen Rank = 12
	King  Rank = 13
	Ace   Rank = 14
)

type handRank int16

// Poker hand ranks mapped to arbitrary values with descending order based on rank
const (
	royalFlushRank    handRank = 100
	straightFlushRank handRank = 99
	fourOfAKindRank   handRank = 98
	flushRank         handRank = 97
	straightRank      handRank = 96
	threeOfAKindRank  handRank = 95
	twoPairRank       handRank = 94
	pairRank          handRank = 93
	highCardRank      handRank = 92
)

// Card represents a playing card.
type Card struct {
	rank Rank
	suit Suit
}

type Hand []Card

// Implement the sort.Interface so that we can sort a hand.
func (h Hand) Len() int           { return len(h) }
func (h Hand) Less(a, b int) bool { return h[a].rank < h[b].rank }
func (h Hand) Swap(a, b int)      { h[a], h[b] = h[b], h[a] }

// Deck represents a deck of cards.
type Deck struct {
	cards []Card
}

// GenerateDeck returns a Deck of 52 shuffled playing cards.
func GenerateDeck() Deck {
	suits := []Suit{Spade, Club, Heart, Diamond}
	ranks := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
	cards := make([]Card, 0)
	for _, s := range suits {
		for _, r := range ranks {
			c := Card{r, s}
			cards = append(cards, c)
		}
	}
	shuffledCards := shuffle(cards)
	deck := Deck{shuffledCards}
	return deck
}

func shuffle(cards []Card) []Card {
	shuffledDeck := []Card{}
	i := len(cards)
	for i > 0 {
		rand_idx := rand.Intn(len(cards))
		// Add randomly selected card to new deck.
		c := cards[rand_idx]
		shuffledDeck = append(shuffledDeck, c)
		// Move element at end of slice to position of the randomly selected card.
		cards[rand_idx] = cards[len(cards)-1]
		// Remove last element from slice since it's been copied to the position of the card we
		// just removed.
		cards = cards[:len(cards)-1]
		i--
	}
	return shuffledDeck
}

// Length returns the number of cards in the deck.
func (deck Deck) Length() int {
	return len(deck.cards)
}

// Draw removes and returns the last card from the deck.
func (deck *Deck) Draw() (Card, error) {
	if deck.Length() == 0 {
		return Card{}, errors.New("Deck is empty.")
	}
	c := deck.cards[len(deck.cards)-1]
	deck.cards = deck.cards[:len(deck.cards)-1]
	return c, nil
}

func (deck Deck) GetCards() []Card {
	return deck.cards
}

func getHandRank(hand []Card) handRank {
	if hasRoyalFlush, rank := royalFlush(hand); hasRoyalFlush {
		return rank
	} else if hasStraightFlush, rank := straightFlush(hand); hasStraightFlush {
		return rank
	} else if hasFourOfAKind, rank := fourOfAKind(hand); hasFourOfAKind {
		return rank
	} else if hasFlush, rank := flush(hand); hasFlush {
		return rank
	} else if hasStraight, rank := straight(hand); hasStraight {
		return rank
	} else if hasThreeOfAKind, rank := threeOfAKind(hand); hasThreeOfAKind {
		return rank
	} else if hasTwoPair, rank := twoPair(hand); hasTwoPair {
		return rank
	} else if hasPair, rank := pair(hand); hasPair {
		return rank
	} else {
		return highCardRank
	}
}

func royalFlush(hand []Card) (bool, handRank) {
	tens := getCardsByRank(hand, Ten)
	for _, t := range tens {
		hasRoyalFlush := royalFlushSearch(hand, t.suit, t)
		if hasRoyalFlush {
			return true, royalFlushRank
		}
	}
	return false, 0
}

// royalFlushSearch performs a depth first search for a Royal Flush in a slice of cards.
// hand is the deck of cards to search, current is the card to start the search at (root node),
// and suit is the suit of the starting card.
func royalFlushSearch(hand []Card, suit Suit, current Card) bool {
	if current.rank == Ten && current.suit == suit {
		idx, jack := cardSearchByRankAndSuit(hand, Jack, suit)
		if idx == -1 {
			return false
		}
		return royalFlushSearch(hand, suit, jack)
	} else if current.rank == Jack && current.suit == suit {
		idx, queen := cardSearchByRankAndSuit(hand, Queen, suit)
		if idx == -1 {
			return false
		}
		return royalFlushSearch(hand, suit, queen)
	} else if current.rank == Queen && current.suit == suit {
		idx, king := cardSearchByRankAndSuit(hand, King, suit)
		if idx == -1 {
			return false
		}
		return royalFlushSearch(hand, suit, king)
	} else if current.rank == King && current.suit == suit {
		idx, ace := cardSearchByRankAndSuit(hand, Ace, suit)
		if idx == -1 {
			return false
		}
		return royalFlushSearch(hand, suit, ace)
	} else if current.rank == Ace && current.suit == suit {
		return true
	}
	return false
}

func straightFlush(hand []Card) (bool, handRank) {
	// map of suits to array of bools that indicate whether or not a card exists.
	// index 0 represents an ace.
	cardMapAceLow := createCardMap()
	for _, card := range hand {
		if card.rank == Ace {
			cardMapAceLow[card.suit][0] = true
		} else {
			// When representing ace as the low card, use rank -1 as card position
			cardMapAceLow[card.suit][card.rank-1] = true
		}
	}
	cardMapAceHigh := createCardMap()
	for _, card := range hand {
		// When representing ace as the high card, use rank -2 as card position
		cardMapAceHigh[card.suit][card.rank-2] = true
	}

	if fiveInARow(cardMapAceLow) || fiveInARow(cardMapAceHigh) {
		return true, straightFlushRank
	}
	return false, 0
}

// FiveInARow iterates through a map of boolean arrays, returning true if any of the arrays contain 5 consecutive 'true' values.
// Otherwise it returns false.
func fiveInARow(cardMap map[Suit]*[13]bool) bool {
	for _, row := range cardMap {
		count := 0
		for i := 0; i < len(row); i++ {
			if row[i] {
				count++
			} else {
				count = 0
			}

			if count == 5 {
				return true
			}
		}
	}
	return false
}

func createCardMap() map[Suit]*[13]bool {
	cardMap := make(map[Suit]*[13]bool)
	for _, s := range suits {
		var row [13]bool
		cardMap[s] = &row
	}
	return cardMap
}

func fourOfAKind(hand []Card) (bool, handRank) {
	cardCounts := cardCountsByRank(hand)
	for _, v := range cardCounts {
		if v == 4 {
			return true, fourOfAKindRank
		}
	}
	return false, 0
}

func flush(hand []Card) (bool, handRank) {
	suitCounts := cardCountsBySuit(hand)
	for _, v := range suitCounts {
		if v == 5 {
			return true, flushRank
		}
	}
	return false, 0
}

func straight(hand []Card) (bool, handRank) {
	if len(hand) < 5 {
		return false, 0
	}
	// Check for straight with Ace as low card
	orderedHandAceLow := orderByRank(hand, true)
	consecutiveCount := 1
	prev := orderedHandAceLow[0]
	i := 1
	for i < len(hand) {
		currentRank := orderedHandAceLow[i].rank
		prevRank := prev.rank
		// Treat ace as low card
		if orderedHandAceLow[i].rank == Ace {
			currentRank = 1
		}
		if prev.rank == Ace {
			prevRank = 1
		}

		// Increment count, reset count, or do nothing (The case when currentRank == prevRank)
		if currentRank == prevRank+1 {
			consecutiveCount++
		} else if currentRank > prevRank+1 {
			consecutiveCount = 1
		}
		prev = orderedHandAceLow[i]
		i += 1

		if consecutiveCount == 5 {
			return true, straightRank
		}
	}

	// Check for straight with Ace as high card
	orderedHandAceHigh := orderByRank(hand, false)
	consecutiveCount = 1
	prev = orderedHandAceHigh[0]
	i = 1
	for i < len(hand) {
		currentRank := orderedHandAceHigh[i].rank
		if currentRank == prev.rank+1 {
			consecutiveCount++
		} else if currentRank > prev.rank+1 {
			consecutiveCount = 1
		}
		prev = orderedHandAceHigh[i]
		i += 1

		if consecutiveCount == 5 {
			return true, straightRank
		}
	}
	return false, 0
}

func threeOfAKind(hand []Card) (bool, handRank) {
	cardCounts := cardCountsByRank(hand)
	for _, count := range cardCounts {
		if count == 3 {
			return true, threeOfAKindRank
		}
	}
	return false, 0
}

func twoPair(hand []Card) (bool, handRank) {
	cardCounts := cardCountsByRank(hand)
	pairCount := 0
	for _, count := range cardCounts {
		if count == 2 {
			pairCount++
		}
		if pairCount == 2 {
			return true, twoPairRank
		}
	}
	return false, 0
}

func pair(hand []Card) (bool, handRank) {
	cardCounts := cardCountsByRank(hand)
	for _, count := range cardCounts {
		if count == 2 {
			return true, pairRank
		}
	}
	return false, 0
}

func highCard(hand []Card) Rank {
	high := Rank(-1)
	for _, card := range hand {
		if card.rank > high {
			high = card.rank
		}
	}
	return high
}

// copyAndRemoveCard returns a copy of the cards passed to the function
// minus the card at the specified index. The calling slice
// of cards is unaffected.
func copyAndRemoveCard(cards []Card, idx int) ([]Card, error) {
	cardsCopy := make([]Card, len(cards))
	copy(cardsCopy, cards)
	cardsCopy, err := removeCard(cardsCopy, idx)
	if err != nil {
		return []Card{}, err
	}
	return cardsCopy, nil
}

func removeCard(cards []Card, idx int) ([]Card, error) {
	if idx >= len(cards) || idx < 0 {
		return []Card{}, fmt.Errorf("No card at index %v. cards: %v", idx, cards)
	}
	// copy elements 1 to the right of deletion index into deletion index.
	copy(cards[idx:], cards[idx+1:])
	// Clear the card at the end of the slice, since it is now a duplicate of the card to its left.
	cards[len(cards)-1] = Card{}
	// Shrink slice by 1
	cards = cards[:len(cards)-1]
	return cards, nil
}

func getCardsByRank(cards []Card, rank Rank) []Card {
	matches := []Card{}
	for _, c := range cards {
		if c.rank == rank {
			matches = append(matches, c)
		}
	}
	return matches
}

func cardSearchByRank(cards []Card, rank Rank) (int, Card) {
	for idx, c := range cards {
		if c.rank == rank {
			return idx, c
		}
	}
	return -1, Card{}
}

func cardSearchByRankAndSuit(cards []Card, rank Rank, suit Suit) (int, Card) {
	for idx, c := range cards {
		if c.rank == rank && c.suit == suit {
			return idx, c
		}
	}
	return -1, Card{}
}

func orderByRank(cards []Card, aceLow bool) []Card {
	// Copy contents of calling slice into new slice so that the original is unaffected.
	cardsCopy := make([]Card, len(cards))
	copy(cardsCopy, cards)
	ordered := []Card{}
	for len(cardsCopy) > 0 {
		// Set min to a fake card with super rank that is higher than the possible ranks to start.
		min := Card{Rank(99), Heart}
		minIdx := -1
		for idx, card := range cardsCopy {
			if aceLow {
				if card.rank == Ace {
					min = card
					minIdx = idx
				} else if card.rank <= min.rank && min.rank != Ace {
					min = card
					minIdx = idx
				}
			} else {
				if card.rank <= min.rank {
					min = card
					minIdx = idx
				}
			}
		}
		// add smallest card to ordered slice of cards
		ordered = append(ordered, min)
		// remove the card we just added from the original slice of cards
		cardsCopy[minIdx] = cardsCopy[len(cardsCopy)-1]
		cardsCopy[len(cardsCopy)-1] = Card{}
		cardsCopy = cardsCopy[:len(cardsCopy)-1]
	}
	return ordered
}

func cardCountsByRank(cards []Card) map[Rank]int {
	counts := make(map[Rank]int)
	for _, c := range cards {
		if _, exists := counts[c.rank]; exists {
			counts[c.rank]++
		} else {
			counts[c.rank] = 1
		}
	}
	return counts
}

func cardCountsBySuit(cards []Card) map[Suit]int {
	counts := make(map[Suit]int)
	for _, c := range cards {
		if _, exists := counts[c.suit]; exists {
			counts[c.suit]++
		} else {
			counts[c.suit] = 1
		}
	}
	return counts
}
