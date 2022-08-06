package game

import (
	"errors"
	"fmt"

	"github.com/Chris-Behan/gopoker/cards"
)

type player struct {
	id               int // id of the player which is the same as where they are seated at the table
	hand             [2]cards.Card
	money            int
	alive            bool // whether or not the player is still in the game
	amountBetInRound int  // amount the player has bet in the current round
}

type gamePhase int8

const (
	preFlop  gamePhase = 0
	flop     gamePhase = 1
	turn     gamePhase = 2
	river    gamePhase = 3
	showdown gamePhase = 4
)

type GameState struct {
	table             []player // players playing at the table
	bigBlindAmount    int
	smallBlindAmount  int
	bigBlindPos       int // index of table where the big blind is
	smallBlindPos     int // index of table where the small blind is
	pot               int // Amount of money in the pot
	highestBetInRound int // Highest betting amount of the current round
	whoseTurn         int // id of the player whose turn it is
	phase             gamePhase
	participating     []int // id of players participating in the round
	betInCurrentRound bool  // whether or not there has been a bet in the current round (round being preflop, flop, turn, etc)
}

func NewGame(numPlayers int, playerCash int, bigBlindAmt int) GameState {
	game := GameState{[]player{}, bigBlindAmt, bigBlindAmt / 2, 1, 0, 0, 0, 0, preFlop, []int{}, false}
	for i := 0; i < numPlayers; i++ {
		p := player{i, [2]cards.Card{}, playerCash, true, 0}
		game.table = append(game.table, p)
	}

	return game
}

func gameLoop() {

}

func (g *GameState) newRound() {
	g.phase = preFlop
	g.addAllPlayers()
	g.dealCards()
	g.handleBlinds()
	g.whoseTurn = g.participantClockwiseToPlayer(g.bigBlindPos)
}

// Adds all players to the GameState.participating slice.
func (g *GameState) addAllPlayers() {
	ids := []int{}
	for _, p := range g.alivePlayers() {
		ids = append(ids, p.id)
	}
	g.participating = ids
}

// Returns a slice of players still in the game.
func (g *GameState) alivePlayers() []player {
	players := []player{}
	for _, p := range g.table {
		if p.alive {
			players = append(players, p)
		}
	}
	return players
}

// Deal cards to all alive players. Assumes that every alive player is in the participating slice
// and that ONLY alive players are in the participating slice.
func (g *GameState) dealCards() {
	playerIDs := g.participating
	deck := cards.GenerateDeck()
	numPlayers := len(playerIDs)
	cardsDealt := 0
	cardsToDeal := numPlayers * 2
	playerIdx := 0
	for cardsDealt < cardsToDeal {
		playerToDealTo := playerIDs[playerIdx]
		card, err := deck.Draw()
		if err != nil {
			panic(err)
		}
		cardIdx := 0
		if cardsDealt >= numPlayers {
			cardIdx = 1
		}
		g.table[playerToDealTo].hand[cardIdx] = card
		if playerIdx == numPlayers-1 {
			playerIdx = 0
		} else {
			playerIdx++
		}
		cardsDealt++
	}
}

func (g *GameState) handleBlinds() {
	// deduct blinds from players and add to pot
	g.table[g.smallBlindPos].money -= g.smallBlindAmount
	g.pot += g.smallBlindAmount
	g.table[g.bigBlindPos].money -= g.bigBlindAmount
	g.pot += g.bigBlindAmount
}

// Returns the next participating player clockwise to the specified player.
func (g GameState) participantClockwiseToPlayer(playerID int) int {
	id := g.getClockwisePlayerID(playerID)
	for {
		if !intInSlice(id, g.participating) {
			id = g.getClockwisePlayerID(id)
		} else {
			return id
		}
	}
}

// Returns the id of the player clockwise to the player ID provided.
func (g GameState) getClockwisePlayerID(from int) int {
	if from+1 == len(g.table) {
		return 0
	}
	return from + 1
}

// Check checks for the specified player or returns an error if the player cannot check.
func (g *GameState) Check(playerID int) error {
	err := g.validateCheck(playerID)
	if err != nil {
		return fmt.Errorf("error checking: %v", err)
	}
	// handle turn end
	return nil
}

// Fold removes the specified player from the current round.
func (g *GameState) Fold(playerID int) error {
	if playerID != g.whoseTurn {
		return fmt.Errorf("error folding player %v because it is player %v's turn", playerID, g.whoseTurn)
	}
	newParticipating, err := removeIntFromSlice(g.participating, playerID)
	if err != nil {
		return fmt.Errorf("error folding for player %v: %v", playerID, err)
	}
	g.participating = newParticipating
	// handle turn end
	return nil
}

// Bet makes the first wager of the round. Only possible during the flop, turn, or river.
func (g *GameState) Bet(playerID int, amount int) error {
	err := g.validateBet(playerID, amount)
	if err != nil {
		return fmt.Errorf("error betting: %v", err)
	}

	g.table[playerID].money -= amount
	g.table[playerID].amountBetInRound += amount
	g.pot += amount
	g.betInCurrentRound = true
	g.highestBetInRound = amount

	g.whoseTurn = g.getNextPlayersTurn()
	// handle turn end
	return nil
}

