package solvent

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"sort"
	"testing"
)

const listTitle0 = "list0"

const itemTitle0 = "item0"
const itemTitle1 = "item1"
const itemTitle2 = "item2"

func TestNewToDoList(t *testing.T) {
	list, err := NewToDoList(listTitle0)

	assertEquals(t, nil, err, "NewToDoList error")
	assertEquals(t, listTitle0, list.Title, "list.Title")
	assertEquals(t, 0, len(list.liveSet), "list.liveSet length")
	assertEquals(t, 0, len(list.tombstoneSet), "list.tombstoneSet length")
}

func TestAddItem(t *testing.T) {
	list, _ := NewToDoList(listTitle0)

	id, err := list.AddItem(itemTitle0)
	assertEquals(t, nil, err, "list.AddItem error")

	item, err := list.GetItem(id)
	assertEquals(t, nil, err, "list.GetItem error")
	assertEquals(t, itemTitle0, item.Title, "item.Title")
	assertEquals(t, false, item.Checked, "item.Checked")
}

func TestRemoveItem(t *testing.T) {
	list, _ := NewToDoList(listTitle0)
	id, _ := list.AddItem(itemTitle0)

	list.RemoveItem(id)

	_, err := list.GetItem(id)
	expected := &NotFoundError{
		ID:      id,
		message: fmt.Sprintf("item with ID %v could not be found", id),
	}
	assertEquals(t, expected, err, "list.GetItem error")
}

func TestCheckItem(t *testing.T) {
	list, _ := NewToDoList(listTitle0)
	id, _ := list.AddItem(itemTitle0)

	id1, err := list.CheckItem(id)
	assertEquals(t, nil, err, "list.CheckItem error")
	assertEquals(t, id, id1, "list.CheckItem id")

	item, _ := list.GetItem(id1)
	assertEquals(t, itemTitle0, item.Title, "item.Title")
	assertEquals(t, true, item.Checked, "item.Checked")
}

func TestUncheckItem(t *testing.T) {
	list, _ := NewToDoList(listTitle0)
	id0, _ := list.AddItem(itemTitle0)
	id1, _ := list.CheckItem(id0)

	id2, err := list.UncheckItem(id0)
	assertEquals(t, nil, err, "list.UncheckItem error")
	assertNotEquals(t, id1, id2, "list.UncheckItem id")

	item, _ := list.GetItem(id2)
	assertEquals(t, itemTitle0, item.Title, "item.Title")
	assertEquals(t, false, item.Checked, "item.Checked")
}

func TestGetItems(t *testing.T) {
	list, _ := NewToDoList(listTitle0)
	id0, _ := list.AddItem(itemTitle0)
	id1, _ := list.AddItem(itemTitle1)

	items := orderedItems(&list)
	item0 := items[0]
	item1 := items[1]
	assertEquals(t, id0, item0.ID, "item0.ID")
	assertEquals(t, id1, item1.ID, "item1.ID")
}

func TestMoveItem(t *testing.T) {
	list, _ := NewToDoList(listTitle0)
	id0, _ := list.AddItem(itemTitle0)
	id1, _ := list.AddItem(itemTitle1)
	id2, _ := list.AddItem(itemTitle2)

	ids := itemIDs(orderedItems(&list))
	expected := []uuid.UUID{id0, id1, id2}
	assertEquals(t, expected, ids, "Initial item ordering")

	err := list.MoveItem(id2, 1)
	assertEquals(t, nil, err, "list.MoveItem error")
	ids = itemIDs(orderedItems(&list))
	expected = []uuid.UUID{id0, id2, id1}
	assertEquals(t, expected, ids, "First move item ordering")

	err = list.MoveItem(id2, -10)
	assertEquals(t, nil, err, "list.MoveItem error")
	ids = itemIDs(orderedItems(&list))
	expected = []uuid.UUID{id2, id0, id1}
	assertEquals(t, expected, ids, "Second move item ordering")

	err = list.MoveItem(id2, 10)
	assertEquals(t, nil, err, "list.MoveItem error")
	ids = itemIDs(orderedItems(&list))
	expected = []uuid.UUID{id0, id1, id2}
	assertEquals(t, expected, ids, "Third move item ordering")
}

func orderedItems(tdl *ToDoList) []ToDoItem {
	items := tdl.GetItems()
	sort.Slice(items, func(i, j int) bool { return items[i].OrderValue < items[j].OrderValue })

	return items
}

func itemIDs(list []ToDoItem) []uuid.UUID {
	ids := make([]uuid.UUID, len(list))
	for i, v := range list {
		ids[i] = v.ID
	}

	return ids
}

func assertEquals(t *testing.T, expected interface{}, actual interface{}, title string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v should be '%v' but was '%v'", title, expected, actual)
	}
}

func assertNotEquals(t *testing.T, expected interface{}, actual interface{}, title string) {
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("%v should not be '%v' but was '%v'", title, expected, actual)
	}
}
