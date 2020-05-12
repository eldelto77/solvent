package persistence

import (
	"fmt"

	"github.com/eldelto/solvent"
	"github.com/google/uuid"
)

type ToDoListStore  map[uuid.UUID]solvent.ToDoList

type Repository struct {
	store ToDoListStore
}

func NewRepository() Repository {
	return Repository{
		store: ToDoListStore{},
	}
}

func (r *Repository) Store(list *solvent.ToDoList) error {
	r.store[list.ID] = *list
	return nil
}

func (r *Repository) Update(list *solvent.ToDoList) error {
	_, ok := r.store[list.ID]
	if !ok {
		return fmt.Errorf("ToDoList with ID '%v' could not be found", list.ID)
	}
	r.store[list.ID] = *list

	return nil
}

func (r *Repository) Fetch(id uuid.UUID) (*solvent.ToDoList, error) {
	list, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("ToDoList with ID '%v' could not be found", id)
	}

	return &list, nil
}