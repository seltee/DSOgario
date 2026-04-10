package main

import (
	"fmt"
	"log"
	"net/http"
	"test/auth"
	"test/game"
	"test/handlers"
	"test/status"
	"test/tokens"

	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("Starting server app ...")

	router := chi.NewRouter()

	tokenMgr := tokens.New()
	defer tokenMgr.Stop()

	game := game.New()
	defer game.Stop()

	status.RouteInit(router)
	handlers.RouteInit(router, tokenMgr, game)
	auth.RouteInit(router, tokenMgr)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<h1>Home Page</h1>
					<p>Status: 200 OK</p>
				</body>
			</html>
		`))
	})

	// Custom 404 handler
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`
			<html>
				<body>
					<h1>404 - Page Not Found</h1>
				</body>
			</html>
		`))
	})

	fmt.Println("Server is started")
	error := http.ListenAndServe(":8080", router)
	if error != nil {
		log.Fatal("Error running server ", error)
	}

	fmt.Println("Closing server app ...")
}
