package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"test/names"
	"test/tokens"

	"github.com/go-chi/chi/v5"
)

type Response struct {
	Token string `json:"token"`
}

func RouteInit(router *chi.Mux, manager *tokens.Manager) {
	router.Get("/auth/{value}", func(w http.ResponseWriter, r *http.Request) {
		value := chi.URLParam(r, "value")
		var adv, name string
		var colorIndex uint8 = 0

		parts := strings.SplitN(value, ":", 3)

		if len(parts) == 2 {
			name = parts[0]
			u64, _ := strconv.ParseUint(parts[1], 0, 8)
			colorIndex = uint8(u64)
		} else if len(parts) == 3 {
			adv = parts[0]
			name = parts[1]
			u64, _ := strconv.ParseUint(parts[2], 0, 8)
			colorIndex = uint8(u64)
		} else {
			http.Error(w, "bad name", http.StatusBadRequest)
			fmt.Println("Bad name", adv, name, colorIndex)
			return
		}
		fmt.Println("Name", adv, name, colorIndex)

		var isNameCorrect bool
		if adv == "" {
			isNameCorrect = names.CheckName(name)
		} else {
			isNameCorrect = names.CheckAdvName(adv, name)
		}
		if !isNameCorrect {
			http.Error(w, "bad name", http.StatusBadRequest)
			return
		}
		var fullName string
		if adv == "" {
			fullName = name
		} else {
			fullName = adv + " " + name
		}

		token := manager.AddNewUser(fullName, colorIndex)

		resp := Response{
			Token: token,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
}
