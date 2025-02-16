package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sayanli/calculator/internal/entity"
)

func InstructionsHandler(w http.ResponseWriter, r *http.Request) {
	var instructions []entity.Instruction

	err := json.NewDecoder(r.Body).Decode(&instructions)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for i, instr := range instructions {
		fmt.Println(i, instr.Type)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Instructions processed successfully"))
}

func Run() {
	fmt.Println("Calculator is running...")
	r := chi.NewRouter()
	r.Post("/", InstructionsHandler)
	http.ListenAndServe(":8080", r)
}
