package service

import (
	"github.com/eldelto/solvent"
	"github.com/google/uuid"
)

type Repository interface {
	Store(list *solvent.ToDoList) error
	Update(list *solvent.ToDoList) error
	Fetch(id uuid.UUID) (*solvent.ToDoList, error)
	FetchAll() []solvent.ToDoList
	BulkUpdate(lists []solvent.ToDoList) error
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

func (s *Service) FetchAll() []solvent.ToDoList {
	return s.repository.FetchAll()
}

func (s *Service) Update(list *solvent.ToDoList) (*solvent.ToDoList, error) {
	oldList, err := s.Fetch(list.ID)
	if err != nil {
		return nil, err
	}

	merged, err := list.Merge(oldList)
	if err != nil {
		return nil, err
	}
	mergedList := merged.(*solvent.ToDoList)

	err = s.repository.Update(mergedList)
	if err != nil {
		return nil, err
	}

	return mergedList, nil
}

func (s *Service) BulkUpdate(lists []solvent.ToDoList) ([]solvent.ToDoList, error) {
	updateList := make([]solvent.ToDoList, len(lists))
	for i, toDoList := range lists {
		oldList, err := s.Fetch(toDoList.ID)

		if err == nil {
			merged, err := toDoList.Merge(oldList)
			if err != nil {
				return nil, err
			}
			mergedList := merged.(*solvent.ToDoList)
			updateList[i] = *mergedList
		} else {
			updateList[i] = toDoList
		}
	}

	err := s.repository.BulkUpdate(updateList)
	if err != nil {
		return nil, err
	}

	return updateList, nil
}

// TODO: Handle archived ToDoLists
