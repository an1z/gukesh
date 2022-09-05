package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to connected players.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from a player.
	pongWait = 60 * time.Second

	// Send pings to players with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from a player.
	maxMessageSize = 512
)

// PlayerId uniquely identifies a player.
type PlayerId string

// Player stores the state of a player.
type player struct {
	id        PlayerId
	gameId    GameId // Empty string if not yet connected to a game
	isWhite   bool
	nextFrame chan *Frame // The next frame to be rendered on player's screen

	// Connections.
	gm   *gameManager
	conn *websocket.Conn
}

// PieceView is the player facing view of a piece.
type PieceView struct {
	IsWhite bool   `json:"is_white"`
	Name    string `json:"name"`
}

// Frame holds the contents of a players current screen.
type Frame struct {
	GameId        GameId        `json:"game_id"`
	IsGameStarted bool          `json:"is_game_started"`
	IsWhite       bool          `json:"is_white"`
	IsTurn        bool          `json:"is_turn"`
	Board         [][]PieceView `json:"board"`
	Msg           string        `json:"msg"`
}

// PlayerAction is an action performed by a player during a turn. The fields are
// evaluated top to bottom by the game loop to resolve which action was
// performed. Note that actions are mutually exclusive.
type PlayerAction struct {
	StartGame bool   `json:"start_game"`
	JoinGame  GameId `json:"join_game"`
	Move      Move   `json:"move"`
}

// Move denotes a move by a player.
type Move struct {
	PlayerId PlayerId `json:"-"`
	From     Pos      `json:"from"`
	To       Pos      `json:"to"`
}

// Pos denotes a piece's position on the chess board.
type Pos struct {
	Row int8 `json:"row"`
	Col int8 `json:"col"`
}

// newPlayer creates a new player.
func newPlayer(gm *gameManager, conn *websocket.Conn) *player {
	return &player{
		id:        PlayerId(uuid.NewString()),
		gm:        gm,
		conn:      conn,
		nextFrame: make(chan *Frame),
	}
}

// gameReadLoop handles incoming actions from a player.
func (p *player) gameReadLoop() {
	defer func() {
		p.gm.leaveGame(p.id)
		p.conn.Close()
	}()
	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error {
		p.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("conn closed: %v", err)
			}
			break
		}
		var playerAction PlayerAction
		if err := json.Unmarshal(message, &playerAction); err != nil {
			log.Fatalf("error: %v", err)
		}
		if playerAction.StartGame {
			p.isWhite = true // the host player gets white
			p.gm.startNewGame(p)
		} else if playerAction.JoinGame != "" {
			p.gm.addSecondPlayer(p, playerAction.JoinGame)
		} else {
			// Tell game manager to handle the move.
			playerAction.Move.PlayerId = p.id
			p.gm.incomingMoves <- &playerAction.Move
		}
	}
}

// gameWriteLoop pushes new frames to a player when game state changes.
func (p *player) gameWriteLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.conn.Close()
	}()
	for {
		select {
		case frame, ok := <-p.nextFrame:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Game manager closed the channel.
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Fatal("error acquiriting writer", err)
			}
			bytes, err := json.Marshal(frame)
			if err != nil {
				log.Fatal("serialization failed")
			}
			w.Write(bytes)
			if err := w.Close(); err != nil {
				log.Fatal(err)
			}

		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("error sending ping", err)
				return
			}
		}
	}
}

// StartGameSession registers a new player and initializes a websocket
// connection for it.
func startGameSession(gm *gameManager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	player := newPlayer(gm, conn)
	go player.gameReadLoop()
	go player.gameWriteLoop()
}
