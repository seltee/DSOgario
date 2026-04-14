package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

	fs := http.FileServer(http.Dir("../frontend/build/web"))
	router.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Started server, port :%s", port)

	error := http.ListenAndServe(":8080", router)
	if error != nil {
		log.Fatal("Error running server ", error)
	}

	fmt.Println("Closing server app ...")
}
