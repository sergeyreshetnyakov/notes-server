package notes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

const url string = "http://localhost:8080"

func TestNotes(t *testing.T) {
	t.Run("[GET] notes", func(t *testing.T) {
		_, err := http.Get(url)
		if err != nil {
			t.Error(err.Error())
			t.Fail()
		}
	})

	req := bytes.NewBufferString(`{
		"header": "wash the basement",
		"content": "immediatly"
	}`)

	var addedNote struct {
		Id int `json:"id"`
	}

	t.Run("[ADD] note", func(t *testing.T) {
		res, err := http.Post(url, "application/json", req)
		if err != nil {
			t.Error(err.Error())
			t.Fail()
		}
		json.NewDecoder(res.Body).Decode(&addedNote)
	})
	t.Log(addedNote.Id)

	req = bytes.NewBufferString(`{
			"header": "immediatly",
			"content": "wash the basement",
			"id": `)

	req.Write([]byte(strconv.Itoa(addedNote.Id)))
	req.WriteString(`}`)

	t.Run("[PATCH] note", func(t *testing.T) {
		patchReq, err := http.NewRequest(http.MethodPatch, url, req)
		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}

		_, err = http.DefaultClient.Do(patchReq)
		if err != nil {
			t.Error(err.Error())
			t.Fail()
		}
	})

	t.Run("[DELETE] note", func(t *testing.T) {
		req := bytes.NewBufferString(`{
		"id":
	`)
		req.Write([]byte(strconv.Itoa(addedNote.Id)))
		req.WriteString(`}`)
		deleteReq, err := http.NewRequest(http.MethodDelete, url, req)
		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}
		http.DefaultClient.Do(deleteReq)
		if err != nil {
			t.Error(err.Error())
			t.Fail()
		}
	})
}
