package game

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const entitySize = 20
const MaxPlayers = 256

const TypePlayer = 1
const playerViewDistance = 120.0

const MessageTypeFrame = 1
const MessageTypeScores = 2

const MessageInputTypeMove = 1
const MessageInputTypeDivide = 2

type Player struct {
	Name             string
	Token            string
	ID               uint32
	Conn             *websocket.Conn
	sendChan         chan []byte
	Chunks           []*PlayerChunk
	MoveTo           Position
	PlayerCenter     Position
	ColorIndex       uint8
	Speed            float64
	Eaten            bool
	EatenTime        time.Time
	DisconnectedTime time.Time
	MarkedForRemoval bool
	Size             uint16
	MergeBlock       float64
}

type PlayerChunk struct {
	Position  Position
	ShiftTo   Position
	Size      uint16
	SizeTimer float64
	Radius    float64
	ID        uint32
}

type PlayerJoin struct {
	Token      string
	Name       string
	ColorIndex uint8
	Conn       *websocket.Conn
}

type PlayerInput struct {
	Token     string
	RelTarget Position
	Divide    bool
}

type PlayerScoreItem struct {
	ID         uint32
	ColorIndex uint8
	Name       string
	Size       uint16
}

type WSEntity struct {
	Type       uint8
	ColorIndex uint8
	Size       uint16
	ID         uint32
	OwnerID    uint32
	RelPosX    int16
	RelPosY    int16
	RelMoveToX int16
	RelMoveToY int16
}

func (game *Game) AddPlayer(p *PlayerJoin) {
	game.addPlayerChan <- p
}

// method on Player
func (player *Player) readPump(game *Game) {
	defer func() {
		if player.Conn != nil {
			player.Conn.Close()
			player.Conn = nil
			player.DisconnectedTime = time.Now()
		}
	}()

	token := player.Token

	for {
		// Read message
		_, message, err := player.Conn.ReadMessage()
		if err != nil {
			return // client disconnected
		}

		if len(message) == 6 {
			// 0-2 message type
			messageType := int16(binary.BigEndian.Uint16(message[0:2]))

			relX := float64(int16(binary.BigEndian.Uint16(message[2:4]))) / Precision
			relY := float64(int16(binary.BigEndian.Uint16(message[4:6]))) / Precision

			if messageType == MessageInputTypeMove {
				input := PlayerInput{
					Token: token,
					RelTarget: Position{
						X: relX,
						Y: relY,
					},
					Divide: false,
				}

				select {
				case game.inputChan <- &input:
				default: // input queue full, dropping the input
				}
			}

			if messageType == MessageInputTypeDivide {
				input := PlayerInput{
					Token: token,
					RelTarget: Position{
						X: relX,
						Y: relY,
					},
					Divide: true,
				}

				select {
				case game.inputChan <- &input:
				default: // input queue full, dropping the input
				}
			}
		}
	}
}

func (player *Player) writePump() {
	defer func() {
		if player.Conn != nil {
			player.Conn.Close()
		}
	}()

	for msg := range player.sendChan {
		if player.Conn != nil {
			player.Conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))

			err := player.Conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				fmt.Println("Write error:", err, "player: ", player.Name)
				player.Conn.Close()
				player.Conn = nil
				player.DisconnectedTime = time.Now()
				return
			}
		}
	}
}

func (game *Game) buildFrameFor(player *Player) []byte {
	visible := game.getVisibleEntitiesFor(player)
	count := len(visible)

	frame := make([]byte, 12+count*entitySize)
	binary.BigEndian.PutUint16(frame[0:2], MessageTypeFrame) // message type
	binary.BigEndian.PutUint16(frame[2:4], uint16(count))    // entity count
	binary.BigEndian.PutUint32(frame[4:8], player.ID)
	if player.Eaten { // is eaten
		frame[8] = 1
	} else {
		frame[8] = 0
	}

	offset := 12
	for _, entity := range visible {
		frame[offset+0] = entity.Type
		frame[offset+1] = entity.ColorIndex
		binary.BigEndian.PutUint16(frame[offset+2:offset+4], entity.Size)
		binary.BigEndian.PutUint32(frame[offset+4:offset+8], entity.ID)
		binary.BigEndian.PutUint32(frame[offset+8:offset+12], entity.OwnerID)
		binary.BigEndian.PutUint16(frame[offset+12:offset+14], uint16(entity.RelPosX))
		binary.BigEndian.PutUint16(frame[offset+14:offset+16], uint16(entity.RelPosY))
		offset += entitySize
	}

	return frame
}

