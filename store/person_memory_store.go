package store

import (
	"fmt"
	"sync"
	"time"

	. "github.com/genepg/simple-match-api/models"
	"github.com/google/uuid"
)

type RecordNotFoundError struct{}

func (e *RecordNotFoundError) Error() string {
	return "record not found"
}

type DuplicateKeyError struct {
	Id uuid.UUID
}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("duplicate movie id: %v", e.Id)
}

type CreatePersonParams struct {
	Name        string `json:"name"`
	Height      int    `json:"height"`
	Gender      string `json:"gender"`
	WantedDates int    `json:"wanted_dates"`
}

type PersonMemoryStore struct {
	people map[uuid.UUID]Person
	mu     sync.RWMutex
}

func NewPersonMemoryStore() *PersonMemoryStore {
	return &PersonMemoryStore{
		people: map[uuid.UUID]Person{},
	}
}

func (s *PersonMemoryStore) GetAll() ([]Person, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var people []Person
	for _, m := range s.people {
		people = append(people, m)
	}
	return people, nil
}

func (s *PersonMemoryStore) GetByID(id uuid.UUID) (Person, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, ok := s.people[id]
	if !ok {
		return Person{}, &RecordNotFoundError{}
	}

	return m, nil
}

func (s *PersonMemoryStore) Create(createPersonParams CreatePersonParams) (Person, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newPersonId := uuid.New()

	if p, ok := s.people[newPersonId]; ok {
		return p, &DuplicateKeyError{Id: newPersonId}
	}

	person := Person{
		Id:          newPersonId,
		Name:        createPersonParams.Name,
		Height:      createPersonParams.Height,
		Gender:      createPersonParams.Gender,
		WantedDates: createPersonParams.WantedDates,
		CreatedAt:   time.Now().UTC(),
	}

	s.people[person.Id] = person
	return person, nil
}

func (s *PersonMemoryStore) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.people, id)
	return nil
}
