package main

import (
	"github.com/google/uuid"
)

const (
	_invalidMoveMsg = "invalid move"
	_notYourTurnMsg = "not your turn"
)

// GameId uniquely identifies a game.
type GameId string

// game stores the game state.
type game struct {
	id               GameId
	board            [8][8]piece
	playerIdToPlayer map[PlayerId]*player
	whoseTurn        *player
}

// newGame constructs a new game and adds the first player.
func newGame(p *player) *game {
	game := &game{
		id:               GameId(uuid.NewString()),
		board:            newBoard(),
		playerIdToPlayer: map[PlayerId]*player{p.id: p},
	}
	if p.isWhite {
		game.whoseTurn = p
	}
	game.redrawAllScreens()
	return game
}

// addSecondPlayer adds the second player to the game.
func (g *game) addSecondPlayer(p *player) {
	g.playerIdToPlayer[p.id] = p
	if p.isWhite {
		g.whoseTurn = p
	}
	g.redrawAllScreens()
}

// createNextFrame creates a next frame to be pushed to a player's screen.
func (g *game) createNextFrame(p *player, msg string) *Frame {
	frame := &Frame{
		GameId:        g.id,
		IsGameStarted: len(g.playerIdToPlayer) == 2,
		IsTurn:        g.whoseTurn.id == p.id,
		IsWhite:       p.isWhite,
		Msg:           msg,
	}
	// Transform internal board for display.
	for _, row := range g.board {
		var pieceViewRow []PieceView
		for _, piece := range row {
			// Empty cell.
			if piece == nil {
				pieceViewRow = append(pieceViewRow, PieceView{})
				continue
			}
			pieceViewRow = append(pieceViewRow, PieceView{
				IsWhite: piece.isWhite(),
				Name:    piece.name(),
			})
		}
		frame.Board = append(frame.Board, pieceViewRow)
	}
	return frame
}

// redrawAllScreens redraws the screen for all players.
func (g *game) redrawAllScreens() {
	for _, p := range g.playerIdToPlayer {
		p.nextFrame <- g.createNextFrame(p, "" /*msg*/)
	}
}

// drawScreen redraws a player's screen.
func (g *game) drawScreen(p *player, msg string) {
	p.nextFrame <- g.createNextFrame(p, msg)
}

// handleMove handles a move and redraws player(s) screen with the outcome.
func (g *game) handleMove(m *Move) {
	var player, other_player *player
	for _, p := range g.playerIdToPlayer {
		if p.id == m.PlayerId {
			player = p
		} else {
			other_player = p
		}
	}

	// Handle playing out of turn.
	if g.whoseTurn.id != player.id {
		g.drawScreen(player, _notYourTurnMsg)
		return
	}
	// Handle moving a non existed piece.
	piece := g.board[m.From.Row][m.From.Col]
	if piece == nil {
		g.drawScreen(player, _invalidMoveMsg)
		return
	}
	piece.move(m, &g.board)    // Make the move
	g.whoseTurn = other_player // Switch the turn
	g.redrawAllScreens()
}

// newBoard constructs a new board.
func newBoard() [8][8]piece {
	var board [8][8]piece
	for i := 0; i < 8; i++ {
		board[1][i] = newPawn(false /*white*/)
		board[6][i] = newPawn(true /*white*/)
	}
	board[0][0] = newRook(false /*white*/)
	board[0][1] = newKnight(false /*white*/)
	board[0][2] = newBishop(false /*white*/)
	board[0][3] = newQueen(false /*white*/)
	board[0][4] = newKing(false /*white*/)
	board[0][5] = newBishop(false /*white*/)
	board[0][6] = newKnight(false /*white*/)
	board[0][7] = newRook(false /*white*/)

	board[7][0] = newRook(true /*white*/)
	board[7][1] = newKnight(true /*white*/)
	board[7][2] = newBishop(true /*white*/)
	board[7][3] = newQueen(true /*white*/)
	board[7][4] = newKing(true /*white*/)
	board[7][5] = newBishop(true /*white*/)
	board[7][6] = newKnight(true /*white*/)
	board[7][7] = newRook(true /*white*/)
	return board
}
