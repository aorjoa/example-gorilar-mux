package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

type insertStub struct{}

func TestCheckin_should_be_ok(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(&Check{ID: 1, PlaceID: 212})
	req := httptest.NewRequest("POST", "http://www.example.com/req", payload)
	w := httptest.NewRecorder()

	var fn InFunc = func(ID, placeID int64) error {
		return nil
	}

	resp := w.Result
	body, _ := ioutil.ReadAll(resp)
	fmt.Println(resp.StatusCode)
}
