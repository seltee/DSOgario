package game

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

type Game struct {
	players       map[string]*Player
	addPlayerChan chan *PlayerJoin
	inputChan     chan *PlayerInput
	crumbs        []*Crumb
	crumbsTimer   float64
	fieldSize     float64
}

func New() *Game {
	game := &Game{
		players:       make(map[string]*Player),
		addPlayerChan: make(chan *PlayerJoin, 4),
		inputChan:     make(chan *PlayerInput, MaxPlayers),
	}

	go game.Run()
	return game
}

func (game *Game) Run() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	lastTime := time.Now()
	game.fieldSize = 160.0

	for {
		now := time.Now()
		delta := now.Sub(lastTime).Seconds()
		lastTime = now

		game.processPendingPlayers()

		game.processNewCrumbs(delta)

		game.processInputs()

		game.updateWorld(delta)

		game.broadcastFrame()

		<-ticker.C
	}
}

func (game *Game) processPendingPlayers() {
	for {
		select {
		case newPlayer := <-game.addPlayerChan:
			position := genFieldPosition(game.fieldSize)
			player := &Player{
				Name:       newPlayer.Name,
				Token:      newPlayer.Token,
				Conn:       newPlayer.Conn,
				ColorIndex: newPlayer.ColorIndex,
				sendChan:   make(chan []byte, 1024),
				Position:   position,
				MoveTo:     position,
				Size:       10,
				Radius:     sizeToRadius(10),
				Speed:      sizeToSpeed(10),
				ID:         GetNextID(),
			}
			game.players[newPlayer.Token] = player

			go player.readPump(game)
			go player.writePump()

			fmt.Println("PLAYER ADDED", newPlayer.Name)
		default:
			return
		}
	}
}

func (game *Game) processNewCrumbs(delta float64) {
	game.crumbsTimer -= delta * float64(crumbsPerSecond)
	if game.crumbsTimer < 0.0 {
		game.crumbsTimer += 1.0
		size := uint16(rand.IntN(2) + 1)
		game.crumbs = append(game.crumbs, &Crumb{
			Position: genFieldPosition(game.fieldSize),
			Size:     size,
			Radius:   sizeToRadius(size),
			ID:       GetNextID(),
		})
	}
}

func (game *Game) processInputs() {
	for {
		select {
		case input := <-game.inputChan:
			if player, ok := game.players[input.Token]; ok {
				player.MoveTo = Position{
					X: player.Position.X + input.RelMoveTo.X,
					Y: player.Position.Y + input.RelMoveTo.Y,
				}
			}

		default:
			return
		}
	}
}

func (game *Game) updateWorld(delta float64) {
	for _, player := range game.players {
		distToMoveX := player.MoveTo.X - player.Position.X
		distToMoveY := player.MoveTo.Y - player.Position.Y
		dist2 := distToMoveX*distToMoveX + distToMoveY*distToMoveY
		if dist2 > 0.001 {
			dist := math.Sqrt(dist2)
			if dist < player.Speed*delta {
				player.Position.X = player.MoveTo.X
				player.Position.Y = player.MoveTo.Y
			} else {
				player.Position.X += (distToMoveX / dist) * player.Speed * delta
				player.Position.Y += (distToMoveY / dist) * player.Speed * delta
			}
		}
	}

	i := 0
	for i < len(game.crumbs) {
		crumb := game.crumbs[i]
		eaten := false

		for _, player := range game.players {
			diffX := crumb.Position.X - player.Position.X
			diffY := crumb.Position.Y - player.Position.Y
			distSq := diffX*diffX + diffY*diffY
			eatDist := math.Max(player.Radius, crumb.Radius)

			if distSq < eatDist*eatDist {
				// Crumb eaten!
				player.Size += uint16(crumb.Size)
				player.Radius = sizeToRadius(player.Size)
				eaten = true
				break
			}
		}

		if eaten {
			// Remove crumb
			last := len(game.crumbs) - 1
			game.crumbs[i] = game.crumbs[last]
			game.crumbs = game.crumbs[:last]
		} else {
			i++
		}
	}
}

func (game *Game) broadcastFrame() {
	for _, player := range game.players {
		frame := game.buildFrameFor(player)
		if frame == nil {
			continue
		}

		select {
		case player.sendChan <- frame:
			// data frame sent
		default:
			// data frame dropped
		}
	}
}

func (game *Game) removeCrumb(index int) {
	last := len(game.crumbs) - 1
	game.crumbs[index] = game.crumbs[last]
	game.crumbs = game.crumbs[:last]
}

func (game *Game) Stop() {
}
