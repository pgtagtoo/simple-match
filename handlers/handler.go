package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/genepg/simple-match-api/models"
	"github.com/genepg/simple-match-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	personStore *store.PersonMemoryStore
}

func NewHandler(memoryStore *store.PersonMemoryStore) *Handler {
	return &Handler{
		personStore: memoryStore,
	}
}

func (h *Handler) AddSinglePersonAndMatch(w http.ResponseWriter, r *http.Request) {
	var personParams store.CreatePersonParams

	err := json.NewDecoder(r.Body).Decode(&personParams)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if personParams.WantedDates == 0 {
		http.Error(w, "WantedDates for the new person is zero", http.StatusBadRequest)
		return
	}

	newPerson, err := h.personStore.Create(personParams)

	fmt.Println("created new person", newPerson)

	people, err := h.personStore.GetAll()

	if err != nil {
		http.Error(w, "Can't get the person data", http.StatusBadRequest)
		return
	}

	var matches []models.Person
	for _, p := range people {
		maleMatchRule := newPerson.Gender == "male" && p.Gender == "female" && p.Height < newPerson.Height
		femaleMatchRule := newPerson.Gender == "female" && p.Gender == "male" && p.Height > newPerson.Height
		if maleMatchRule || femaleMatchRule {
			matches = append(matches, p)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

func (h *Handler) RemoveSinglePerson(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	err = h.personStore.Delete(id)
	if err != nil {
		http.Error(w, "Fail to delete the person", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Person with ID '%s' removed from the system", id)
}

func (h *Handler) QuerySinglePeople(w http.ResponseWriter, r *http.Request) {
	nParam := r.URL.Query().Get("n")
	n, err := strconv.Atoi(nParam)
	if err != nil {
		http.Error(w, "Invalid parameter for 'n'", http.StatusBadRequest)
		return
	}

	var possibleMatches []models.Person
	people, err := h.personStore.GetAll()

	for _, p := range people {
		if p.WantedDates > 0 {
			possibleMatches = append(possibleMatches, p)
		}
		if len(possibleMatches) >= n {
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(possibleMatches)
}
