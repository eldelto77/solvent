package service

import (
	"github.com/eldelto/solvent"
	"github.com/google/uuid"
)

type Repository interface {
	Store(list *solvent.ToDoList) error
	Update(list *solvent.ToDoList) error
	Fetch(id uuid.UUID) (*solvent.ToDoList, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

// TODO: Wrap returned errors with custom ones
func (s *Service) Create(title string) (*solvent.ToDoList, error) {
	list, err := solvent.NewToDoList(title)
	if err != nil {
		return nil, err
	}

	err = s.repository.Store(&list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (s *Service) Fetch(id uuid.UUID) (*solvent.ToDoList, error) {
	return s.repository.Fetch(id)
}

func (s *Service) Update(list *solvent.ToDoList) (*solvent.ToDoList, error) {
	oldList, err := s.Fetch(list.ID)
	if err != nil {
		return nil, err
	}

	mergedList, err := list.Merge(oldList)
	if err != nil {
		return nil, err
	}

	err = s.repository.Update(&mergedList)
	if err != nil {
		return nil, err
	}

	return &mergedList, nil
}

// TODO: Handle archived ToDoLists