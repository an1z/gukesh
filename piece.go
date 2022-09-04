package main

import "errors"

// piece is a chess piece.
type piece interface {
	// isWhite returns true if the piece is white.
	isWhite() bool

	// Name returns the name of the piece.
	name() string

	// move moves the piece on the board. Returns an error if the move was
	// unsuccessful.
	move(m *Move, board *[8][8]piece) error
}

// checkBounds returns an error if move is out of bounds.
func (m *Move) checkBounds() error {
	if m.To.Row < 0 || m.To.Row > 7 || m.To.Col < 0 || m.To.Col > 7 {
		return errors.New(_invalidMoveMsg)
	}
	return nil
}

type pawn struct {
	white   bool
	nameVal string
}

type king struct {
	white   bool
	nameVal string
}

type queen struct {
	white   bool
	nameVal string
}

type knight struct {
	white   bool
	nameVal string
}

type bishop struct {
	white   bool
	nameVal string
}

type rook struct {
	white   bool
	nameVal string
}

func newPawn(white bool) piece {
	return &pawn{white: white, nameVal: "pawn"}
}

func newKing(white bool) piece {
	return &king{white: white, nameVal: "king"}
}
func newQueen(white bool) piece {
	return &queen{white: white, nameVal: "queen"}
}

func newKnight(white bool) piece {
	return &knight{white: white, nameVal: "knight"}
}

func newBishop(white bool) piece {
	return &bishop{white: white, nameVal: "bishop"}
}

func newRook(white bool) piece {
	return &rook{white: white, nameVal: "rook"}
}

func (p *pawn) isWhite() bool {
	return p.white
}

func (p *pawn) name() string {
	return p.nameVal
}

func (p *queen) isWhite() bool {
	return p.white
}

func (p *queen) name() string {
	return p.nameVal
}

func (p *king) isWhite() bool {
	return p.white
}

func (p *king) name() string {
	return p.nameVal
}

func (p *rook) isWhite() bool {
	return p.white
}

func (p *rook) name() string {
	return p.nameVal
}

func (p *knight) isWhite() bool {
	return p.white
}

func (p *knight) name() string {
	return p.nameVal
}

func (p *bishop) isWhite() bool {
	return p.white
}

func (p *bishop) name() string {
	return p.nameVal
}

func (p *pawn) move(m *Move, board *[8][8]piece) error {
	if err := m.checkBounds(); err != nil {
		return err
	}
	board[m.From.Row][m.From.Col] = nil
	board[m.To.Row][m.To.Col] = p
	return nil
}

func (p *king) move(m *Move, board *[8][8]piece) error {
	if err := m.checkBounds(); err != nil {
		return err
	}
	board[m.From.Row][m.From.Col] = nil
	board[m.To.Row][m.To.Col] = p
	return nil
}

func (p *queen) move(m *Move, board *[8][8]piece) error {
	if err := m.checkBounds(); err != nil {
		return err
	}
	board[m.From.Row][m.From.Col] = nil
	board[m.To.Row][m.To.Col] = p
	return nil
}

func (p *rook) move(m *Move, board *[8][8]piece) error {
	if err := m.checkBounds(); err != nil {
		return err
	}
	board[m.From.Row][m.From.Col] = nil
	board[m.To.Row][m.To.Col] = p
	return nil
}

func (p *bishop) move(m *Move, board *[8][8]piece) error {
	if err := m.checkBounds(); err != nil {
		return err
	}
	board[m.From.Row][m.From.Col] = nil
	board[m.To.Row][m.To.Col] = p
	return nil
}

func (p *knight) move(m *Move, board *[8][8]piece) error {
	if err := m.checkBounds(); err != nil {
		return err
	}
	board[m.From.Row][m.From.Col] = nil
	board[m.To.Row][m.To.Col] = p
	return nil
}
