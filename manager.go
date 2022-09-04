package main

import (
	"log"
	"sync"
)

// gameManager stores the state of all games, and mediates moves for all games.
type gameManager struct {
	// Protects access to the maps.
	mu sync.Mutex

	gameIdToGame   map[GameId]*game
	playerIdToGame map[PlayerId]*game
	incomingMoves  chan *Move // Moves queued by the players
}

// newGameManager creates a new gameManager.
func newGameManager() *gameManager {
	return &gameManager{
		gameIdToGame:   make(map[GameId]*game),
		playerIdToGame: make(map[PlayerId]*game),
		incomingMoves:  make(chan *Move),
	}
}

// startNewGame starts a new game and adds the first player to it.
func (gm *gameManager) startNewGame(p *player) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game := newGame(p)
	gm.gameIdToGame[game.id] = game
	gm.playerIdToGame[p.id] = game
}

// addSecondPlayer adds a second player to a game.
func (gm *gameManager) addSecondPlayer(p *player, gameId GameId) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, ok := gm.gameIdToGame[gameId]
	if !ok {
		log.Fatal("map lookup failed")
	}
	game.addSecondPlayer(p)
	gm.playerIdToGame[p.id] = game
}

// leaveGame closes the game. It closes the game for both players.
func (gm *gameManager) leaveGame(playerId PlayerId) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, ok := gm.playerIdToGame[playerId]
	if !ok {
		// Game is already closed.
		return
	}
	for _, player := range game.playerIdToPlayer {
		// Ends the game write loop for player.
		close(player.nextFrame)
		delete(gm.playerIdToGame, player.id)
	}
	delete(gm.gameIdToGame, game.id)
}

func (gm *gameManager) runGames() {
	// Single threaded game manager for now.
	for move := range gm.incomingMoves {
		game, ok := gm.playerIdToGame[move.PlayerId]
		if !ok {
			log.Fatal("map lookup failed")
		}
		game.handleMove(move)
	}
}
