package cards

import "testing"

// Tests that highCard returns the rank of the highest card in a hand.
func TestHighCard(t *testing.T) {
	hand := []Card{{Two, Heart}, {Ace, Spade}, {King, Diamond}, {Queen, Club}, {Five, Spade}}
	rank := highCard(hand)
	if rank != Ace {
		t.Errorf("Expected rank to be %v but instead it was %v.", Ace, rank)
	}
}

// Tests that highCard returns -1 of the input hand is empty.
func TestHighCardEmptyHand(t *testing.T) {
	rank := highCard([]Card{})
	if rank != -1 {
		t.Errorf("Expected the rank of an empty hand to be -1 but instead it was %v.", rank)
	}
}

func TestRemoveCard(t *testing.T) {
	tests := []struct {
		inputCards []Card
		inputIdx   int
		expected   []Card
	}{
		{[]Card{{Ten, Club}, {Two, Heart}, {Three, Club}, {Four, Diamond}, {Five, Spade}},
			2,
			[]Card{{Ten, Club}, {Two, Heart}, {Four, Diamond}, {Five, Spade}},
		},
		{[]Card{{Ten, Club}},
			0,
			[]Card{},
		},
		{[]Card{{Ten, Club}, {Two, Heart}, {Three, Club}, {Four, Diamond}, {Five, Spade}},
			5,
			[]Card{},
		},
	}

	for _, test := range tests {
		cards, _ := removeCard(test.inputCards, test.inputIdx)
		if !cardsEqual(cards, test.expected) {
			t.Errorf("Expected: %v Actual: %v", test.expected, cards)
		}
	}
}

func TestRemoveCardOutOfBounds(t *testing.T) {
	tests := []struct {
		inputCards []Card
		inputIdx   int
	}{
		{[]Card{{Ten, Club}, {Two, Heart}, {Three, Club}, {Four, Diamond}, {Five, Spade}},
			7,
		},
		{[]Card{{Ten, Club}, {Two, Heart}, {Three, Club}, {Four, Diamond}, {Five, Spade}},
			-1,
		},
	}
	for _, test := range tests {
		_, err := removeCard(test.inputCards, test.inputIdx)
		if err == nil {
			t.Errorf("Expected an error to be returned but there wasn't.")
		}
	}
}

func TestCopyAndRemoveCard(t *testing.T) {
	inputCards := []Card{{Ten, Club}, {Two, Heart}, {Three, Club}, {Four, Diamond}, {Five, Spade}}
	expectedCards := []Card{{Ten, Club}, {Three, Club}, {Four, Diamond}, {Five, Spade}}
	cards, _ := copyAndRemoveCard(inputCards, 1)
	if !cardsEqual(cards, expectedCards) {
		t.Errorf("Expected: %v Actual: %v", expectedCards, cards)
	}

	// Test that change to original cards does not affect the copy
	cardsLength := len(cards)
	// Remove all cards in original slice
	inputCardsLength := len(inputCards)
	for i := 0; i < inputCardsLength; i++ {
		inputCards, _ = removeCard(inputCards, 0)
	}
	if cardsLength != len(cards) {
		t.Errorf("Expected the copied slice of cards to be unaffected by modifications to the original.")
	}
}

func TestOrderByRankAceLow(t *testing.T) {
	unordered := []Card{{King, Heart}, {Ace, Heart}, {Queen, Heart}, {Jack, Heart}, {Ten, Heart}, {Nine, Heart}}
	ordered := []Card{{Ace, Heart}, {Nine, Heart}, {Ten, Heart}, {Jack, Heart}, {Queen, Heart}, {King, Heart}}
	result := orderByRank(unordered, true)
	if !cardsEqual(result, ordered) {
		t.Errorf("Expected: %v Actual: %v", ordered, result)
	}
}

func TestOrderByRankAceHigh(t *testing.T) {
	unordered := []Card{{King, Heart}, {Queen, Heart}, {Jack, Heart}, {Ten, Heart}, {Nine, Heart}, {Ace, Heart}}
	ordered := []Card{{Nine, Heart}, {Ten, Heart}, {Jack, Heart}, {Queen, Heart}, {King, Heart}, {Ace, Heart}}
	result := orderByRank(unordered, false)
	if !cardsEqual(result, ordered) {
		t.Errorf("Expected: %v Actual: %v", ordered, result)
	}
}

