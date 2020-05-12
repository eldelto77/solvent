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

// TODO: Test failure cases

func TestNewToDoList(t *testing.T) {
	list, err := NewToDoList(listTitle0)

	assertEquals(t, nil, err, "NewToDoList error")
	assertEquals(t, listTitle0, list.Title, "list.Title")
	assertEquals(t, 0, len(list.LiveSet), "list.LiveSet length")
	assertEquals(t, 0, len(list.TombstoneSet), "list.TombstoneSet length")
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
		message: fmt.Sprintf("item with ID '%v' could not be found", id),
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

func TestMerge(t *testing.T) {
	list0, _ := NewToDoList(listTitle0)
	_, _ = list0.AddItem(itemTitle0)
	id1, _ := list0.AddItem(itemTitle1)

	list1, _ := NewToDoList(listTitle0)
	list1.ID = list0.ID
	_, _ = list1.AddItem(itemTitle2)

	item1, _ := list0.GetItem(id1)
	item1.Checked = true
	item1.OrderValue = 5.0
	list1.LiveSet[id1] = item1

	mergedList, err := list0.Merge(&list1)
	assertEquals(t, nil, err, "list0.Merge error")

	// TODO: Handle equal sort order assigned from item creation
	/*ids := itemIDs(orderedItems(&mergedList))
	expected := []uuid.UUID{id1, id0, id2}
	assertEquals(t, expected, ids, "Item ordering")*/

	item1, _ = mergedList.GetItem(id1)
	assertEquals(t, 5.0, item1.OrderValue, "item0.OrderValue")
	assertEquals(t, true, item1.Checked, "item0.Checked")
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
