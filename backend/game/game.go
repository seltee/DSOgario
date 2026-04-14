package game

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

type Game struct {
	players              map[string]*Player
	addPlayerChan        chan *PlayerJoin
	inputChan            chan *PlayerInput
	crumbs               []*Crumb
	crumbsTimer          float64
	fieldSize            float64
	broadcastScoreTicker int
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

		game.broadcastScoreTicker++
		if game.broadcastScoreTicker > 200 {
			game.broadcastScoreTicker = 0
			game.broadcastScore()
		}

		<-ticker.C
	}
}

func (game *Game) processPendingPlayers() {
	for {
		select {
		case newPlayer := <-game.addPlayerChan:
			position := genFieldPosition(game.fieldSize)
			player := &Player{
				Name:             newPlayer.Name,
				Token:            newPlayer.Token,
				Conn:             newPlayer.Conn,
				ColorIndex:       newPlayer.ColorIndex,
				sendChan:         make(chan []byte, 1024),
				Eaten:            false,
				MarkedForRemoval: false,
				Chunks:           make([]*PlayerChunk, 0, 8),
				MoveTo:           position,
				Speed:            sizeToSpeed(10),
				Size:             10,
				PlayerCenter:     position,
				ID:               GetNextID(),
			}
			player.Chunks = append(player.Chunks, &PlayerChunk{
				Position:  position,
				ShiftTo:   Position{X: 0, Y: 0},
				Radius:    sizeToRadius(10),
				Size:      10,
				SizeTimer: 0,
				ID:        GetNextID(),
			})

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
		if len(game.crumbs) < 800 {
			size := uint16(rand.IntN(2) + 1)
			game.crumbs = append(game.crumbs, &Crumb{
				Position: genFieldPosition(game.fieldSize),
				Size:     size,
				Radius:   sizeToRadius(size),
				ID:       GetNextID(),
			})
		}
	}
}

func (game *Game) processInputs() {
	for {
		select {
		case input := <-game.inputChan:
			if player, ok := game.players[input.Token]; ok {
				if input.Divide {
					forDivision := len(player.Chunks)
					for i := 0; i < forDivision; i++ {
						chunk := player.Chunks[i]
						size := chunk.Size
						if size >= 20 {
							chunk.Size = size / 2
							chunk.Radius = sizeToRadius(chunk.Size)
							newChunkSize := size - chunk.Size

							distToMoveX := player.MoveTo.X - chunk.Position.X
							distToMoveY := player.MoveTo.Y - chunk.Position.Y
							dist2 := distToMoveX*distToMoveX + distToMoveY*distToMoveY
							if dist2 > 0.001 {
								dist := math.Sqrt(dist2)
								shiftX := (distToMoveX / dist)
								shiftY := (distToMoveY / dist)

								player.Chunks = append(player.Chunks, &PlayerChunk{
									Position: Position{
										X: chunk.Position.X,
										Y: chunk.Position.Y,
									},
									ShiftTo: Position{
										X: shiftX,
										Y: shiftY,
									},
									Size:   newChunkSize,
									Radius: sizeToRadius(newChunkSize),
									ID:     GetNextID(),
								})
							}
						}
					}
					player.MergeBlock = 1.0
				} else {
					for _, chunk := range player.Chunks {
						player.MoveTo = Position{
							X: chunk.Position.X + input.RelTarget.X,
							Y: chunk.Position.Y + input.RelTarget.Y,
						}
					}
				}
			}

		default:
			return
		}
	}
}

func (game *Game) updateWorld(delta float64) {
	timeNow := time.Now()

	for _, player := range game.players {
		if player.Eaten {
			continue
		}

		if player.MergeBlock > 0 {
			player.MergeBlock -= delta
		} else {
			for i := 0; i < len(player.Chunks); i++ {
				chunk := player.Chunks[i]
				for check := i + 1; check < len(player.Chunks); check++ {
					chunkCheck := player.Chunks[check]

					diffX := chunk.Position.X - chunkCheck.Position.X
					diffY := chunk.Position.Y - chunkCheck.Position.Y
					dist2 := diffX*diffX + diffY*diffY
					maxRadius := math.Max(chunk.Radius, chunkCheck.Radius)
					if dist2 < maxRadius*maxRadius {
						chunk.Size += chunkCheck.Size
						chunk.Radius = sizeToRadius(chunk.Size)
						last := len(player.Chunks) - 1
						player.Chunks[check] = player.Chunks[last]
						player.Chunks = player.Chunks[:last]
						break
					}
				}
			}
		}

		// merge with itself
		player.Size = 0
		for _, chunk := range player.Chunks {
			// reduce chunk size
			chunk.SizeTimer += delta
			reduceCount := 0.5 + 200.0/float64(chunk.Size)
			if chunk.SizeTimer > reduceCount {
				chunk.SizeTimer = 0
				if chunk.Size > 10 {
					chunk.Size--
					chunk.Radius = sizeToRadius(chunk.Size)
				}
			}

			// move chunk
			distToMoveX := player.MoveTo.X - chunk.Position.X
			distToMoveY := player.MoveTo.Y - chunk.Position.Y
			dist2 := distToMoveX*distToMoveX + distToMoveY*distToMoveY
			if dist2 > 0.001 {
				dist := math.Sqrt(dist2)
				if dist < player.Speed*delta {
					chunk.Position.X = player.MoveTo.X
					chunk.Position.Y = player.MoveTo.Y
				} else {
					chunk.Position.X += (distToMoveX / dist) * player.Speed * delta
					chunk.Position.Y += (distToMoveY / dist) * player.Speed * delta
				}
			}
			if math.Abs(chunk.ShiftTo.X) > 0.001 {
				chunk.Position.X += chunk.ShiftTo.X * delta * player.Speed * 10
				chunk.ShiftTo.X = chunk.ShiftTo.X - chunk.ShiftTo.X*delta*4
			}
			if math.Abs(chunk.ShiftTo.Y) > 0.001 {
				chunk.Position.Y += chunk.ShiftTo.Y * delta * player.Speed * 10
				chunk.ShiftTo.Y = chunk.ShiftTo.Y - chunk.ShiftTo.Y*delta*4
			}

			// eat other players
			for _, playerCheck := range game.players {
				if !playerCheck.Eaten && playerCheck != player {

					i := 0
					for i < len(playerCheck.Chunks) {
						chunkCheck := playerCheck.Chunks[i]
						eaten := false

						if chunk.Size > uint16(float64(chunkCheck.Size)*1.1) {
							distToX := chunkCheck.Position.X - chunk.Position.X
							distToY := chunkCheck.Position.Y - chunk.Position.Y
							dist2 := distToX*distToX + distToY*distToY
							if dist2 < chunk.Radius*chunk.Radius {
								eaten = true
								chunk.Size += chunkCheck.Size
								chunk.Radius = sizeToRadius(chunk.Size)
								if len(playerCheck.Chunks) <= 1 {
									playerCheck.Eaten = true
									playerCheck.EatenTime = timeNow
								} else {
									last := len(playerCheck.Chunks) - 1
									game.crumbs[i] = game.crumbs[last]
									game.crumbs = game.crumbs[:last]
								}
							}
						}

						if eaten {
							last := len(playerCheck.Chunks) - 1
							playerCheck.Chunks[i] = playerCheck.Chunks[last]
							playerCheck.Chunks = playerCheck.Chunks[:last]
						} else {
							i++
						}
					}
				}
			}

			// calc player size
			player.Size += chunk.Size
		}

		for _, player := range game.players {
			centerX := 0.0
			centerY := 0.0
			for _, chunk := range player.Chunks {
				centerX += chunk.Position.X
				centerY += chunk.Position.Y
			}
			centerX /= float64(len(player.Chunks))
			centerY /= float64(len(player.Chunks))
			player.PlayerCenter.X = centerX
			player.PlayerCenter.Y = centerY
		}
	}

	for key, player := range game.players {
		// Mark to remove when eaten and not watching
		if player.Eaten {
			if player.Conn == nil {
				player.MarkedForRemoval = true

			}
			if timeNow.Sub(player.EatenTime).Seconds() > 30 {
				player.MarkedForRemoval = true
			}
		}

		// Mark to remove when disconnected for too long
		if player.Conn == nil {
			if timeNow.Sub(player.DisconnectedTime).Seconds() > 180 {
				player.MarkedForRemoval = true
			}
		}

		// remove player and close connection
		if player.MarkedForRemoval {
			if player.Conn != nil {
				player.Conn.Close()
				player.Conn = nil
			}
			delete(game.players, key)
		}
	}

	i := 0
	for i < len(game.crumbs) {
		crumb := game.crumbs[i]
		eaten := false

		for _, player := range game.players {
			if !player.Eaten {
				for _, chunk := range player.Chunks {
					diffX := crumb.Position.X - chunk.Position.X
					diffY := crumb.Position.Y - chunk.Position.Y
					distSq := diffX*diffX + diffY*diffY
					eatDist := math.Max(chunk.Radius, crumb.Radius)

					if distSq < eatDist*eatDist {
						// Crumb eaten!
						chunk.Size += uint16(crumb.Size)
						chunk.Radius = sizeToRadius(player.Size)
						eaten = true
						break
					}
				}
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

func (game *Game) broadcastScore() {
	score := game.buildScore()
	for _, player := range game.players {
		select {
		case player.sendChan <- score:
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