func TestCardCountsByRank(t *testing.T) {
	cards := []Card{{King, Heart}, {King, Diamond}, {Ace, Spade}, {Two, Spade}, {Two, Diamond}, {Two, Club}}
	expectedCounts := map[Rank]int{
		King: 2,
		Ace:  1,
		Two:  3,
	}
	counts := cardCountsByRank(cards)
	for k, v := range expectedCounts {
		if counts[k] != v {
			t.Errorf("Expected cardCountsByRank(%v) to return %v, but instead it returned %v.", k, v, counts[k])
		}
	}
}

func TestCardCountsBySuit(t *testing.T) {
	cards := []Card{{King, Heart}, {Three, Heart}, {Two, Spade}, {Ace, Spade}, {Four, Spade}, {Five, Diamond}}
	expectedCounts := map[Suit]int{
		Heart:   2,
		Spade:   3,
		Diamond: 1,
	}
	counts := cardCountsBySuit(cards)
	for k, v := range expectedCounts {
		if counts[k] != v {
			t.Errorf("Expected cardCountsBySuit(%v) to return %v, but instead it returned %v.", k, v, counts[k])
		}
	}
}

func TestRoyalFlush(t *testing.T) {
	tests := []struct {
		hand          []Card
		hasRoyalFlush bool
	}{
		{
			[]Card{{Ten, Club}, {Ten, Heart}, {Three, Club}, {Jack, Heart}, {Queen, Heart}, {King, Heart}, {Ace, Heart}},
			true,
		},
		{
			[]Card{{Ten, Club}, {Ten, Heart}, {Three, Club}, {Jack, Heart}, {Queen, Heart}, {King, Heart}, {Ace, Club}},
			false,
		},
		{
			[]Card{},
			false,
		},
	}
	for _, test := range tests {
		hasRoyalFlush, _ := royalFlush(test.hand)
		if hasRoyalFlush != test.hasRoyalFlush {
			t.Errorf("Expected royalFlush(%v) to return %v, but instead it returned %v.",
				test.hand,
				test.hasRoyalFlush,
				hasRoyalFlush)
		}
	}
}

func TestStraightFlush(t *testing.T) {
	tests := []struct {
		hand             []Card
		hasStraightFlush bool
	}{
		{
			[]Card{{Nine, Club}, {Nine, Heart}, {Three, Club}, {Ten, Heart}, {Jack, Heart}, {Queen, Heart}, {King, Heart}},
			true,
		},
		{
			[]Card{{Ace, Diamond}, {Two, Diamond}, {Three, Heart}, {Three, Diamond}, {Four, Diamond}, {Five, Diamond}},
			true,
		},
		{
			[]Card{{Two, Spade}, {Ace, Club}, {Three, Spade}, {Four, Spade}, {Jack, Heart}, {Five, Spade}},
			false,
		},
		{
			[]Card{{Queen, Club}, {King, Club}, {Ace, Club}, {Two, Club}, {Three, Club}},
			false,
		},
		{
			[]Card{{Two, Spade}, {Ace, Spade}, {Three, Spade}, {Four, Spade}, {Jack, Heart}, {Six, Spade}},
			false,
		},
	}

	for _, test := range tests {
		hasStraightFlush, _ := straightFlush(test.hand)
		if hasStraightFlush != test.hasStraightFlush {
			t.Errorf("Expected straightFlush(%v) to return %v, but instead it returned %v.",
				test.hand,
				test.hasStraightFlush,
				hasStraightFlush)
		}
	}
}

func TestFourOfAKind(t *testing.T) {
	tests := []struct {
		hand           []Card
		hasFourOfAKind bool
	}{
		{
			[]Card{{Nine, Club}, {Nine, Heart}, {Three, Club}, {Ten, Heart}, {Nine, Spade}, {Queen, Heart}, {Nine, Diamond}},
			true,
		},
		{
			[]Card{{Two, Spade}, {Ace, Spade}, {Three, Spade}, {Four, Spade}, {Jack, Heart}, {Five, Spade}},
			false,
		},
		{
			[]Card{},
			false,
		},
	}
	for _, test := range tests {
		hasFourOfAKind, _ := fourOfAKind(test.hand)
		if hasFourOfAKind != test.hasFourOfAKind {
			t.Errorf("Expected fourOfAKind(%v) to return %v, but instead it returned %v.",
				test.hand,
				test.hasFourOfAKind,
				hasFourOfAKind)
		}
	}
}

