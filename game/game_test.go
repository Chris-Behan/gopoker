package game

import (
	"testing"
)

func TestNewRound(t *testing.T) {
	gameState := NewGame(5, 100, 4)
	gameState.newRound()
	// fmt.Printf("GameState: %v", gameState)
	t.Logf("GameState: %v", gameState)
	gameState.newRound()
	t.Logf("GameState: %v", gameState)
}
