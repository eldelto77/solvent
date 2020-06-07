package persistence

import (
	"fmt"
	"sync"

	"github.com/eldelto/solvent"
	"github.com/google/uuid"
)

type NotebookStore map[uuid.UUID]solvent.Notebook

type InMemoryRepository struct {
	store NotebookStore
	mutex sync.Mutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		store: NotebookStore{},
		mutex: sync.Mutex{},
	}
}

func (r *InMemoryRepository) Store(notebook *solvent.Notebook) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.store[notebook.ID] = *notebook
	return nil
}

func (r *InMemoryRepository) Update(notebook *solvent.Notebook) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, ok := r.store[notebook.ID]
	if !ok {
		return fmt.Errorf("Notebook with ID '%v' could not be found", notebook.ID)
	}
	r.store[notebook.ID] = *notebook

	return nil
}

func (r *InMemoryRepository) Fetch(id uuid.UUID) (*solvent.Notebook, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	notebook, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("Notebook with ID '%v' could not be found", id)
	}

	return &notebook, nil
}

func (r *InMemoryRepository) Remove(id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.store, id)

	return nil
}