func TestFlush(t *testing.T) {
	tests := []struct {
		hand     []Card
		hasFlush bool
	}{
		{
			[]Card{{Nine, Club}, {Two, Club}, {Three, Club}, {Ten, Club}, {Jack, Spade}, {Queen, Club}, {Nine, Diamond}},
			true,
		},
		{
			[]Card{{Two, Heart}, {Ace, Heart}, {Three, Heart}, {Four, Heart}, {Jack, Diamond}, {Five, Spade}, {Seven, Diamond}},
			false,
		},
		{
			[]Card{},
			false,
		},
	}
	for _, test := range tests {
		hasFlush, _ := flush(test.hand)
		if hasFlush != test.hasFlush {
			t.Errorf("Expected flush(%v) to return %v, but it instead it returned %v.",
				test.hand,
				test.hasFlush,
				hasFlush)
		}
	}
}

func TestStraight(t *testing.T) {
	tests := []struct {
		hand        []Card
		hasStraight bool
	}{
		{
			[]Card{{Ace, Club}, {Jack, Diamond}, {Two, Club}, {Three, Club}, {Queen, Diamond}, {Four, Spade}, {Five, Diamond}},
			true,
		},
		{
			[]Card{{Ten, Club}, {Jack, Heart}, {Queen, Diamond}, {King, Spade}, {Ace, Heart}, {Three, Spade}, {Seven, Club}},
			true,
		},
		{
			[]Card{{Two, Heart}, {Ace, Heart}, {Three, Heart}, {Four, Heart}, {Jack, Diamond}, {Six, Spade}, {Seven, Diamond}},
			false,
		},
		{
			[]Card{},
			false,
		},
	}
	for _, test := range tests {
		hasStraight, _ := straight(test.hand)
		if hasStraight != test.hasStraight {
			t.Errorf("Expected straight(%v) to return %v, but instead it returned %v.",
				test.hand,
				test.hasStraight,
				hasStraight)
		}
	}
}

func TestThreeOfAKind(t *testing.T) {
	tests := []struct {
		hand            []Card
		hasThreeOfAKind bool
	}{
		{
			[]Card{{Two, Diamond}, {Two, Club}, {Four, Spade}, {Two, Heart}},
			true,
		},
		{
			[]Card{{Two, Diamond}, {Two, Club}, {Four, Spade}, {Four, Heart}},
			false,
		},
		{
			[]Card{},
			false,
		},
	}
	for _, test := range tests {
		hasThreeOfAKind, _ := threeOfAKind(test.hand)
		if hasThreeOfAKind != test.hasThreeOfAKind {
			t.Errorf("Expected threeOfAKind(%v) to to return %v, but instead it returned %v.",
				test.hand,
				test.hasThreeOfAKind,
				hasThreeOfAKind)
		}
	}
}

func TestTwoPair(t *testing.T) {
	tests := []struct {
		hand       []Card
		hasTwoPair bool
	}{
		{
			[]Card{{Two, Diamond}, {Two, Club}, {Four, Spade}, {Five, Heart}, {Four, Club}},
			true,
		},
		{
			[]Card{{Two, Diamond}, {Two, Club}, {Four, Spade}, {Six, Heart}, {Seven, Diamond}},
			false,
		},
		{
			[]Card{},
			false,
		},
	}
	for _, test := range tests {
		hasTwoPair, _ := twoPair(test.hand)
		if hasTwoPair != test.hasTwoPair {
			t.Errorf("Expected twoPair(%v) to return %v, but instead it returned %v.",
				test.hand,
				test.hasTwoPair,
				hasTwoPair)
		}
	}
}

func TestPair(t *testing.T) {
	tests := []struct {
		hand    []Card
		hasPair bool
	}{
		{
			[]Card{{Two, Diamond}, {Two, Club}, {Four, Spade}, {Five, Heart}, {Six, Club}, {Jack, Club}, {Queen, Heart}},
			true,
		},
		{
			[]Card{{Two, Diamond}, {Ace, Club}, {Four, Spade}, {Eight, Spade}, {Ten, Heart}, {Six, Heart}, {Six, Diamond}},
			true,
		},
		{
			[]Card{{Jack, Heart}, {King, Diamond}, {Ten, Spade}, {Nine, Heart}, {Three, Club}, {Five, Diamond}, {Seven, Club}},
			false,
		},
		{
			[]Card{},
			false,
		},
	}
	for _, test := range tests {
		hasPair, _ := pair(test.hand)
		if hasPair != test.hasPair {
			t.Errorf("Expected twoPair(%v) to return %v, but instead it returned %v.",
				test.hand,
				test.hasPair,
				hasPair)
		}
	}
}

func cardsEqual(a, b []Card) bool {
	if len(a) != len(b) {
		return false
	}
	for i, cardA := range a {
		cardB := b[i]
		if cardA.rank != cardB.rank || cardA.suit != cardB.suit {
			return false
		}
	}
	return true
}
