package solvent

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
)

type ToDoItem struct {
	ID         uuid.UUID
	Title      string
	Checked    bool
	OrderValue float64
}

type ToDoItemMap map[uuid.UUID]ToDoItem

type ToDoList struct {
	ID           uuid.UUID
	Title        string
	liveSet      ToDoItemMap
	tombstoneSet ToDoItemMap
}

func NewToDoList(title string) (ToDoList, error) {
	// TODO: validate input string
	id, err := randomUuid()
	if err != nil {
		return ToDoList{}, err
	}

	toDoList := ToDoList{
		ID:           id,
		Title:        title,
		liveSet:      ToDoItemMap{},
		tombstoneSet: ToDoItemMap{},
	}

	return toDoList, nil
}

func (tdl *ToDoList) AddItem(title string) (uuid.UUID, error) {
	// TODO: validate input string

	id, err := randomUuid()
	if err != nil {
		return uuid.Nil, err
	}

	item := ToDoItem{
		ID:         id,
		Title:      title,
		Checked:    false,
		OrderValue: tdl.nextOrderValue(),
	}
	tdl.liveSet[id] = item

	return id, nil
}

func (tdl *ToDoList) GetItem(id uuid.UUID) (ToDoItem, error) {
	item, ok := tdl.liveView()[id]
	if ok == false {
		return item, newNotFoundError(id)
	}

	return item, nil
}

func (tdl *ToDoList) RemoveItem(id uuid.UUID) {
	if item, ok := tdl.liveView()[id]; ok {
		tdl.tombstoneSet[id] = item
	}
}

func (tdl *ToDoList) CheckItem(id uuid.UUID) (uuid.UUID, error) {
	item, err := tdl.GetItem(id)
	if err == nil {
		item.Checked = true
		tdl.liveSet[id] = item
	}

	return item.ID, err
}

func (tdl *ToDoList) UncheckItem(id uuid.UUID) (uuid.UUID, error) {
	item, err := tdl.GetItem(id)
	if err != nil {
		return uuid.Nil, newNotFoundError(id)
	}
	tdl.RemoveItem(item.ID)

	newID, err := randomUuid()
	if err != nil {
		return newID, err
	}
	newItem := ToDoItem{
		ID:         newID,
		Title:      item.Title,
		Checked:    false,
		OrderValue: item.OrderValue,
	}
	tdl.liveSet[newID] = newItem

	return newID, nil
}

func (tdl *ToDoList) GetItems() []ToDoItem {
	// TODO: Benchmark pre-allocation
	items := []ToDoItem{}
	for _, item := range tdl.liveView() {
		items = append(items, item)
	}

	return items
}

func (tdl *ToDoList) MoveItem(id uuid.UUID, targetIndex int) error {
	item, err := tdl.GetItem(id)
	if err != nil {
		return err
	}

	items := tdl.GetItems()
	index := clampIndex(targetIndex, items)
	sort.Slice(items, func(i, j int) bool { return items[i].OrderValue < items[j].OrderValue })

	orderValueMid := items[index].OrderValue
	orderValueAdjacent := orderValueMid
	if orderValueMid < item.OrderValue {
		// Moving item up
		orderValueAdjacent = 0.0
		if (index - 1) >= 0 {
			orderValueAdjacent = items[index-1].OrderValue
		}
	} else if orderValueMid > item.OrderValue {
		// Moving item down
		orderValueAdjacent = tdl.nextOrderValue()
		if (index + 1) < len(items) {
			orderValueAdjacent = items[index+1].OrderValue
		}
	} else {
		// Already on correct position
		return nil
	}

	newOrderValue := (orderValueMid + orderValueAdjacent) / 2
	item.OrderValue = newOrderValue
	tdl.liveSet[item.ID] = item

	return nil
}

func (tdl *ToDoList) Merge(other *ToDoList) (ToDoList, error) {
	return ToDoList{}, nil
}

func randomUuid() (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		err = &UnknownError{
			message: "item creation failed with nested error",
			err:     err,
		}
	}

	return id, err
}

func (tdl *ToDoList) liveView() ToDoItemMap {
	// TODO: Pass expected length?
	liveView := ToDoItemMap{}

	for key, value := range tdl.liveSet {
		_, deleted := tdl.tombstoneSet[key]
		if !deleted {
			liveView[key] = value
		}
	}

	return liveView
}

func (tdl *ToDoList) nextOrderValue() float64 {
	orderValue := 10.0
	for _, item := range tdl.liveView() {
		if item.OrderValue > orderValue {
			orderValue = item.OrderValue
		}
	}

	return orderValue + 10
}

func clampIndex(index int, list []ToDoItem) int {
	max := len(list) - 1
	if index < 0 {
		return 0
	} else if index > max {
		return max
	} else {
		return index
	}
}

// Errors

type NotFoundError struct {
	ID      uuid.UUID
	message string
}

func newNotFoundError(id uuid.UUID) *NotFoundError {
	return &NotFoundError{
		ID:      id,
		message: fmt.Sprintf("item with ID %v could not be found", id),
	}
}

func (e *NotFoundError) Error() string {
	return e.message
}

type InvalidTitleError struct {
	title   string
	message string
}

func (e *InvalidTitleError) Error() string {
	return e.message
}

type CannotBeMergedError struct {
	ID      uuid.UUID
	message string
}

func (e *CannotBeMergedError) Error() string {
	return e.message
}

type UnknownError struct {
	err     error
	message string
}

func (e *UnknownError) Error() string {
	return e.message
}

func (e *UnknownError) Unwrap() error {
	return e.err
}
