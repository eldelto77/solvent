package dto

import (
	"github.com/eldelto/solvent"
	"github.com/google/uuid"
)

// TODO: Write custom Unmarshal functions to check for required fields

type OrderValueDto struct {
	Value     float64 `json:"value"`
	UpdatedAt int64   `json:"updatedAt"`
}

func orderValueToDto(orderValue solvent.OrderValue) OrderValueDto {
	return OrderValueDto{
		Value:     orderValue.Value,
		UpdatedAt: orderValue.UpdatedAt,
	}
}

func orderValueFromDto(orderValue OrderValueDto) solvent.OrderValue {
	return solvent.OrderValue{
		Value:     orderValue.Value,
		UpdatedAt: orderValue.UpdatedAt,
	}
}

type TitleDto struct {
	Value     string `json:"value"`
	UpdatedAt int64  `json:"updatedAt"`
}

func titleToDto(title solvent.Title) TitleDto {
	return TitleDto{
		Value:     title.Value,
		UpdatedAt: title.UpdatedAt,
	}
}

func titleFromDto(title TitleDto) solvent.Title {
	return solvent.Title{
		Value:     title.Value,
		UpdatedAt: title.UpdatedAt,
	}
}

// ToDoItemDto is a DTO representing a ToDoItem as JSON"
type ToDoItemDto struct {
	ID         uuid.UUID     `json:"id"`
	Title      string        `json:"title"`
	Checked    bool          `json:"checked"`
	OrderValue OrderValueDto `json:"orderValue"`
}

// ToDoItemToDto converts a ToDoItem to its DTO representation
func ToDoItemToDto(item solvent.ToDoItem) ToDoItemDto {
	return ToDoItemDto{
		ID:         item.ID,
		Title:      item.Title,
		Checked:    item.Checked,
		OrderValue: orderValueToDto(item.OrderValue),
	}
}

// ToDoItemFromDto converts a DTO representation to an actual ToDoItem
func ToDoItemFromDto(item ToDoItemDto) solvent.ToDoItem {
	return solvent.ToDoItem{
		ID:         item.ID,
		Title:      item.Title,
		Checked:    item.Checked,
		OrderValue: orderValueFromDto(item.OrderValue),
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

// ToDoListDto is a DTO representing a ToDoList as JSON"
type ToDoListDto struct {
	ID           uuid.UUID     `json:"id"`
	Title        TitleDto      `json:"title"`
	LiveSet      []ToDoItemDto `json:"liveSet"`
	TombstoneSet []ToDoItemDto `json:"tombstoneSet"`
	UpdatedAt    int64         `json:"updatedAt"`
	CreatedAt    int64         `json:"createdAt"`
}

// ToDoListToDto converts a ToDoList to its DTO representation
func ToDoListToDto(list *solvent.ToDoList) ToDoListDto {
	return ToDoListDto{
		ID:           list.ID,
		Title:        titleToDto(list.Title),
		LiveSet:      toDoItemMapToDto(list.LiveSet),
		TombstoneSet: toDoItemMapToDto(list.TombstoneSet),
		CreatedAt:    list.CreatedAt,
	}
}

// ToDoListFromDto converts a DTO representation to an actual ToDoList
func ToDoListFromDto(list *ToDoListDto) solvent.ToDoList {
	return solvent.ToDoList{
		ID:           list.ID,
		Title:        titleFromDto(list.Title),
		LiveSet:      toDoItemMapFromDto(list.LiveSet),
		TombstoneSet: toDoItemMapFromDto(list.TombstoneSet),
		CreatedAt:    list.CreatedAt,
	}
}
