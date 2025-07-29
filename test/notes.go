package notes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

func TestNotes(t *testing.T) {
	url := "http://localhost:8080"

	t.Run("GET notes", func(t *testing.T) {
		_, err := http.Get(url)
		if err != nil {
			t.Error(err.Error())
		}
	})
	var addedNote struct {
		Id int64
	}
	t.Run("ADD note", func(t *testing.T) {
		var jsonBody = []byte(`{
				"header": "wash the basement",
				"content": "immediatly"
			}`)

		res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Error(err)
		}
		if err := json.NewDecoder(res.Body).Decode(&addedNote); err != nil {
			t.Error(err)
		}
	})
	t.Run("PATCH note", func(t *testing.T) {
		var jsonBody = bytes.NewBuffer([]byte(`{
				"header": "immediatly",
				"content": "wash the basement",
				"id": 
			}`))
		jsonBody.Write([]byte(strconv.Itoa(int(addedNote.Id))))
		req, err := http.NewRequest(http.MethodPatch, url, jsonBody)
		if err != nil {
			t.Error(err)
		}
		if _, err := http.DefaultClient.Do(req); err != nil {
			t.Error(err)
		}
	})

	t.Run("DELETE note", func(t *testing.T) {
		var jsonBody = bytes.NewBuffer([]byte(`{
			"id":
		}`))
		jsonBody.Write([]byte(strconv.Itoa(int(addedNote.Id))))

		req, err := http.NewRequest(http.MethodDelete, url, jsonBody)
		if err != nil {
			t.Error(err)
		}

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
		}
	})
}
