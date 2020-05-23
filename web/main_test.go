package main

import (
	"fmt"
	"testing"

	. "github.com/eldelto/solvent/internal/testutils"
	"github.com/google/uuid"
)

const listTitle0 = "list0"
const listTitle1 = "list1"

const itemTitle0 = "item0"

func TestCreateToDoList(t *testing.T) {
	ts := NewTestServer(t, MainController.Handler)
	defer ts.Close()

	body := fmt.Sprintf(`{"title":"%s"}`, listTitle0)
	response := ts.POST("/api/to-do-list", body)
	AssertEquals(t, 200, response.StatusCode, "response.StatusCode")

	title := response.Body()["title"]
	AssertEquals(t, listTitle0, title, "title")

	id := response.Body()["id"].(string)
	_, err := uuid.Parse(id)
	AssertEquals(t, nil, err, "uuid.Parse error")
}

func TestFetchToDoLists(t *testing.T) {
	ts := NewTestServer(t, MainController.Handler)
	defer ts.Close()

	body := fmt.Sprintf(`{"title":"%s"}`, listTitle0)
	response := ts.POST("/api/to-do-list", body)
	rawID := response.Body()["id"].(string)
	id0, _ := uuid.Parse(rawID)

	body = fmt.Sprintf(`{"title":"%s"}`, listTitle1)
	response = ts.POST("/api/to-do-list", body)
	rawID = response.Body()["id"].(string)
	id1, _ := uuid.Parse(rawID)

	response = ts.GET("/api/to-do-list")

	type toDoListJSON struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	type responseJSON struct {
		ToDoLists []toDoListJSON `json:"toDoLists"`
	}

	var responseBody responseJSON
	err := response.Decode(&responseBody)
	AssertEquals(t, nil, err, "response.Decode error")

	expectedList0 := toDoListJSON{
		ID:    id0.String(),
		Title: listTitle0,
	}
	expectedList1 := toDoListJSON{
		ID:    id1.String(),
		Title: listTitle1,
	}

	untypedSlice := make([]interface{}, len(responseBody.ToDoLists))
	for i := range responseBody.ToDoLists {
		untypedSlice[i] = responseBody.ToDoLists[i]
	}

	AssertContains(t, expectedList0, untypedSlice, "responseBody.ToDoLists")
	AssertContains(t, expectedList1, untypedSlice, "responseBody.ToDoLists")
}

func TestFetchToDoList(t *testing.T) {
	ts := NewTestServer(t, MainController.Handler)
	defer ts.Close()

	body := fmt.Sprintf(`{"title":"%s"}`, listTitle0)
	response := ts.POST("/api/to-do-list", body)

	id := response.Body()["id"].(string)
	parsedID, _ := uuid.Parse(id)

	response = ts.GET("/api/to-do-list/" + parsedID.String())
	AssertEquals(t, 200, response.StatusCode, "response.StatusCode")

	title := response.Body()["title"]
	AssertEquals(t, listTitle0, title, "title")

	id = response.Body()["id"].(string)
	parsedID1, err := uuid.Parse(id)
	AssertEquals(t, nil, err, "uuid.Parse error")
	AssertEquals(t, parsedID, parsedID1, "UUID")
}

func TestUpdateToDoList(t *testing.T) {
	ts := NewTestServer(t, MainController.Handler)
	defer ts.Close()

	body := fmt.Sprintf(`{"title":"%s"}`, listTitle0)
	response := ts.POST("/api/to-do-list", body)

	id := response.Body()["id"].(string)
	parsedID, _ := uuid.Parse(id)

	body = fmt.Sprintf(`{
		"id": "%s",
		"title": "%s",
		"liveSet": [
			{
				"id": "%s",
				"title": "%s",
				"checked": false,
				"orderValue": 10.0
			}
		],
		"tombstoneSet": []
	}`, parsedID, listTitle0, parsedID, itemTitle0)

	response = ts.PUT("/api/to-do-list", body)
	AssertEquals(t, 200, response.StatusCode, "response.StatusCode")

	title := response.Body()["title"]
	AssertEquals(t, listTitle0, title, "title")

	id = response.Body()["id"].(string)
	parsedID1, err := uuid.Parse(id)
	AssertEquals(t, nil, err, "uuid.Parse error")
	AssertEquals(t, parsedID, parsedID1, "UUID")

	itemTitle := response.Body()["liveSet"].([]interface{})[0].(map[string]interface{})["title"].(string)
	AssertEquals(t, itemTitle0, itemTitle, "itemTitle")
}
