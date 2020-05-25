package solvent

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

type ToDoItem struct {
	ID         uuid.UUID
	Title      string
	Checked    bool
	OrderValue float64
	UpdatedAt  int64
}

type ToDoItemMap map[uuid.UUID]ToDoItem

type ToDoList struct {
	ID           uuid.UUID
	Title        string
	LiveSet      ToDoItemMap
	TombstoneSet ToDoItemMap
	UpdatedAt    int64
	CreatedAt    int64
}

func NewToDoList(title string) (ToDoList, error) {
	// TODO: validate input string
	id, err := randomUUID()
	if err != nil {
		return ToDoList{}, err
	}

	toDoList := ToDoList{
		ID:           id,
		Title:        title,
		LiveSet:      ToDoItemMap{},
		TombstoneSet: ToDoItemMap{},
		UpdatedAt:    time.Now().Local().UnixNano(),
		CreatedAt:    time.Now().Local().UnixNano(),
	}

	return toDoList, nil
}

// TODO: Use types for ToDoListID and ToDoItemID
func (tdl *ToDoList) Rename(title string) (uuid.UUID, error) {
	tdl.Title = title
	tdl.UpdatedAt = time.Now().UnixNano()

	return tdl.ID, nil
}

func (tdl *ToDoList) AddItem(title string) (uuid.UUID, error) {
	// TODO: validate input string

	id, err := randomUUID()
	if err != nil {
		return uuid.Nil, err
	}

	item := ToDoItem{
		ID:         id,
		Title:      title,
		Checked:    false,
		OrderValue: tdl.nextOrderValue(),
		UpdatedAt:  time.Now().UnixNano(),
	}
	tdl.LiveSet[id] = item

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
		tdl.TombstoneSet[id] = item
	}
}

func (tdl *ToDoList) CheckItem(id uuid.UUID) (uuid.UUID, error) {
	item, err := tdl.GetItem(id)
	if err == nil {
		item.Checked = true
		tdl.LiveSet[id] = item
	}

	return item.ID, err
}

func (tdl *ToDoList) UncheckItem(id uuid.UUID) (uuid.UUID, error) {
	item, err := tdl.GetItem(id)
	if err != nil {
		return uuid.Nil, newNotFoundError(id)
	}
	tdl.RemoveItem(item.ID)

	newID, err := randomUUID()
	if err != nil {
		return newID, err
	}
	newItem := ToDoItem{
		ID:         newID,
		Title:      item.Title,
		Checked:    false,
		OrderValue: item.OrderValue,
		UpdatedAt:  item.UpdatedAt,
	}
	tdl.LiveSet[newID] = newItem

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
	var orderValueAdjacent float64
	if orderValueMid < item.OrderValue {
		// Moving item up
		if (index - 1) >= 0 {
			orderValueAdjacent = items[index-1].OrderValue
		} else {
			orderValueAdjacent = 0.0
		}
	} else if orderValueMid > item.OrderValue {
		// Moving item down
		if (index + 1) < len(items) {
			orderValueAdjacent = items[index+1].OrderValue
		} else {
			orderValueAdjacent = tdl.nextOrderValue()
		}
	} else {
		// Already on correct position
		return nil
	}

	newOrderValue := (orderValueMid + orderValueAdjacent) / 2
	item.OrderValue = newOrderValue
	tdl.LiveSet[item.ID] = item

	return nil
}

func (tdl *ToDoList) Merge(other *ToDoList) (ToDoList, error) {
	if tdl.ID != other.ID {
		return ToDoList{}, newCannotBeMergedError(tdl.ID, other.ID)
	}

	var updatedAt int64
	var title string
	if other.UpdatedAt > tdl.UpdatedAt {
		updatedAt = other.UpdatedAt
		title = other.Title
	} else {
		updatedAt = tdl.UpdatedAt
		title = tdl.Title
	}

	mergedLiveSet, err := mergeToDoItemMaps(tdl.LiveSet, other.LiveSet)
	if err != nil {
		return ToDoList{}, err
	}

	mergedTombstoneSet, err := mergeToDoItemMaps(tdl.TombstoneSet, other.TombstoneSet)
	if err != nil {
		return ToDoList{}, err
	}

	mergedToDoList := ToDoList{
		ID:           tdl.ID,
		Title:        title,
		LiveSet:      mergedLiveSet,
		TombstoneSet: mergedTombstoneSet,
		UpdatedAt:    updatedAt,
		CreatedAt:    tdl.CreatedAt,
	}
	return mergedToDoList, nil
}

func randomUUID() (uuid.UUID, error) {
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

	for key, value := range tdl.LiveSet {
		_, deleted := tdl.TombstoneSet[key]
		if !deleted {
			liveView[key] = value
		}
	}

	return liveView
}

func (tdl *ToDoList) nextOrderValue() float64 {
	orderValue := 0.0
	for _, item := range tdl.liveView() {
		if item.OrderValue > orderValue {
			orderValue = item.OrderValue
		}
	}

	return orderValue + 10
}

func mergeToDoItems(this, other ToDoItem) (ToDoItem, error) {
	if this.ID != other.ID {
		return this, newCannotBeMergedError(this.ID, other.ID)
	}

	if other.Checked {
		this.Checked = true
	}

	if other.UpdatedAt > this.UpdatedAt {
		this.OrderValue = other.OrderValue
		this.UpdatedAt = other.UpdatedAt
	}

	return this, nil
}

func mergeToDoItemMaps(thisMap, otherMap ToDoItemMap) (ToDoItemMap, error) {
	mergedMap := ToDoItemMap{}
	for k, v := range thisMap {
		mergedMap[k] = v
	}
	for k, otherItem := range otherMap {
		thisItem, ok := thisMap[k]
		if ok {
			mergedItem, err := mergeToDoItems(thisItem, otherItem)
			if err != nil {
				return ToDoItemMap{}, newCannotBeMergedError(thisItem.ID, otherItem.ID)
			}
			mergedMap[k] = mergedItem
		}
	}

	return mergedMap, nil
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
		message: fmt.Sprintf("item with ID '%v' could not be found", id),
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
	thisID  uuid.UUID
	otherID uuid.UUID
	message string
}

func newCannotBeMergedError(thisID, otherID uuid.UUID) *CannotBeMergedError {
	return &CannotBeMergedError{
		thisID:  thisID,
		otherID: otherID,
		message: fmt.Sprintf("item with ID '%v' cannot be merged with item with ID '%v'", thisID, otherID),
	}
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