func (game *Game) buildScore() []byte {
	list := make([]*PlayerScoreItem, 0, 128)
	for _, listPlayer := range game.players {
		if !listPlayer.Eaten {
			list = append(list, &PlayerScoreItem{
				ID:         listPlayer.ID,
				ColorIndex: listPlayer.ColorIndex,
				Name:       listPlayer.Name,
				Size:       listPlayer.Size,
			})
		}
	}
	count := len(list)

	byteSize := 4 // 0 - 2 message type, 2 - 4 entity count
	for _, info := range list {
		byteSize += 8 // ID - 4, Size - 2, Color - 1, NameLength - 1
		byteSize += len(info.Name)
	}

	frame := make([]byte, byteSize)
	binary.BigEndian.PutUint16(frame[0:2], MessageTypeScores) // message type
	binary.BigEndian.PutUint16(frame[2:4], uint16(count))     // entity count

	offset := 4
	for _, entity := range list {
		binary.BigEndian.PutUint32(frame[offset+0:offset+4], entity.ID)
		binary.BigEndian.PutUint16(frame[offset+4:offset+6], entity.Size)
		frame[offset+6] = entity.ColorIndex
		frame[offset+7] = uint8(len(entity.Name))

		// Write the name bytes starting at offset+8
		nameBytes := []byte(entity.Name)
		copy(frame[offset+8:offset+8+len(nameBytes)], nameBytes)

		// Advance the offset for the next entity
		offset += 8 + len(nameBytes)
	}

	return frame
}

func (game *Game) getVisibleEntitiesFor(player *Player) []*WSEntity {
	out := make([]*WSEntity, 0, 128)
	baseX := player.PlayerCenter.X
	baseY := player.PlayerCenter.Y

	viewDistSq := playerViewDistance * playerViewDistance

	for _, listPlayer := range game.players {
		if !listPlayer.Eaten {
			for _, chunk := range listPlayer.Chunks {
				diffX := chunk.Position.X - baseX
				diffY := chunk.Position.Y - baseY
				distSq := diffX*diffX + diffY*diffY

				if distSq < viewDistSq {
					out = append(out, &WSEntity{
						Type:       TypePlayer,
						ColorIndex: listPlayer.ColorIndex,
						Size:       chunk.Size,
						ID:         chunk.ID,
						OwnerID:    listPlayer.ID,
						RelPosX:    int16(diffX * Precision),
						RelPosY:    int16(diffY * Precision),
						RelMoveToX: int16((listPlayer.MoveTo.X - baseX) * Precision),
						RelMoveToY: int16((listPlayer.MoveTo.Y - baseY) * Precision),
					})
				}
			}
		}
	}

	for _, listCrumb := range game.crumbs {
		diffX := listCrumb.Position.X - baseX
		diffY := listCrumb.Position.Y - baseY
		distSq := diffX*diffX + diffY*diffY

		if distSq < viewDistSq {
			out = append(out, &WSEntity{
				Type:       TypeCrumb,
				ColorIndex: 0,
				Size:       listCrumb.Size,
				ID:         listCrumb.ID,
				OwnerID:    listCrumb.ID,
				RelPosX:    int16((listCrumb.Position.X - baseX) * Precision),
				RelPosY:    int16((listCrumb.Position.Y - baseY) * Precision),
				RelMoveToX: 0,
				RelMoveToY: 0,
			})
		}
	}

	return out
}
