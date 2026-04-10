package status

import (
	"encoding/json"
	"net/http"
	"test/names"
	"time"

	"github.com/go-chi/chi/v5"
)

type Response struct {
	ServerRunningMin int      `json:"serverRunningMin"`
	AdvList          []string `json:"advList"`
	NameList         []string `json:"nameList"`
}

func RouteInit(router *chi.Mux) {
	start := time.Now()

	router.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		currentTyme := time.Now()
		diff := currentTyme.Sub(start)

		resp := Response{
			ServerRunningMin: int(diff.Minutes()),
			AdvList:          names.Adv,
			NameList:         names.Name,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
}
