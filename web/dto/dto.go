package dto

import (
	"github.com/eldelto/solvent"
	"github.com/google/uuid"
)

// TODO: Write custom Unmarshal functions to check for required fields

// ToDoItemDto is a DTO representing a ToDoItem as JSON"
type ToDoItemDto struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Checked    bool      `json:"checked"`
	OrderValue float64   `json:"orderValue"`
	UpdatedAt  int64     `json:"updatedAt"`
}

// ToDoListDto is a DTO representing a ToDoList as JSON"
type ToDoListDto struct {
	ID           uuid.UUID     `json:"id"`
	Title        string        `json:"title"`
	LiveSet      []ToDoItemDto `json:"liveSet"`
	TombstoneSet []ToDoItemDto `json:"tombstoneSet"`
	UpdatedAt    int64         `json:"updatedAt"`
	CreatedAt    int64         `json:"createdAt"`
}

// ToDoItemToDto converts a ToDoItem to its DTO representation
func ToDoItemToDto(item solvent.ToDoItem) ToDoItemDto {
	return ToDoItemDto{
		ID:         item.ID,
		Title:      item.Title,
		Checked:    item.Checked,
		OrderValue: item.OrderValue,
		UpdatedAt:  item.UpdatedAt,
	}
}

// ToDoItemFromDto converts a DTO representation to an actual ToDoItem
func ToDoItemFromDto(item ToDoItemDto) solvent.ToDoItem {
	return solvent.ToDoItem{
		ID:         item.ID,
		Title:      item.Title,
		Checked:    item.Checked,
		OrderValue: item.OrderValue,
		UpdatedAt:  item.UpdatedAt,
	}
}

func toDoItemMapToDto(items solvent.ToDoItemMap) []ToDoItemDto {
	itemsDto := make([]ToDoItemDto, len(items))

	i := 0
	for _, v := range items {
		itemsDto[i] = ToDoItemToDto(v)
		i++
	}

	return itemsDto
}

func toDoItemMapFromDto(itemsDto []ToDoItemDto) solvent.ToDoItemMap {
	items := solvent.ToDoItemMap{}
	for _, v := range itemsDto {
		items[v.ID] = ToDoItemFromDto(v)
	}

	return items
}

// ToDoListToDto converts a ToDoList to its DTO representation
func ToDoListToDto(list *solvent.ToDoList) ToDoListDto {
	return ToDoListDto{
		ID:           list.ID,
		Title:        list.Title,
		LiveSet:      toDoItemMapToDto(list.LiveSet),
		TombstoneSet: toDoItemMapToDto(list.TombstoneSet),
		UpdatedAt:    list.UpdatedAt,
		CreatedAt:    list.CreatedAt,
	}
}

// ToDoListFromDto converts a DTO representation to an actual ToDoList
func ToDoListFromDto(list *ToDoListDto) solvent.ToDoList {
	return solvent.ToDoList{
		ID:           list.ID,
		Title:        list.Title,
		LiveSet:      toDoItemMapFromDto(list.LiveSet),
		TombstoneSet: toDoItemMapFromDto(list.TombstoneSet),
		UpdatedAt:    list.UpdatedAt,
		CreatedAt:    list.CreatedAt,
	}
}
