package dto

import (
	"github.com/eldelto/solvent"
	"github.com/google/uuid"
)

type ToDoItemDto struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Checked    bool      `json:"checked"`
	OrderValue float64   `json:"orderValue"`
}

type ToDoItemDtoMap map[uuid.UUID]ToDoItemDto

type ToDoListDto struct {
	ID           uuid.UUID      `json:"id"`
	Title        string         `json:"title"`
	LiveSet      ToDoItemDtoMap `json:"liveSet"`
	TombstoneSet ToDoItemDtoMap `json:"tombstoneSet"`
}

func ToDoItemToDto(item solvent.ToDoItem) ToDoItemDto {
	return ToDoItemDto{
		ID:         item.ID,
		Title:      item.Title,
		Checked:    item.Checked,
		OrderValue: item.OrderValue,
	}
}

func ToDoItemFromDto(item ToDoItemDto) solvent.ToDoItem {
	return solvent.ToDoItem{
		ID:         item.ID,
		Title:      item.Title,
		Checked:    item.Checked,
		OrderValue: item.OrderValue,
	}
}

func toDoItemMapToDto(items solvent.ToDoItemMap) ToDoItemDtoMap {
	itemsDto := ToDoItemDtoMap{}
	for k, v := range items {
		itemsDto[k] = ToDoItemToDto(v)
	}

	return itemsDto
}

func toDoItemMapFromDto(itemsDto ToDoItemDtoMap) solvent.ToDoItemMap {
	items := solvent.ToDoItemMap{}
	for k, v := range itemsDto {
		items[k] = ToDoItemFromDto(v)
	}

	return items
}

func ToDoListToDto(list *solvent.ToDoList) *ToDoListDto {
	return &ToDoListDto{
		ID:           list.ID,
		Title:        list.Title,
		LiveSet:      toDoItemMapToDto(list.LiveSet),
		TombstoneSet: toDoItemMapToDto(list.TombstoneSet),
	}
}

func ToDoListFromDto(list *ToDoListDto) *solvent.ToDoList {
	return &solvent.ToDoList{
		ID:           list.ID,
		Title:        list.Title,
		LiveSet:      toDoItemMapFromDto(list.LiveSet),
		TombstoneSet: toDoItemMapFromDto(list.TombstoneSet),
	}
}