// Call matches the current bet.
func (g *GameState) Call(playerID int) error {
	err := g.validateCall(playerID)
	if err != nil {
		return fmt.Errorf("error calling: %v", err)
	}

	callAmount := g.callAmount(playerID)
	g.table[playerID].money -= callAmount
	g.table[playerID].amountBetInRound += callAmount
	g.pot += callAmount

	// handle turn end
	return nil
}

// Raise increases the current bet.
func (g *GameState) Raise(playerID int, amount int) error {
	err := g.validateRaise(playerID, amount)
	if err != nil {
		return fmt.Errorf("error raising: %v", err)
	}
	// amount player is betting is call + raise
	betAmount := g.callAmount(playerID) + amount
	g.table[playerID].money -= betAmount
	g.table[playerID].amountBetInRound += betAmount
	g.pot += betAmount
	g.highestBetInRound = g.table[playerID].amountBetInRound

	// handle turn end
	return nil
}

func (g GameState) validateCheck(playerID int) error {
	if playerID != g.whoseTurn {
		return errors.New(notYourTurnMsg(playerID, g.whoseTurn))
	}
	if g.highestBetInRound > g.table[playerID].amountBetInRound {
		return fmt.Errorf("cannot check, instead you must call or raise the current betting amount of %v",
			g.highestBetInRound)
	}
	return nil
}

func (g GameState) validateBet(playerID int, amount int) error {
	if playerID != g.whoseTurn {
		return errors.New(notYourTurnMsg(playerID, g.whoseTurn))
	}
	if g.betInCurrentRound {
		return fmt.Errorf("can only Bet if there hasn't been a bet this round. If you wish to increase the bet, call Raise")
	}
	minBet := g.minimumBet()
	if amount < minBet {
		return fmt.Errorf("minimum bet is $%v", minBet)
	}
	playersMoney := g.table[playerID].money
	if playersMoney < amount {
		return fmt.Errorf("player %v does not have enough money to bet $%v (they only have $%v)",
			playerID, amount, playersMoney)
	}

	return nil
}

func (g GameState) validateCall(playerID int) error {
	if playerID != g.whoseTurn {
		return errors.New(notYourTurnMsg(playerID, g.whoseTurn))
	}
	player := g.table[playerID]
	amountToCall := g.callAmount(playerID)
	if amountToCall > player.money {
		return fmt.Errorf("player %v does not have enough money to call $%v (they only have $%v)",
			player.id,
			amountToCall,
			player.money)
	}

	return nil
}

func (g GameState) validateRaise(playerID int, amount int) error {
	if playerID != g.whoseTurn {
		return errors.New(notYourTurnMsg(playerID, g.whoseTurn))
	}
	minRaise := g.minimumRaise()
	if amount < minRaise {
		return fmt.Errorf("minimum raise is $%v", minRaise)
	}
	// amount to call + raise
	callAmount := g.callAmount(playerID)
	totalAmount := callAmount + amount
	player := g.table[playerID]
	if totalAmount > player.money {
		return fmt.Errorf("player %v does not have enough money to raise, needs $%v ($%v to call plus $%v raise), but only has $%v",
			playerID, totalAmount, callAmount, amount, player.money)
	}
	return nil
}

func (g GameState) callAmount(playerID int) int {
	return g.highestBetInRound - g.table[playerID].amountBetInRound
}

func notYourTurnMsg(playerWhoTriedToMakeMove int, whoseTurn int) string {
	return fmt.Sprintf("it is not player %v's turn, it is player %v's turn",
		playerWhoTriedToMakeMove,
		whoseTurn)
}

// Updates game state to the next players turn.
func (g *GameState) nextPlayersTurn() {
	// prevPlayersTurn := g.getTablePos(g.whoseTurn)
	if g.getTablePos(g.whoseTurn) == len(g.table)-1 {
		g.whoseTurn = g.table[0].id
	} else {
		g.whoseTurn++
	}
}

func (g GameState) getNextPlayersTurn() int {
	return g.participantClockwiseToPlayer(g.whoseTurn)
}

// Returns the player whose turn it is.
func (g GameState) getWhoseTurn() player {
	tablePos := g.getTablePos(g.whoseTurn)
	return g.table[tablePos]
}

func (g GameState) minimumBet() int {
	return g.bigBlindAmount
}

func (g GameState) minimumRaise() int {
	return g.highestBetInRound
}

// Returns the index of the table where the player with the specified id is or -1 if no player is found.
func (g GameState) getTablePos(player_id int) int {
	for idx, p := range g.table {
		if p.id == player_id {
			return idx
		}
	}
	return -1
}

func intInSlice(i int, s []int) bool {
	for _, j := range s {
		if i == j {
			return true
		}
	}
	return false
}

// Removes the specified int from a slice of ints and returns the result.
// If the int is not found in the slice, an error is returned.
// Note that the resulting slice may not be in the same order as it was passed in.
func removeIntFromSlice(nums []int, n int) ([]int, error) {
	idxOfInt := -1
	for idx, val := range nums {
		if val == n {
			idxOfInt = idx
			break
		}
	}
	if idxOfInt == -1 {
		return []int{}, fmt.Errorf("couldn't find specified int in the slice: %v is not in %v", n, nums)
	}

	nums[idxOfInt] = nums[len(nums)-1]
	nums[len(nums)-1] = 0
	nums = nums[:len(nums)-1]
	return nums, nil
}
