package handlers

import (
	"fmt"
	"log"
	"net/http"
	"test/game"
	"test/tokens"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO add domain restrictions
		return true
	},
}

func RouteInit(router *chi.Mux, tokenManager *tokens.Manager, gameManager *game.Game) {
	router.Get("/ws/{token}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CONNECTION ATTEMPT - Method:", r.Method)
		fmt.Println("Upgrade header:", r.Header.Get("Upgrade"))
		fmt.Println("Connection header:", r.Header.Get("Connection"))

		token := chi.URLParam(r, "token")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		info, ok := tokenManager.Validate(token)
		if !ok {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Upgrade http to WS
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("websocket upgrade failed: %v", err)
			return
		}

		tokenManager.Remove(token)

		player := &game.PlayerJoin{
			Token:      token,
			Name:       info.Name,
			ColorIndex: info.ColorIndex,
			Conn:       conn,
		}
		gameManager.AddPlayer(player)
	})
}
