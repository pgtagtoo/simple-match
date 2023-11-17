package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/genepg/simple-match-api/api"
	"github.com/genepg/simple-match-api/models"
	"github.com/genepg/simple-match-api/store"

	"github.com/stretchr/testify/require"
)

func init() {
	fmt.Println("test is starting!!!!")
}

func TestRoot(t *testing.T) {
	s := api.CreateNewServer()
	s.MountHandlers()

	// Create a New Request
	req, _ := http.NewRequest("GET", "/", nil)

	// Execute Request
	response := executeRequest(req, s)

	if http.StatusOK != response.Code {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	// We can use testify/require to assert values, as it is more convenient
	require.Equal(t, "This is simple match", response.Body.String())
}

func TestAddSinglePersonAndMatch(t *testing.T) {
	s := api.CreateNewServer()
	s.MountHandlers()

	// Create a user
	newPerson := store.CreatePersonParams{
		Name:        "Mandy",
		Height:      155,
		Gender:      "female",
		WantedDates: 1,
	}
	_, err := s.PersonStore.Create(newPerson)

	// create a request for creating a person and find matches
	payload := `{
		"name": "John",
		"height": 175,
		"gender": "male",
		"wanted_dates": 2
	}`

	req, err := http.NewRequest("POST", "/addPersonAndMatch", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(req, s)

	// Check the status code
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("Response body: %s", response.Body.String())
	}

	// Check the matches
	var matches []models.Person
	err = json.NewDecoder(response.Body).Decode(&matches)
	if err != nil {
		t.Fatal(err)
	}
	require.Len(t, matches, 1)
}

func TestRemoveSinglePerson(t *testing.T) {
	s := api.CreateNewServer()
	s.MountHandlers()

	// Create a user
	newPerson := store.CreatePersonParams{
		Name:        "Sam",
		Height:      180,
		Gender:      "male",
		WantedDates: 3,
	}
	createdNewPerson, err := s.PersonStore.Create(newPerson)

	// Create a request with a specific ID to delete
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/removePerson/%s", createdNewPerson.Id), nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(req, s)

	// Check the status code
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("Response body: %s", response.Body.String())
	}
}

func TestQuerySinglePeople(t *testing.T) {
	s := api.CreateNewServer()
	s.MountHandlers()

	// Create a user
	newPerson := store.CreatePersonParams{
		Name:        "Sandy",
		Height:      160,
		Gender:      "female",
		WantedDates: 3,
	}
	_, err := s.PersonStore.Create(newPerson)

	// Create a request for querying people
	req, err := http.NewRequest("GET", "/queryPeople?n=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(req, s)

	// Check the status code
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("Response body: %s", response.Body.String())
	}

	// Check the length of people
	var people []models.Person
	err = json.NewDecoder(response.Body).Decode(&people)
	if err != nil {
		t.Fatal(err)
	}
	require.Len(t, people, 1)
}

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *api.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}
